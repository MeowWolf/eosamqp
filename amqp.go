package eosamqp

import (
	"fmt"
	"strings"

	log "github.com/MeowWolf/eoslog"
	"github.com/streadway/amqp"
)

// Connection is the Connection type from the streadway module
// to be used by clients without having to imporrt streadway
type Connection = amqp.Connection

// Channel is the Channel type from the streadway module
// to be used by clients without having to imporrt streadway
type Channel = amqp.Channel

// Delivery is the Delivery type from the streadway module
// to be used by clients without having to imporrt streadway
type Delivery = amqp.Delivery

// Error is the Error type from the streadway module
// to be used by clients without having to imporrt streadway
type Error = amqp.Error

// ExchangeConfig holds config data for an amqp exchange
type ExchangeConfig struct {
	Name        string
	Type        string
	Durable     bool
	AutoDeleted bool
	Internal    bool
	NoWait      bool
	Arguments   amqp.Table
}

// QueueConfig holds config data for an amqp queue
type QueueConfig struct {
	Name       string
	RoutingKey string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	NoAck      bool
	Arguments  amqp.Table
}

// NewConnection creates and returns a new amqp connection
func NewConnection(brokerURL string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		log.Error.Printf("Failed to connect to RabbitMQ: %s", err)
		return nil, err
	}

	return conn, nil
}

// NewChannel creates and returns a new amqp channel
func NewChannel(
	conn *amqp.Connection,
	exConfig ExchangeConfig,
) (*amqp.Channel, error) {
	if conn == nil {
		return nil, fmt.Errorf("could not create channel, no connection to broker")
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Error.Printf("Failed to open a channel: %s", err)
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		exConfig.Name,
		exConfig.Type,
		exConfig.Durable,
		exConfig.AutoDeleted,
		exConfig.Internal,
		exConfig.NoWait,
		exConfig.Arguments,
	); err != nil {
		return ch, err
	}

	return ch, nil
}

// Consume binds a queue to an exchange and sets up a message consumer on a given channel
func Consume(
	exchangeName string, ch *amqp.Channel, qc QueueConfig) <-chan amqp.Delivery {
	q := declareQueue(ch, qc)
	bindQueue(exchangeName, ch, &q, qc)
	return consume(ch, &q, qc)
}

// PublishAmqpMessage publishes an amqp message
func PublishAmqpMessage(
	exchangeName string,
	ch *amqp.Channel,
	qc QueueConfig,
	payload []byte,
) {
	err := ch.Publish(
		exchangeName,  // exchange
		qc.RoutingKey, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        payload,
		})
	if err != nil {
		log.Error.Printf("Failed to publish a message: %s", err)
	}
}

func declareQueue(ch *amqp.Channel, qc QueueConfig) amqp.Queue {
	q, err := ch.QueueDeclare(
		qc.Name,       // name
		qc.Durable,    // durable
		qc.AutoDelete, // delete when unused
		qc.Exclusive,  // exclusive
		qc.NoWait,     // no-wait
		qc.Arguments,  // arguments
	)
	if err != nil {
		log.Error.Printf("Failed to declare a queue: %s", err)
	}

	return q
}

func bindQueue(exchangeName string, ch *amqp.Channel, q *amqp.Queue, qc QueueConfig) {
	err := ch.QueueBind(
		q.Name,        // queue name
		qc.RoutingKey, // routing key
		exchangeName,  // exchange
		qc.NoWait,
		qc.Arguments,
	)
	if err != nil {
		log.Error.Printf("Failed to bind a queue: %s", err)
	}
}

func consume(ch *amqp.Channel, q *amqp.Queue, qc QueueConfig) <-chan amqp.Delivery {
	messageChan, err := ch.Consume(
		q.Name,       // queue
		"",           // consumer
		qc.NoAck,     // auto ack
		qc.Exclusive, // exclusive
		false,        // no local
		qc.NoWait,    // no wait
		qc.Arguments, // args
	)
	if err != nil {
		log.Error.Printf("Failed to register a consumer: %s", err)
	}

	log.Info.Printf(" [*] Waiting for messages on '" + qc.RoutingKey + "'. To exit press CTRL+C")
	return messageChan
}

// GetBrokerURL constructs an amqp broker url used to create a connection
func GetBrokerURL(useTLS, username, password, host, port string) string {
	var protocol string
	if useTLS != "false" {
		protocol = "amqps://"
	} else {
		protocol = "amqp://"
	}

	return protocol + username + ":" + password + "@" + host + ":" + port + "/"
}

// GetRoutingKeyValueByIndex returns a function that can be used to get specific parts of an amqp topic exchange routing key
func GetRoutingKeyValueByIndex(routingKey string) func(int) string {
	routingKeyParts := strings.Split(routingKey, ".")
	return func(index int) string {
		if index >= len(routingKeyParts) {
			return ""
		}
		return routingKeyParts[index]
	}
}
