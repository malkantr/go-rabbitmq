package main

import (
	"context"
	"log"
	"time"

	"eventdrivenrabbit/internal"

	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

func main(){
	conn, err := internal.ConnectRabbitMQ("dave", "1234", "localhost:5671", "customers",
	"/home/path/tls-gen/basic/result/ca_certificate.pem",
	"/home/path/tls-gen/basic/result/client_PCNAME_certificate.pem",
	"/home/path/tls-gen/basic/result/client_PCNAME_key.pem",
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	publishConn, err := internal.ConnectRabbitMQ("dave","1234", "localhost:5671", "customers",
	"/home/path/tls-gen/basic/result/ca_certificate.pem",
	"/home/path/tls-gen/basic/result/client_PCNAME_certificate.pem",
	"/home/path/tls-gen/basic/result/client_PCNAME_key.pem",
	)
	if err != nil {
		panic(err)
	}
	defer publishConn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	publishClient, err := internal.NewRabbitMQClient(publishConn)
	if err != nil {
		panic(err)
	}
	defer publishClient.Close()

	queue, err := client.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	if err := client.CreateBinding(queue.Name, "", "customer_events"); err != nil {
		panic(err)
	}

	messageBus, err := client.Consume(queue.Name, "email-service", false)
	if err != nil {
		panic(err)
	}
	// set a timeout for 15 secs
	ctx := context.Background()

	var blocking chan struct{}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	
	g, ctx :=errgroup.WithContext(ctx)

	// Apply a Hard limit on the server
	if err := client.ApplyQos(10, 0, true); err != nil {
		panic(err)
	}

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

				if err := publishClient.Send(ctx, "customer_callbacks", msg.ReplyTo,amqp091.Publishing{
					ContentType: "text/plain",
					DeliveryMode: amqp091.Persistent,
					Body: []byte("RPC COMPLETE"),
					CorrelationId: msg.CorrelationId,
				}); err != nil{
					panic(err)
				} 

				log.Printf("Acknowledge message %s\n", message.MessageId)
				return nil
			})
		}
	}()

	log.Println("Consuming, use CTRL+C to exit")
	<-blocking
}