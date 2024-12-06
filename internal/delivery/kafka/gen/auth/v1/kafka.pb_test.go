package v1

import (
	"testing"

	ck "github.com/Karzoug/meower-common-go/kafka"
)

func TestChangedEvent_Reset(t *testing.T) {
	authChangedEventFngpnt := ck.MessageTypeHeaderValue(&ChangedEvent{})
	t.Log(authChangedEventFngpnt)
}
