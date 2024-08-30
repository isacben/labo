package main

import (
	"log"
	"net/http"

	"github.com/isacben/labo/pkg/api"
)

func main() {
	// Navigation
	http.HandleFunc("/financialTransactions", api.ViewFinancialTransactions)
	http.HandleFunc("/financialReports", api.ViewFinancialReports)
	http.HandleFunc("/downloadFinancialReport", api.DownloadFinancialReport)

	// internal API endpoints
	http.HandleFunc("/getFinancialTransactions", api.GetFinancialTransactions)
	http.HandleFunc("/createFinancialReport", api.CreateFinancialReport)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
