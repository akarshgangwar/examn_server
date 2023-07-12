package handlers

import (
	// "fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
var jwtKey = []byte("ExamnProject")
type Claims struct {
	Username string `json:"username"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func LoginHandler(c *gin.Context,db *gorm.DB){
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Username == "" && req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Email is required"})
		return
	}

	// if creadential is correct
	accessToken, refreshToken, err := GenerateTokens(req.Email,req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}
func GenerateTokens(username string,email string) (string, string, error) {
	accessExpirationTime := time.Now().Add(15 * time.Minute)
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Generate access token
	accessClaims := Claims{
		Username: username,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpirationTime.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshClaims := &Claims{
		Username: username,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpirationTime.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
