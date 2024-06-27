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
	GetOccurences() (*Occurrences, error)
	GetWords() (*Words, error)
	AddWord(word string) (*Words, error)
}

type Api struct {
	client  http.Client
	baseUrl string
}

func (api Api) GetOccurences() (*Occurrences, error) {
	response, err := doRequest(api.client, api.baseUrl+"/occurrence")
	if err != nil {
		return nil, err
	}

	occurrences, ok := response.(Occurrences)
	if !ok {
		return nil, fmt.Errorf("expected Response type of Occurrences but got something else")
	}

	return &occurrences, nil
}

func (api Api) GetWords() (*Words, error) {
	response, err := doRequest(api.client, api.baseUrl+"/words")
	if err != nil {
		return nil, err
	}

	words, ok := response.(Words)
	if !ok {
		return nil, fmt.Errorf("expected Response type of Words but got something else")
	}

	return &words, nil
}

func (api Api) AddWord(word string) (*Words, error) {
	response, err := doRequest(api.client, api.baseUrl+"/words?input="+word)
	if err != nil {
		return nil, err
	}

	words, ok := response.(Words)
	if !ok {
		return nil, fmt.Errorf("expected Response type of Words but got something else")
	}

	return &words, nil
}

func New(options Options) (*Api, error) {

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

		client.Transport = MyJwtTransport{
			token:     token,
			transport: http.DefaultTransport,
		}
	}
	return &Api{
		client:  client,
		baseUrl: options.BaseUrl,
	}, nil
}