package api

import (
	"net/http"
)

func NewRequest(method string, apiURL string, r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, apiURL, r.Body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Access-Control-Allow-Origin", "*")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
