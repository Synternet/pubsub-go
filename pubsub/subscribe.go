package pubsub

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
)

func (sn *NatsService) subscribeHandlers(ctx context.Context) error {
	for subject, msgHandler := range sn.nats.handlers {
		sub, err := sn.subscribe(ctx, subject, msgHandler)
		if err != nil {
			return fmt.Errorf("failed to subscribe handler to subject %s: %s", subject, err.Error())
		}
		sn.nats.subscriptions = append(sn.nats.subscriptions, sub)
	}
	return nil
}

func (sn *NatsService) unsubscribeHandlers() error {
	var errors MultiError
	for _, subscription := range sn.nats.subscriptions {
		err := subscription.Unsubscribe()
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to unsubscribe: %w", err))
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (sn *NatsService) subscribe(ctx context.Context, subject string, msgHandler any) (*nats.Subscription, error) {
	switch handler := msgHandler.(type) {
	case Handler:
		return sn.nats.conn.Subscribe(subject, func(msg *nats.Msg) {
			err := handler(msg.Data)
			if err != nil {
				return
			}
		})
	case HandlerWithSubject:
		return sn.nats.conn.Subscribe(subject, func(msg *nats.Msg) {
			err := handler(msg.Data, msg.Subject)
			if err != nil {
				return
			}
		})
	}
	return nil, errors.New("invalid handler function")
}
