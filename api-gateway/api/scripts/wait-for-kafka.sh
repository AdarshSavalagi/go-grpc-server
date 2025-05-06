#!/bin/sh
KAFKA_HOST=${KAFKA_HOST:-kafka}
KAFKA_PORT=${KAFKA_PORT:-9092}

echo "Waiting for Kafka at $KAFKA_HOST:$KAFKA_PORT..."

while ! nc -z "$KAFKA_HOST" "$KAFKA_PORT"; do
  sleep 1
  echo "Still waiting for Kafka..."
done

echo "Kafka is up!"
exec "$@"