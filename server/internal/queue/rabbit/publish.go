package rabbit

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Publish struct {
	r *rabbit
}

func NewPublish(dsn, exchangeName, queue string) *Publish {
	p := &Publish{
		r: &rabbit{
			exchangeName: exchangeName,
			queue:        queue,
			dsn:          dsn,
			closed:       false,
		},
	}

	return p
}

func (p *Publish) Connect() error {
	if err := p.r.connect(); err != nil {
		return fmt.Errorf("connect error: %v", err)
	}

	_, err := p.r.declare()
	if err != nil {
		return err
	}

	go func() {
		for {
			log.Info("starting reconnecting gorutine...")

			if <-p.r.done != nil {
				err := p.r.reconnect()
				if err != nil {
					log.Errorf("reconnect error (stop reconnecting): %v", err)

					return
				}
			}
			if p.r.closed {
				log.Info("rabbit connection closed as excpected")

				return
			}

			log.Info("...reconnected")
		}
	}()

	return nil
}

func (p *Publish) Publish(body interface{}, routingKey string) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	log.Debugf("publish body: %s", b)

	ap := amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	}

	err = p.r.channel.Publish(
		p.r.exchangeName, // exchange
		routingKey,       // routing key
		false,            // mandatory
		false,            // immediate
		ap,
	)

	return err
}

func (p *Publish) Close() {
	log.Info("trying to close publisher")
	p.r.close()
}
