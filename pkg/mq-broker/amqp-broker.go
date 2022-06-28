package mq_broker

import (
	"encoding/json"
	"fmt"
	cfg "github.com/advancemg/vimb-loader/internal/config"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/valinurovam/garagemq/config"
	"github.com/valinurovam/garagemq/metrics"
	"github.com/valinurovam/garagemq/server"
	"math"
	"net"
	"time"
)

type Config struct {
	MqHost     string           `json:"Host"`
	MqPort     string           `json:"Port"`
	MqUsername string           `json:"Username"`
	MqPassword string           `json:"Password"`
	conn       *amqp.Connection `json:"-"`
}

type Error struct {
	amqp.Error
}

type QInfo struct {
	Messages  int // count of messages not awaiting acknowledgment
	Consumers int // number of consumers receiving deliveries
}

func InitConfig() *Config {
	return &Config{
		MqHost:     cfg.Config.Amqp.MqHost,
		MqPort:     cfg.Config.Amqp.MqPort,
		MqUsername: cfg.Config.Amqp.MqUsername,
		MqPassword: cfg.Config.Amqp.MqPassword,
	}
}

func (c *Config) Close() error {
	return c.conn.Close()
}

func (c *Config) Ping() bool {
	for {
		endpoint := fmt.Sprintf("%s:%s", c.MqHost, c.MqPort)
		conn, err := net.DialTimeout("tcp", endpoint, time.Second*1)
		if err != nil {
			time.Sleep(time.Second * 2)
			log.PrintLog("vimb-loader", "amqp", "info", "ping amqp endpoint", endpoint, "...")
			continue
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	return true
}

func (c *Config) ServerStart() error {
	cfg, err := config.CreateDefault()
	if err != nil {
		return err
	}
	cfg.Queue.MaxMessagesInRAM = 1000
	cfg.Connection.ChannelsMax = 2000
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
	c, err := c.connection()
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
	c, err := c.connection()
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
	c, err := c.connection()
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
	defer c.conn.Close()
	defer ch.Close()
	return &declareQueue, nil
}

func (c *Config) PublishJson(queue string, msg interface{}) error {
	c, err := c.connection()
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
	//defer c.conn.Close()
	defer ch.Close()
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
	c, err := c.connection()
	if err != nil {
		return err
	}
	//defer c.conn.Close()
	var args = make(amqp.Table)
	args["x-message-ttl"] = int64(math.MaxInt64)
	_, err = c.QueueDeclare(name, true, false, false, false, args)
	if err != nil {
		return err
	}
	return nil
}
