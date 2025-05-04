package main

import (
	"context"
	"dog-service/pubsub"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleSignals(cancel)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	subscriber := pubsub.NewSubscriber(rdb, "logs")

	fmt.Println("🚀 Log listener started. Waiting for events...")

	err := subscriber.Subscribe(ctx, func(eventType string, payload map[string]any) {
		p, _ := json.MarshalIndent(payload, "", "\t")
		log.Printf("📨 Received event: %s\n📦 Payload: %+v\n", eventType, string(p))
	})
	if err != nil {
		log.Fatalf("❌ Subscriber error: %v", err)
	}

	fmt.Println("👋 Log listener stopped.")
}

func handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println("\n🛑 Received shutdown signal")
	cancel()
}
