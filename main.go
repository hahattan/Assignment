package main

import (
	"context"
	"sync"

	"github.com/hahattan/Appaegis/pubsub"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	s := pubsub.NewPubSubService(ctx, &wg)

	s1 := pubsub.NewSubscriber("sub1", "topic1")
	s2 := pubsub.NewSubscriber("sub2", "topic2")
	s3 := pubsub.NewSubscriber("sub3", "topic3")
	s4 := pubsub.NewSubscriber("sub4", "topic1")

	s.Subscribe("topic1", s1)
	s.Subscribe("topic2", s2)
	s.Subscribe("topic3", s3)
	s.Subscribe("topic1", s4)

	s.Publish("topic1", []byte("hello1"))
	s.Publish("topic2", []byte("hello2"))
	s.Publish("topic3", []byte("hello3"))

	cancel()
	wg.Wait()
}
