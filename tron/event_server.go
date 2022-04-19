package tron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type EventServer struct {
	url    string
	client *http.Client
}

func NewEventServer(url string) *EventServer {
	return &EventServer{
		url: url,
		client: &http.Client{
			Timeout: time.Millisecond * time.Duration(10*1000),
		},
	}
}

func (e *EventServer) GetContractsEventsParams(contractAddress string, params url.Values, res interface{}) error {

	path := fmt.Sprintf("/v1/contracts/%s/events", contractAddress)
	err := e.get(e.url+path, params, res)
	if err != nil {
		return err
	}
	return nil
}

// get 请求
func (e *EventServer) get(urlStr string, params url.Values, res interface{}) error {

	Url, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		// handle error
		return err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		// handle error
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}

	return nil
}

// post 请求
func (e *EventServer) post(urlStr string, data map[string]interface{}, res interface{}) error {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(bytesData))
	if err != nil {
		// handle error
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.client.Do(req)
	if err != nil {
		// handle error
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}

	return nil
}
