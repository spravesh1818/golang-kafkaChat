#!/bin/bash

# Check if the number of arguments is correct
if [ $# -ne 1 ]; then
  echo "Usage: $0 <topic_name>"
  exit 1
fi

# Set the Kafka container name or host
KAFKA_CONTAINER_NAME="kafka"  # Replace with your Kafka container name or host

# Extract the topic name from the command-line argument
TOPIC_NAME=$1

# Execute the kafka-topics.sh script to create the topic
docker exec -it $KAFKA_CONTAINER_NAME kafka-topics.sh \
  --create \
  --zookeeper zookeeper:2181 \
  --topic $TOPIC_NAME \
  --partitions 1 \
  --replication-factor 1