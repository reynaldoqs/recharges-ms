package msgbroker

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type rabbitmqBroker struct {
	conex *amqp.Connection
}

// NewRabbitMqBroker creates a new instance for msgBroker for rabbitmq
func NewRabbitMqBroker(connectionURL string) (MessageBroker, error) {
	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		return nil, errors.Wrap(err, "msgbroker.rabbitmq.NewRabbitMqBroker")
	}

	broker := rabbitmqBroker{
		conex: conn,
	}
	return &broker, nil
}

func (r *rabbitmqBroker) PublishOnQueue(body []byte, queueName string) error {
	if r.conex == nil {
		return fmt.Errorf("connection isn't initialized")
	}
	ch, err := r.conex.Channel()
	defer ch.Close()

	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "msgbroker.rabbitmq.PublishOnQueue")
	}

	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body, // Our JSON body as []byte
	})
	if err != nil {
		return errors.Wrap(err, "msgbroker.PublishOnQueue")
	}

	log.Printf("A message was sent to queue %v: %v", queueName, string(body))
	return nil
}

func (r *rabbitmqBroker) Subscribe(exchangeName string, handlerFunc func(data []byte)) error {
	amqpChannel, err := r.conex.Channel()
	if err != nil {
		return errors.Wrap(err, "msgbroker.rabbitmq.Subscribe")
	}
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare(exchangeName, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "msgbroker.rabbitmq.Subscribe")
	}

	err = amqpChannel.Qos(1, 0, false)
	if err != nil {
		return errors.Wrap(err, "msgbroker.rabbitmq.Subscribe")
	}
	messageChannel, err := amqpChannel.Consume(queue.Name, "", false, false, false, false, nil)

	if err != nil {
		return errors.Wrap(err, "msgbroker.rabbitmq.Subscribe")
	}
	stopChan := make(chan bool)

	go readMessages(messageChannel, handlerFunc)
	// Stop for program termination
	<-stopChan
	return nil
}

func readMessages(messageChannel <-chan amqp.Delivery, handler func(data []byte)) {
	log.Printf("Consumer ready, PID: %d", os.Getpid())
	for d := range messageChannel {
		handler(d.Body)
		if err := d.Ack(false); err != nil {
			log.Printf("Error acknowledging message : %s", err)
		} else {
			log.Printf("Acknowledged message")
		}

	}
}

func (r *rabbitmqBroker) Close() {
	r.conex.Close()
}
