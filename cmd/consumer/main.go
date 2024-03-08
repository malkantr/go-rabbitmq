package main

import (
	"context"
	"log"
	"time"

	"eventdrivenrabbit/internal"

	"golang.org/x/sync/errgroup"
)

func main(){
	conn, err := internal.ConnectRabbitMQ("dave", "1234", "localhost:5672", "customers")

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	messageBus, err := client.Consume("customers_created", "email-service", false)
	if err != nil {
		panic(err)
	}

	// set a timeout for 15 secs
	ctx := context.Background()

	var blocking chan struct{}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	
	g, ctx :=errgroup.WithContext(ctx)

	// errgroup allows us concurrent tasks
	g.SetLimit(10)

	go func() {
		for message := range messageBus {
			// Spawn a worker
			msg := message
			g.Go(func() error{
				log.Printf("New Message: %v", msg)
				time.Sleep(10*time.Second)
				if err := msg.Ack(false); err !=nil{
					log.Println("Ack message failed")
					return err
				}
				log.Printf("Acknowledge message %s\n", message.MessageId)
				return nil
			})
		}
	}()

	log.Println("Consuming, use CTRL+C to exit")
	<-blocking
}