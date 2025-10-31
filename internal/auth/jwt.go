package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTManager struct {
	secretKey string
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

// Creates a JWT manager
func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
	}
}

// Creates a new JWT token
func (j *JWTManager) GenerateToken(userID uint, email, role string) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	// Sign the token
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validates a token and returns claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// Convert claims to our Claims type
		return &Claims{
			UserID: uint(claims["user_id"].(float64)),
			Email:  claims["email"].(string),
			Role:   claims["role"].(string),
		}, nil
	}

	return nil, errors.New("invalid token")
}

// Hashes a password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compares password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
