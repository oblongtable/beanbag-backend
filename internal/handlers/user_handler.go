package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oblongtable/beanbag-backend/internal/services"
	"github.com/oblongtable/beanbag-backend/middleware"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUserRequest represents the request body for creating a user.
// @Description User details
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the given details
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User details"
// @Success 201 {object} db.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(ctx, req.Name, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} db.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserById(ctx, int32(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Define the request struct specifically for the sync endpoint
type SyncUserRequest struct {
	Name  string `json:"name"` // Name might not always be present from Auth0, handle potential empty string
	Email string `json:"email" binding:"required,email"`
}

// SyncUser godoc
// @Summary Sync Auth0 user with backend database via Email
// @Description Creates a user if they don't exist based on email, or updates existing user's name.
// @Tags users
// @Accept json
// @Produce json
// @Param user body SyncUserRequest true "User details (name, email) from Auth0"
// @Success 200 {object} db.User "User found and potentially updated"
// @Success 201 {object} db.User "User created"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized (JWT missing/invalid)"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/sync [post]
// @Security BearerAuth
func (h *UserHandler) SyncUser(ctx *gin.Context) {
	var req SyncUserRequest

	// --- Retrieve authenticated user's email from Gin context ---
	// Use the key defined in the middleware package
	jwtEmail := ctx.GetString(middleware.GinContextKeyUserEmail)

	// The ExtractAndSetClaims middleware already checked for existence and emptiness
	// If we reach here, jwtEmail should be valid and non-empty.

	// Bind the request body JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// --- Compare JWT email with request body email (Case-Insensitive) ---
	if strings.ToLower(jwtEmail) != strings.ToLower(req.Email) {
		log.Printf("Forbidden SyncUser attempt: JWT email '%s' != Request email '%s'", jwtEmail, req.Email)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Request email does not match authenticated user"})
		return
	}
	// --- End email comparison ---

	// Basic validation (email format is already validated by binding)
	if req.Email == "" { // Should be caught by binding:"required" but double-check
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	// Use email as name if name is empty
	if req.Name == "" {
		req.Name = req.Email
	}

	// Emails match, proceed with the service call
	// Pass the validated email from the JWT to the service layer
	user, created, err := h.userService.SyncUser(ctx.Request.Context(), req.Name, jwtEmail) // Use jwtEmail
	if err != nil {
		// Log the internal error for debugging
		log.Printf("Error syncing user %s: %v", jwtEmail, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync user data"}) // Keep error message generic for client
		return
	}

	// Return different status codes based on whether user was created or found/updated
	if created {
		ctx.JSON(http.StatusCreated, user) // 201 Created
	} else {
		ctx.JSON(http.StatusOK, user) // 200 OK
	}
}