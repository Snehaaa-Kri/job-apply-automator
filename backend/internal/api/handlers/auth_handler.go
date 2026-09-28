package handlers

import (
	"github.com/Snehaaa-Kri/job-apply-automator/internal/service"
	"github.com/Snehaaa-Kri/job-apply-automator/pkg/response"
	"github.com/Snehaaa-Kri/job-apply-automator/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler holds the auth service dependency.
// Node analogy: class AuthController { constructor(private authService: AuthService) {} }
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler is the constructor.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/v1/auth/register
// Node analogy: authController.register = async (req, res) => { ... }
func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput

	// 1. Parse JSON body into struct
	// Node analogy: const body = req.body (but with type checking)
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	// 2. Validate the struct fields
	// Node analogy: const { error } = schema.validate(body)
	if err := validator.Validate.Struct(input); err != nil {
		// Return field-level errors: { "Email": "Email is required" }
		c.JSON(400, gin.H{
			"success": false,
			"errors":  validator.FormatErrors(err),
		})
		return
	}

	// 3. Call service
	result, err := h.authService.Register(c.Request.Context(), input)
	if err != nil {
		// Specific known errors get 400, unexpected ones get 500
		if err.Error() == "email already in use" {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "registration failed")
		return
	}

	// 4. Send response
	response.Created(c, "registration successful", result)
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := validator.Validate.Struct(input); err != nil {
		c.JSON(400, gin.H{"success": false, "errors": validator.FormatErrors(err)})
		return
	}

	result, err := h.authService.Login(c.Request.Context(), input)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, "login successful", result)
}

// GetMe handles GET /api/v1/auth/me (protected route)
// The userID was set by AuthRequired middleware
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Read userID from context — set by JWT middleware
	// Node analogy: req.user.userId
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c, "invalid user identity")
		return
	}

	user, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "user fetched", user)
}
