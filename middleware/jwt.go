package middleware

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/oblongtable/beanbag-backend/initializers"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
	Email string `json:"email"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func VerifyToken() func(next http.Handler) http.Handler {
	config := initializers.GetConfig()
	issuerURL, err := url.Parse("https://" + config.AuthDomain + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{config.AuthAudience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Access Denied"}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}

const (
	GinContextKeyUserEmail = "user_email"
	GinContextKeyUserSub   = "user_sub"
)

// ExtractAndSetClaims is a Gin middleware that should run *after*
// the JWT validation middleware (like VerifyToken adapted).
// It extracts claims set by go-jwt-middleware/v2 into the request context
// and puts specific ones (email, sub) into the Gin context.
func ExtractAndSetClaims() gin.HandlerFunc {
	return func(c *gin.Context) {
		// go-jwt-middleware/v2 stores the validated claims in the request context
		token := c.Request.Context().Value(jwtmiddleware.ContextKey{})
		
		if token == nil {
			// This should ideally be caught by VerifyToken, but check defensively
			log.Println("Middleware Error: No token found in request context after VerifyToken")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not found in context"})
			return
		}

		validatedClaims, ok := token.(*validator.ValidatedClaims)
		if !ok {
			log.Printf("Middleware Error: Token in context is not *validator.ValidatedClaims: %T", token)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid token format in context"})
			return
		}

		// --- Extract Subject (sub) ---
		sub := validatedClaims.RegisteredClaims.Subject
		if sub == "" {
			log.Println("Middleware Error: Subject (sub) claim missing in validated token")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Token missing required user identifier"})
			return
		}

		// --- Extract Email from CustomClaims ---
		// Cast the CustomClaims interface{} to our specific *CustomClaims type
		customClaims, ok := validatedClaims.CustomClaims.(*CustomClaims)
		if !ok || customClaims == nil {
			log.Printf("Middleware Error: CustomClaims in context is not *CustomClaims or is nil")
			// This might happen if the token validation succeeded but custom claims parsing failed,
			// or if the CustomClaims field wasn't populated correctly by the validator.
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not process custom claims from token"})
			return
		}

		email := customClaims.Email // Access the Email field we added
		if email == "" {
			log.Println("Middleware Error: Email claim missing or empty in validated token's custom claims")
			// This is added by a custom action post-login in our Auth0 setup
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Token missing required user email information"})
			return
		}

		// --- Set values in Gin context for the handler ---
		c.Set(GinContextKeyUserEmail, email)
		c.Set(GinContextKeyUserSub, sub)
		log.Printf("Middleware: Set user_email=%s, user_sub=%s in Gin context", email, sub) // Debug log

		c.Next() // Proceed to the next handler (e.g., the actual API handler)
	}
}


// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}
