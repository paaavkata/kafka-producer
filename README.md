# kafka-producer
GoLang function that produces message to a given topic in a given Kafka cluster with provided payload - base64 encoded JSON

# Usage

The image expects 3 env variables:
- KAFKA_BROKER
- KAFKA_TOPIC
- PAYLOAD - base64 encoded JSON
```bash
payload='{"status": "started", "job_id": 1234}'
docker run infra/kafka-producer \
    -e KAFKA_BROKER=127.0.0.1:9092 \
    -e KAFKA_TOPIC=job-status \
    -e PAYLOAD=$(echo -n $payload | base64)
```

The script automatically adds a `timestamp` field in the JSON.