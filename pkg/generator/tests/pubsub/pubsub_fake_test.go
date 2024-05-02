package pubsub_test

import (
	"testing"
	"time"

	"github.com/zarldev/ffakes/pkg/generator/tests/pubsub"
)

func TestPubSubFake(t *testing.T) {
	pub := pubsub.NewFakePublisher(t)
	sub := pubsub.NewFakeSubscriber(t)
	broker := pubsub.NewFakeBroker(t)
	resChan := make(chan *pubsub.Message)
	id := 1
	topic := "test-topic"
	data := []byte("test-data")
	timestamp := time.Now()
	message := pubsub.Message{
		ID:        id,
		Topic:     topic,
		Data:      data,
		Timestamp: timestamp,
	}
	t.Run("when publishing a message", func(t *testing.T) {
		// Setup
		pub.OnPublishMessage(func(m *pubsub.Message) error {
			return nil
		})
		// Test
		err := pub.PublishMessage(&message)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
		if pub.PublishMessageCount != 1 {
			t.Errorf("expected PublishMessage to be called once but got %d", pub.PublishMessageCount)
		}
	})
	t.Run("when subscribing to messages", func(t *testing.T) {
		go func() {
			resChan <- &message
		}()
		// Setup
		sub.OnSubscribeToMessages(func() (chan *pubsub.Message, error) {
			return resChan, nil
		})
		// Test
		_, err := sub.SubscribeToMessages()
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
		if sub.SubscribeToMessagesCount != 1 {
			t.Errorf("expected SubscribeToMessages to be called once but got %d", sub.SubscribeToMessagesCount)
		}
		select {
		case msg := <-resChan:
			if msg != &message {
				t.Errorf("expected to receive message %v but got %v", message, msg)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("expected to receive a message but got none")
		}
	})
	t.Run("when publishing and subscribing to messages", func(t *testing.T) {
		// Setup
		broker.OnPublishMessage(func(m *pubsub.Message) error {
			return nil
		})
		broker.OnSubscribeToMessages(func() (chan *pubsub.Message, error) {
			return make(chan *pubsub.Message), nil
		})
		// Test
		err := broker.PublishMessage(&message)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
		_, err = broker.SubscribeToMessages()
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}

	})

}
