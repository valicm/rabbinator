package mandrill

import (
	"encoding/json"
	"fmt"
	"github.com/bkway/gochimp/mandrill"
	"os"
	"rabbinator/cmd/utility"
)

var queueStatus utility.QueueStatus

// Definition for Mandrill queue item.
type QueueItem struct {
	Data struct {
		// Specifics for Drupal module mandrill output
		// Otherwise we could directly map mandrill.Message struct.
		Id     string `json:"id"`
		Module string `json:"module, omitempty"`
		mandrill.Message
	} `json:"message"`
}

// Process queue item. Unmarshal data to Mandrill struct
// Preform API calls and return allowed string for status.
func ProcessItem(QueueBody []byte, apiKey string, defaultTemplate string, moduleTemplates map[string]string) string{
	var data QueueItem

	err := json.Unmarshal(QueueBody, &data)
	// If we have mapping issue, just print an error and continue.
	if err != nil {
		fmt.Println("There was an error in data mapping:", err)
	}

	// We should not reach here, but if we are.
	// Exit from rabbinator. No point of constant requeue
	// item if no api key is provided.
	if apiKey == "" {
		fmt.Println("Missing Mandrill Api key. Exiting...")
		os.Exit(1)
	}

	client := mandrill.NewClient(apiKey)

	var templateContent [] mandrill.Variable

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
		fmt.Println("There was an error:", err)
		return queueStatus.Retry
	}

	// Get received status from Mandrill
	var sentStatus = send[0].Status

	// Reject or requeue messages depending on status received from Mandrill.
	switch sentStatus {
	case "rejected":
	case "invalid":
	case "error":
		return queueStatus.Reject
	}

	// Mark message as delivered.
	return queueStatus.Success

}
