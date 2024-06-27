package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type loginRequest struct {
	Password string `json:"password"`
}
type loginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client http.Client, baseUrl *url.URL, password string) (string, error) {
	loginRequest := loginRequest{
		Password: password,
	}

	requestBody, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("unmarshal error: %s", err)
	}
	loginUrl := baseUrl.Scheme + "://" + baseUrl.Host
	loginUrl += baseUrl.Path
	loginUrl += "login"

	response, err := client.Post(loginUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", RequestError{
			Url:      loginUrl,
			HTTPCode: response.StatusCode,
			Body:     string(requestBody),
			Err:      fmt.Sprintf("request error: %s", err),
		}
	}

	if response.StatusCode != 200 {
		return "", RequestError{
			Url:      loginUrl,
			HTTPCode: response.StatusCode,
			Body:     string(requestBody),
			Err:      fmt.Sprintf("request error: %s", err),
		}
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", RequestError{
			Url:      loginUrl,
			HTTPCode: response.StatusCode,
			Body:     string(requestBody),
			Err:      fmt.Sprintf("error reading body: %s", err),
		}
	}

	if !json.Valid(body) {
		return "", RequestError{
			Url:      loginUrl,
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("invalid JSON: %s", err),
		}
	}

	var loginResponse loginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return "", RequestError{
			Url:      loginUrl,
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("unmarshall error: %s", err),
		}
	}

	return loginResponse.Token, nil
}
