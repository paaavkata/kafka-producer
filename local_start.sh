# /bin/bash

payload='{"status": "started", "job_id": 1234}'

export KAFKA_BROKER=127.0.0.1:9092
export KAFKA_TOPIC=job-status
export PAYLOAD=$(echo -n $payload | base64)

go run main.go