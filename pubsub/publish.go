package pubsub

import (
	"context"
	"encoding/json"
)

func (sn *NatsService) Publish(ctx context.Context, subject string, msg interface{}) error {
	err := sn.nats.conn.Publish(subject, msg.([]byte))
	if err != nil {
		return err
	}
	return nil
}

func (sn *NatsService) PublishAsJSON(ctx context.Context, subject string, msg interface{}) error {
	var err error
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = sn.nats.conn.Publish(subject, jsonData)
	if err != nil {
		return err
	}
	return nil
}
