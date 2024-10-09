package api

import (
	"encoding/json"
	"bytes"
	"fmt"
	"log"
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

	codeVerifier, err := verifier()
	if err != nil {
		log.Println("error: authorize: could not calculate code verifier: ", err)
	}

	codeChallange := codeVerifier.CodeChallengeS256()
	awxBenRes, err := authorize3(codeChallange)

	if err != nil {
		fmt.Printf("error: could not authorize: %s\n", err)
	}

	if awxBenRes.StatusCode != 200 {
		log.Println("error: authorize: ", string(awxBenRes.Body))
	}

	awxAuthRes := struct {
		AuthorizationCode string `json:"authorization_code"`
	}{}

	aerr := json.Unmarshal(awxBenRes.Body, &awxAuthRes)
	if aerr != nil {
		fmt.Printf("error: card details: %s\n", aerr)
	}

	c := Component {
		codeVerifier.Value,
		awxAuthRes.AuthorizationCode,
		os.Getenv("clientId"),
	}

	data := struct {
		Customer string
		Component Component
		Error    string
	}{
		Customer: os.Getenv("customer"),
		Component: c,
	}

	t.Execute(w, data)
}

func authorize3 (codeChallange string) (AwxResponse, error) {
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
	}{
		codeChallange,
		[]string{"w:awx_action:transfers_edit"},
	}

	body, _ := json.Marshal(data)

	url := "/api/v1/authentication/authorize"
	awxRes, err := sendRequest("POST", url, bytes.NewBuffer(body), header)
	if err != nil {
		fmt.Printf("connector: %v\n", err)
	}

	return awxRes, nil
}
