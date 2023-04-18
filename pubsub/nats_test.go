package pubsub

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

const (
	testSubject = "test.subject"
	NSHost      = "localhost"
)

func NewNatsServer() (*server.Server, error) {
	opts := server.Options{
		Host:   NSHost,
		Port:   -1,
		NoLog:  true,
		NoSigs: true,
	}

	nServer := test.RunServer(&opts)
	if nServer == nil {
		panic("could not start NATS test server")
	}

	return nServer, nil
}

func TestNatsServiceWithNewNatsServer(t *testing.T) {
	// Run the NATS server
	natsServer, err := NewNatsServer()
	if err != nil {
		t.Fatalf("Error creating NATS server: %v", err)
	}
	defer natsServer.Shutdown()

	// Create a NatsService instance
	cfg := Config{
		URI:  natsServer.ClientURL(),
		Opts: []nats.Option{nats.Name("NATS Test Client")},
	}
	service := MustConnect(cfg)

	// Add a handler function
	handlerCalled := make(chan bool, 1)
	service.AddHandler(testSubject, func(data []byte) error {
		handlerCalled <- true
		return nil
	})

	// Serve the NatsService with a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := service.Serve(ctx); err != nil && err != context.DeadlineExceeded && err != context.Canceled {
			t.Errorf("Error serving NATS: %v", err)
		}
	}()

	// Wait for the service to be ready
	time.Sleep(500 * time.Millisecond)

	// Connect to the NATS server
	nc, err := nats.Connect(natsServer.ClientURL())
	if err != nil {
		t.Fatalf("Error connecting to NATS server: %v", err)
	}
	defer nc.Close()

	// Publish a message
	err = nc.Publish(testSubject, []byte("test message"))
	if err != nil {
		t.Fatalf("Error publishing message: %v", err)
	}

	// Wait for the handler to be called
	select {
	case <-handlerCalled:
	case <-time.After(2 * time.Second):
		t.Error("Handler was not called")
	}
}
