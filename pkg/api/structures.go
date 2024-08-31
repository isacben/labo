package api

// Financial transactions

type Items struct {
	Id                 string  `json:"id"`
	BatchId            string  `json:"batch_id"`
	SourceId           string  `json:"source_id"`
	FundingSourceId    string  `json:"funding_source_id"`
	SourceType         string  `json:"source_type"`
	TransactionType    string  `json:"transaction_type"`
	Currency           string  `json:"currency"`
	Amount             float64 `json:"amount"`
	ClientRate         float64 `json:"client_rate"`
	CurrencyPair       string  `json:"currency_pair"`
	Net                float64 `json:"net"`
	Fee                float64 `json:"fee"`
	EstimatedSettledAt string  `json:"estimated_settled_at"`
	SettledAt          string  `json:"settled_at"`
	Description        string  `json:"description"`
	Status             string  `json:"status"`
	CreatedAt          string  `json:"created_at"`
}

type Report struct {
	HasMore bool    `json:"has_more"`
	Items   []Items `json:"items"`
}

// Financial Reports

type FinancialReportItems struct {
	Id              string `json:"id"`
	FileName        string `json:"file_name"`
	FileFormat      string `json:"file_format"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	ReportExpiresAt string `json:"report_expires_at"`
}

type FinancialReport struct {
	HasMore bool                   `json:"has_more"`
	Items   []FinancialReportItems `json:"items"`
}

// General

type AwxResponse struct {
	Body       []byte
	Header     map[string][]string
	StatusCode int
}
