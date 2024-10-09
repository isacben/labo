package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
)

func ViewFinancialReports(w http.ResponseWriter, r *http.Request) {
	var report FinancialReport

	financialReport, err := getListOfFinancialReports(100)
	if err != nil {
		fmt.Printf("error: could not get financial reports: %s\n", err)
	}

	jerr := json.Unmarshal(financialReport, &report)
	if jerr != nil {
		fmt.Println("financial reports: ", jerr)
	}

	t, _ := template.ParseFiles(
		"templates/financialReports.html",
		"templates/navbar.html",
		"templates/header.html",
	)

	data := struct {
		Customer string
		Report   FinancialReport
		Error    string
	}{
		Customer: os.Getenv("customer"),
		Report:   report,
	}

	t.Execute(w, data)
}

func DownloadFinancialReport(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	format := r.URL.Query().Get("file_format")
	name := r.URL.Query().Get("file_name")

	fileContent, err := getFinancialReportContent(id)
	if err != nil {
		fmt.Printf("error: could not get financial report: %s\n", err)
	}

	switch format {
	case "csv":
		format = "text/csv"
	case "pdf":
		format = "application/pdf"
	case "xlsx":
		format = "application/vnd.ms-excel"
	}

	w.Header().Add("Content-Type", format)
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	w.Header().Add("Content-Length", fmt.Sprint(len(fileContent)))

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("Cannot use flusher")
	}

	w.Write(fileContent)
	flusher.Flush()
}

func getFinancialReportContent(id string) ([]byte, error) {
	token, err := getToken()
	if err != nil {
		return nil, fmt.Errorf("authentication error: %s", err)
	}

	url := fmt.Sprintf("https://api-demo.airwallex.com/api/v1/finance/financial_reports/%s/content", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http error: could not create request: %s", err)
	}

	authHeader := fmt.Sprintf("Bearer %s", token)
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {authHeader},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: could not send request: %s", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io error: could not read response body: %s", err)
	}

	return body, nil
}

func CreateFinancialReport(w http.ResponseWriter, r *http.Request) {
	token, err := getToken()
	if err != nil {
		//return nil, fmt.Errorf("authentication error: %s", err)
	}

	formData := struct {
		From        string `json:"from_created_at"`
		To          string `json:"to_created_at"`
		ReportType  string `json:"type"`
		File_format string `json:"file_format"`
	}{
		r.FormValue("from_created_at"),
		r.FormValue("to_created_at"),
		r.FormValue("type"),
		r.FormValue("file_format"),
	}

	data, _ := json.Marshal(formData)

	url := "https://api-demo.airwallex.com/api/v1/finance/financial_reports/create"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		//return nil, fmt.Errorf("http error: could not create request: %s", err)
	}

	authHeader := fmt.Sprintf("Bearer %s", token)
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {authHeader},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//return nil, fmt.Errorf("http error: could not send request: %s", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		//return nil, fmt.Errorf("io error: could not read response body: %s", err)
	}

	fmt.Print(string(body))
	//return body, nil

	var report FinancialReport

	financialReport, err := getListOfFinancialReports(1)
	if err != nil {
		fmt.Printf("error: could not get financial reports: %s\n", err)
	}

	jerr := json.Unmarshal(financialReport, &report)
	if jerr != nil {
		fmt.Println("financial reports: ", jerr)
	}

	t, _ := template.ParseFiles(
		"templates/financialReports.html",
		"templates/navbar.html",
		"templates/header.html",
	)
	t.ExecuteTemplate(w, "report-table", report)
}

func getListOfFinancialReports(pageSize int64) ([]byte, error) {
	token, err := getToken()
	if err != nil {
		return nil, fmt.Errorf("authentication error: %s", err)
	}

	url := fmt.Sprintf("https://api-demo.airwallex.com/api/v1/finance/financial_reports?page_size=%v", pageSize)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("http error: could not create request: %s", err)
	}

	openId := os.Getenv("openId")
	authHeader := fmt.Sprintf("Bearer %s", token)
	req.Header = http.Header{
		"Content-Type":   {"application/json"},
		"Authorization":  {authHeader},
		"x-on-behalf-of": {openId},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error: could not send request: %s", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io error: could not read response body: %s", err)
	}

	return body, nil
}
