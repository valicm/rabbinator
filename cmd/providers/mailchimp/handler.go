package mailchimp

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

// Definition for mailchimp queue item.
type Mailchimp struct {
	Args struct {
		Email       string   `json:"email"`
		ListId      string   `json:"list_id"`
		DoubleOptin bool     `json:"double_optin"`
		Format      string   `json:"format, omitempty"`
		MergeVars   []string `json:"merge_vars, omitempty"`
		Interests   []string `json:"interests, omitempty"`
	} `json:"args"`
}

// Process queue item. Unmarshal data to Mailchimp struct
// Preform API calls and set Delivery.Acknowledger status.
func ProcessItem(Delivery amqp.Delivery, apiKey string) {
	var data Mailchimp

	err := json.Unmarshal(Delivery.Body, &data)
	if err != nil {
		fmt.Println("There was an error:", err)
		//Delivery.Acknowledger.Reject()
	}


}