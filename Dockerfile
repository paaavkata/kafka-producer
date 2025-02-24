# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY main.go .
RUN go mod init producer && go mod tidy && go build -o producer main.go

# Runtime stage
FROM gcr.io/distroless/base-debian10
COPY --from=builder /app/producer .
ENTRYPOINT ["/producer"]
