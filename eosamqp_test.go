package eosamqp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MeowWolf/eosamqp/mocks"
)

/*************************************
* Connect
*************************************/
func TestAmqpConnection(t *testing.T) {
	d := Deps{
		logError: func(format string, v ...interface{}) {
		},
		dial: mocks.MockDial,
	}
	amqp := New(&d)
	_, err := amqp.Connect("broker.url")
	if err != nil {
		t.Errorf("NewConnection() We should not have an error here")
	}
}

func TestBadAmqpConnection(t *testing.T) {
	want := "failed to connect"
	got := ""

	d := Deps{
		logError: func(format string, v ...interface{}) {
			got = fmt.Sprintf(format, v...)
		},
		dial: mocks.MockDialError,
	}

	amqp := New(&d)
	amqp.Connect("hello")
	if !strings.Contains(got, want) {
		t.Errorf("NewConnection() error - want: '%s', got: '%s'", want, got)
	}
}

/*************************************
* NewChannel
*************************************/
func TestNewChannelWithNilConnection(t *testing.T) {
	want := "could not create channel"
	got := ""

	amqp := New(nil)
	_, err := amqp.NewChannel(nil)
	got = err.Error()
	if !strings.Contains(got, want) {
		t.Errorf("NewChannel() error - want: '%s', got: '%s'", want, got)
	}
}

/*************************************
* DeclareExchange
*************************************/
func TestDeclareExchange(t *testing.T) {
	amqp := New(nil)
	err := amqp.DeclareExchange(&mocks.ChannelMock{}, ExchangeConfig{})
	if err != nil {
		t.Errorf("DeclareExchange() err should be nil here: %s", err)
	}
}

func TestBadDeclareExchange(t *testing.T) {
	want := "could not declare exchange"
	got := ""
	d := Deps{
		logError: func(format string, v ...interface{}) {
			got = fmt.Sprintf(format, v...)
		},
	}
	amqp := New(&d)
	err := amqp.DeclareExchange(&mocks.ChannelWithBadExchangeDeclareMock{}, ExchangeConfig{})
	if !strings.Contains(got, want) {
		t.Errorf("DeclareExchange() error - want: '%s', got: '%s'", want, got)
	}
	if err == nil {
		t.Errorf("DeclareExchange() err should not be nil here")
	}
}

/*************************************
* Consume
*************************************/
func TestConsume(t *testing.T) {
	want := "Waiting for"
	got := ""

	d := Deps{
		logInfo: func(format string, v ...interface{}) {
			got = fmt.Sprintf(format, v...)
		},
	}

	amqp := New(&d)
	deliveryChan, err := amqp.Consume("exchangeName", &mocks.ChannelMock{}, QueueConfig{})
	if !strings.Contains(got, want) {
		t.Errorf("Consume() error - want: '%s', got: '%s'", want, got)
	}
	if err != nil {
		t.Errorf("Consume() err should be nil here: %s", err)
	}
	if deliveryChan == nil {
		t.Errorf("Consume() returned delivery chan should not be nil")
	}
}

func TestConsumeWithBadQueueDeclare(t *testing.T) {
	amqp := New(nil)
	deliveryChan, err := amqp.Consume("exchangeName", &mocks.ChannelWithBadQueueDeclareMock{}, QueueConfig{})
	if err == nil {
		t.Errorf("Consume() err should not be nil here")
	}
	if deliveryChan != nil {
		t.Errorf("Consumer) returned delivery chan that should be nil")
	}
}

func TestConsumeWithBadChannelConsume(t *testing.T) {
	amqp := New(nil)
	deliveryChan, err := amqp.Consume("exchangeName", &mocks.ChannelWithBadConsumeMock{}, QueueConfig{})
	if err == nil {
		t.Errorf("Consume() err should not be nil here")
	}
	if deliveryChan != nil {
		t.Errorf("Consume() returned delivery chan that should be nil")
	}
}

