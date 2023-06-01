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
go get -u gitlab.com/syntropynet/amberdm/sdk/pubsub-go
```

## Usage
Here is a simple example demonstrating how to subscribe to a data stream and republish the received data to another stream:

### The preferred method of authentication is using an access token from the developer portal.
```go
package main

import (
	"context"
	"github.com/nats-io/nats.go"
	"gitlab.com/syntropynet/amberdm/sdk/pubsub-go/pubsub"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	natsUrl                 = "nats://127.0.0.1"
	accessToken             = "SAAGNJOZTRPYYXG2NJX3ZNGXYUSDYX2BWO447W3SHG6XQ7U66RWHQ3JUXM"
	exampleSubscribeSubject = "example.sub.subject"
)

// RepublishData receives a message on a given subject and republishes it to another subject.
// It takes a context, the service instance, and the data (message) as input arguments.
func PrintData(ctx context.Context, service *pubsub.NatsService, data []byte) error {
	log.Println("Received message on", exampleSubscribeSubject, "subject")
	return nil
}

func main() {
	jwt, _ := pubsub.CreateAppJwt(accessToken)
	println(jwt)
	// Set user credentials and options for NATS connection.
	opts := []nats.Option{}
	opts = append(opts, nats.UserJWTAndSeed(jwt, accessToken))

	// Connect to the NATS server using the provided options.
	service := pubsub.MustConnect(
		pubsub.Config{
			URI:  natsUrl,
			Opts: opts,
		})
	log.Println("Connected to NATS server.")

	// Create a context with a cancel function to control the cancellation of ongoing operations.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Add a handler function to process messages received on the exampleSubscribeSubject.
	service.AddHandler(exampleSubscribeSubject, func(data []byte) error {
		return PrintData(ctx, service, data)
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