package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var pub Publisher

type Publisher struct {
	client *redis.Client
	topic  string
}

func InitPublisher(redisURL, topic string) error {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return err
	}

	pub = Publisher{
		client: client,
		topic:  topic,
	}
	return nil
}

func Publish(ctx context.Context, eventType string, payload any) error {
	message := map[string]any{
		"type":    eventType,
		"payload": payload,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := pub.client.Publish(ctx, pub.topic, data).Err(); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
