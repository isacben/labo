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

func ViewCard(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"templates/cardDetails.html",
		"templates/navbar.html",
		"templates/header.html",
	)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

func DisplayCard(w http.ResponseWriter, r *http.Request) {
	cardId := "6a1b3d9e-4fb0-464e-bef8-32abb9bfb4c6"

	var panToken PanToken

	awxRes, err := getPanToken(cardId)
	if err != nil {
		fmt.Printf("error: could not get pan token: %s\n", err)
	}

	fmt.Println(string(awxRes.Body))

	t, _ := template.ParseFiles(
		"templates/header.html",
		"templates/cardDetails.html",
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
		awxScaRes, err := authorize2(codeChallange)

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
			fmt.Printf("error: card details: %s\n", aerr)
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

	jerr := json.Unmarshal(awxRes.Body, &panToken)
	if jerr != nil {
		fmt.Printf("error: card details: %s\n", jerr)
	}

	t.ExecuteTemplate(w, "pan-delegation", panToken)
}

// request to Airwallex

func getPanToken(card_id string) (AwxResponse, error) {
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

	log.Println(card_id)
	data := struct {
		CardId string `json:"card_id"`
	}{
		card_id,
	}

	body, _ := json.Marshal(data)
	log.Println(string(body))
	url := "/api/v1/issuing/pantokens/create"
	awxRes, err := sendRequest("POST", url, bytes.NewBuffer(body), header)
	if err != nil {
		fmt.Printf("connector: %v\n", err)
	}

	return awxRes, nil
}

func authorize2(codeChallange string) (AwxResponse, error) {
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
