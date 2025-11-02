package main

import (
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func main() {
	kafkaURL := os.Getenv("KAFKA_BROKERS")
	if kafkaURL == "" {
		log.Fatal("KAFKA_BROKERS is not set")
	}

	kafkaProducer, err := sarama.NewSyncProducer([]string{kafkaURL}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaProducer.Close()

	handlers := NewHandlers(kafkaProducer)

	consumers, err := NewConsumers(kafkaURL)
	if err != nil {
		log.Fatal(err)
	}
	defer consumers.Shutdown()

	go func() {
		err := consumers.Start()
		if err != nil {
			log.Printf("Error starting consumers: %v", err)
		}
	}()

	mux := gin.New()
	mux.GET("/api/events/health", handlers.GetEventsServiceHealth)
	mux.POST("/api/events/movie", handlers.CreateMovieEvent)
	mux.POST("/api/events/user", handlers.CreateUserEvent)
	mux.POST("/api/events/payment", handlers.CreatePaymentEvent)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Starting events service on port %s", port)
	mux.Run("0.0.0.0:" + port)
}
