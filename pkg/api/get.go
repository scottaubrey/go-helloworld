package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type page struct {
	Name string `json:"page"`
}

type words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (words words) GetResponse() string {
	formattedWords := "Words\n"
	formattedWords += "-----\n\n"
	formattedWords += strings.Join(words.Words, "\n")
	return formattedWords
}

type occurrences struct {
	Words map[string]int `json:"words"`
}

func (occurrences occurrences) GetResponse() string {
	formattedOccurrences := "Word\tCount\n"
	formattedOccurrences += "----\t-----\n\n"

	for word, count := range occurrences.Words {
		formattedOccurrences += fmt.Sprintf("%s\t%d\n", word, count)
	}
	return formattedOccurrences
}

type Response interface {
	GetResponse() string
}

func doRequest(client http.Client, baseUrl *url.URL, path string) (Response, error) {

	requestUrl := baseUrl.Scheme + "://" + baseUrl.Host
	if baseUrl.Path == "" {
		requestUrl += "/"
	} else {
		requestUrl += baseUrl.Path
	}
	requestUrl += path

	response, err := client.Get(requestUrl)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET Error: %s", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP GET Error: Status code is %d not 200. Body: %s", response.StatusCode, body)
	}

	if !json.Valid(body) {
		return nil, RequestError{
			Url:      requestUrl,
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      "Invalid Json",
		}
	}

	var page page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			Url:      requestUrl,
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("JSON unmarshall Error: %s", err),
		}
	}

	switch page.Name {
	case "words":
		var words words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				Url:      requestUrl,
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("JSON unmarshall Error: %s", err),
			}
		}

		return words, nil
	case "occurrence":
		var occurrences occurrences
		err = json.Unmarshal(body, &occurrences)
		if err != nil {
			return nil, RequestError{
				Url:      requestUrl,
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("JSON unmarshall Error: %s", err),
			}
		}

		if _, ok := occurrences.Words["Scott"]; ok {
			fmt.Println("\n> Hey! I found a Scott! 👋\n")
		}

		return occurrences, nil
	}

	return nil, RequestError{
		Url:      requestUrl,
		HTTPCode: response.StatusCode,
		Body:     string(body),
		Err:      fmt.Sprintf("Unknown Error: %s", err),
	}
}
