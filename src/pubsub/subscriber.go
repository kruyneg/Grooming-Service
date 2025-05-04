package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type EventHandler func(eventType string, payload map[string]interface{})

type Subscriber struct {
	client *redis.Client
	topic  string
}

func NewSubscriber(client *redis.Client, topic string) *Subscriber {
	return &Subscriber{
		client: client,
		topic:  topic,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context, handler EventHandler) error {
	sub := s.client.Subscribe(ctx, s.topic)
	ch := sub.Channel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-ch:
			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &raw); err != nil {
				fmt.Println("invalid message format:", err)
				continue
			}

			eventType, _ := raw["type"].(string)
			payload, _ := raw["payload"].(map[string]interface{})

			go handler(eventType, payload)
		}
	}
}
