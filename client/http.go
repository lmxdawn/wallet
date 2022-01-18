package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
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
	coinName          string
	rechargeNotifyUrl string
	withdrawNotifyUrl string
	timeout           int // 超时时间，毫秒单位
}

// NewHttpClient 创建
func NewHttpClient(protocol, coinName, rechargeNotifyUrl, withdrawNotifyUrl string) *HttpClient {
	return &HttpClient{
		protocol,
		coinName,
		rechargeNotifyUrl,
		withdrawNotifyUrl,
		1000 * 10,
	}
}

// RechargeSuccess 充值成功通知
func (h *HttpClient) RechargeSuccess(hash string, status uint, address string, value int64) error {

	data := make(map[string]interface{})
	data["protocol"] = h.protocol
	data["coinName"] = h.coinName
	data["hash"] = hash
	data["status"] = status
	data["address"] = address
	data["value"] = value

	var res Response
	err := h.post(h.withdrawNotifyUrl, data, &res)
	if err != nil {
		return err
	}

	if !res.Success() {
		return errors.New(res.Message)
	}

	return nil
}

// WithdrawSuccess 提现成功通知
func (h *HttpClient) WithdrawSuccess(hash string, status uint, orderId string, address string, value int64) error {

	data := make(map[string]interface{})
	data["protocol"] = h.protocol
	data["coinName"] = h.coinName
	data["hash"] = hash
	data["status"] = status
	data["orderId"] = orderId
	data["address"] = address
	data["value"] = value

	var res Response
	err := h.post(h.withdrawNotifyUrl, data, &res)

	if err != nil {
		return err
	}

	if !res.Success() {
		return errors.New(res.Message)
	}

	return nil
}

// get 请求
func (h *HttpClient) get(urlStr string, params url.Values, res interface{}) error {

	Url, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()

	client := &http.Client{
		Timeout: time.Millisecond * time.Duration(h.timeout),
	}
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		// handle error
		return err
	}
	resp, err := client.Do(req)
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
func (h *HttpClient) post(urlStr string, data map[string]interface{}, res interface{}) error {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: time.Millisecond * time.Duration(h.timeout),
	}
	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(bytesData))
	if err != nil {
		// handle error
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
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
