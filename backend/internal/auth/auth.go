package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/krunpy0/todo-premium-ver/internal/streak"
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
		fmt.Println(err)
		return
	}

	if ex, err := user.QueryUserBool(userStruct.Username); err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "",
			"error": "unexpected error",
		})
		fmt.Println(err)
		return
	} else if ex {
		c.JSON(http.StatusConflict, gin.H{
			"data": "",
			"error": "user already exists",
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(userStruct.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"data": "", "error": "unexpected error"})
		fmt.Println(err)
		return
	}

	createdUser, err := user.CreateUser(userStruct.Username, string(hashed))
	if err != nil {
		c.JSON(500, gin.H{"data": "", "error": "unexpected error"})
		fmt.Println(err)
		return
	}

	err = streak.CreateStreak(createdUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"data": "", "error": "unexpected error"})
		fmt.Println(err)
		return
	}

	tokenStr, err := SignToken(createdUser.ID)
	if err != nil {
    c.JSON(500, gin.H{"data": "", "error": "unexpected error"})
		fmt.Println(err)
    return
}

	c.SetCookie("token", tokenStr, 3600 * 24 * 30, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"data": "register successful, cookie set",
		"error": "",
	})
}
 
func comparePassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func SignToken(userID string) (string, error) {

	claims:= &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 30 * time.Hour)),
		Subject: userID,
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
	var userStruct = user.User{}
	if err:=c.ShouldBindJSON(&userStruct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "",
			"error": "unexpected error",
		})
		fmt.Println(err)
		return
	} 
	queriedUser, err := user.QueryUser(userStruct.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "",
			"error": "unexpected error",
		})	
		fmt.Println(err)
		return
	}
	if err := comparePassword(userStruct.Password, queriedUser.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": "",
			"error": "invalid password",
		})
		
		return
	}
	tokenStr, err := SignToken(queriedUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": "",
			"error": "unexpected error",
		})
		fmt.Println(err)
		return
	}
	c.SetCookie("token", tokenStr, 3600 * 24 * 30, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"data": "login successful",
		"error": "",
	})

}