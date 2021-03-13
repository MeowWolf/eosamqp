package mocks

import (
	"fmt"

	"github.com/streadway/amqp"
)

func MockDial(url string) (*amqp.Connection, error) {
	return &amqp.Connection{}, nil
}

func MockDialError(url string) (*amqp.Connection, error) {
	return &amqp.Connection{}, fmt.Errorf("an error occurred")
}

/*************************************
* Channel
*************************************/
type ChannelMock struct{}

func (c *ChannelMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *ChannelMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (c *ChannelMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return nil
}

func (c *ChannelMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan amqp.Delivery, error) {
	deliveryChan := make(chan amqp.Delivery)
	return deliveryChan, nil
}

func (c *ChannelMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

func (ch *ChannelMock) NotifyClose(c chan *amqp.Error) chan *amqp.Error {
	return c
}

/*************************************
* Channel with bad ExchangeDeclare
*************************************/
type ChannelWithBadExchangeDeclareMock struct{ ChannelMock }

func (c *ChannelWithBadExchangeDeclareMock) ExchangeDeclare(name, etype string, durable, autoDeleted, internal, noWait bool, arguments amqp.Table) error {
	return fmt.Errorf("an error occurred")
}

/*************************************
* Channel with bad QueueDeclare
*************************************/
type ChannelWithBadQueueDeclareMock struct{ ChannelMock }

func (c *ChannelWithBadQueueDeclareMock) QueueDeclare(name string, durable, autoDeleted, exclusive, noWait bool, arguments amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, fmt.Errorf("an error occurred")
}

/*************************************
* Channel with bad Consume
*************************************/
type ChannelWithBadConsumeMock struct{ ChannelMock }

func (c *ChannelWithBadConsumeMock) Consume(name, consumer string, noAck, exclusive, noLocal, noWait bool, arguments amqp.Table) (<-chan amqp.Delivery, error) {
	deliveryChan := make(chan amqp.Delivery)
	return deliveryChan, fmt.Errorf("an error occurred")
}

/*************************************
* Channel with bad QueueBind
*************************************/
type ChannelWithBadQueueBindMock struct{ ChannelMock }

func (c *ChannelWithBadQueueBindMock) QueueBind(name, routingKey, exchangeName string, noWait bool, arguments amqp.Table) error {
	return fmt.Errorf("an error occurred")
}

/*************************************
* Channel with bad Publish
*************************************/
type ChannelWithBadPublishMock struct{ ChannelMock }

func (c *ChannelWithBadPublishMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return fmt.Errorf("an error occurred")
}
