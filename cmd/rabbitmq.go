package cmd

import (
	"github.com/streadway/amqp"
	"log"
	"rabbinator/cmd/providers/mailchimp"
	"rabbinator/cmd/providers/mandrill"
	"rabbinator/cmd/utility"
)

// Stored configuration for processing queue.
var config utility.Config

// Initialize all task necessarily for establishing connection.
func Initialize(consumer string, configDir string)  {

	// Initialize and set config.
	config = utility.ConfigSetup(consumer, configDir)

	// Make connection to RabbitMQ.
	connectRabbitMQ()
}

// Rabbit connection handler and processing items.
func connectRabbitMQ() {

	// Start connection.
	conn, err := amqp.Dial(config.Client.Uri)
	utility.ErrorLog("Failed to connect to RabbitMQ", err)
	defer conn.Close()

	ch, err := conn.Channel()
	utility.ErrorLog("Failed to open a channel", err)
	defer ch.Close()

	// Declare queue.
	q, err := ch.QueueDeclare(
		config.QueueName,
		config.Client.Queue.Durable,
		config.Client.Queue.AutoDelete,
		config.Client.Queue.Exclusive,
		config.Client.Queue.NoWait,
		nil,
	)
	utility.ErrorLog("Failed to declare a queue", err)

	err = ch.Qos(
		config.Client.Prefetch.Count,
		config.Client.Prefetch.Size,
		config.Client.Prefetch.Global,
	)

	utility.ErrorLog("Failed to set QoS", err)

	msgs, err := ch.Consume(
		config.QueueName,
		config.Consumer + q.Name,
		config.Client.Consume.AutoAck,
		config.Client.Consume.Exclusive,
		config.Client.Consume.NoLocal,
		config.Client.Consume.NoWait,
		nil,
	)
	utility.ErrorLog("Failed to register a consumer", err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Process queue items.
			processQueueItem(d)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// Process queue item.
// TODO: make it dynamic?
func processQueueItem(Delivery amqp.Delivery) {

	switch config.Type {
	case "mandrill":
		mandrill.ProcessItem(Delivery, config.ApiKey)
	case "mailchimp":
		mailchimp.ProcessItem(Delivery, config.ApiKey)

	default:
		// TODO: Reject item and write syslog?
		//Delivery.Acknowledger.Reject()
	}

}
