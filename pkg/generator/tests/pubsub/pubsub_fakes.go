// Code generated by ffakes v0.0.1 DO NOT EDIT.

package pubsub

import (
	"testing"
)

type FakePublisher struct {
	t                   *testing.T
	PublishMessageCount int
	FPublishMessage     []func(m *Message) error
}

type PublishMessageFunc = func(m *Message) error
type PublisherOption func(f *FakePublisher)

func OnPublishMessage(fn ...PublishMessageFunc) PublisherOption {
	return func(f *FakePublisher) {
		f.FPublishMessage = append(f.FPublishMessage, fn...)
	}
}

func (f *FakePublisher) OnPublishMessage(fns ...PublishMessageFunc) {
	for _, fn := range fns {
		f.FPublishMessage = append(f.FPublishMessage, fn)
	}
}

func NewFakePublisher(t *testing.T, opts ...PublisherOption) *FakePublisher {
	f := &FakePublisher{t: t}
	for _, opt := range opts {
		opt(f)
	}
	t.Cleanup(func() {
		if f.PublishMessageCount != len(f.FPublishMessage) {
			t.Fatalf("expected PublishMessage to be called %d times but got %d", len(f.FPublishMessage), f.PublishMessageCount)
		}
	})
	return f
}

func (f *FakePublisher) PublishMessage(m *Message) error {
	var idx = f.PublishMessageCount
	if f.PublishMessageCount >= len(f.FPublishMessage) {
		idx = len(f.FPublishMessage) - 1
	}
	if len(f.FPublishMessage) != 0 {
		o1 := f.FPublishMessage[idx](m)
		f.PublishMessageCount++
		return o1
	}
	return nil
}

type FakeSubscriber struct {
	t                        *testing.T
	SubscribeToMessagesCount int
	FSubscribeToMessages     []func() (chan *Message, error)
}

type SubscribeToMessagesFunc = func() (chan *Message, error)
type SubscriberOption func(f *FakeSubscriber)

func OnSubscribeToMessages(fn ...SubscribeToMessagesFunc) SubscriberOption {
	return func(f *FakeSubscriber) {
		f.FSubscribeToMessages = append(f.FSubscribeToMessages, fn...)
	}
}

func (f *FakeSubscriber) OnSubscribeToMessages(fns ...SubscribeToMessagesFunc) {
	for _, fn := range fns {
		f.FSubscribeToMessages = append(f.FSubscribeToMessages, fn)
	}
}

func NewFakeSubscriber(t *testing.T, opts ...SubscriberOption) *FakeSubscriber {
	f := &FakeSubscriber{t: t}
	for _, opt := range opts {
		opt(f)
	}
	t.Cleanup(func() {
		if f.SubscribeToMessagesCount != len(f.FSubscribeToMessages) {
			t.Fatalf("expected SubscribeToMessages to be called %d times but got %d", len(f.FSubscribeToMessages), f.SubscribeToMessagesCount)
		}
	})
	return f
}

func (f *FakeSubscriber) SubscribeToMessages() (chan *Message, error) {
	var idx = f.SubscribeToMessagesCount
	if f.SubscribeToMessagesCount >= len(f.FSubscribeToMessages) {
		idx = len(f.FSubscribeToMessages) - 1
	}
	if len(f.FSubscribeToMessages) != 0 {
		o1, o2 := f.FSubscribeToMessages[idx]()
		f.SubscribeToMessagesCount++
		return o1, o2
	}
	return nil, nil
}

type FakeBroker struct {
	t                        *testing.T
	BrokerConnectionCount    int
	PublishMessageCount      int
	SubscribeToMessagesCount int
	FBrokerConnection        []func() error
	FPublishMessage          []func(m *Message) error
	FSubscribeToMessages     []func() (chan *Message, error)
}

type BrokerConnectionFunc = func() error
type BrokerOption func(f *FakeBroker)

func BrokerOnBrokerConnection(fn ...BrokerConnectionFunc) BrokerOption {
	return func(f *FakeBroker) {
		f.FBrokerConnection = append(f.FBrokerConnection, fn...)
	}
}

func BrokerOnPublishMessage(fn ...PublishMessageFunc) BrokerOption {
	return func(f *FakeBroker) {
		f.FPublishMessage = append(f.FPublishMessage, fn...)
	}
}

func BrokerOnSubscribeToMessages(fn ...SubscribeToMessagesFunc) BrokerOption {
	return func(f *FakeBroker) {
		f.FSubscribeToMessages = append(f.FSubscribeToMessages, fn...)
	}
}

func (f *FakeBroker) OnBrokerConnection(fns ...BrokerConnectionFunc) {
	for _, fn := range fns {
		f.FBrokerConnection = append(f.FBrokerConnection, fn)
	}
}

func (f *FakeBroker) OnPublishMessage(fns ...PublishMessageFunc) {
	for _, fn := range fns {
		f.FPublishMessage = append(f.FPublishMessage, fn)
	}
}

func (f *FakeBroker) OnSubscribeToMessages(fns ...SubscribeToMessagesFunc) {
	for _, fn := range fns {
		f.FSubscribeToMessages = append(f.FSubscribeToMessages, fn)
	}
}

func NewFakeBroker(t *testing.T, opts ...BrokerOption) *FakeBroker {
	f := &FakeBroker{t: t}
	for _, opt := range opts {
		opt(f)
	}
	t.Cleanup(func() {
		if f.BrokerConnectionCount != len(f.FBrokerConnection) {
			t.Fatalf("expected BrokerConnection to be called %d times but got %d", len(f.FBrokerConnection), f.BrokerConnectionCount)
		}
		if f.PublishMessageCount != len(f.FPublishMessage) {
			t.Fatalf("expected PublishMessage to be called %d times but got %d", len(f.FPublishMessage), f.PublishMessageCount)
		}
		if f.SubscribeToMessagesCount != len(f.FSubscribeToMessages) {
			t.Fatalf("expected SubscribeToMessages to be called %d times but got %d", len(f.FSubscribeToMessages), f.SubscribeToMessagesCount)
		}
	})
	return f
}

func (f *FakeBroker) BrokerConnection() error {
	var idx = f.BrokerConnectionCount
	if f.BrokerConnectionCount >= len(f.FBrokerConnection) {
		idx = len(f.FBrokerConnection) - 1
	}
	if len(f.FBrokerConnection) != 0 {
		o1 := f.FBrokerConnection[idx]()
		f.BrokerConnectionCount++
		return o1
	}
	return nil
}

func (f *FakeBroker) PublishMessage(m *Message) error {
	var idx = f.PublishMessageCount
	if f.PublishMessageCount >= len(f.FPublishMessage) {
		idx = len(f.FPublishMessage) - 1
	}
	if len(f.FPublishMessage) != 0 {
		o1 := f.FPublishMessage[idx](m)
		f.PublishMessageCount++
		return o1
	}
	return nil
}

func (f *FakeBroker) SubscribeToMessages() (chan *Message, error) {
	var idx = f.SubscribeToMessagesCount
	if f.SubscribeToMessagesCount >= len(f.FSubscribeToMessages) {
		idx = len(f.FSubscribeToMessages) - 1
	}
	if len(f.FSubscribeToMessages) != 0 {
		o1, o2 := f.FSubscribeToMessages[idx]()
		f.SubscribeToMessagesCount++
		return o1, o2
	}
	return nil, nil
}