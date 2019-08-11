package mandrill

import (
	"encoding/json"
	"fmt"
	"github.com/bkway/gochimp/mandrill"
	"github.com/streadway/amqp"
)

// Definition for Mandrill queue item.
type QueueItem struct {
	Message struct {
		Id     string `json:"id"`
		Module string `json:"module, omitempty"`
		mandrill.Message
	} `json:"message"`
}

// Process queue item. Unmarshal data to Mandrill struct
// Preform API calls and set Delivery.Acknowledger status.
func ProcessItem(Delivery amqp.Delivery, apiKey string, defaultTemplate string, moduleTemplates map[string]string) bool{
	var data QueueItem

	//var queueTag = Delivery.DeliveryTag

	err := json.Unmarshal(Delivery.Body, &data)
	if err != nil {
		fmt.Println("There was an error:", err)
		//Delivery.Acknowledger.Reject(queueTag, true)
	}

	if apiKey == "" {
		fmt.Println("Missing api key")
	}

	client := mandrill.NewClient(apiKey)

	var templateContent [] mandrill.Variable

	templateContent = append(templateContent, mandrill.Variable{
		Name:    "body",
		Content: data.Message.Message.Html,
	})

	// Specifics for usage with Drupal mandrill module,
	// but could be reused elsewhere if needed.
	var templateId = moduleTemplates[data.Message.Id]

	// We don't have specifics. Use default template.
	if templateId == "" {
		templateId = defaultTemplate
	}

	send, err := client.MessagesSendTemplate(templateId, templateContent, &data.Message.Message, true, map[string]string{})
	if err != nil {
		fmt.Println("There was an error:", err)
	}

	fmt.Println(send)

	// Mark message as delivered.
	return true

}
