package types

// Response ...
type Response struct {
	Code    int    `json:"code"`    // 错误code码
	Message string `json:"message"` // 错误信息
}

func (br Response) Success() bool {
	return br.Code == 0
}

type AddressCountRes struct {
	Response `json:"response"`
	Data     int `json:"data,omitempty"` // 生成地址的数量
}

type WithdrawRes struct {
	OrderId string `json:"orderId,omitempty"` // 订单号
	To      string `json:"to,omitempty"`      // 提现地址
	Value   string `json:"value,omitempty"`   // 提现值
}

type WithdrawListRes struct {
	Response `json:"response"`
	Data     []WithdrawRes `json:"data,omitempty"` // 提现列表数据
}
