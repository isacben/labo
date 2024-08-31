package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AuthJson struct {
	Expires_at string `json:"expires_at"`
	Token      string `json:"token"`
}

func getToken() (string, error) {
	clientId := os.Getenv("clientId")
	apiKey := os.Getenv("apiKey")

	header := http.Header{"x-client-id": {clientId}, "x-api-key": {apiKey}}
	awxRes, err := sendRequest("POST", "/api/v1/authentication/login", nil, header)
	if err != nil {
		fmt.Println(err)
	}

	var authJson AuthJson
	jerr := json.Unmarshal(awxRes.Body, &authJson)
	if jerr != nil {
		fmt.Printf("authorization: %v\n", jerr)
	}
	return authJson.Token, nil
}

func sendRequest(method string, url string, body io.Reader, header http.Header) (AwxResponse, error) {
	reqUrl := fmt.Sprintf("https://api-demo.airwallex.com%s", url)
	req, err := http.NewRequest(method, reqUrl, body)

	req.Header = header
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return AwxResponse{}, fmt.Errorf("could not create http request: %s", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return AwxResponse{}, fmt.Errorf("error making http request: %s", err)
	}

	defer res.Body.Close()

	var awxRes AwxResponse

	awxRes.Body, err = io.ReadAll(res.Body)
	if err != nil {
		return AwxResponse{}, fmt.Errorf("io error: could not read response body: %s", err)
	}

	awxRes.Header = res.Header
	awxRes.StatusCode = res.StatusCode

	return awxRes, nil
}
