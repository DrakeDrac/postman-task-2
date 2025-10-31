package users

import (
	"postman-task/internal/auth"
	"postman-task/internal/core"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db  *gorm.DB
	jwt *auth.JWTManager
}

// Creates a new user handler
func NewUserHandler(db *gorm.DB, jwt *auth.JWTManager) *UserHandler {
	return &UserHandler{
		db:  db,
		jwt: jwt,
	}
}

// Register a new user
func (h *UserHandler) Register(c *gin.Context) {
	// Get data from request
	var data core.RegisterRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}
	// Checking for admin
	if data.Role == "admin" {
		c.JSON(403, gin.H{"error": "Admin user cannot be registered via API"})
		return
	}

	// If role registered is faculty or warden, ensure requester is an admin
	if data.Role == "faculty" || data.Role == "warden" {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(403, gin.H{"error": "Admin token required to register faculty/warden"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		claims, err := h.jwt.ValidateToken(tokenParts[1])
		if err != nil || claims.Role != "admin" {
			c.JSON(403, gin.H{"error": "Only admin can register faculty/warden"})
			return
		}
	}

	// Check if user exists
	var existingUser core.User
	h.db.Where("email = ?", data.Email).First(&existingUser)
	if existingUser.ID != 0 {
		c.JSON(400, gin.H{"error": "Email already in use"})
		return
	}

	// Hash password
	hash, err := auth.HashPassword(data.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}

	// Create user
	user := core.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: hash,
		Role:     data.Role,
		Dept:     data.Dept,
	}

	// Save to database
	result := h.db.Create(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Could not create user"})
		return
	}

	// Return success
	c.JSON(200, gin.H{
		"message": "User created",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// Login user
func (h *UserHandler) Login(c *gin.Context) {
	// Get login data
	var data core.LoginRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	// Find user
	var user core.User
	result := h.db.Where("email = ?", data.Email).First(&user)
	if result.Error != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !auth.CheckPasswordHash(data.Password, user.Password) {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := h.jwt.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}

	// Return token
	c.JSON(200, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Get all users, admin only
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Get pagination parameters
	page := 1
	p := c.Query("page")
	if p != "" {
		pn, err := strconv.ParseInt(p, 10, 32) // parse int in base 10
		if err == nil && pn > 0 {
			page = int(pn)
		}
	}
	pageSize := 10

	// Get total count
	var total int64
	h.db.Model(&core.User{}).Count(&total)

	// Get page of users
	var users []core.User
	offset := (page - 1) * pageSize
	h.db.Offset(offset).
		Limit(pageSize).
		Find(&users)

	// Remove passwords from data
	for i := range users {
		users[i].Password = ""
	}

	// Return paginated result
	result := core.PageResult{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Items:    users,
	}

	c.JSON(200, result)
}

// Get user by ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Get user ID from URL
	id := c.Param("id")

	// Find user
	var user core.User
	result := h.db.First(&user, id)

	if result.Error != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Remove password from data
	user.Password = ""

	c.JSON(200, user)
}
