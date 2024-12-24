package p2p

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Message struct {
	Id              string  `json:"id"`
	Message         string  `json:"message"`
	StatusCode      int     `json:"status_code"`
	Signature       string  `json:"signature"`
	TransactionData string  `json:"transaction_data"`
	Status          *string `json:"status"`
}

type Validator interface {
	ReceiveQuery(m Message) (Message, error)
	SendQuery(m Message, destinationId string) (Message, error)
	SelectRandomValidator() []string
	SelectRandomValidatorV2(excludeId string) []string
	VerifyMessage(m Message) (bool, error)
	CreateTransaction(tranData string) (string, error)

	GetID() string
}

func CommunicateHTTP[T any](url string, method string, headers map[string]string, body []byte, result *T, timeout time.Duration) (int, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	client := &http.Client{
		Timeout: time.Second * 100,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("failed to read response body: %v", err)
	}
	if err := json.Unmarshal(respBody, result); err != nil {
		return resp.StatusCode, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return resp.StatusCode, nil
}
