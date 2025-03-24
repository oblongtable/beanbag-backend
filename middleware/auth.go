package middleware

import (
	"net/http"

	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	"github.com/oblongtable/beanbag-backend/initializers"
	"gopkg.in/square/go-jose.v2"
)

func AuthMiddlewareAuth0(ctx *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		config := initializers.GetConfig()

		var auth0Domain = "https://" + config.AuthDomain + "/"
		client := auth0.NewJWKClient(
			auth0.JWKClientOptions{
				URI: auth0Domain + ".well-known/jwks.json",
			},
			nil,
		)
		configuration := auth0.NewConfiguration(
			client,
			[]string{config.AuthAudience},
			auth0Domain,
			jose.RS256,
		)
		validator := auth0.NewValidator(configuration, nil)
		_, err := validator.ValidateRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token",
			})
			c.Abort()
			return
		}

		ctx.Next()
	}
}
