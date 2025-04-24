FROM golang:1.24.0

RUN apt-get update && \
    apt-get install -y librdkafka-dev && \
    apt-get clean

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y netcat-openbsd

RUN go mod tidy
RUN go build -o grpc-server ./cmd/server/
ENV ENV_FILE=/app/.env.development
ENV CONFIG_PATH=/app/config/config.development.yaml

COPY .env.development /app/.env.development
COPY config/config.development.yaml /app/config/config.development.yaml

COPY wait-for-kafka.sh /usr/bin/wait-for-kafka.sh
RUN chmod +x /usr/bin/wait-for-kafka.sh

ENTRYPOINT ["/usr/bin/wait-for-kafka.sh"]
CMD ["./grpc-server"]