package eosamqp_test

import "github.com/MeowWolf/eosamqp"

func Example() {

	amqp := eosamqp.Amqp{}

	/*************************************
	 * Connect to broker
	 *************************************/
	brokerURL := eosamqp.GetBrokerURL(true, "username", "password", "host", "port")
	conn, err := amqp.Connect(brokerURL)
	if err != nil {
		// handle error
	}

	/*************************************
	 * Create new channel
	 *************************************/
	ch, err := amqp.NewChannel(conn)
	if err != nil {
		// handle error
	}

	/*************************************
	 * Declare exchange
	 *************************************/
	exchangeConfig := eosamqp.ExchangeConfig{
		Name: "exchange-name",
		Type: "topic",
	}
	err = amqp.DeclareExchange(ch, exchangeConfig)
	if err != nil {
		// handle error
	}

	/*************************************
	 * Publish
	 *************************************/
	queuePublishConfig := eosamqp.QueueConfig{
		Name:       "publish-queue",
		RoutingKey: "a.routing.key",
	}
	message := []byte("outgoing message")
	amqp.Publish(exchangeConfig.Name, ch, queuePublishConfig, message)

	/*************************************
	 * Consume
	 *************************************/
	queueConsumeConfig := eosamqp.QueueConfig{
		Name:       "consume-queue",
		RoutingKey: "a.routing.key",
	}
	messageStream, err := amqp.Consume(exchangeConfig.Name, ch, queueConsumeConfig)
	if err != nil {
		// handle error
	}

	for amqpMessage := range messageStream {

		// do something with amqpMessage

		// acknowledge the message
		if err := amqpMessage.Ack(false); err != nil {
			// handle error
		}
	}
}
