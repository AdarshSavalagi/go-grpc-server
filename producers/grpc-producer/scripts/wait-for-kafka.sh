#!/bin/sh

# Set default Kafka host and port if not already defined as environment variables
KAFKA_HOST=${KAFKA_HOST:-kafka}
KAFKA_PORT=${KAFKA_PORT:-9092}

echo "Waiting for Kafka at $KAFKA_HOST:$KAFKA_PORT..."

# Use a loop to continuously check the network connection to Kafka
until nc -z "$KAFKA_HOST" "$KAFKA_PORT"; do
  echo "Still waiting for Kafka at $KAFKA_HOST:$KAFKA_PORT..."
  sleep 1 # Wait for 1 second before the next attempt
done

echo "Kafka is up and reachable at $KAFKA_HOST:$KAFKA_PORT!"

# Execute the command passed as arguments to this script
echo "Executing command: $*"
exec "$@"