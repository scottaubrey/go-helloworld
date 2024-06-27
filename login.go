package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client http.Client, loginUrl, password string) (string, error) {
	loginRequest := LoginRequest{
		Password: password,
	}

	requestBody, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %s", err)
	}

	response, err := client.Post(loginUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(requestBody),
			Err:      fmt.Sprintf("request error: %s", err),
		}
	}

	if response.StatusCode != 200 {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(requestBody),
			Err:      fmt.Sprintf("request error: %s", err),
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(requestBody),
			Err:      fmt.Sprintf("error reading body: %s", err),
		}
	}

	if !json.Valid(body) {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("invalid JSON: %s", err),
		}
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("unmarshall error: %s", err),
		}
	}

	return loginResponse.Token, nil
}
