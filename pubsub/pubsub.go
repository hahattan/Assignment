package pubsub

import (
	"context"
	"fmt"
	"sync"
)

type PubSubService struct {
	sync.RWMutex
	subs map[string][]*Subscriber
	ctx  context.Context
	wg   *sync.WaitGroup
}

func NewPubSubService(ctx context.Context, wg *sync.WaitGroup) PubSubService {
	return PubSubService{
		subs: make(map[string][]*Subscriber),
		ctx:  ctx,
		wg:   wg,
	}
}

func (s *PubSubService) Publish(topic string, message []byte) {
	s.RLock()
	subs := s.subs[topic]
	s.RUnlock()

	for _, sub := range subs {
		sub.OnMessage(message)
	}
}

func (s *PubSubService) Subscribe(topic string, subscriber *Subscriber) {
	s.Lock()
	defer s.Unlock()
	if s.subs[topic] == nil {
		s.subs[topic] = make([]*Subscriber, 0)
	}

	s.subs[topic] = append(s.subs[topic], subscriber)
	s.wg.Add(1)
	go subscriber.Run(s.ctx, s.wg)
}

type Subscriber struct {
	name     string
	topic    string
	messages chan string
}

func NewSubscriber(name string, topic string) *Subscriber {
	return &Subscriber{
		name:     name,
		topic:    topic,
		messages: make(chan string),
	}
}

func (s *Subscriber) OnMessage(message []byte) {
	s.messages <- string(message)
}

func (s *Subscriber) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-s.messages:
			// handle the received message, print the log here just for demo
			fmt.Printf("%s received %s on %s\n", s.name, msg, s.topic)
		}
	}
}
