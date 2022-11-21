package rabbitmq

import (
	"context"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type Client struct {
	log             *zap.SugaredLogger
	rabbitmqClient  *amqp.Connection
	rabbitmqChannel *amqp.Channel
	rabbitmqIsConn  atomic.Bool

	address      string
	queue        string
	kind         string
	exchange     string
	routingKey   string
	contentType  string
	consumerFunc func(msg amqp.Delivery)
}

func New(
	address string,
	kind string,
	queue string,
	exchange string,
	routingKey string,
	contentType string,
	consumerFunc func(msg amqp.Delivery),
) *Client {
	client := Client{
		log:          zap.S().With("module", "pkg.rabbitmq"),
		address:      address,
		kind:         kind,
		queue:        queue,
		exchange:     exchange,
		routingKey:   routingKey,
		contentType:  contentType,
		consumerFunc: consumerFunc,
	}
	client.rabbitmqIsConn.Store(false)
	go client.Connect()

	return &client
}

func (c *Client) Connect() {
	defer func() {
		if err := recover(); err != nil {
			c.log.Errorf("consume err: %+v", err)
			time.Sleep(3 * time.Second)
			go c.Connect()
			// for !c.rabbitmqIsConn.Load() {
			// 	c.log.Infof("rabbitmq reconnect")
			// 	time.Sleep(3 * time.Second)
			// 	c.Connect()
			// }
		}
	}()

	var err error
	c.rabbitmqClient, err = amqp.Dial(c.address)
	if err != nil {
		c.log.Errorf("rabbitmq connect error: %+v", err)
		return
	}

	c.rabbitmqChannel, err = c.rabbitmqClient.Channel()
	if err != nil {
		c.log.Errorf("rabbitmq open channel error: %+v", err)
		return
	}

	err = c.rabbitmqChannel.Qos(200, 0, false)
	if err != nil {
		c.log.Errorf("rabbitmq Qos error: %+v", err)
		return
	}

	c.log.Infof("rabbitmq ExchangeDeclare exchange=%s kind=%s", c.exchange, c.kind)
	err = c.rabbitmqChannel.ExchangeDeclare(c.exchange, c.kind, true, false, false, false, nil)
	if err != nil {
		c.log.Errorf("rabbitmq ExchangeDeclare error: %+v", err)
		return
	}

	c.log.Infof("rabbitmq QueueDeclare queue=%s", c.queue)
	q, err := c.rabbitmqChannel.QueueDeclare(c.queue, true, false, false, false, nil)
	if err != nil {
		c.log.Errorf("rabbitmq QueueDeclare error: %+v", err)
		return
	}

	c.log.Infof("rabbitmq QueueBind name=%s routingKey=%s exchange=%s", q.Name, c.routingKey, c.exchange)
	err = c.rabbitmqChannel.QueueBind(q.Name, c.routingKey, c.exchange, false, nil)
	if err != nil {
		c.log.Errorf("rabbitmq QueueBind error: %+v", err)
		return
	}

	c.rabbitmqIsConn.Store(true)

	var msgs <-chan amqp.Delivery
	if c.consumerFunc != nil {
		msgs, err = c.rabbitmqChannel.Consume(c.queue, "", false, false, false, false, nil)
		if err != nil {
			c.log.Errorf("rabbitmq Consume error: %s", err)
		}
	} else {
		msgs = make(<-chan amqp.Delivery)
	}

	closeChan := make(chan *amqp.Error, 1)
	notifyClose := c.rabbitmqChannel.NotifyClose(closeChan)
	closeFlag := false
	for {
		select {
		case <-notifyClose:
			c.rabbitmqIsConn.Store(false)
			c.log.Errorf("rabbitmq notifyClose")
			close(closeChan)
			closeFlag = true
			c.Connect()
		case msg := <-msgs:
			c.consumerFunc(msg)
		}
		if closeFlag {
			break
		}
	}
}

func (c *Client) Produce(body []byte) error {
	return c.rabbitmqChannel.Publish(c.exchange, c.routingKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  c.contentType,
		Body:         body,
	})
}

type ProduceParam struct {
	Exchange    string
	RoutingKey  string
	Body        []byte
	ContentType string
}

func (c *Client) ProduceWithParam(param *ProduceParam) error {
	return c.rabbitmqChannel.Publish(param.Exchange, param.RoutingKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  param.ContentType,
		Body:         param.Body,
	})
}

type Params struct {
	Kind        string
	Queue       string
	Exchange    string
	RoutingKey  string
	ContentType string
}

func (c *Client) Shutdown(ctx context.Context) error {
	return nil
}
