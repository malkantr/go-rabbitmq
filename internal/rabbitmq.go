package internal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	// The connection used by the client // tcp connection
	conn *amqp.Connection		
	// Channel is used to process / send messages // multiplex connection
	ch *amqp.Channel			
}

func ConnectRabbitMQ(username, password, host, vhost, caCert, clientCert, clientKey string) (*amqp.Connection, error) {
	ca, err := os.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	// Load keypair
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	// Add the RootCA to the cert pool
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(ca)

	tlsCfg := &tls.Config{
		RootCAs: rootCAs,
		Certificates: []tls.Certificate{cert},
	}
	
	return amqp.DialTLS(fmt.Sprintf("amqps://%s:%s@%s/%s", username, password, host, vhost), tlsCfg)
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	if err := ch.Confirm(false); err != nil{
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch: ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error) {
	q, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	if err != nil {
		return amqp.Queue{}, nil
	}
	return q, err
}

// CreateBinding will bind the current channel to the given exchange using the routingkey provided
func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	// leaving nowait false, having nowait set to false will make the channel return an error if its fails to bind
	// name = name of the queue
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

// Send is used to publish payloads onto an exchange with the given routingkey
func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	confirmation, err :=  rc.ch.PublishWithDeferredConfirmWithContext(ctx, 
		exchange, 
		routingKey, 
		// Mandatory is used to determine if an error should be returned upon failure
		true,
		// Immediate
		// immediate Removed in MQ 3 or ...
		false,
		options,
	)

	if err != nil {
		return err
	}
	log.Println(confirmation.Wait())
	return nil
}

// Consume is used to consume a queue
func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

// ApplyQos
// prefetch count - an integer on how many unacknowledge messages the server can send
// prefecth size  - is int of how many bytes
// global		  - Determines if the rule should be applied globally or not
func (rc RabbitClient) ApplyQos(count, size int, global bool) error {
	return rc.ch.Qos(count, size, global)
}