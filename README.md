 docker compose --profile arm up

 protoc --go_out=. --go-grpc_out=. log.proto

docker exec -it kafka bash

 /opt/bitnami/kafka/bin/kafka-topics.sh --create --topic logs --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1