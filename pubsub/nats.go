package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// NatsService represents the main struct for the NATS service.
type NatsService struct {
	nats *natsConnection
}

// Handler represents a function that takes a byte slice as input and returns an error.
type Handler func([]byte) error

// Config is a struct that holds the configuration settings for connecting to the NATS server.
type Config struct {
	URI  string
	Opts []nats.Option
}

// natsConnection represents the underlying connection to NATS, its subscriptions, handlers, and error channel.
type natsConnection struct {
	conn          *nats.Conn
	subscriptions []*nats.Subscription
	handlers      map[string]Handler
	errorCh       <-chan error
}

// MustConnect connects to the NATS server using the provided configuration and panics if it fails to connect.
func MustConnect(cfg Config) *NatsService {
	nats, err := Connect(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to connect nats: %w", err))
	}
	return nats
}

// Connect connects to the NATS server using the provided configuration and returns a new NatsService instance.
func Connect(cfg Config) (*NatsService, error) {
	conn, errorCh, err := connect(cfg)
	if err != nil {
		return &NatsService{}, err
	}

	s := newService(conn, errorCh)

	return s, nil
}

// connect is a helper function that establishes a connection to the NATS server and returns the connection, error channel, and any error encountered.
func connect(cfg Config) (*nats.Conn, <-chan error, error) {
	errorCh := make(chan error, 1)

	// Configure the reconnect settings
	totalWait := 60 * time.Minute
	reconnectDelay := time.Second

	cfg.Opts = append(cfg.Opts, nats.ReconnectWait(reconnectDelay))
	cfg.Opts = append(cfg.Opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	cfg.Opts = append(cfg.Opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		reason := fmt.Errorf("disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
		errorCh <- reason
	}))
	cfg.Opts = append(cfg.Opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		reason := fmt.Errorf("reconnected [%s]", nc.ConnectedUrl())
		errorCh <- reason
	}))
	cfg.Opts = append(cfg.Opts, nats.ClosedHandler(func(nc *nats.Conn) {
		reason := fmt.Errorf("exiting: %v", nc.LastError())
		errorCh <- reason
	}))

	natsConn, err := nats.Connect(cfg.URI, cfg.Opts...)
	if err != nil {
		errorCh <- err
		return natsConn, errorCh, err
	}
	return natsConn, errorCh, err
}

// newService is a helper function that creates a new NatsService instance for testing purposes.
func newService(conn *nats.Conn, errorCh <-chan error) *NatsService {
	return &NatsService{
		nats: &natsConnection{
			conn:     conn,
			handlers: map[string]Handler{},
			errorCh:  errorCh,
		},
	}
}

// AddHandler registers a new handler function for the given subject in the NatsService instance.
func (sn *NatsService) AddHandler(subject string, handlerFn Handler) {
	// Check if a handler with the same subject is already registered, and panic if it is
	if _, ok := sn.nats.handlers[subject]; ok {
		panic(fmt.Errorf("handler with subject %s already registered", subject))
	}

	// Register the handler function for the given subject
	sn.nats.handlers[subject] = handlerFn
}

// CloseConnection flushes the connection and closes it.
func (sn *NatsService) CloseConnection() {
	sn.nats.conn.Flush()
	sn.nats.conn.Close()
}

// Serve subscribes to the registered handlers and waits for the context to be done or an error to occur.
func (sn *NatsService) Serve(ctx context.Context) (err error) {
	// Subscribe to the registered handlers
	if err = sn.subscribeHandlers(ctx); err != nil {
		return err
	}

	// Wait for the context to be done or an error to occur
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-sn.nats.errorCh:
	}
	return err
}
