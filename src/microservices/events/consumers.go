package main

import (
	"log"

	"github.com/IBM/sarama"
)

type Consumers struct {
	UserEventConsumer    *UserEventConsumer
	MovieEventConsumer   *MovieEventConsumer
	PaymentEventConsumer *PaymentEventConsumer
}

func NewConsumers(kafkaURL string) (*Consumers, error) {
	userEventConsumer, err := NewUserEventConsumer(kafkaURL)
	if err != nil {
		return nil, err
	}

	movieEventConsumer, err := NewMovieEventConsumer(kafkaURL)
	if err != nil {
		return nil, err
	}

	paymentEventConsumer, err := NewPaymentEventConsumer(kafkaURL)
	if err != nil {
		return nil, err
	}

	return &Consumers{UserEventConsumer: userEventConsumer, MovieEventConsumer: movieEventConsumer, PaymentEventConsumer: paymentEventConsumer}, nil
}

func (c *Consumers) Start() error {
	go func() {
		err := c.UserEventConsumer.Consume("user-events", c.UserEventConsumer.HandleMessage)
		if err != nil {
			log.Printf("Error consuming user-events: %v", err)
		}
	}()

	go func() {
		err := c.MovieEventConsumer.Consume("movie-events", c.MovieEventConsumer.HandleMessage)
		if err != nil {
			log.Printf("Error consuming movie-events: %v", err)
		}
	}()

	go func() {
		err := c.PaymentEventConsumer.Consume("payment-events", c.PaymentEventConsumer.HandleMessage)
		if err != nil {
			log.Printf("Error consuming payment-events: %v", err)
		}
	}()

	return nil
}

func (c *Consumers) Shutdown() error {
	err := c.UserEventConsumer.Shutdown()
	if err != nil {
		return err
	}

	err = c.MovieEventConsumer.Shutdown()
	if err != nil {
		return err
	}

	err = c.PaymentEventConsumer.Shutdown()
	if err != nil {
		return err
	}

	return nil
}

type UserEventConsumer struct {
	kafkaConsumer sarama.Consumer
}

func NewUserEventConsumer(kafkaURL string) (*UserEventConsumer, error) {
	kafkaConsumer, err := sarama.NewConsumer([]string{kafkaURL}, nil)
	if err != nil {
		return nil, err
	}
	return &UserEventConsumer{kafkaConsumer: kafkaConsumer}, nil
}

func (c *UserEventConsumer) HandleMessage(message *sarama.ConsumerMessage) {
	log.Println("Received user event:", string(message.Value))
}

func (c *UserEventConsumer) Consume(topic string, callback func(message *sarama.ConsumerMessage)) error {
	partitionConsumer, err := c.kafkaConsumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		callback(message)
	}

	return nil
}

func (c *UserEventConsumer) Shutdown() error {
	return c.kafkaConsumer.Close()
}

type MovieEventConsumer struct {
	kafkaConsumer sarama.Consumer
}

func NewMovieEventConsumer(kafkaURL string) (*MovieEventConsumer, error) {
	kafkaConsumer, err := sarama.NewConsumer([]string{kafkaURL}, nil)
	if err != nil {
		return nil, err
	}
	return &MovieEventConsumer{kafkaConsumer: kafkaConsumer}, nil
}

func (c *MovieEventConsumer) HandleMessage(message *sarama.ConsumerMessage) {
	log.Println("Received movie event:", string(message.Value))
}

func (c *MovieEventConsumer) Consume(topic string, callback func(message *sarama.ConsumerMessage)) error {
	partitionConsumer, err := c.kafkaConsumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		callback(message)
	}
	return nil
}

func (c *MovieEventConsumer) Shutdown() error {
	return c.kafkaConsumer.Close()
}

type PaymentEventConsumer struct {
	kafkaConsumer sarama.Consumer
}

func NewPaymentEventConsumer(kafkaURL string) (*PaymentEventConsumer, error) {
	kafkaConsumer, err := sarama.NewConsumer([]string{kafkaURL}, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentEventConsumer{kafkaConsumer: kafkaConsumer}, nil
}

func (c *PaymentEventConsumer) HandleMessage(message *sarama.ConsumerMessage) {
	log.Println("Received payment event:", string(message.Value))
}

func (c *PaymentEventConsumer) Consume(topic string, callback func(message *sarama.ConsumerMessage)) error {
	partitionConsumer, err := c.kafkaConsumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		callback(message)
	}
	return nil
}

func (c *PaymentEventConsumer) Shutdown() error {
	return c.kafkaConsumer.Close()
}
