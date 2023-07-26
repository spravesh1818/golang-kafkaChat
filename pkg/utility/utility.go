package utility

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func Ulid() string {
	t := time.Now().UTC()
	id := ulid.MustNew(ulid.Timestamp(t), rand.Reader)
	return id.String()
}
