package medium

import (
	"context"
	"encoding/json"
	chatinterfaces "golang_kafka_chat_application/pkg/interfaces"
	"golang_kafka_chat_application/pkg/utility"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

type producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) chatinterfaces.Producer {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientID,
	}

	p := kafka.WriterConfig{
		Brokers:          brokers,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}

	return &producer{kafka.NewWriter(p)}

}

func (p *producer) Publish(ctx context.Context, payload interface{}) error {
	message, err := p.encodeMessage(payload)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, message)
}

func (p *producer) encodeMessage(payload interface{}) (kafka.Message, error) {
	m, err := json.Marshal(payload)
	if err != nil {
		return kafka.Message{}, err
	}

	key := utility.Ulid()
	return kafka.Message{
		Key:   []byte(key),
		Value: m,
	}, nil
}
