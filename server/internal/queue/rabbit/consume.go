package rabbit

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Consume struct {
	r           *rabbit
	consumerTag string
	qos         int
}

func NewConsume(dsn, consumerTag, exchangeName, queue string) *Consume {
	c := &Consume{
		consumerTag: consumerTag,
		qos:         10,
		r: &rabbit{
			exchangeName: exchangeName,
			queue:        queue,
			dsn:          dsn,
			closed:       false,
		},
	}

	return c
}

func (c *Consume) Handle(fn func(<-chan amqp.Delivery), threads int) error {
	var err error
	if err = c.r.connect(); err != nil {
		return fmt.Errorf("connect error: %v", err)
	}
	msgCh, err := c.consume()
	if err != nil {
		return fmt.Errorf("consume error: %v", err)
	}

	go func() {
		for {
			log.Info("starting consumers...")
			for i := 0; i < threads; i++ {
				go fn(msgCh)
			}

			if <-c.r.done != nil {
				msgCh, err = c.reconnectConsume()
				if err != nil {
					log.Errorf("reconnect error: %v, still working gorutines: %d", err, runtime.NumGoroutine())

					return
				}
			}
			if c.r.closed {
				log.Info("rabbit connection closed as excpected")

				return
			}

			log.Info("...reconnected")
		}
	}()

	return nil
}

func (c *Consume) consume() (<-chan amqp.Delivery, error) {
	queueName, err := c.r.declare()
	if err != nil {
		return nil, err
	}

	// Число сообщений, которые можно подтвердить за раз.
	err = c.r.channel.Qos(c.qos, 0, false)
	if err != nil {
		return nil, fmt.Errorf("setting qos error: %v", err)
	}

	msgCh, err := c.r.channel.Consume(
		queueName,
		c.consumerTag,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("queue consume error: %s", err)
	}

	return msgCh, err
}

func (c *Consume) reconnectConsume() (<-chan amqp.Delivery, error) {
	err := c.r.reconnect()
	if err != nil {
		return nil, err
	}
	msgCh, err := c.consume()
	if err != nil {
		return nil, err
	}

	return msgCh, nil
}

func (c *Consume) Close() {
	log.Info("trying to close consumer")
	c.r.close()
}
