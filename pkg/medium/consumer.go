package medium

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	chatapp "golang_kafka_chat_application/pkg"
	chatinterfaces "golang_kafka_chat_application/pkg/interfaces"
	utility "golang_kafka_chat_application/pkg/utility"
)

type consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string) chatinterfaces.Consumer {
	c := kafka.ReaderConfig{
		Brokers:         brokers,
		Topic:           topic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
		GroupID:         utility.Ulid(),
		StartOffset:     kafka.LastOffset,
	}

	return &consumer{kafka.NewReader(c)}
}

func (c *consumer) Read(ctx context.Context, msg chan chatapp.Message, chErr chan error) {
	defer c.reader.Close()

	for {
		m, err := c.reader.ReadMessage(ctx)

		if err != nil {
			chErr <- fmt.Errorf("ERROR WHILE READING MESSAGE %v", err)
			continue
		}

		var message chatapp.Message
		err = json.Unmarshal(m.Value, &message)

		if err != nil {
			chErr <- err
		}

		msg <- message
	}
}
