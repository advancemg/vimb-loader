package mq_broker

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"github.com/streadway/amqp"
	"math"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	conn     *amqp.Connection
}

type Error struct {
	amqp.Error
}

func InitConfig() *Config {
	return &Config{
		Host:     utils.GetEnv("RABBITMQ_HOST", "localhost"),
		Port:     utils.GetEnv("RABBITMQ_PORT", "5555"),
		Username: utils.GetEnv("RABBITMQ_USER", "guest"),
		Password: utils.GetEnv("RABBITMQ_PWD", "guest"),
	}
}

func (c *Config) connection() (*Config, error) {
	dialConfig := amqp.Config{}
	dialConfig.FrameSize = 0
	dialConfig.Dial = amqp.DefaultDial(time.Second * 10)
	dialConfig.Heartbeat = 10 * time.Second
	connection, err := amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Username, c.Password, c.Host, c.Port), dialConfig)
	if err != nil {
		return nil, err
	}
	c.conn = connection
	return c, nil
}

func (c *Config) Channel() (*amqp.Channel, error) {
	_, err := c.connection()
	if err != nil {
		return nil, err
	}
	return c.conn.Channel()
}

func (c *Config) ConsumerChannel() (*amqp.Channel, *amqp.Connection, error) {
	_, err := c.connection()
	if err != nil {
		return nil, nil, err
	}
	channel, err := c.conn.Channel()
	return channel, c.conn, err
}

func (c *Config) GetQueueInfo(queue string) (*amqp.Queue, error) {
	_, err := c.connection()
	if err != nil {
		return nil, err
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()
	inspect, err := ch.QueueInspect(queue)
	return &inspect, err
}

func (c *Config) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (*amqp.Queue, error) {
	_, err := c.connection()
	if err != nil {
		return nil, err
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	declareQueue, err := ch.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
	if err != nil {
		return nil, err
	}
	defer ch.Close()
	defer c.conn.Close()
	return &declareQueue, nil
}

func (c *Config) PublishJson(queue string, msg interface{}) error {
	_, err := c.connection()
	if err != nil {
		return err
	}
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	defer c.conn.Close()
	err = ch.Publish(``, queue, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         msgByte,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) DeclareSimpleQueue(name string) error {
	_, err := c.connection()
	if err != nil {
		return err
	}
	defer c.conn.Close()
	var args = make(amqp.Table)
	args["x-message-ttl"] = int64(math.MaxInt64)
	_, err = c.QueueDeclare(name, true, false, false, false, args)
	if err != nil {
		return err
	}
	return nil
}
