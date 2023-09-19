package v1

import "github.com/streadway/amqp"

type RabbitInterface interface {
	Send(data interface{}, queue ...string) error
	initRabbitChannel() error
	getChannel() (*amqp.Channel, error)
}
