package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	broker := os.Getenv("KAFKA_BROKER")   // e.g., "kafka-broker:9092"
	topic := os.Getenv("KAFKA_TOPIC")     // e.g., "job-status"
	payloadBase64 := os.Getenv("PAYLOAD") // Base64-encoded JSON payload

	if payloadBase64 == "" {
		fmt.Println("Error: PAYLOAD environment variable is required.")
		os.Exit(1)
	}

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
	defer client.Close()

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
	fmt.Println("Message sent successfully:", string(payloadBytes))
}
