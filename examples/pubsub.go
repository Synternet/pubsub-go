package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Synternet/pubsub-go/pubsub"
	"github.com/nats-io/nats.go"
)

const (
	natsUrl                 = "nats://127.0.0.1:4222"
	userCredsJWT            = "USER_JWT"
	accessToken             = "EXAMPLE_ACCESS_TOKEN"
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
	opts = append(opts, nats.UserJWTAndSeed(userCredsJWT, accessToken))

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
