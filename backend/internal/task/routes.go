package task

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krunpy0/todo-premium-ver/internal/user"
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

func GetTaskByIDRoute(c *gin.Context) {
	taskID := c.Param("taskID")

	userID, ex := c.Get("userID")
	if !ex {
		c.JSON(500, gin.H{
			"data": "",
			"err": "userID not found",
		})
		return
	}
	task, err := GetTaskById(taskID, userID.(string))
	if err != nil {
		c.JSON(500, gin.H{
			"data": "",
			"err": "unexpected server error",
		})
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"data": task,
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

func xpByDifficulty(difficulty string) (int, error) {
	switch difficulty {
	case "easy":
			return 20, nil
	case "medium":
			return 40, nil
	case "hard":
			return 50, nil
	default:
			return 0, fmt.Errorf("unknown difficulty: %s", difficulty)
	}
}

func CompleteTaskRoute(c *gin.Context) {
	taskID := c.Param("taskID")

	userID, ex := c.Get("userID")
	if !ex {
		c.JSON(500, gin.H{
			"data": "",
			"err": "userID not found",
		})
		return
	}

	task, err := CompleteTask(userID.(string), taskID)
	if err != nil {
		c.JSON(500, gin.H{
			"data": "",
			"err": "task already completed, failed or not found",
		})
		fmt.Println(err)
		return
	}

	xpAmount, err := xpByDifficulty(task.Difficulty)

	if err != nil {
		c.JSON(500, gin.H{
			"data": "",
			"err": "unexpected server error",
		})
		fmt.Println(err)
		return
	}
	updatedXP, err := user.UpdateUserXP(userID.(string), xpAmount); if err != nil {
		c.JSON(500, gin.H{
			"data": "",
			"err": "unexpected server error",
		})
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{
		"data": gin.H{
			"updatedXP": updatedXP,
			"updatedTask": task,
		},
		"err": "",
	})
}

func FailTaskRoute(c *gin.Context) {
	taskID := c.Param("taskID")

	userID, ex := c.Get("userID")
	if !ex {
		c.JSON(500, gin.H{
			"data": "",
			"err": "userID not found",
		})
		return
	}
	task, err := FailTask(userID.(string), taskID)
	if err != nil {
		c.JSON(500, gin.H{
			"data": "",
			"err": "task already completed, failed or not found",
		})
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"data": task,
		"err": "",
	})
}

func CancelTaskRoute(c *gin.Context) {
	taskID := c.Param("taskID")

	userID, ex := c.Get("userID")
	if !ex {
		c.JSON(500, gin.H{
			"data": "",
			"err": "userID not found",
		})
		return
	}
	task, err := CancelTask(userID.(string), taskID)
	if err != nil {
		c.JSON(500, gin.H{
			"data": "",
			"err": "task cannot be cancelled or not found",
		})
		fmt.Println(err)
		return
	}
	c.JSON(200, gin.H{
		"data": task,
		"err": "",
	})
}