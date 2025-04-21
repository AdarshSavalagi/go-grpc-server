FROM golang:1.24.0

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o grpc-server ./grpc-kafka-server/cmd/server/
ENV ENV_FILE=/app/.env.development
ENV CONFIG_PATH=/app/config/config.development.yaml

COPY .env.development /app/.env.development
COPY config/config.development.yaml /app/config/config.development.yaml
CMD ["./grpc-server"]
