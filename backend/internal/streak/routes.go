package streak

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func QueryStreakRoute(c *gin.Context) {
	userID, ex := c.Get("userID")
	if !ex{
		c.JSON(500, gin.H{
			"data":"",
			"err": "unexpected server error",
		})
		return
	}
	
	streak, err := QueryStreak(userID.(string))
	if err != nil {
		c.JSON(500, gin.H{
			"data":"",
			"err": "unexpected server error",
		})
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"data":streak,
		"err": "",
	})
}
