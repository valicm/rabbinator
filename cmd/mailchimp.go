package cmd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/bkway/gochimp"
	"log"
	"reflect"
	"strings"
)

const memberStatusSubscribed gochimp.SubscriptionStatus = "subscribed"
const memberStatusPending gochimp.SubscriptionStatus = "pending"

// Definition for Mailchimp queue item.
type queueItemMailchimp struct {
	Args struct {
		Email       string                 `json:"email"`
		EmailType   string                 `json:"email_type,omitempty"`
		ListId      string                 `json:"list_id"`
		DoubleOptin bool                   `json:"double_optin"`
		Format      string                 `json:"format,omitempty"`
		MergeVars   map[string]interface{} `json:"merge_vars,omitempty"`
		Interests   map[string]bool        `json:"interests,omitempty"`
	} `json:"args"`
}

// processMailchimpItem Unmarshal data to Mailchimp struct
// Preform API calls and return string response.
func processMailchimpItem(QueueBody []byte, apiKey string) string {
	var data queueItemMailchimp

	err := json.Unmarshal(QueueBody, &data)
	// If we have mapping issue, just print an error in the log and continue.
	// Probably could be minor / not blocking mapping, so we can let it hopefully.
	if err != nil {
		log.Println("There was an error in data mapping: ", err)
	}

	// We should not reach here, but if we are.
	// Exit from rabbinator. No point of constant requeue
	// item if no api key is provided.
	if apiKey == "" {
		log.Fatalf("%s: %s", "Missing Mailchimp Api key. Exiting...", err)
	}

	// Start Mailchimp client.
	client := gochimp.NewClient(apiKey)

	var memberStatus = memberStatusSubscribed

	// If double opt-in is required, set member status to 'pending',
	// but only if the user isn't already subscribed.
	if data.Args.DoubleOptin {
		memberInfo, err := client.Member(data.Args.ListId, data.Args.Email)
		if err == nil {
			if memberInfo.Status == memberStatusSubscribed {
				// If member is already subscribed, we don't need to send
				// it again.
				return queueSuccess
			}

			memberStatus = memberStatusPending

		}
	}

	// Construct our local member variable.
	var memberData = gochimp.Member{
		Id:           generateUserId(data.Args.Email),
		EmailAddress: data.Args.Email,
		EmailType:    gochimp.EmailType(data.Args.Format),
		Status:       memberStatus,
		MergeFields:  data.Args.MergeVars,
		Interests:    data.Args.Interests,
		ListId:       data.Args.ListId,
	}

	// Use method for adding/updating members.
	subscribe, err := client.UpsertMember(data.Args.ListId, &memberData)
	if err != nil {
		log.Print("mailchimp unable to make subscription api call due to error: ", err)
		return queueReject
	}

	// Why this: If we on response get proper mapping to SubscriptionStatus type
	// we are sure it is done. If you send malformed request, Mailchimp will return
	// integer as status, and gocimp library would map it to Member type, causing error.
	// TODO: explore gochimp logging/error improvements possibility.
	if reflect.TypeOf(subscribe.Status).Name() == "SubscriptionStatus" {
		return queueSuccess
	}

	// Retry item. If we reached here, some strange error occurred.
	log.Print("mailchimp returns faulty response:", subscribe)
	return queueRetry

}

// Mailchimp requires user id - md5 hash of the email (lowercase).
func generateUserId(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(s))))
}
