package mq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go-recruitment-spider/config"
	"log"
)

type MyReceiver struct {
	queueName string
	conn      *amqp.Connection
}

func NewMyReceiver(queueName string) *MyReceiver {
	conn, err := amqp.Dial(config.GetTomlConfig().Rabbit.Addr)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}
	return &MyReceiver{queueName: queueName, conn: conn}
}

func (m *MyReceiver) GetQueue() (*amqp.Queue, error) {
	ch, err := m.conn.Channel()
	if err != nil {
		return nil, errors.New("Failed to open a channel")
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		m.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	//log.Println(q.Messages)
	if err != nil {
		return nil, errors.New("Failed to declare a queue")
	}
	return &q, nil
}

func (m *MyReceiver) GetDelivery() (<-chan amqp.Delivery, error) {
	//url := genUrl(m.cfg)
	//conn, err := amqp.Dial(url)
	//if err != nil {
	//	return nil, errors.New("Failed to connect to RabbitMQ")
	//}
	//defer conn.Close()
	ch, err := m.conn.Channel()
	if err != nil {
		return nil, errors.New("Failed to open a channel")
	}
	//defer ch.Close()
	q, err := ch.QueueDeclare(
		m.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	//log.Println(q.Messages)
	if err != nil {
		return nil, errors.New("Failed to declare a queue")
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, errors.New("Failed to set QoS")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		return nil, errors.New("Failed to get Consume")
	}
	return msgs, nil
}
