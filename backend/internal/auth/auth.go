package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/krunpy0/todo-premium-ver/internal/user"
	"golang.org/x/crypto/bcrypt"
)



func Register (c *gin.Context) {
	var userStruct = user.User{}

	if err:= c.ShouldBindJSON(&userStruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "",
			"error": "unexpected error",
		})
	}

	if ex, err := user.QueryUserBool(userStruct.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "",
			"error": "unexpected error",
		})
	} else if ex {
		c.JSON(http.StatusConflict, gin.H{
			"data": "",
			"error": "user already exists",
		})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(userStruct.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"data": "", "error": "unexpected error"})
	}

	userId, err := user.CreateUser(userStruct.Username, string(hashed))
	if err != nil {
		c.JSON(500, gin.H{"data": "", "error": "unexpected error"})
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId": userId,
	})
}

func SignToken(username string, password string, hashedPassword string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedPassword))
	if err != nil {
		return "", err
	}

	claims:= &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 30 * time.Hour)),
		Subject: username,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Login(c *gin.Context) {
	var user = user.User{}
	if err:=c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unexpected error",
		})
	}
}