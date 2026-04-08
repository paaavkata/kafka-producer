# kafka-producer
GoLang utility that produces a message to a given topic in a given Kafka cluster with provided payload - base64 encoded JSON

## Usage

The application accepts parameters in two ways (command-line arguments take precedence over environment variables):

### Option 1: Command-line Arguments
```bash
payload='{"status": "started", "job_id": 1234}'
kafka-producer 127.0.0.1:9092 job-status $(echo -n $payload | base64)
```

### Option 2: Environment Variables
```bash
payload='{"status": "started", "job_id": 1234}'
export KAFKA_BROKER=127.0.0.1:9092
export KAFKA_TOPIC=job-status
export PAYLOAD=$(echo -n $payload | base64)
kafka-producer
```

### Docker Usage
```bash
payload='{"status": "started", "job_id": 1234}'
docker run infra/kafka-producer \
    127.0.0.1:9092 \
    job-status \
    $(echo -n $payload | base64)
```

Or with environment variables:
```bash
payload='{"status": "started", "job_id": 1234}'
docker run \
    -e KAFKA_BROKER=127.0.0.1:9092 \
    -e KAFKA_TOPIC=job-status \
    -e PAYLOAD=$(echo -n $payload | base64) \
    infra/kafka-producer
```

## Features

- Accepts both command-line arguments and environment variables
- Automatically adds a `timestamp` field (Unix timestamp) to the JSON payload
- Validates that the payload is valid base64-encoded JSON
- Uses franz-go Kafka client for reliable message delivery