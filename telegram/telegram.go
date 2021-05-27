package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Producer struct {
	chat_id  string
	token    string
	endpoint string
}

type MessageResult struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

func NewProducer(chat_id, token string) *Producer {
	endpoint := "https://api.telegram.org/bot"
	return &Producer{chat_id: chat_id, token: token, endpoint: endpoint}
}

func (p Producer) SendMessage(message string) (MessageResult, error) {
	reqURL := fmt.Sprintf("%s%s/sendMessage?chat_id=%s&text=%s", p.endpoint, p.token, p.chat_id, message)
	log.Println("reqURL:", reqURL)
	client := &http.Client{}
	var result MessageResult
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return result, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("status not ok ", resp.StatusCode)
		return result, fmt.Errorf("status not ok")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}
	return result, nil
}
