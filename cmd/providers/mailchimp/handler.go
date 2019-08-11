package mailchimp

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/bkway/gochimp"
	"os"
	"rabbinator/cmd/utility"
	"strings"
)

const MemberStatusSubscribed gochimp.SubscriptionStatus = "subscribed"
const MemberStatusPending gochimp.SubscriptionStatus = "pending"

var queueStatus utility.QueueStatus

// Definition for mailchimp queue item.
type QueueItem struct {
	Args struct {
		Email       string                 `json:"email"`
		EmailType   string                 `json:"email_type, omitempty"`
		ListId      string                 `json:"list_id"`
		DoubleOptin bool                   `json:"double_optin"`
		Format      string                 `json:"format, omitempty"`
		MergeVars   map[string]interface{} `json:"merge_vars, omitempty"`
		Interests   map[string]bool        `json:"interests,omitempty"`
	} `json:"args"`
}

// Process queue item. Unmarshal data to Mailchimp struct
// Preform API calls and return string response.
func ProcessItem(QueueBody []byte, apiKey string) string {
	var data QueueItem

	err := json.Unmarshal(QueueBody, &data)
	if err != nil {
		fmt.Println("There was an error in data mapping:", err)
	}

	// We should not reach here, but if we are.
	// Exit from rabbinator. No point of constant requeue
	// item if no api key is provided.
	if apiKey == "" {
		fmt.Println("Missing Mailchimp Api key. Exiting...")
		os.Exit(1)
	}

	// Start Mailchimp client.
	client := gochimp.NewClient(apiKey)

	var memberStatus = MemberStatusSubscribed

	// If double opt-in is required, set member status to 'pending',
	// but only if the user isn't already subscribed.
	if data.Args.DoubleOptin {
		memberInfo, err :=client.Member(data.Args.ListId, data.Args.Email)
		if err == nil {
			if memberInfo.Status == MemberStatusSubscribed {
				// If member is already subscribed, we don't need to send
				// it again.
				return queueStatus.Success
			} else {
				memberStatus = MemberStatusPending
			}
		}
	}

	// Construct our local member variable.
	var memberData = gochimp.Member{
		Id : generateUserId(data.Args.Email),
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
		fmt.Println("There was an error:", err)
		return queueStatus.Reject
	}

	if subscribe != nil {
		return queueStatus.Success
	}

	// Retry item.
	return queueStatus.Retry

}

// Mailchimp requires user id - md5 hash of the email (lowercase).
func generateUserId(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(s))))
}