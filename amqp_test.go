package eosamqp

// go test . -coverprofile=c.out
// go tool cover -html=c.out
// go test -v -run="none" -bench="Benchmark"

import (
	"fmt"
	"strings"
	"testing"
)

/*************************************
* NewConnection
*************************************/
func TestAmqpConnection(t *testing.T) {
	d := deps{
		logError: func(format string, v ...interface{}) {
		},
		dial: mockDial,
	}
	amqp := New(&d)
	err := amqp.NewConnection("broker.url")
	if err != nil {
		t.Errorf("NewConnection() We should not have an error here")
	}
}

func TestBadAmqpConnection(t *testing.T) {
	want := "failed to connect"
	got := ""

	d := deps{
		logError: func(format string, v ...interface{}) {
			got = fmt.Sprintf(format, v...)
		},
		dial: mockDialError,
	}

	amqp := New(&d)
	amqp.NewConnection("hello")
	if !strings.Contains(got, want) {
		t.Errorf("NewConnection() error - want: %s, got: %s", want, got)
	}
}

/*************************************
* NewChannel
*************************************/
func TestNewChannelWithNilConnection(t *testing.T) {
	want := "could not create channel"
	got := ""

	amqp := New(nil)
	err := amqp.NewChannel(nil, ExchangeConfig{})
	got = err.Error()
	if !strings.Contains(got, want) {
		t.Errorf("NewChannel() error - want: %s, got: %s", want, got)
	}
}

// func TestNewChannelThatWillNotOpen(t *testing.T) {
// 	want := "failed to open"
// 	got := ""

// 	d := deps{
// 		logError: func(format string, v ...interface{}) {
// 			got = fmt.Sprintf(format, v...)
// 		},
// 		dial: mockDialError,
// 	}

// 	amqp := New(&d)
// 	amqp.NewChannel(&badConnectionMock{}, ExchangeConfig{})
// 	if !strings.Contains(got, want) {
// 		t.Errorf("NewChannel() error - want: %s, got: %s", want, got)
// 	}
// }

// func TestNewChanneThatReturnsAnError(t *testing.T) {
// 	want := "could not declare"
// 	got := ""

// 	d := deps{
// 		logError: func(format string, v ...interface{}) {
// 			got = fmt.Sprintf(format, v...)
// 		},
// 		dial: mockDialError,
// 	}

// 	amqp := New(&d)
// 	amqp.NewChannel(&connectionMockWithBadChannel{}, ExchangeConfig{})
// 	if !strings.Contains(got, want) {
// 		t.Errorf("NewChannel() error - want: %s, got: %s", want, got)
// 	}
// }

// func TestNewChanne(t *testing.T) {

// 	d := deps{
// 		dial: mockDialError,
// 	}

// 	amqp := New(&d)
// 	err := amqp.NewChannel(&connectionMock{}, ExchangeConfig{})
// 	if err != nil {
// 		t.Errorf("NewChannel() should not have an error here: %s", err)
// 	}
// }

/*************************************
* Consume
*************************************/
func TestConsume(t *testing.T) {
	amqp := New(nil)
	deliveryChan, err := amqp.Consume("exchangeName", &channelMock{}, QueueConfig{})
	if err != nil {
		t.Errorf("Consumer() err should be nil here: %s", err)
	}
	if deliveryChan == nil {
		t.Errorf("Consumer() returned delivery chan should not be nil")
	}
}

func TestConsumeWithBadQueueDeclare(t *testing.T) {
	amqp := New(nil)
	deliveryChan, err := amqp.Consume("exchangeName", &channelWithBadQueueDeclareMock{}, QueueConfig{})
	if err == nil {
		t.Errorf("Consumer() err should not be nil here")
	}
	if deliveryChan != nil {
		t.Errorf("Consumer() returned delivery chan that should be nil")
	}
}
