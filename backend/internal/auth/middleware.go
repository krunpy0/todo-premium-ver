package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func Auth(c *gin.Context) {
	godotenv.Load()
	token, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	claims := &Claims{}
	strToken, err := jwt.ParseWithClaims(token, claims, func (token *jwt.Token) (any, error)  {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if err != nil || !strToken.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	userID :=claims.UserID
	fmt.Println(userID)
	c.Set("userID", int(userID))
	c.Next()
}