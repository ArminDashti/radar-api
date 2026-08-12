package handlers

import (
	"net/http"

	radarauth "github.com/ArminDashti/radar-api/internal/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}
	var passwordHash string
	if err := s.Pool.QueryRow(requestContext(c),
		`SELECT password_hash FROM users WHERE username = $1`, input.Username,
	).Scan(&passwordHash); err != nil ||
		bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	token, err := radarauth.IssueToken(s.JWTSecret, input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "username": input.Username})
}
