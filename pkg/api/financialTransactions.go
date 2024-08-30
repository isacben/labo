package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

func ViewFinancialTransactions(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"templates/financialTransactions.html",
		"templates/navbar.html",
		"templates/header.html",
	)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

func GetFinancialTransactions(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from_created_at")
	to := r.FormValue("to_created_at")

	var report Report

	financialTransactions, err := financialTransactions(from, to)
	if err != nil {
		fmt.Printf("error: could not get financial transactions: %s\n", err)
	}

	jerr := json.Unmarshal(financialTransactions, &report)
	if jerr != nil {
		fmt.Printf("error: financial transactions: %s\n", jerr)
	}

	t, _ := template.ParseFiles(
		"templates/header.html",
		"templates/financialTransactions.html",
	)
	t.ExecuteTemplate(w, "report-table", report)

	fmt.Printf("financialTransactions: %s", string(financialTransactions))
}

// request to Airwallex

func financialTransactions(from string, to string) ([]byte, error) {
	token, err := getToken()
	if err != nil {
		return nil, fmt.Errorf("authentication error: %s", err)
	}

	openId := os.Getenv("openId")
	authHeader := fmt.Sprintf("Bearer %s", token)
	header := http.Header{
		"Authorization":  {authHeader},
		"x-on-behalf-of": {openId},
	}

	url := fmt.Sprintf("/api/v1/financial_transactions?from_created_at=%s&to_created_at=%s", from, to)
	resBody, err := sendRequest("GET", url, nil, header)
	if err != nil {
		fmt.Printf("connector: %v\n", err)
	}

	return resBody, nil
}
