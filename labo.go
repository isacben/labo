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
	http.HandleFunc("/transfers", api.ViewTransfers)
	http.HandleFunc("/viewCard", api.ViewCard)
	http.HandleFunc("/beneficiaryComponent", api.BeneficiaryComponent)

	// internal API endpoints
	http.HandleFunc("/getFinancialTransactions", api.GetFinancialTransactions)
	http.HandleFunc("/createFinancialReport", api.CreateFinancialReport)
	http.HandleFunc("/getTransfers", api.GetTransfers)
	http.HandleFunc("/displayCard", api.DisplayCard)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
