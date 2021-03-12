package eosamqp

import (
	"fmt"

	"github.com/streadway/amqp"
)

func mockDial(url string) (*amqp.Connection, error) {
	return &amqp.Connection{}, nil
}

func mockDialError(url string) (*amqp.Connection, error) {
	return &amqp.Connection{}, fmt.Errorf("an error occurred")
}

/*************************************
* Channel
*************************************/
type channelMock struct{}

func (c *channelMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *channelMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (c *channelMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *channelMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan Delivery, error) {
	deliveryChan := make(chan Delivery)
	return deliveryChan, nil
}

func (c *channelMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

/*************************************
* Channel with bad QueueDeclare
*************************************/
type channelWithBadQueueDeclareMock struct{}

func (c *channelWithBadQueueDeclareMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *channelWithBadQueueDeclareMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, fmt.Errorf("an error occurred")
}

func (c *channelWithBadQueueDeclareMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *channelWithBadQueueDeclareMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan Delivery, error) {
	deliveryChan := make(chan Delivery)
	return deliveryChan, nil
}
func (c *channelWithBadQueueDeclareMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

/*************************************
* Channel with bad Consume
*************************************/
type channelWithBadConsumeMock struct{}

func (c *channelWithBadConsumeMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *channelWithBadConsumeMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (c *channelWithBadConsumeMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *channelWithBadConsumeMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan Delivery, error) {
	deliveryChan := make(chan Delivery)
	return deliveryChan, fmt.Errorf("an error occurred")
}
func (c *channelWithBadConsumeMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

/*************************************
* Channel with bad QueueBind
*************************************/
type channelWithBadQueueBindMock struct{}

func (c *channelWithBadQueueBindMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *channelWithBadQueueBindMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (c *channelWithBadQueueBindMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return fmt.Errorf("an error occurred")
}

func (c *channelWithBadQueueBindMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan Delivery, error) {
	deliveryChan := make(chan Delivery)
	return deliveryChan, nil
}

func (c *channelWithBadQueueBindMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

/*************************************
* Channel with bad Publish
*************************************/
type channelWithBadPublishMock struct{}

func (c *channelWithBadPublishMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *channelWithBadPublishMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (c *channelWithBadPublishMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *channelWithBadPublishMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan Delivery, error) {
	deliveryChan := make(chan Delivery)
	return deliveryChan, nil
}

func (c *channelWithBadPublishMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return fmt.Errorf("an error occurred")
}
