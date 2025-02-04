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
	accessToken             = "EXAMPLE_ACCESS_TOKEN"
	exampleSubscribeSubject = "synternet.example.subject"
)

// RepublishData receives a message on a given subject and republishes it to another subject.
// It takes a context, the service instance, and the data (message) as input arguments.
func PrintData(ctx context.Context, service *pubsub.NatsService, data []byte) error {
	log.Printf("Received message on subject %s: %s\n", exampleSubscribeSubject, string(data))
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
