package msgbroker

// MessageBroker defines our interface for connecting, producing and consuming messages
type MessageBroker interface {
	PublishOnQueue(body []byte, queueName string) error
	Subscribe(exchangeName string, handlerFunc func(data []byte)) error
	Close()
}
