FROM golang:1.24.0

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o grpc-server ./grpc-kafka-server/cmd/server/

CMD ["./grpc-server"]
