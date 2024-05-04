package pubsub

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

type Subscription struct {
	Name            string
	Topic           string
	DeadLetterTopic string
}

type Message struct {
	ID        string
	Data      []byte
	Timestamp time.Time
	Ack       func()
}

func (p *Pubsub) Subscriptions() ([]Subscription, error) {
	ctx := context.Background()

	it := p.client.Subscriptions(ctx)

	var subs []Subscription
	for {
		sub, err := it.NextConfig()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		p.subscriptions[sub.ID()] = sub

		dlt := ""
		if sub.DeadLetterPolicy != nil {
			dlt = sub.DeadLetterPolicy.DeadLetterTopic
			dlt = dlt[strings.LastIndex(dlt, "/")+1:]
		}

		subs = append(subs, Subscription{
			Name:            sub.ID(),
			Topic:           sub.Topic.ID(),
			DeadLetterTopic: dlt,
		})
	}

	return subs, nil
}

func (p *Pubsub) Messages(ctx context.Context, subscription string) chan Message {
	// TODO: Handle errors
	sub := p.client.Subscription(subscription)
	sub.ReceiveSettings.MinExtensionPeriod = 20 * time.Second
	sub.ReceiveSettings.MaxExtensionPeriod = 20 * time.Second
	sub.ReceiveSettings.MaxExtension = 20 * time.Second
	sub.ReceiveSettings.NumGoroutines = 1

	existingMsgs := map[string]bool{}
	msg := make(chan Message, 1)
	go sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		if _, ok := existingMsgs[m.ID]; !ok {
			existingMsgs[m.ID] = true
			msg <- Message{
				ID:        m.ID,
				Data:      m.Data,
				Timestamp: m.PublishTime,
				Ack:       m.Ack,
			}
		}
	})

	return msg
}

func (p *Pubsub) Publish(subscription string, content string) error {
	sub := p.subscriptions[subscription]

	result := sub.Topic.Publish(context.Background(), &pubsub.Message{
		Data: []byte(content),
	})

	_, err := result.Get(context.Background())

	return err
}
