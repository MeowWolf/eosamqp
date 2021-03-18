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

// Publishing is the Error type from the streadway module
// to be used by clients without having to import streadway
type Publishing = amqp.Publishing

// Queue is the Error type from the streadway module
// to be used by clients without having to import streadway
type Queue = amqp.Queue

// Table is the Error type from the streadway module
// to be used by clients without having to import streadway
type Table = amqp.Table

// EOSChannel ...
type EOSChannel interface {
	ExchangeDeclare(string, string, bool, bool, bool, bool, Table) error
	QueueDeclare(string, bool, bool, bool, bool, Table) (Queue, error)
	QueueBind(string, string, string, bool, Table) error
	Consume(string, string, bool, bool, bool, bool, Table) (<-chan Delivery, error)
	Publish(string, string, bool, bool, Publishing) error
	NotifyClose(chan *Error) chan *Error
}

// ExchangeConfig holds config data for an amqp exchange
type ExchangeConfig struct {
	Name        string
	Type        string
	Durable     bool
	AutoDeleted bool
	Internal    bool
	NoWait      bool
	Arguments   Table
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
	Arguments  Table
}

// Amqp contains a pointer to the amqp connection
type Amqp struct {
	dial     func(url string) (*Connection, error)
	logError func(format string, v ...interface{})
	logInfo  func(format string, v ...interface{})

	conn *Connection
}

// Deps is the optional dependencies struct that can be injected into a New Amqp struct
type Deps struct {
	Dial     func(url string) (*Connection, error)
	LogError func(format string, v ...interface{})
	LogInfo  func(format string, v ...interface{})
}

// New creates a new instance of Amqp struct
func New(d *Deps) Amqp {
	a := Amqp{}
	if d != nil && d.Dial != nil {
		a.dial = d.Dial
	} else {
		a.dial = amqp.Dial
	}
	if d != nil && d.LogError != nil {
		a.logError = d.LogError
	} else {
		a.logError = log.Error.Printf
	}
	if d != nil && d.LogInfo != nil {
		a.logInfo = d.LogInfo
	} else {
		a.logInfo = log.Info.Printf
	}

	return a
}

// Connect creates a new amqp connection
func (a *Amqp) Connect(brokerURL string) (*Connection, error) {
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
func (a *Amqp) NewChannel(conn *Connection) (EOSChannel, error) {
	if conn == nil {
		return nil, fmt.Errorf("could not create channel, no connection to broker")
	}
	ch, err := conn.Channel()
	if err != nil {
		a.logError("failed to open a channel: %s", err)
		return nil, err
	}

	return ch, nil
}

// DeclareExchange declares an exchange on a channel
func (a *Amqp) DeclareExchange(
	ch EOSChannel,
	config ExchangeConfig,
) error {
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
		return err
	}

	return nil
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
		Publishing{
			ContentType: "text/json",
			Body:        payload,
		})
	if err != nil {
		a.logError("failed to publish a message: %s", err)
		return err
	}
	return nil
}

func (a *Amqp) declareQueue(ch EOSChannel, qc QueueConfig) (*Queue, error) {
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

func (a *Amqp) bindQueue(exchangeName string, ch EOSChannel, q *Queue, qc QueueConfig) error {
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
