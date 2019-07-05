package main

import (
	"log"
	"math/rand"

	"github.com/streadway/amqp"

	"github.com/iamgoangle/rabbit-go/internal/rabbitmq"
)

func randBool() bool {
	return rand.Intn(2) == 0
}

func retryProducer(conn *amqp.Connection, c rabbitmq.ConfigProducer) (rabbitmq.Producer, error) {
	producer, err := rabbitmq.NewProducer(conn, c)
	if err != nil {
		log.Fatalf("[main]: unable to initial producer RabbitMQ %+v", err)
	}

	return producer, nil
}
func main() {
	conn, err := rabbitmq.NewConnection(rabbitmq.ConfigConnection{
		Type: "standalone",
		Url:  "amqp://admin:1234@localhost:5672/",
	})
	if err != nil {
		log.Fatalln("[main]: unable to connect RabbitMQ %+v", err)
	}

	rbMqConfig := rabbitmq.ConfigConsumer{
		Exchange: rabbitmq.ConfigExchange{
			Type:       rabbitmq.ExchangeDirect,
			Name:       "work.exchange",
			RoutingKey: "work.routing",
		},
		Queue: rabbitmq.ConfigQueue{
			Name: "work.queue",
			Bind: rabbitmq.ConfigQueueBind{
				RoutingKey:   "work.routing",
				ExchangeName: "work.exchange",
			},
		},
	}

	consumer, err := rabbitmq.NewConsumer(conn, rbMqConfig)
	if err != nil {
		log.Panic(err)
	}

	producerConfig := rabbitmq.ConfigProducer{
		Exchange: rabbitmq.ConfigExchange{
			Type:       rabbitmq.ExchangeDirect,
			Name:       "retry.exchange",
			RoutingKey: "retry.routing",
		},
		Queue: rabbitmq.ConfigQueue{
			Name: "retry.queue",
		},
	}
	retryProducer, err := retryProducer(conn, producerConfig)
	if err != nil {
		log.Panic(err)
	}

	msgs, err := consumer.WorkerProcessor()
	if err != nil {
		log.Println(err)
	}

	forever := make(chan bool)

	go func() {
		success := 0

		for m := range msgs {
			if randBool() {
				if err := m.Ack(false); err != nil {
					log.Printf("Error acknowledging message : %s", err)
				}
				success = success + 1
				log.Printf("success %v of / 10 \n", success)
				log.Println("Received a message: %s", string(m.Body))
			} else {
				err = retryProducer.Publish(m.Body)
				if err != nil {
					log.Printf("unable to send my event %+v", err)
				}

				m.Reject(false)
				log.Printf("Random Reject %+v", string(m.Body))
			}
		}
	}()

	<-forever

	consumer.Close()
}
