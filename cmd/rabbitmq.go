package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"rabbinator/cmd/consumers/mailchimp"
	_ "rabbinator/cmd/consumers/mailchimp"
	"rabbinator/cmd/consumers/mandrill"
	_ "rabbinator/cmd/consumers/mandrill"
)

// Rabbit connection handler and processing items.
func Connect() {

	var clientChannel string = viper.GetString("channel")
	var queueType string = viper.GetString("type")
	var acknowledge bool
	//var message string

	// Start connection.
	conn, err := amqp.Dial(viper.GetString("client.uri"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	fmt.Println(clientChannel);

	// Declare queue.
	q, err := ch.QueueDeclare(
		clientChannel,
		viper.GetBool("client.queue.durable"),
		viper.GetBool("client.queue.autodelete"),
		viper.GetBool("client.queue.exclusive"),
		viper.GetBool("client.queue.nowait"),
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		viper.GetInt("client.prefetch.count"),
		viper.GetInt("client.prefetch.size"),
		viper.GetBool("client.prefetch.global"),
	)

	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		clientChannel,
		q.Name,
		viper.GetBool("client.consume.autoack"),
		viper.GetBool("client.consume.exclusive"),
		viper.GetBool("client.consume.nolocal"),
		viper.GetBool("client.consume.nowait"),
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			switch queueType {
			case "mandrill":
				acknowledge  = mandrill.Mandrill(d.Body)
			case "mailchimp":
				acknowledge = mailchimp.Mailchimp(d.Body)
			}
			log.Printf("Received a message: %s", d.Body)

			d.Ack(acknowledge)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
