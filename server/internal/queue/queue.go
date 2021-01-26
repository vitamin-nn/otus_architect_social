package queue

import "github.com/streadway/amqp"

type Consumer interface {
	Handle(fn func(<-chan amqp.Delivery), threads int) error
	Close()
}

type Publisher interface {
	Connect() error
	Publish(body interface{}, routingKey string) error
	Close()
}
