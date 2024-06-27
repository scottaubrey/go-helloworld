package api

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Options struct {
	BaseUrl  string
	Password string
}

type API interface {
	GetOccurences() (*occurrences, error)
	GetWords() (*words, error)
	AddWord(word string) (*words, error)
}

type api struct {
	client  http.Client
	baseUrl *url.URL
}

func (api api) GetOccurences() (*occurrences, error) {
	response, err := doRequest(api.client, api.baseUrl, "/occurrence")
	if err != nil {
		return nil, err
	}

	occurrences, ok := response.(occurrences)
	if !ok {
		return nil, fmt.Errorf("expected Response type of Occurrences but got something else")
	}

	return &occurrences, nil
}

func (api api) GetWords() (*words, error) {
	response, err := doRequest(api.client, api.baseUrl, "/words")
	if err != nil {
		return nil, err
	}

	words, ok := response.(words)
	if !ok {
		return nil, fmt.Errorf("expected Response type of Words but got something else")
	}

	return &words, nil
}

func (api api) AddWord(word string) (*words, error) {
	response, err := doRequest(api.client, api.baseUrl, "/words?input="+word)
	if err != nil {
		return nil, err
	}

	words, ok := response.(words)
	if !ok {
		return nil, fmt.Errorf("expected Response type of Words but got something else")
	}

	return &words, nil
}

func New(options Options) (API, error) {

	parsedUrl, err := url.ParseRequestURI(options.BaseUrl)
	if err != nil {
		fmt.Printf("URL Parse Error: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	client := http.Client{}

	if options.Password != "" {
		token, err := doLoginRequest(client, parsedUrl, options.Password)
		if err != nil {
			return nil, err
		}

		client.Transport = myJwtTransport{
			token:     token,
			transport: http.DefaultTransport,
		}
	}
	return &api{
		client:  client,
		baseUrl: parsedUrl,
	}, nil
}
