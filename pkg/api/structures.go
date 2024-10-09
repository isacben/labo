package api

import "net/http"

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

type TransfersItems struct {
	Id        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Fee       float64 `json:"fee"`
	CreatedAt string  `json:"created_at"`
}

type Transfers struct {
	HasMore bool             `json:"has_more"`
	Items   []TransfersItems `json:"items"`
}

type PanToken struct {
	Expires_at string `json:"expires_at"`
	Token      string `json:"token"`
}

// General

type AwxResponse struct {
	Body       []byte
	Header     http.Header
	StatusCode int
}

type Sca struct {
	CodeVerifier      string
	AuthorizationCode string
	ClientId          string
	Email             string
	SessionCode       string
}

type Component struct {
	CodeVerifier      string
	AuthorizationCode string
	ClientId          string
}
