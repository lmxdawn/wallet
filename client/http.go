package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Response ...
type Response struct {
	Code    int    `json:"code"`    // 错误code码
	Message string `json:"message"` // 错误信息
}

func (br Response) Success() bool {
	return br.Code == 0
}

type HttpClient struct {
	protocol          string
	rechargeNotifyUrl string
	withdrawNotifyUrl string
}

// NewHttpClient 创建
func NewHttpClient(protocol, rechargeNotifyUrl, withdrawNotifyUrl string) (*HttpClient, error) {
	return &HttpClient{
		protocol,
		rechargeNotifyUrl,
		withdrawNotifyUrl,
	}, nil
}

// RechargeSuccess 充值成功通知
func (h *HttpClient) RechargeSuccess(hash string, address string, value int64) error {

	data := make(map[string]interface{})
	data["protocol"] = h.protocol
	data["hash"] = hash
	data["address"] = address
	data["value"] = value

	var res Response
	err := post(h.withdrawNotifyUrl, data, &res)
	if err != nil {
		return err
	}

	if !res.Success() {
		return errors.New(res.Message)
	}

	return nil
}

// WithdrawSuccess 提现成功通知
func (h *HttpClient) WithdrawSuccess(hash string, orderId string, address string, value int64) error {

	data := make(map[string]interface{})
	data["protocol"] = h.protocol
	data["hash"] = hash
	data["orderId"] = orderId
	data["address"] = address
	data["value"] = value

	var res Response
	err := post(h.withdrawNotifyUrl, data, &res)
	if err != nil {
		return err
	}

	if !res.Success() {
		return errors.New(res.Message)
	}

	return nil
}

// get 请求
func get(urlStr string, params url.Values, res interface{}) error {

	Url, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}

	return nil
}

// post 请求
func post(urlStr string, data map[string]interface{}, res interface{}) error {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(urlStr, "application/json", bytes.NewReader(bytesData))
	if err != nil {
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
