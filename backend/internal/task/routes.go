package task

import (
	"fmt"
	"net/http"

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
			"err": "unexpected server error",
		})
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"data": userTasks,
		"err": "",
	})
}

func CreateTaskRoute(c *gin.Context) {
	var task = Task{}
	err := c.ShouldBindJSON(&task); if err != nil {
		c.JSON(500, gin.H{
			"data": "",
			"err": "unexpected server error",
		})
		fmt.Println(err)
		return
	}
	if task.Difficulty != "easy" && task.Difficulty != "medium" && task.Difficulty != "hard" {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "",
			"err":"invalid difficulty type. allowed: 'easy', 'medium', 'hard' ",
		})
		return
	}
	userID, ex := c.Get("userID")
	if !ex {
		c.JSON(500, gin.H{
			"data": "",
			"err": "userID not found",
		})
		return
	}

	taskID, err := CreateTask(userID.(string), task.Title, task.Difficulty, task.Due)
	if err != nil {
		c.JSON(500, gin.H{
			"data": "",
			"err": "unexpected server error",
		})
		fmt.Println(err)
		return	
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": taskID,
		"err": "",
	})

}