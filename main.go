package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	logger "github.com/paaavkata/go-logger"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	// First check the input argouments
	logger.Init(
		"info",
		"json",
		"kafka-producer",
		"dev",
		false,
		true,
		false,
		nil,
		nil)

	if len(os.Args) != 4 {
		fmt.Println("Usage: kafka-producer <broker> <topic> <payload>")
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]
	payloadBase64 := os.Args[3]

	if payloadBase64 == "" {
		fmt.Println("Error: PAYLOAD environment variable is required.")
		os.Exit(1)
	}

	if broker == "" || topic == "" || payloadBase64 == "" {
		broker = os.Getenv("KAFKA_BROKER")   // e.g., "kafka-broker:9092"
		topic = os.Getenv("KAFKA_TOPIC")     // e.g., "job-status"
		payloadBase64 = os.Getenv("PAYLOAD") // Base64-encoded JSON payload

		if payloadBase64 == "" {
			fmt.Println("Error: PAYLOAD environment variable is required.")
			os.Exit(1)
		}
	}

	logger.Debugf("broker: %s", broker)
	logger.Debugf("topic: %s", topic)
	logger.Debugf("payload: %s", string(payloadBase64))

	// Decode base64 payload
	payloadBytes, err := base64.StdEncoding.DecodeString(payloadBase64)
	if err != nil {
		fmt.Println("Error decoding base64 payload:", err)
		os.Exit(1)
	}

	// Validate JSON structure
	var jsonPayload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &jsonPayload); err != nil {
		fmt.Println("Error parsing JSON payload:", err)
		os.Exit(1)
	}

	// Add timestamp to payload
	jsonPayload["timestamp"] = time.Now().Unix()

	// Re-encode JSON payload
	payloadBytes, err = json.Marshal(jsonPayload)
	if err != nil {
		fmt.Println("Error encoding JSON payload:", err)
		os.Exit(1)
	}

	// Create Kafka producer
	client, err := kgo.NewClient(
		kgo.SeedBrokers(broker),
	)
	if err != nil {
		fmt.Println("Failed to create Kafka client:", err)
		os.Exit(1)
	}

	// Produce message
	record := &kgo.Record{
		Topic: topic,
		Value: payloadBytes,
	}
	ctx := context.Background()
	err = client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		fmt.Println("Failed to send message:", err)
		os.Exit(1)
	}

	// Ensure the producer exits after the first successful send.
	client.Close()
	fmt.Println("Message sent successfully:", string(payloadBytes))
	os.Exit(0)
}
