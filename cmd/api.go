package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Api struct {
	Host  string
	Token string
	Port  int
}

type NotificationSetting struct {
	Conditions       ConditionEntry `json:"Conditions"`
	SendNotification bool           `json:"Notification"`
}

type ConditionEntry struct {
	Prefix   string `json:"Prefix"`
	Protocol string `json:"Protocol"`
	DstIP    string `json:"DstIP"`
	DstPort  int    `json:"DstPort"`
	SrcIP    string `json:"SrcIP"`
	SrcPort  int    `json:"SrcPort"`
	InDev    string `json:"InDev"`
	OutDev   string `json:"OutDev"`
}

func NewApi(host string, token string) Api {
	port := 80
	return Api{host, token, port}
}

func (api Api) Get(url string) (string, error) {
	requestURL := fmt.Sprintf("http://%s:%d%s", api.Host, api.Port, url)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(body), nil
}

func (api Api) Put(url string, body string) (string, error) {
	data := strings.NewReader(body)
	requestURL := fmt.Sprintf("http://%s:%d%s", api.Host, api.Port, url)

	req, err := http.NewRequest("PUT", requestURL, data)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(bodyResp), nil
}

func (api Api) Delete(url string, body string) (string, error) {
	data := strings.NewReader(body)
	requestURL := fmt.Sprintf("http://%s:%d%s", api.Host, api.Port, url)

	req, err := http.NewRequest("DELETE", requestURL, data)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(bodyResp), nil
}
