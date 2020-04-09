package mq

import (
	"errors"
	"github.com/streadway/amqp"
	"go-recruitment-spider/config"
	"log"
)

type MySender struct {
	queueName string
	conn      *amqp.Connection
}

func NewMySender(queueName string) *MySender {
	conn, err := amqp.Dial(config.GetTomlConfig().Rabbit.Addr)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}
	return &MySender{queueName: queueName, conn: conn}
}

func (m *MySender) GetQueue() (*amqp.Queue, error) {
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

func (m *MySender) Send(msgCh <-chan string) {
	ch, err := m.conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", err.Error())
	}
	defer ch.Close()
	q, _ := m.GetQueue()
	for {
		select {
		case msg := <-msgCh:
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msg),
				})
			//log.Printf(" [x] Sent %s", msg)
			if err != nil {
				log.Println("Failed to publish message: ", msg)
			}
		}
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
