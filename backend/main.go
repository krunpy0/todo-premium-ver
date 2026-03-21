package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/krunpy0/todo-premium-ver/db"
	"github.com/krunpy0/todo-premium-ver/internal/auth"
	"github.com/krunpy0/todo-premium-ver/internal/task"
)

func main() {
	godotenv.Load()
	CONN_STR := os.Getenv("CONN_STR")
	if CONN_STR == "" {
		log.Fatal("CONN_STR is not set")
	}
	
	if err := db.Init(CONN_STR); err != nil {
		log.Fatal(err)	
	}
	defer db.DB.Close()

	router := gin.Default()
	api := router.Group("/api")
	api.Use(auth.Auth)
	{
		api.GET("/protected", func (c *gin.Context) {
			c.JSON(200, gin.H{
				"data": "this route is protected",
				"err": "",
			})
		})
		api.GET("/tasks",task.GetTasksRoute) 
		api.POST("/tasks", task.CreateTaskRoute)
	}
	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	router.Run(":" + os.Getenv("PORT"))
}
