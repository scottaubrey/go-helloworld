package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (words Words) GetResponse() string {
	formattedWords := "Words\n"
	formattedWords += "-----\n\n"
	formattedWords += strings.Join(words.Words, "\n")
	return formattedWords
}

type Occurrences struct {
	Words map[string]int `json:"words"`
}

func (occurrences Occurrences) GetResponse() string {
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

func main() {
	var requestUrl string
	var password string
	var parsedUrl *url.URL

	flag.StringVar(&requestUrl, "url", "", "Url to make a request to")
	flag.StringVar(&password, "password", "", "password")

	flag.Parse()

	parsedUrl, err := url.ParseRequestURI(requestUrl)
	if err != nil {
		fmt.Printf("URL Parse Error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	client := http.Client{}

	if password != "" {
		token, err := doLoginRequest(client, parsedUrl.Scheme+"://"+parsedUrl.Host+"/login", password)
		if err != nil {
			if requestErr, ok := err.(RequestError); ok {
				log.Fatalf(
					"Error with login request with status code %d: %s\n\nBody:\n%s",
					requestErr.HTTPCode,
					requestErr.Err,
					requestErr.Body,
				)
			}
			log.Fatal(err)
		}

		client.Transport = MyJwtTransport{
			token:     token,
			transport: http.DefaultTransport,
		}
	}

	response, err := doRequest(client, parsedUrl.String())
	if err != nil {
		if requestErr, ok := err.(RequestError); ok {
			log.Fatalf(
				"Error with request to %s with status code %d: %s\n\nBody:\n%s",
				requestUrl,
				requestErr.HTTPCode,
				requestErr.Err,
				requestErr.Body,
			)
		}
		log.Fatal(err)
	}

	fmt.Printf("Response: \n%v\n", response.GetResponse())
}

func doRequest(client http.Client, requestUrl string) (Response, error) {
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
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      "Invalid Json",
		}
	}

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("JSON unmarshall Error: %s", err),
		}
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("JSON unmarshall Error: %s", err),
			}
		}

		return words, nil
	case "occurrence":
		var occurrences Occurrences
		err = json.Unmarshal(body, &occurrences)
		if err != nil {
			return nil, RequestError{
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
		HTTPCode: response.StatusCode,
		Body:     string(body),
		Err:      fmt.Sprintf("Unknown Error: %s", err),
	}
}
