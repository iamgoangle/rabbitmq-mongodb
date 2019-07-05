package rabbitmq

import (
	"github.com/streadway/amqp"
)

type ConfigConnection struct {
	Url  string
	Type string
}

// NewConnection factory connection to RabbitMQ
// TODO: refactor interface for testing
func NewConnection(c ConfigConnection) (conn *amqp.Connection, err error) {
	switch c.Type {
	case "cluster":
		return newClientCluster(c.Url)
	default:
		return newClient(c.Url)
	}
}

func newClientCluster(host string) (conn *amqp.Connection, err error) {
	return nil, nil
	// h := strings.Split(host, ",")
	// conn, err := amqp.Dial(h)
	// failOnError(err, "Failed to connect to RabbitMQ Cluster")
	// return conn, err
}

func newClient(host string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(host)
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn, err
}
