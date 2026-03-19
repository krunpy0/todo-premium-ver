package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/krunpy0/todo-premium-ver/db"
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
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	router.Run(":" + os.Getenv("PORT"))
}
