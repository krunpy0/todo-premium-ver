package task

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetTasksRoute(c *gin.Context) {
	val, ex := c.Get("userID")
	if !ex {
		c.JSON(500, gin.H{
			"data": "",
			"err": "userID not found",
		})
		return
	}

	userTasks, err := GetUserTasks(val.(string))
	if err != nil {
		c.JSON(500, gin.H{
			"data": "",
			"err": "some tesyt err",
		})
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"data": userTasks,
		"err": "",
	})
}