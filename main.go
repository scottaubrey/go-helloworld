package main

import (
	"encoding/json"
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
	args := os.Args

	if len(args) < 2 {
		log.Fatalf("Usage: %v <url>\n", args[0])
	}

	response, err := doRequest(args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: \n%v\n", response.GetResponse())
}

func doRequest(requestUrl string) (Response, error) {
	if _, err := url.ParseRequestURI(requestUrl); err != nil {
		return nil, fmt.Errorf("URL Parse Error: %s", err)
	}

	response, err := http.Get(requestUrl)
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

	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshall Error: %s", err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, fmt.Errorf("JSON unmarshall Error: %s", err)
		}
		return words, nil
	case "occurrence":
		var occurrences Occurrences
		err = json.Unmarshal(body, &occurrences)
		if err != nil {
			return nil, fmt.Errorf("JSON unmarshall Error: %s", err)
		}

		if _, ok := occurrences.Words["Scott"]; ok {
			fmt.Println("\n> Hey! I found a Scott! 👋\n")
		}

		return occurrences, nil
	}

	return nil, fmt.Errorf("unknown Page: %s", page.Name)
}
