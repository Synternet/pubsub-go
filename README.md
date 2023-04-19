# Syntropy PubSub-Go

`syntropy-pubsub-go` is a Golang library for the Syntropy DataMesh project that allows you to subscribe to existing data streams or publish new ones. This library is built on top of the NATS messaging system and provides a convenient way to integrate your Golang applications with the Syntropy DataMesh platform.

## Features

- Subscribe to existing data streams
- Publish new data streams
- Support for JSON messages
- Customizable connection options

## Installation

To install the library, use the following command:

```bash
go get -u github.com/SyntropyNet/pubsub-go
```

## Usage
Here is a simple example demonstrating how to subscribe to a data stream and republish the received data to another stream:

```go
package main

import (
	"context"
	"github.com/SyntropyNet/pubsub-go/pubsub"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	NatsUrl                 = "nats://127.0.0.1:4222"
	UserCredsJWT            = "USER_JWT"
	UserCredsSeed           = "CREDS_SEED"
	exampleSubscribeSubject = "example.sub.subject"
	examplePublishSubject   = "example.pub.subject"
)

// RepublishData receives a message on a given subject and republishes it to another subject.
// It takes a context, the service instance, and the data (message) as input arguments.
func RepublishData(ctx context.Context, service *pubsub.NatsService, data []byte) error {
	log.Println("Received message on", exampleSubscribeSubject, "subject")
	err := service.Publish(ctx, examplePublishSubject, data)
	if err != nil {
		return err
	}
	log.Println("Published message on", examplePublishSubject, "subject")
	return nil
}

func main() {
	// Set user credentials and options for NATS connection.
	opts := []nats.Option{}
	opts = append(opts, nats.UserJWTAndSeed(UserCredsJWT, UserCredsSeed))

	// Connect to the NATS server using the provided options.
	service := pubsub.MustConnect(
		pubsub.Config{
			URI:  NatsUrl,
			Opts: opts,
		})
	log.Println("Connected to NATS server.")

	// Create a context with a cancel function to control the cancellation of ongoing operations.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Add a handler function to process messages received on the exampleSubscribeSubject.
	service.AddHandler(exampleSubscribeSubject, func(data []byte) error {
		return RepublishData(ctx, service, data)
	})

	// Set up signal handling to gracefully shut down when receiving SIGINT or SIGTERM signals.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	// Start serving messages and processing them using the registered handler function.
	service.Serve(ctx)
}
```

This example demonstrates how to connect to a NATS server, subscribe to a subject, and republish received messages to another subject using the Syntropy PubSub-Go library.

## License
This project is licensed under the MIT License.