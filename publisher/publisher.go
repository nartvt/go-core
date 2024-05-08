package publisher

import (
	"encoding/json"
	"errors"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	topic   string
	host    string
	natsCli *nats.Conn
	opts    []nats.Option
}

func NewPublisher(host string, topic string, opts ...nats.Option) (*Publisher, error) {
	conn, err := nats.Connect(host, opts...)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		topic:   topic,
		host:    host,
		natsCli: conn,
		opts:    opts,
	}, nil
}

func (p *Publisher) Publish(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = p.natsCli.Publish(p.topic, b)
	if err == nats.ErrConnectionClosed {
		// retry if conn is close
		conn, err := nats.Connect(p.host, p.opts...)
		if err != nil {
			return errors.New("RESTART_NATS_CONN_FAILED")
		}
		p.natsCli = conn
		return p.natsCli.Publish(p.topic, b)
	}

	return err
}

func (p *Publisher) PublishWithTopic(topic string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = p.natsCli.Publish(topic, b)
	if err == nats.ErrConnectionClosed {
		// retry if conn is close
		conn, err := nats.Connect(p.host, p.opts...)
		if err != nil {
			return errors.New("RESTART_NATS_CONN_FAILED")
		}
		p.natsCli = conn
		return p.natsCli.Publish(p.topic, b)
	}

	return err
}
