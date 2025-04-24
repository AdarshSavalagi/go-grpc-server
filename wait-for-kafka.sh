#!/bin/bash

echo "Waiting for Kafka to be ready on kafka:9092..."
until nc -z kafka 9092; do
  echo "Kafka not ready yet..."
  sleep 2
done

echo "Kafka is up. Starting gRPC server..."
exec "$@"