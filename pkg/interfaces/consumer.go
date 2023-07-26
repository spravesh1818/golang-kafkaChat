package chatinterfaces

import (
	"context"
	chatapp "golang_kafka_chat_application/pkg"
)

type Consumer interface {
	Read(ctx context.Context, msg chan chatapp.Message, err chan error)
}
