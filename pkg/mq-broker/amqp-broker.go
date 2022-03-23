package mq_broker

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/valinurovam/garagemq/config"
	"github.com/valinurovam/garagemq/metrics"
	"github.com/valinurovam/garagemq/server"
	"math"
	"os"
	"time"
)

type Config struct {
	MqHost     string `json:"mqHost"`
	MqPort     string `json:"mqPort"`
	MqUsername string `json:"mqUsername"`
	MqPassword string `json:"mqPassword"`
	conn       *amqp.Connection
}

type Error struct {
	amqp.Error
}

func InitConfig() *Config {
	return loadConfig()
}

func loadConfig() *Config {
	var config Config
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config
}

func (c *Config) ServerStart() error {
	cfg, _ := config.CreateDefault()
	metrics.NewTrackRegistry(15, time.Second, false)
	srv := server.NewServer(c.MqHost, c.MqPort, cfg.Proto, cfg)
	srv.Start()
	return errors.New("Amqp server - stop")
}

func (c *Config) connection() (*Config, error) {
	dialConfig := amqp.Config{}
	dialConfig.FrameSize = 0
	dialConfig.Dial = amqp.DefaultDial(time.Second * 10)
	dialConfig.Heartbeat = 10 * time.Second
	connection, err := amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.MqUsername, c.MqPassword, c.MqHost, c.MqPort), dialConfig)
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
