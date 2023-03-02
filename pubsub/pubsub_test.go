package pubsub

import (
	"context"
	"sync"
	"testing"
)

func TestPubSubService_Publish(t *testing.T) {
	sub := NewSubscriber("testSub", "topic1")
	t.Run("Publish", func(t *testing.T) {
		pubsub := NewPubSubService(context.Background(), &sync.WaitGroup{})
		pubsub.Subscribe("topic1", sub)
		pubsub.Publish("topic1", []byte("test"))
		// test if the message is printed
	})
}

func TestPubSubService_Subscribe(t *testing.T) {
	t.Run("Subscribe", func(t *testing.T) {
		pubsub := NewPubSubService(context.Background(), &sync.WaitGroup{})
		pubsub.Subscribe("test1", &Subscriber{})
		pubsub.Subscribe("test1", &Subscriber{})
		pubsub.Subscribe("test2", &Subscriber{})

		if len(pubsub.subs) != 2 {
			t.Errorf("got %d, want %d", len(pubsub.subs), 2)
		}
		if len(pubsub.subs["test1"]) != 2 {
			t.Errorf("got %d, want %d", len(pubsub.subs["test1"]), 2)
		}
	})
}
