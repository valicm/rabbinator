package cmd

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"rabbinator/cmd/providers/mailchimp"
	"rabbinator/cmd/providers/mandrill"
	"rabbinator/cmd/utility"
)

// Initialize all task necessarily for establishing connection.
func Initialize(channel string, configFile string)  {

	// Initialize and set config.
	utility.ConfigSetup(channel, configFile)

	// Make connection to rabbitmq.
	connectRabbitMQ()
}

// Rabbit connection handler and processing items.
func connectRabbitMQ() {

	var clientChannel = viper.GetString("channel")
	var queueType = viper.GetString("type")

	// Start connection.
	conn, err := amqp.Dial(viper.GetString("client.uri"))
	utility.ErrorLog("Failed to connect to RabbitMQ", err)
	defer conn.Close()

	ch, err := conn.Channel()
	utility.ErrorLog("Failed to open a channel", err)
	defer ch.Close()

	// Declare queue.
	q, err := ch.QueueDeclare(
		clientChannel,
		viper.GetBool("client.queue.durable"),
		viper.GetBool("client.queue.autodelete"),
		viper.GetBool("client.queue.exclusive"),
		viper.GetBool("client.queue.nowait"),
		nil,
	)
	utility.ErrorLog("Failed to declare a queue", err)

	err = ch.Qos(
		viper.GetInt("client.prefetch.count"),
		viper.GetInt("client.prefetch.size"),
		viper.GetBool("client.prefetch.global"),
	)

	utility.ErrorLog("Failed to set QoS", err)

	msgs, err := ch.Consume(
		clientChannel,
		q.Name,
		viper.GetBool("client.consume.autoack"),
		viper.GetBool("client.consume.exclusive"),
		viper.GetBool("client.consume.nolocal"),
		viper.GetBool("client.consume.nowait"),
		nil,
	)
	utility.ErrorLog("Failed to register a consumer", err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Process queue items.
			processQueueItem(d, queueType)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// Process queue item.
// TODO: make it dynamic?
func processQueueItem(Delivery amqp.Delivery, queueType string) {

	switch queueType {
	case "mandrill":
		mandrill.ProcessItem(Delivery)
	case "mailchimp":
		mailchimp.ProcessItem(Delivery)

	default:
		// TODO: Reject item and write syslog?
		//Delivery.Acknowledger.Reject()
	}

}
