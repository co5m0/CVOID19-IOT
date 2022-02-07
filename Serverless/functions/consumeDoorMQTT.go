// This function read the values posted on a mqtt queue.
// Then define an AMQP queue and post there the message received
// with additional info.

// MQTT JSON message structure:
// {
//    "status":"open/close"
// }
// 
// AMQP JSON message structure:
// {
//      "door":"enter/exit",
//      "status":"open/close",
//      "timestamp":"xxxxxxxxxxxxx"
// }
//
// ENV vars:
//  RabbitMQURL: rabbitmq endpoint amqp://username:password@endpoint:port eg. amqp://guest:guest@192.168.68.112:5672
//  RabbitMQEnterQ: amqp queue name eg. iot/rmq/enter
//

package main

import (
	"encoding/json"
	"fmt"
	"os"
    "strconv"
	"time"

	"github.com/nuclio/nuclio-sdk-go"
	"github.com/streadway/amqp"
)

type mqttdoormsg struct {
	Status string `json:"status"`
}

type DoorMQ struct {
	Door      string `json:"door"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	fmt.Println("Starting Go application", os.Getenv("RabbitMQEnterQ"))

	msg := &mqttdoormsg{}
	json.Unmarshal([]byte(event.GetBody()), msg)
	conn, err := amqp.Dial(os.Getenv("RabbitMQURL"))

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Successfully Connected to RabbitMQ")

	ch, err := conn.Channel()

	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		os.Getenv("RabbitMQEnterQ"),
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Printing queue info %s\n", q)
	fmt.Printf("MESSAGE %s\n", msg.Status)

    msgjsonenc, err := json.Marshal(DoorMQ{
				Door:      os.Getenv("Door"),
				Status:    msg.Status,
				Timestamp: strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond))),
	})
    
    if err != nil {
		panic(err)
	}
    
	// publish message
	err = ch.Publish(
		"",
		os.Getenv("RabbitMQEnterQ"),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:msgjsonenc,
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully Published to RabbitMQ")
	return nil, nil
}


