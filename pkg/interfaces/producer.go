package chatinterfaces

import "context"

type Producer interface {
	Publish(ctx context.Context, payload interface{}) error
}
