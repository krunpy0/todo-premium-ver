package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user "github.com/krunpy0/todo-premium-ver/internal/user"
)

func Register (c *gin.Context) {
	var user = user.User{}

	if err:= c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unexpected error",
		})
	}

	

}