package mq

import (
	"errors"
	"github.com/streadway/amqp"
	"go-recruitment-spider/config"
	"log"
)

type mqSession struct {
	Conn  *amqp.Connection
	Ch    *amqp.Channel
	Queue *amqp.Queue
}

func NewSession(queueName string) *mqSession {
	conn, err := amqp.Dial(config.GetTomlConfig().Rabbit.Addr)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel")
	}
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	//log.Println(q.Messages)
	if err != nil {
		log.Fatal("Failed to declare a queue")
	}
	return &mqSession{
		Conn:  conn,
		Ch:    ch,
		Queue: &q,
	}
}

func (s *mqSession) GetDelivery() (<-chan amqp.Delivery, error) {
	err := s.Ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, errors.New("Failed to set QoS")
	}
	msgs, err := s.Ch.Consume(
		s.Queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, errors.New("Failed to get Consume")
	}
	return msgs, nil
}
