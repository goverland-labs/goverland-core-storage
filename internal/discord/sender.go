package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Message struct {
	Content string `json:"content"`
}

type Sender struct {
	WebhookURL string
}

func NewSender(webhookURL string) *Sender {
	return &Sender{WebhookURL: webhookURL}
}

func (s *Sender) SendMessage(content string) error {
	message := Message{
		Content: content,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
