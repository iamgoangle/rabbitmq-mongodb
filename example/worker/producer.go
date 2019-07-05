package main

import (
	"log"

	"github.com/iamgoangle/rabbitmq-mongodb/internal/rabbitmq"
)

func main() {
	rbMqConfig := rabbitmq.ConfigProducer{
		Conn: rabbitmq.ConfigConnection{
			Type: "standalone",
			Url:  "amqp://admin:1234@localhost:5672/",
		},
		// direct type
		// Exchange: rabbitmq.ConfigExchange{
		// 	Type: "",
		// 	Name: "Campaign-test",
		// },

		// topic type
		Exchange: rabbitmq.ConfigExchange{
			Type:    rabbitmq.ExchangeTopic,
			Name:    "Campaign-test-topic",
			Routing: "campaign.pool.channelID4",
		},
		Queue: rabbitmq.ConfigQueue{
			Name:            "hello-direct",
			DeleteWhenUnuse: false,
			Durable:         false,
			Exclusive:       false,
			NoWait:          false,
		},
	}

	producer, err := rabbitmq.NewProducer(rbMqConfig)
	if err != nil {
		log.Println("[main]: unable to connect RabbitMQ %+v", err)
	}

	body := []byte(`{"eventType":"newCampaign","data":"Hello GGEZ GGEEEE","Time":1294706395881547000}`)
	err = producer.Publish(body)
	if err != nil {
		log.Println("unable to send my hello world %+v", err)
	}

	// body := []byte(`{"eventType":"newCampaign","data":"Hello Mook","Time":1294706395881547000}`)
	// var mid = rand.Intn(100)
	// bMessage, err := json.Marshal(map[string]interface{}{
	// 	"eventType": "newCampaign",
	// 	"data": map[string]interface{}{
	// 		"Mid":    mid,
	// 		"Thread": "thanakorn",
	// 		"Time":   time.Now().Unix(),
	// 	},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = producer.PublishDirectExchange(bMessage)
	// if err != nil {
	// 	log.Printf("unable to send my hello world %+v", err)
	// }
}
