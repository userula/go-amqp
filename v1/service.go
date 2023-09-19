package v1

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type client struct {
	conn    *amqp.Connection
	configs *RabbitConfig
}

func NewRabbitService(conf *RabbitConfig) RabbitInterface {
	rmqClient := &client{
		configs: conf,
	}

	err := conf.Validate()
	if err != nil {
		log.Errorf("%v", err)
		return nil
	}

	err = rmqClient.initRabbitChannel()
	if err != nil {
		log.Errorf("%v", err)
		return nil
	}
	return rmqClient
}

func (c *client) initRabbitChannel() error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.configs.Username, c.configs.Password, c.configs.Host, c.configs.Port))
	if err != nil {
		return fmt.Errorf("[RABBIT] connect fail: %v", err)
	}
	log.Info("[RABBIT] connected!")
	c.conn = conn
	return nil
}

func (c *client) getChannel() (*amqp.Channel, error) {
	if c.conn.IsClosed() {
		log.Error("[RABBIT]: Channel is closed")
		if c.configs != nil {
			log.Info("connecting to rabbit..")
			err := c.initRabbitChannel()
			if err != nil {
				return nil, err
			}
			log.Info("rabbit - done")
		} else {
			return nil, fmt.Errorf("[RABBIT]: configs is nil")
		}
	}
	channel, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (c *client) Send(data interface{}, queue ...string) error {
	channel, err := c.getChannel()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	q := c.configs.QueueName
	if len(queue) > 0 {
		q = queue[0]
	}
	err = publish(data, channel, q)
	if err != nil {
		return err
	}
	return nil
}

func publish(data interface{}, channel *amqp.Channel, queue string) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("[RABBIT] could not to marshal rabbit message: %v", err)
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}

	err = channel.Publish("", queue, false, false, msg)
	if err != nil {
		return fmt.Errorf("[RABBIT] publish message error: %v", err)
	}
	log.Infof("[RABBIT] sent data %v", msg)
	return nil
}
