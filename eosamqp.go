package eosamqp

import (
	"fmt"
	"strings"

	log "github.com/MeowWolf/eoslog"
	"github.com/streadway/amqp"
)

// Connection is the Connection type from the streadway module
// to be used by clients without having to import streadway
type Connection = amqp.Connection

// Channel is the Channel type from the streadway module
// to be used by clients without having to import streadway
type Channel = amqp.Channel

// Delivery is the Delivery type from the streadway module
// to be used by clients without having to import streadway
type Delivery = amqp.Delivery

// Error is the Error type from the streadway module
// to be used by clients without having to import streadway
type Error = amqp.Error

// EOSChannel ...
type EOSChannel interface {
	ExchangeDeclare(string, string, bool, bool, bool, bool, amqp.Table) error
	QueueDeclare(string, bool, bool, bool, bool, amqp.Table) (amqp.Queue, error)
	QueueBind(string, string, string, bool, amqp.Table) error
	Consume(string, string, bool, bool, bool, bool, amqp.Table) (<-chan Delivery, error)
	Publish(string, string, bool, bool, amqp.Publishing) error
}

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

type Amqp struct {
	dial     func(url string) (*Connection, error)
	logError func(format string, v ...interface{})
	logInfo  func(format string, v ...interface{})

	conn *amqp.Connection
}

type deps struct {
	dial     func(url string) (*Connection, error)
	logError func(format string, v ...interface{})
	logInfo  func(format string, v ...interface{})
}

// New creates a new instance of Amqp struct
func New(d *deps) Amqp {
	a := Amqp{}
	if d != nil && d.dial != nil {
		a.dial = d.dial
	} else {
		a.dial = amqp.Dial
	}
	if d != nil && d.logError != nil {
		a.logError = d.logError
	} else {
		a.logError = log.Error.Printf
	}
	if d != nil && d.logInfo != nil {
		a.logInfo = d.logInfo
	} else {
		a.logInfo = log.Info.Printf
	}

	return a
}

// Connect creates a new amqp connection
func (a *Amqp) Connect(brokerURL string) (*amqp.Connection, error) {
	conn, err := a.dial(brokerURL)
	if err != nil {
		a.conn = nil
		a.logError("failed to connect to RabbitMQ: %s", err)
		return nil, err
	}

	a.conn = conn
	return conn, nil
}

// NewChannel creates and returns a new amqp channel
func (a *Amqp) NewChannel(
	conn *amqp.Connection,
	config ExchangeConfig,
) (EOSChannel, error) {
	if conn == nil {
		return nil, fmt.Errorf("could not create channel, no connection to broker")
	}

	ch, err := conn.Channel()
	if err != nil {
		a.logError("failed to open a channel: %s", err)
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		config.Name,
		config.Type,
		config.Durable,
		config.AutoDeleted,
		config.Internal,
		config.NoWait,
		config.Arguments,
	); err != nil {
		a.logError("could not declare exchange: %s", err)
		return nil, err
	}

	return ch, nil
}

// Consume binds a queue to an exchange and sets up a message consumer on a given channel
func (a *Amqp) Consume(
	exchangeName string,
	ch EOSChannel,
	qc QueueConfig,
) (<-chan Delivery, error) {
	q, err := a.declareQueue(ch, qc)
	if err != nil {
		a.logError("could not declare queue: %s", err)
		return nil, err
	}

	if err := a.bindQueue(exchangeName, ch, q, qc); err != nil {
		a.logError("failed to register a consumer: %s", err)
		return nil, err
	}

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
		a.logError("failed to register a consumer: %s", err)
		return nil, err
	}

	a.logInfo(" [*] Waiting for messages on '" + qc.RoutingKey + "'. To exit press CTRL+C")
	return messageChan, nil
}

// Publish publishes an amqp message
func (a *Amqp) Publish(
	exchangeName string,
	ch EOSChannel,
	qc QueueConfig,
	payload []byte,
) error {
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
		a.logError("failed to publish a message: %s", err)
		return err
	}
	return nil
}

func (a *Amqp) declareQueue(ch EOSChannel, qc QueueConfig) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		qc.Name,       // name
		qc.Durable,    // durable
		qc.AutoDelete, // delete when unused
		qc.Exclusive,  // exclusive
		qc.NoWait,     // no-wait
		qc.Arguments,  // arguments
	)
	if err != nil {
		a.logError("failed to declare a queue: %s", err)
		return nil, err
	}

	return &q, nil
}

func (a *Amqp) bindQueue(exchangeName string, ch EOSChannel, q *amqp.Queue, qc QueueConfig) error {
	err := ch.QueueBind(
		q.Name,        // queue name
		qc.RoutingKey, // routing key
		exchangeName,  // exchange
		qc.NoWait,
		qc.Arguments,
	)

	if err != nil {
		a.logError("failed to bind a queue: %s", err)
		return err
	}

	return nil
}

// GetBrokerURL constructs an amqp broker url used to create a connection
func GetBrokerURL(useTLS bool, username, password, host, port string) string {
	var protocol string
	if useTLS {
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
