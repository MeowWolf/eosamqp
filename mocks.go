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
* Connection
*************************************/
type connectionMock struct{}

func (c *connectionMock) Channel() (EOSChannel, error) {
	return &channelMock{}, nil
}

type connectionMockWithBadChannel struct{}

func (c *connectionMockWithBadChannel) Channel() (EOSChannel, error) {
	return &badChannelMock{}, nil
}

type badConnectionMock struct{}

func (c *badConnectionMock) Channel() (EOSChannel, error) {
	return &amqp.Channel{}, fmt.Errorf("an error occurred")
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

/*************************************
* Bad Channel
*************************************/
type badChannelMock struct{}

func (c *badChannelMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return fmt.Errorf("an error occurred")
}

func (c *badChannelMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (c *badChannelMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *badChannelMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan Delivery, error) {
	deliveryChan := make(chan Delivery)
	return deliveryChan, nil
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
