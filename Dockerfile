FROM gcr.io/distroless/base-debian12
ARG TARGETARCH
COPY kafka-producer-${TARGETARCH} /kafka-producer
ENTRYPOINT ["/kafka-producer"]
