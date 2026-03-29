package user

import "github.com/gin-gonic/gin"

func QueryMe(c *gin.Context) {
	userID, ex:= c.Get("userID")
	if !ex {
		c.JSON(500, gin.H{
			"data":"",
			"err":"unexpected server error",
		})
	}

	user, err := queryUserByID(userID.(string))
	if err != nil {
		c.JSON(500, gin.H{
			"data":"",
			"err":"unexpected server error",
		})
	}

	c.JSON(200, gin.H{
		"data":user,
		"err":"",
	})

}
