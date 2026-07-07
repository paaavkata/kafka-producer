package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	logger "github.com/paaavkata/go-logger"
	gonats "github.com/paaavkata/go-nats"
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
		fmt.Println("Usage: kafka-producer <nats-url> <topic> <payload>")
		os.Exit(1)
	}

	url := os.Args[1]
	topic := os.Args[2]
	payloadBase64 := os.Args[3]

	if payloadBase64 == "" {
		fmt.Println("Error: PAYLOAD environment variable is required.")
		os.Exit(1)
	}

	logger.Debugf("url: %s", url)
	logger.Debugf("topic: %s", topic)
	logger.Debugf("payload: %s", string(payloadBase64))
	logger.Infof("Starting producer with url=%s topic=%s payload_b64_len=%d", url, topic, len(payloadBase64))

	// Optional timeout (seconds) for the produce request.
	// Defaults to 20s to avoid hanging init containers indefinitely.
	// NATS_PRODUCE_TIMEOUT_SEC preferred; KAFKA_PRODUCE_TIMEOUT_SEC accepted
	// for backward compatibility with existing pod templates.
	produceTimeoutSec := 20
	timeoutEnv := os.Getenv("NATS_PRODUCE_TIMEOUT_SEC")
	if timeoutEnv == "" {
		timeoutEnv = os.Getenv("KAFKA_PRODUCE_TIMEOUT_SEC")
	}
	if timeoutEnv != "" {
		if parsed, err := strconv.Atoi(timeoutEnv); err == nil && parsed > 0 {
			produceTimeoutSec = parsed
		} else {
			logger.Errorf("Invalid produce timeout %q, using default=%ds", timeoutEnv, produceTimeoutSec)
		}
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
	logger.Infof("Payload decoded and validated as JSON")

	// Best-effort network diagnostics to quickly identify DNS/connectivity issues
	hostPort := strings.TrimPrefix(strings.TrimPrefix(url, "nats://"), "tls://")
	host, port, splitErr := net.SplitHostPort(hostPort)
	if splitErr == nil {
		ips, err := net.LookupHost(host)
		if err != nil {
			logger.Errorf("DNS lookup failed for host=%s: %v", host, err)
		} else {
			logger.Infof("DNS lookup for host=%s resolved to: %v", host, ips)
		}

		dialTimeout := 3 * time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), dialTimeout)
		if err != nil {
			logger.Errorf("TCP dial check failed to %s:%s (timeout=%s): %v", host, port, dialTimeout, err)
		} else {
			logger.Infof("TCP dial check succeeded to %s:%s", host, port)
			_ = conn.Close()
		}
	} else {
		logger.Errorf("Could not parse host:port from %q: %v", hostPort, splitErr)
	}

	// Add timestamp to payload
	jsonPayload["timestamp"] = time.Now().Unix()

	producer, err := gonats.NewProducer(&gonats.ProducerConfig{
		URLs:     []string{url},
		ClientID: "kafka-producer",
		Topic:    topic,
		Timeout:  time.Duration(produceTimeoutSec) * time.Second,
	})
	if err != nil {
		fmt.Println("Failed to create NATS producer:", err)
		os.Exit(1)
	}
	logger.Infof("NATS producer created successfully")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(produceTimeoutSec)*time.Second)
	defer cancel()
	logger.Infof("Producing message with timeout=%ds", produceTimeoutSec)
	if err := producer.SendMessageWithContext(ctx, "", jsonPayload); err != nil {
		fmt.Println("Failed to send message:", err)
		os.Exit(1)
	}
	logger.Infof("Message produced successfully")

	// Ensure the producer exits after the first successful send.
	producer.Close()
	payloadBytes, _ = json.Marshal(jsonPayload)
	fmt.Println("Message sent successfully:", string(payloadBytes))
	os.Exit(0)
}
