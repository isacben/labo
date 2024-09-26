package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

func ViewTransfers(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"templates/transfers.html",
		"templates/navbar.html",
		"templates/header.html",
	)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

func GetTransfers(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from_created_at")
	to := r.FormValue("to_created_at")

	var report Transfers

	awxRes, err := transfers(from, to)
	if err != nil {
		fmt.Printf("error: could not get transfers: %s\n", err)
	}

	t, _ := template.ParseFiles(
		"templates/header.html",
		"templates/transfers.html",
	)

	if awxRes.StatusCode != 200 {
		msg := fmt.Sprintf("Error: HTTP %v: %s", awxRes.StatusCode, string(awxRes.Body))
		log.Println(msg)
		log.Println(awxRes.Header)
		t.ExecuteTemplate(w, "error-msg", msg)
	}

	if awxRes.Header.Get("X-Sca-Session-Code") != "" {

		codeVerifier, err := verifier()
		if err != nil {
			log.Println("error: authorize: could not calculate code verifier: ", err)
		}

		codeChallange := codeVerifier.CodeChallengeS256()
		awxScaRes, err := authorize(codeChallange)

		if err != nil {
			fmt.Printf("error: could not authorize: %s\n", err)
		}

		if awxScaRes.StatusCode != 200 {
			log.Println("error: authorize: ", string(awxScaRes.Body))
		}

		awxAuthRes := struct {
			AuthorizationCode string `json:"authorization_code"`
		}{}

		aerr := json.Unmarshal(awxScaRes.Body, &awxAuthRes)
		if aerr != nil {
			fmt.Printf("error: transfers: %s\n", aerr)
		}

		scaInfo := struct {
			CodeVerifier      string
			AuthorizationCode string
			ClientId          string
			Email             string
			SessionCode       string
		}{
			codeVerifier.Value,
			awxAuthRes.AuthorizationCode,
			os.Getenv("clientId"),
			"test1231@airwallex.com",
			awxRes.Header.Get("X-Sca-Session-Code"),
		}
		t.ExecuteTemplate(w, "sca-component", scaInfo)
	}

	jerr := json.Unmarshal(awxRes.Body, &report)
	if jerr != nil {
		fmt.Printf("error: transfers: %s\n", jerr)
	}

	t.ExecuteTemplate(w, "report-table", report)
}

// request to Airwallex

func transfers(from string, to string) (AwxResponse, error) {
	token, err := getToken()
	if err != nil {
		return AwxResponse{}, fmt.Errorf("authentication error: %s", err)
	}

	openId := os.Getenv("openId")
	authHeader := fmt.Sprintf("Bearer %s", token)
	header := http.Header{
		"Authorization":  {authHeader},
		"x-on-behalf-of": {openId},
	}

	url := fmt.Sprintf("/api/v1/transfers?from_created_at=%s&to_created_at=%s", from, to)
	awxRes, err := sendRequest("GET", url, nil, header)
	if err != nil {
		fmt.Printf("connector: %v\n", err)
	}

	return awxRes, nil
}

func authorize(codeChallange string) (AwxResponse, error) {
	token, err := getToken()
	if err != nil {
		return AwxResponse{}, fmt.Errorf("authentication error: %s", err)
	}

	openId := os.Getenv("openId")
	authHeader := fmt.Sprintf("Bearer %s", token)
	header := http.Header{
		"Authorization":  {authHeader},
		"x-on-behalf-of": {openId},
	}

	data := struct {
		CodeChallange string   `json:"code_challenge"`
		Scope         []string `json:"scope"`
		Identity      string   `json:"identity"`
	}{
		codeChallange,
		[]string{"w:awx_action:sca_edit", "r:awx_action:sca_view"},
		"user_1234",
	}

	body, _ := json.Marshal(data)

	url := "/api/v1/authentication/authorize"
	awxRes, err := sendRequest("POST", url, bytes.NewBuffer(body), header)
	if err != nil {
		fmt.Printf("connector: %v\n", err)
	}

	return awxRes, nil
}
