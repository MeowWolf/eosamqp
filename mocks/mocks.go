package mocks

import (
	"fmt"

	"github.com/streadway/amqp"
)

// MockDial ...
func MockDial(url string) (*amqp.Connection, error) {
	return &amqp.Connection{}, nil
}

// MockDialError ...
func MockDialError(url string) (*amqp.Connection, error) {
	return &amqp.Connection{}, fmt.Errorf("an error occurred")
}

/*************************************
* Channel
*************************************/

// ChannelMock ...
type ChannelMock struct{}

// ExchangeDeclare ...
func (c *ChannelMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return nil
}

// QueueDeclare ...
func (c *ChannelMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

// QueueBind ...
func (c *ChannelMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return nil
}

// Consume ...
func (c *ChannelMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan amqp.Delivery, error) {
	deliveryChan := make(chan amqp.Delivery)
	return deliveryChan, nil
}

// Publish ...
func (c *ChannelMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

// NotifyClose ...
func (c *ChannelMock) NotifyClose(ch chan *amqp.Error) chan *amqp.Error {
	return ch
}

// Close ...
func (c *ChannelMock) Close() error {
	return nil
}

/*************************************
* Channel with bad ExchangeDeclare
*************************************/

// ChannelWithBadExchangeDeclareMock ...
type ChannelWithBadExchangeDeclareMock struct{ ChannelMock }

// ExchangeDeclare ...
func (c *ChannelWithBadExchangeDeclareMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return fmt.Errorf("an error occurred")
}

/*************************************
* Channel with bad QueueDeclare
*************************************/

// ChannelWithBadQueueDeclareMock ...
type ChannelWithBadQueueDeclareMock struct{ ChannelMock }

// QueueDeclare ...
func (c *ChannelWithBadQueueDeclareMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, fmt.Errorf("an error occurred")
}

/*************************************
* Channel with bad Consume
*************************************/

// ChannelWithBadConsumeMock ...
type ChannelWithBadConsumeMock struct{ ChannelMock }

// Consume ...
func (c *ChannelWithBadConsumeMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan amqp.Delivery, error) {
	deliveryChan := make(chan amqp.Delivery)
	return deliveryChan, fmt.Errorf("an error occurred")
}

/*************************************
* Channel with bad QueueBind
*************************************/

// ChannelWithBadQueueBindMock ...
type ChannelWithBadQueueBindMock struct{ ChannelMock }

// QueueBind ...
func (c *ChannelWithBadQueueBindMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return fmt.Errorf("an error occurred")
}

/*************************************
* Channel with bad Publish
*************************************/

// ChannelWithBadPublishMock ...
type ChannelWithBadPublishMock struct{ ChannelMock }

// Publish ...
func (c *ChannelWithBadPublishMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return fmt.Errorf("an error occurred")
}
