package mandrill

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

// Definition for Mandrill queue item.
type Mandrill struct{
	Message struct {
		Id                      string `json:"id"`
		Module                  string `json:"module"`
		Html                    string `json:"html"`
		Subject                 string `json:"subject"`
		FromEmail               string `json:"from_email"`
		FromName                string `json:"from_name"`
		To                      []struct{
			Email string `json:"email"`
			Name string `json:"name, omitempty"`
			To string `json:"to, omitempty"`
		} `json:"to"`
		Headers                 map[string]string `json:"headers"`
		TrackOpens              bool `json:"track_opens"`
		TrackClicks             bool `json:"track_clicks"`
		AutoText                bool `json:"auto_text"`
		UrlStripQs              bool `json:"url_strip_qs"`
		BccAddress              string `json:"bcc_address"`
		Tags                    []string `json:"tags,omitempty"`
		GoogleAnalyticsDomains  []int `json:"google_analytics_domains,omitempty"`
		GoogleAnalyticsCampaign string `json:"google_analytics_campaign"`
		Attachments             []struct{
			Type string `json:"name"`
			Name string `json:"name"`
			Content string `json:"name"`
		} `json:"attachments, omitempty"`
		Images             []struct{
			Type string `json:"name"`
			Name string `json:"name"`
			Content string `json:"name"`
		} `json:"attachments, omitempty"`
		ViewContentLink         bool `json:"view_content_link"`
		Metadata                []int `json:"metadata, omitempty"`
	} `json:"message"`

}

// Process queue item.
func ProcessItem(Delivery amqp.Delivery) {
	var data Mandrill

	err := json.Unmarshal(Delivery.Body, &data)
	if err != nil {
		fmt.Println("There was an error:", err)
		//Delivery.Acknowledger.Reject()
	}

}