package tron

// EventResult 事件
type EventResult struct {
	Success bool `json:"success"`
	Data    []struct {
		BlockNumber           int               `json:"block_number"`
		BlockTimestamp        int64             `json:"block_timestamp"`
		CallerContractAddress string            `json:"caller_contract_address"`
		ContractAddress       string            `json:"contract_address"`
		EventIndex            int               `json:"event_index"`
		EventName             string            `json:"event_name"`
		Result                map[string]string `json:"result"`
		ResultType            map[string]string `json:"result_type"`
		Event                 string            `json:"event"`
		TransactionID         string            `json:"transaction_id"`
	} `json:"data"`
	Meta struct {
		At          int64  `json:"at"`
		Fingerprint string `json:"fingerprint"`
		PageSize    int    `json:"page_size"`
	} `json:"meta"`
}
