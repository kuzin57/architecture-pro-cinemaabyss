package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/cinemaabyss/microservices/events/generated"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	kafkaProducer sarama.SyncProducer
}

func NewHandlers(kafkaProducer sarama.SyncProducer) *Handlers {
	return &Handlers{
		kafkaProducer: kafkaProducer,
	}
}

func (h *Handlers) GetEventsServiceHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (h *Handlers) CreateMovieEvent(c *gin.Context) {
	var event generated.MovieEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		errorResp := generated.Error{Error: fmt.Sprintf("Invalid request body: %v", err)}
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	partition, offset, err := h.sendEventToKafka("movie-events", event)
	if err != nil {
		errorResp := generated.Error{Error: fmt.Sprintf("Failed to send event to Kafka: %v", err)}
		c.JSON(http.StatusInternalServerError, errorResp)
		return
	}

	response := generated.EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    int(offset),
		Event: generated.Event{
			Id:        fmt.Sprintf("movie-%d-%s", event.MovieId, event.Action),
			Type:      "movie",
			Timestamp: time.Now(),
			Payload:   serialize(event),
		},
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handlers) CreateUserEvent(c *gin.Context) {
	var event generated.UserEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		errorResp := generated.Error{Error: fmt.Sprintf("Invalid request body: %v", err)}
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	partition, offset, err := h.sendEventToKafka("user-events", event)
	if err != nil {
		errorResp := generated.Error{Error: fmt.Sprintf("Failed to send event to Kafka: %v", err)}
		c.JSON(http.StatusInternalServerError, errorResp)
		return
	}

	response := generated.EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    int(offset),
		Event: generated.Event{
			Id:        fmt.Sprintf("user-%d-%s", event.UserId, event.Action),
			Type:      "user",
			Timestamp: event.Timestamp,
			Payload:   serialize(event),
		},
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handlers) CreatePaymentEvent(c *gin.Context) {
	var event generated.PaymentEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		errorResp := generated.Error{Error: fmt.Sprintf("Invalid request body: %v", err)}
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	partition, offset, err := h.sendEventToKafka("payment-events", event)
	if err != nil {
		errorResp := generated.Error{Error: fmt.Sprintf("Failed to send event to Kafka: %v", err)}
		c.JSON(http.StatusInternalServerError, errorResp)
		return
	}

	response := generated.EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    int(offset),
		Event: generated.Event{
			Id:        fmt.Sprintf("payment-%d", event.PaymentId),
			Type:      "payment",
			Timestamp: event.Timestamp,
			Payload:   serialize(event),
		},
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handlers) sendEventToKafka(topic string, event any) (int, int64, error) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return 0, 0, err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(eventJSON),
	}

	partition, offset, err := h.kafkaProducer.SendMessage(msg)
	if err != nil {
		return 0, 0, err
	}

	return int(partition), offset, nil
}

func serialize(entity any) map[string]any {
	bytes, err := json.Marshal(entity)
	if err != nil {
		return nil
	}

	payload := make(map[string]any)
	if err := json.Unmarshal(bytes, &payload); err != nil {
		return nil
	}

	return payload
}
