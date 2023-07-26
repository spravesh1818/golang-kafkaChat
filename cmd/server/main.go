package main

import (
	"context"
	"encoding/json"
	"fmt"
	chatapp "golang_kafka_chat_application/pkg"
	chatinterfaces "golang_kafka_chat_application/pkg/interfaces"
	"golang_kafka_chat_application/pkg/medium"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type Request struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func joinHandler(producer chatinterfaces.Producer) func(*gin.Context) {
	return func(c *gin.Context) {
		var req Request
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		message := chatapp.SystemMessage(fmt.Sprintf("%s has joined the room!", req.Username))
		if err := producer.Publish(context.Background(), message); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "message published"})
	}
}

func publishHandler(producer chatinterfaces.Producer) func(*gin.Context) {
	return func(c *gin.Context) {
		var req Request
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		message := chatapp.NewMessage(req.Username, req.Message)
		if err := producer.Publish(context.Background(), message); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "message published"})
	}
}

func main() {
	var (
		brokers = os.Getenv("KAFKA_BROKERS")
		topic   = os.Getenv("KAFKA_TOPIC")
	)

	producer := medium.NewProducer(strings.Split(brokers, ","), topic)
	r := gin.Default()
	r.POST("/join", joinHandler(producer))
	r.POST("/publish", publishHandler(producer))

	_ = r.Run()
}
