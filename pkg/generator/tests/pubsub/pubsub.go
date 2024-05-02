package pubsub

import "time"

//go:generate ffakes -vv -i Publisher,Subscriber,Broker
type Message struct {
	ID        int
	Topic     string
	Data      []byte
	Timestamp time.Time
}

type Publisher interface {
	PublishMessage(m *Message) error
}

type Subscriber interface {
	SubscribeToMessages() (chan *Message, error)
}

type Broker interface {
	BrokerConnection() error
	Publisher
	Subscriber
}
