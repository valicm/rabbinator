package cmd

import (
	"encoding/json"
	"github.com/bkway/gochimp/mandrill"
	"log"
)

// Definition for Mandrill queue item.
type queueItemMandrill struct {
	Data struct {
		// Specifics for Drupal module mandrill output
		// Otherwise we could directly map mandrill.Message struct.
		Id     string `json:"id"`
		Module string `json:"module,omitempty"`
		mandrill.Message
	} `json:"message"`
}

// processMandrillItem unmarshal data to Mandrill struct
// Preform API calls and return allowed string for status.
func processMandrillItem(QueueBody []byte, apiKey string, defaultTemplate string, moduleTemplates map[string]string) string {
	var data queueItemMandrill

	err := json.Unmarshal(QueueBody, &data)
	// If we have mapping issue, just print an error in the log and continue.
	// Probably could be minor / not blocking mapping, so we can let it hopefully.
	if err != nil {
		log.Print("There was an error in data mapping: ", err)
	}

	// We should not reach here, but if we are.
	// Exit from rabbinator. No point of constant requeue
	// item if no api key is provided.
	if apiKey == "" {
		log.Fatal("Missing Mandrill Api key. Exiting...")
	}

	client := mandrill.NewClient(apiKey)

	var templateContent []mandrill.Variable

	templateContent = append(templateContent, mandrill.Variable{
		Name:    "body",
		Content: data.Data.Message.Html,
	})

	// Specifics for usage with Drupal mandrill module,
	// but could be reused elsewhere if needed.
	var templateId = moduleTemplates[data.Data.Id]

	// We don't have specifics. Use default template.
	if templateId == "" {
		templateId = defaultTemplate
	}

	send, err := client.MessagesSendTemplate(templateId, templateContent, &data.Data.Message, true, map[string]string{})
	if err != nil {
		log.Print("mandrill is unable to send message due to error: ", err)
		return queueRetry
	}

	// Get received status from Mandrill
	var sentStatus = send[0].Status

	// Reject or requeue messages depending on status received from Mandrill.
	switch sentStatus {
	case "rejected":
		log.Print("mandrill rejected email with following details: ", send[0])
		return queueReject
	case "invalid":
		log.Print("mandrill returned invalid sent status: ", send[0])
		return queueUnknown
	case "error":
		log.Print("mandrill returned error: ", send[0])
		return queueRetry
	}

	// Mark message as delivered.
	return queueSuccess

}
