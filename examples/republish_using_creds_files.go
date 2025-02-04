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
	natsUrl                 = "nats://127.0.0.1"
	subscriberNatsCredFile  = "./subNats.creds"
	publisherNatsCredFile	= "./pubNats.creds"
	exampleSubscribeSubject = "synternet.example.subject"
	examplePublishSubject   = "publisher.example.subject"
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
	subOpts := []nats.Option{}
	subOpts = append(subOpts, nats.UserCredentials(subscriberNatsCredFile))
	pubOpts := []nats.Option{}
	pubOpts = append(pubOpts, nats.UserCredentials(publisherNatsCredFile))

	// Connect to the NATS server using the provided options.
	pubService := pubsub.MustConnect(
		pubsub.Config{
			URI:  natsUrl,
			Opts: pubOpts,
		})
	log.Println("Publisher connected to NATS server.")

	subService := pubsub.MustConnect(
		pubsub.Config{
			URI:  natsUrl,
			Opts: subOpts,
		})
	log.Println("Subscriber connected to NATS server.")

	// Create a context with a cancel function to control the cancellation of ongoing operations.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Add a handler function to process messages received on the exampleSubscribeSubject.
	subService.AddHandler(exampleSubscribeSubject, func(data []byte) error {
		return RepublishData(ctx, pubService, data)
	})

	// Set up signal handling to gracefully shut down when receiving SIGINT or SIGTERM signals.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	// Start serving messages and processing them using the registered handler function.
	subService.Serve(ctx)
}
