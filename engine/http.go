package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/lmxdawn/wallet/types"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Http struct {
	url string // 请求地址
}

// GetAddressCount 获取需要生成地址的数量
func (h *Http) GetAddressCount(name string) (int, error) {

	params := url.Values{}
	params.Set("name", name)

	var res types.AddressCountRes
	err := h.get(params, &res)
	if err != nil {
		return 0, err
	}
	if !res.Success() {
		return 0, errors.New("获取生成地址数量错误")
	}

	return res.Data, nil
}

// PostAddress 发送已经生成的地址列表
func (h *Http) PostAddress(name string, addressArr []string) error {

	data := make(map[string]interface{})
	data["addressArr"] = addressArr

	var res types.Response
	err := h.post(data, &res)
	if err != nil {
		return err
	}

	if !res.Success() {
		return errors.New("发送地址列表失败")
	}

	return nil
}

// PostRechargeSuccess 发送充值成功的数据
func (h *Http) PostRechargeSuccess(name string, from string, to string, value uint64) error {

	data := make(map[string]interface{})
	data["from"] = from
	data["to"] = to
	data["value"] = value

	var res types.Response
	err := h.post(data, &res)
	if err != nil {
		return err
	}

	if !res.Success() {
		return errors.New("发送充值失败")
	}

	return nil
}

// GetWithdraw 获取提现列表
func (h *Http) GetWithdraw(name string) ([]types.WithdrawRes, error) {

	params := url.Values{}
	params.Set("name", name)

	var res types.WithdrawListRes
	err := h.get(params, &res)
	if err != nil {
		return nil, err
	}

	if !res.Success() {
		return nil, errors.New("获取提现列表失败")
	}

	return res.Data, nil
}

// PostWithdrawSuccess 发送提现成功数据
func (h *Http) PostWithdrawSuccess(name string, orderId string, address string) error {

	data := make(map[string]interface{})
	data["name"] = name
	data["orderId"] = orderId
	data["address"] = address

	var res types.Response
	err := h.post(data, &res)
	if err != nil {
		return err
	}

	if !res.Success() {
		return errors.New("发送提现成功失败")
	}

	return nil
}

// get 请求
func (h *Http) get(params url.Values, res interface{}) error {

	Url, err := url.Parse(h.url)
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
func (h *Http) post(data map[string]interface{}, res interface{}) error {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(h.url, "application/json", bytes.NewReader(bytesData))
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
