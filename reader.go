package gomail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var client http.Client

type Email struct {
	Token   string
	Address string
}

type Inbox struct {
	MessageCount int
}

type Messages struct {
	AtId           string         `json:"@id"`
	Type           string         `json:"@type"`
	AccountId      string         `json:"accountId"`
	CreatedAt      string         `json:"createdAt"`
	DownloadUrl    string         `json:"downloadUrl"`
	From           string         `json:"from"`
	HasAttachments bool           `json:"hasAttachments"`
	Id             string         `json:"id"`
	Intro          string         `json:"intro"`
	IsDeleted      bool           `json:"isDeleted"`
	MessageId      string         `json:"msgid"`
	Seen           bool           `json:"seen"`
	Size           int            `json:"size"`
	Subject        string         `json:"subject"`
	To             map[int]string `json:"to"`
	TimeUpdated    string         `json:"updatedAt"`
}

type InboxResponse struct {
	Context       string         `json:"@context"`
	Id            string         `json:"@id"`
	Type          string         `json:"@type"`
	Messages      map[int]string `json:"hydra:member"`
	TotalMessages int            `json:"hydra:totalItems"`
}

func AddRequestHeaders(headers map[string]string, req *http.Request) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func (em *Email) GetInbox() (*InboxResponse, error) {
	inbox_resp := InboxResponse{}
	req, err := http.NewRequest("GET", "https://api.mail.tm/messages", nil)
	if err != nil {
		return &inbox_resp, fmt.Errorf("Request could not be made.")
	}
	headers := map[string]string{
		"content-type":  "application/json",
		"user-agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		"authorization": fmt.Sprintf("Bearer %s", em.Token),
	}
	AddRequestHeaders(headers, req)
	resp, err := client.Do(req)
	if err != nil {
		return &inbox_resp, fmt.Errorf("Client could not issue the request.")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &inbox_resp)
	if err != nil {
		return &inbox_resp, fmt.Errorf("Unmarshal was unsuccessful.")
	}
	return &inbox_resp, nil
}

func (ir *InboxResponse) GrepMessages() (*[]Messages, error) {
	messages := make([]Messages, ir.TotalMessages)
	for i := 0; i < ir.TotalMessages; i++ {
		message := &Messages{}
		err := json.Unmarshal([]byte(ir.Messages[i]), &message)
		if err != nil {
			return &messages, fmt.Errorf("Unmarshal could not parse messages.")
		}
		messages = append(messages, *message)
	}
	return &messages, nil
}

/*
func (ir *InboxResponse) GrepSubjects() ([]string, error) {
	subjects := make([]string, ir.TotalMessages)
	if ir.TotalMessages < 1 {
		return subjects, fmt.Errorf("No messages.")
	} else {
		for i := 0; i < len(subjects); i++ {
			subjects = append(subjects, ir.Messages.Subject)
		}
		return subjects, nil
	}
}
*/