func TestConsumeWithBadQueueBind(t *testing.T) {
	want := "failed to register"
	got := ""

	d := Deps{
		logError: func(format string, v ...interface{}) {
			got = fmt.Sprintf(format, v...)
		},
	}
	amqp := New(&d)
	_, err := amqp.Consume("exchangeName", &mocks.ChannelWithBadQueueBindMock{}, QueueConfig{})
	if !strings.Contains(got, want) {
		t.Errorf("Consume() error - want: '%s', got: '%s'", want, got)
	}
	if err == nil {
		t.Errorf("Consume() err should not be nil here")
	}
}

/*************************************
* Publish
*************************************/
func TestPublish(t *testing.T) {
	amqp := New(nil)
	err := amqp.Publish(
		"exchangeName",
		&mocks.ChannelMock{},
		QueueConfig{},
		[]byte("the message"),
	)
	if err != nil {
		t.Errorf("Publish() err should be nil here")
	}
}

func TestPublishWithBadChannelPublish(t *testing.T) {
	want := "failed to publish"
	got := ""

	d := Deps{
		logError: func(format string, v ...interface{}) {
			got = fmt.Sprintf(format, v...)
		},
	}

	amqp := New(&d)
	err := amqp.Publish(
		"exchangeName",
		&mocks.ChannelWithBadPublishMock{},
		QueueConfig{},
		[]byte("the message"),
	)

	if !strings.Contains(got, want) {
		t.Errorf("Publish() error - want: '%s', got: '%s'", want, got)
	}
	if err == nil {
		t.Errorf("Publish() err should not be nil here")
	}
}

/*************************************
* GetBrokerURL
*************************************/
func TestGetBrokerURLWithoutTLS(t *testing.T) {
	want := "amqp://user:pass@host:5672/"
	got := GetBrokerURL(false, "user", "pass", "host", "5672")

	if !strings.Contains(got, want) {
		t.Errorf("GetBrokerURL() error - want: '%s', got: '%s'", want, got)
	}
}

func TestGetBrokerURLWithTLS(t *testing.T) {
	want := "amqps://user:pass@host:5672/"
	got := GetBrokerURL(true, "user", "pass", "host", "5672")

	if !strings.Contains(got, want) {
		t.Errorf("GetBrokerURL() error - want: '%s', got: '%s'", want, got)
	}
}

/*************************************
* GetRoutingKeyValueByIndex
*************************************/
func TestGetRoutingKeyValueByIndex(t *testing.T) {
	routingKey := "make.sure.your.optics.are.clean"
	byIndex := GetRoutingKeyValueByIndex(routingKey)

	if byIndex(1) != "sure" {
		t.Errorf("GetRoutingKeyValueByIndex() error - want: `sure`, got: `%s`", byIndex(1))
	}

	if byIndex(3) != "optics" {
		t.Errorf("GetRoutingKeyValueByIndex() error - want: 'optics', got: '%s'", byIndex(3))
	}

	if byIndex(5) != "clean" {
		t.Errorf("GetRoutingKeyValueByIndex() error - want: `clean`, got: `%s`", byIndex(5))
	}
}

func TestGetRoutingKeyValueByIndexThatIsTooHigh(t *testing.T) {
	routingKey := "kent!"
	byIndex := GetRoutingKeyValueByIndex(routingKey)

	if byIndex(1) != "" {
		t.Errorf("GetRoutingKeyValueByIndex() error - want: ``, got: `%s`", byIndex(1))
	}
}

/*************************************
* Benchmarks
*************************************/
func BenchmarkPublish(b *testing.B) {
	amqp := New(nil)
	for i := 0; i < b.N; i++ {
		amqp.Publish(
			"exchangeName",
			&mocks.ChannelMock{},
			QueueConfig{},
			[]byte("the message"),
		)
	}
}

func BenchmarkConsume(b *testing.B) {
	d := Deps{
		logInfo: func(format string, v ...interface{}) {},
	}
	amqp := New(&d)

	for i := 0; i < b.N; i++ {
		amqp.Consume("exchangeName", &mocks.ChannelMock{}, QueueConfig{})
	}
}
