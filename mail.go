package gomail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

var client http.Client

type Email struct {
	Token   string
	Address string
	Id      string
}

type ResponseStruct struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}

func AddRequestHeaders(headers map[string]string, req *http.Request) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func CreateAccount() (*Email, error) {
	account := Email{}
	account.Address = fmt.Sprintf("hon.%d@midiharmonica.com", rand.Intn(9999999))
	payload := "{\"address\":\"" + account.Address + "\",\"password\":\"JinYang1@\"}"
	req, err := http.NewRequest("POST", "https://api.mail.tm/accounts", strings.NewReader(payload))
	if err != nil {
		return &account, fmt.Errorf("Could not connect.")
	}
	headers := map[string]string{
		"content-type": "application/json",
		"user-agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
	}
	AddRequestHeaders(headers, req)
	resp, err := client.Do(req)
	if err != nil {
		return &account, fmt.Errorf("Could not request.")
	}
	defer resp.Body.Close()
	req, err = http.NewRequest("POST", "https://api.mail.tm/token", strings.NewReader(payload))
	if err != nil {
		return &account, fmt.Errorf("Could not get account token.")
	}
	AddRequestHeaders(headers, req)
	resp, err = client.Do(req)
	if err != nil {
		return &account, fmt.Errorf("Could not request.")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), "token") {
		var responseStruct ResponseStruct
		err = json.Unmarshal(body, &responseStruct)
		if err != nil {
			return &account, fmt.Errorf("Could not parse token")
		}
		account.Token = responseStruct.Token
		account.Id = responseStruct.Id
	}
	return &account, nil
}

func (em *Email) DeleteAccount() error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.mail.tm/%s", em.Id), nil)
	if err != nil {
		return fmt.Errorf("Could not create request.")
	}
	headers := map[string]string{
		"content-type": "application/json",
		"user-agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
	}
	AddRequestHeaders(headers, req)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Could not do the request.")
	}
	defer resp.Body.Close()
	return nil
}
