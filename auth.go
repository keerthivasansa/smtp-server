package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func getAuthBearerToken(authHeader string) string {
	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

type AuthMiddleware struct {
	passwordHash []byte
}

func createAuthMiddleware() AuthMiddleware {
	password := os.Getenv("API_TOKEN")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		panic(err)
	}
	return AuthMiddleware{
		passwordHash: passwordHash,
	}
}

func (am AuthMiddleware) middleware(c *gin.Context) {
	bearerH := c.Request.Header.Get("Authorization")
	token := getAuthBearerToken(bearerH)
	if token == "" {
		c.String(http.StatusUnauthorized, "Missing bearer token.")
		return
	}
	err := bcrypt.CompareHashAndPassword(am.passwordHash, []byte(token))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		c.String(http.StatusUnauthorized, "Invalid bearer token.")
		return
	} else if err != nil {
		println(err.Error())
		c.String(http.StatusInternalServerError, "Unknown error while checking authorization")
		return
	}
	c.Next()
}
