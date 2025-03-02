package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pliliya111/go_sprint2_project/internal/handler"
)

func main() {
	r := gin.Default()

	r.POST("/api/v1/calculate", handler.AddExpression)
	r.GET("/api/v1/expressions", handler.GetExpressions)
	r.GET("/api/v1/expressions/:id", handler.GetExpressionByID)
	r.GET("/internal/task", handler.GetTask)
	r.POST("/internal/task", handler.SubmitTaskResult)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
