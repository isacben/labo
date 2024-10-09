package api

import (
//	"encoding/json"
	"fmt"
//	"log"
	"net/http"
	"os"
	"text/template"
)

func BeneficiaryComponent(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"templates/beneficiaryComponent.html",
		"templates/navbar.html",
		"templates/header.html",
	)
	if err != nil {
		fmt.Println(err)
	}
	data := struct {
		Customer string
		Report   Report
		Error    string
	}{
		Customer: os.Getenv("customer"),
	}

	t.Execute(w, data)
}

//func GetBeneficia(w http.ResponseWriter, r *http.Request) {
//	from := r.FormValue("from_created_at")
//	to := r.FormValue("to_created_at")
//
//	var report Report
//
//	awxRes, err := financialTransactions(from, to)
//	if err != nil {
//		fmt.Printf("error: could not get financial transactions: %s\n", err)
//	}
//
//	t, _ := template.ParseFiles(
//		"templates/header.html",
//		"templates/financialTransactions.html",
//	)
//
//	if awxRes.StatusCode != 200 {
//		msg := fmt.Sprintf("Error: HTTP %v: %s", awxRes.StatusCode, string(awxRes.Body))
//		log.Println(msg)
//		t.ExecuteTemplate(w, "error-msg", msg)
//	}
//
//	jerr := json.Unmarshal(awxRes.Body, &report)
//	if jerr != nil {
//		fmt.Printf("error: financial transactions: %s\n", jerr)
//	}
//
//	t.ExecuteTemplate(w, "report-table", report)
//}
