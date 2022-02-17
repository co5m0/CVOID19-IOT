package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
    "os"

	firebase "firebase.google.com/go/v4"
	"github.com/streadway/amqp"
	"google.golang.org/api/option"
)

// ENV 
// RABBITMQURL amqp://guest:guest@192.168.68.112:5672
// QUEUENAME iot/rmq/enter
// UID enter-service
// FIREBASEPATH cvoid-bar/door-enter
// 


func main() {

	fmt.Println("Starting Go application")
    fmt.Println("ENV: RabbitMQ", os.Getenv("RABBITMQURL"), "Queue Name: ", os.Getenv("QUEUENAME"), "UID: ", os.Getenv("UID"))
	conn, err := amqp.Dial(os.Getenv("RABBITMQURL"))

	if err != nil {
		fmt.Printf("Error connetting to AMQP %s\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Successfully Connected to RabbitMQ")

	ch, err := conn.Channel()

	if err != nil {
		// fmt.Errorf("Erron connetting to AMQP channel %s\n", err)
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
        os.Getenv("QUEUENAME "),
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		// fmt.Errorf("Erron consuming AMQP queue %s\n", err)
		panic(err)
	}

	fmt.Printf("Printing queue info %v\n", q)

	msgs, err := ch.Consume(
        os.Getenv("QUEUENAME"),
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		// fmt.Errorf("Erron consuming AMQP queue %s\n", err)
		panic(err)
	}

	db, err := NewDatabase("/home/sqlite.db")
	if err != nil {
		fmt.Printf("Erron connetting to Database %s\n", err)
		return
	}
	defer db.Close()

	msgToClient := make(chan DoorMQ, 10)

	// Init Firebase app
	fmt.Println("Init Firebase App ...")
	ctx := context.Background()
	// Initialize the app with a custom auth variable, limiting the server's access
	ao := map[string]interface{}{"uid": os.Getenv("UID")}
	conf := &firebase.Config{
		DatabaseURL:  os.Getenv("FIREBASEDATABASEURL"),
		AuthOverride: &ao,
	}
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile("/home/cvoid19.json")

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	// As an admin, the app has access to read and write all data, regradless of Security Rules
	ref := client.NewRef("cvoid-bar/door-exit")

	// go chan
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			body := &DoorMQ{}
			json.Unmarshal([]byte(msg.Body), body)
			fmt.Println("\t--- Message Recieved ---")
			fmt.Printf("door: %s\t status: %s\t timestamp: %s\n", body.Door, body.Status, body.Timestamp)

			if err = db.Insert(*body); err != nil {
				fmt.Println("Error inserting item")
			}

            if body.Status == "open"{
                msgToClient <- *body
            }
		}
	}()

	go func() {
        fmt.Println("Firebase goroutine launched")
		for msg := range msgToClient {
            msg.Status = "1"

            if _, err := ref.Push(ctx, &msg); err != nil {
                fmt.Println("Error pushing child node:", err)
            }
		}
        fmt.Println("Exiting from firebase goroutine")
	}()

	fmt.Println("Successfully connected to RabbitMQ")
	fmt.Println("Waiting for messages ...")


	<-forever
}
