package rabbit

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	exchangeType = "fanout"
)

type rabbit struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan *amqp.Error
	closed       bool
	queue        string
	exchangeName string
	bindingKey   string
	dsn          string
}

func (r *rabbit) connect() error {
	log.Info("connecting to rabbitmq")
	var err error
	r.conn, err = amqp.Dial(r.dsn)
	if err != nil {
		return fmt.Errorf("rabbit connection error: %s", err.Error())
	}

	log.Debug("connection successful")
	log.Debug("getting Channel")
	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel error: %s", err)
	}
	r.done = make(chan *amqp.Error)
	r.conn.NotifyClose(r.done)

	return nil
}

func (r *rabbit) declare() (string, error) {
	var err error

	if err = r.channel.ExchangeDeclare(
		r.exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return "", fmt.Errorf("exchange declare error: %v", err)
	}

	queue, err := r.channel.QueueDeclare(
		r.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("queue declare error: %s", err)
	}

	if err = r.channel.QueueBind(
		queue.Name,
		r.bindingKey,
		r.exchangeName,
		false,
		nil,
	); err != nil {
		return "", fmt.Errorf("queue bind error: %s", err)
	}

	if err != nil {
		return "", fmt.Errorf("queue declare error: %v", err)
	}

	return queue.Name, nil
}

func (r *rabbit) reconnect() error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = 3 * time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 60 * time.Second

	b := backoff.WithContext(be, context.Background())
LOOP:
	for {
		if r.closed {
			log.Info("breaking reconnect because closed connect gracefully")

			break LOOP
		}

		d := b.NextBackOff()
		if d == backoff.Stop {
			return fmt.Errorf("stop reconnecting")
		}

		<-time.After(d)
		if err := r.connect(); err != nil {
			log.Errorf("could not connect in reconnect call: %+v", err)

			continue LOOP
		}

		break LOOP
	}

	return nil
}

func (r *rabbit) close() {
	r.closed = true
	r.channel.Close()
	r.conn.Close()
}
