// based on https://medium.com/@dhanushgopinath/automatically-recovering-rabbitmq-connections-in-go-applications-7795a605ca59
package rabbitmq

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

// MessageBody is the struct for the body passed in the AMQP message. The type will be set on the Request header
type MessageBody struct {
	Data []byte
	Type string
}

// Message is the amqp request to publish
type Message struct {
	Queue         string
	ReplyTo       string
	ContentType   string
	CorrelationID string
	Priority      uint8
	Body          MessageBody
}

type RabbitMQConnectionInterface interface {
	Publish(m *Message) error
}

// Connection is the connection created
type Connection struct {
	name         string
	conn         *amqp.Connection
	dsn          ConnectionDSN
	channel      *amqp.Channel
	exchange     string
	exchangeType string
	queues       []string
	prefetch     uint16
	err          chan error
}

type ConnectionDSN struct {
	Host     string
	Port     int
	User     string
	Password string
	Vhost    string
}

var connectionPool = make(map[string]*Connection)

// NewConnection returns the new connection object
func NewConnection(name, exchange string, exchangeType string, queues []string, dsn ConnectionDSN, prefetch uint16) *Connection {
	if c, ok := connectionPool[name]; ok {
		return c
	}
	c := &Connection{
		exchange:     exchange,
		exchangeType: exchangeType,
		queues:       queues,
		err:          make(chan error),
		dsn:          dsn,
		prefetch:     prefetch,
	}
	connectionPool[name] = c
	return c
}

// GetConnection returns the connection which was instantiated
func GetConnection(name string) *Connection {
	return connectionPool[name]
}

func (c *Connection) Connect() error {
	var err error
	c.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s", c.dsn.User, c.dsn.Password, c.dsn.Host, c.dsn.Port, c.dsn.Vhost))
	if err != nil {
		return fmt.Errorf("Error in creating rabbitmq connection with %s : %s", c.name, err.Error())
	}
	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error)) // Listen to NotifyClose
		c.err <- errors.New("Connection Closed")
	}()
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	c.channel.Qos(
		int(c.prefetch), // prefetch count
		0,               // prefetch size
		false,           // global
	)
	if err := c.channel.ExchangeDeclare(
		c.exchange,     // name
		c.exchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("Error in Exchange Declare: %s", err)
	}
	return nil
}

func (c *Connection) BindQueue() error {
	for _, q := range c.queues {
		if _, err := c.channel.QueueDeclare(q, true, false, false, false, nil); err != nil {
			return fmt.Errorf("error in declaring the queue %s", err)
		}
		if err := c.channel.QueueBind(q, fmt.Sprintf("%s_%s", c.exchange, q), c.exchange, false, nil); err != nil {
			return fmt.Errorf("Queue  Bind error: %s", err)
		}
	}
	return nil
}

// Reconnect reconnects the connection
func (c *Connection) Reconnect() error {
	if err := c.Connect(); err != nil {
		return err
	}
	if err := c.BindQueue(); err != nil {
		return err
	}
	return nil
}

// Consume consumes the messages from the queues and passes it as map of chan of amqp.Delivery
func (c *Connection) Consume() (map[string]<-chan amqp.Delivery, error) {
	m := make(map[string]<-chan amqp.Delivery)
	for _, q := range c.queues {
		deliveries, err := c.channel.Consume(q, "", false, false, false, false, nil)
		if err != nil {
			return nil, err
		}
		m[q] = deliveries
	}
	return m, nil
}

// HandleConsumedDeliveries handles the consumed deliveries from the queues. Should be called only for a consumer connection
func (c *Connection) HandleConsumedDeliveries(q string, delivery <-chan amqp.Delivery, qtdWorkers int, db *sql.DB, fn func(Connection, string, <-chan amqp.Delivery, *sql.DB, int)) {
	for i := 0; i < qtdWorkers; i++ {
		go fn(*c, q, delivery, db, i)
	}
	for {
		if err := <-c.err; err != nil {
			c.Reconnect()
			deliveries, err := c.Consume()
			if err != nil {
				panic(err) // raising panic if consume fails even after reconnecting
			}
			delivery = deliveries[q]
		}
	}
}

// Publish publishes the message to the queue
func (c *Connection) Publish(m *Message) error {
	if err := c.channel.Publish(
		c.exchange, // publish to an exchange
		"",         // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   m.ContentType,
			CorrelationId: m.CorrelationID,
			ReplyTo:       m.ReplyTo,
			Priority:      m.Priority,
			Body:          m.Body.Data,
		},
	); err != nil {
		return err
	}
	return nil
}
