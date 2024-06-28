package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockClient struct {
	MockedResponse http.Response
}

func (client MockClient) Get(requestUrl string) (resp *http.Response, err error) {
	return &client.MockedResponse, nil
}

func TestDoGetWordsRequest(t *testing.T) {
	words := wordsPage{
		page:  page{"words"},
		words: words{Input: "abc", Words: []string{"abc", "def"}},
	}

	wordsJson, err := json.Marshal(words)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	client := MockClient{
		MockedResponse: http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(wordsJson)),
		},
	}

	baseUrl, err := url.Parse("http://localhost/")
	if err != nil {
		t.Fatalf("Url parse error: %v", err)
	}

	response, err := doRequest(client, baseUrl, "words")
	if err != nil {
		t.Fatalf("response error: %v", err)
	}

	responseString := response.GetResponse()
	if responseString != "Words\n-----\n\nabc\ndef" {
		t.Fatalf("Response did not match expected output: %s", responseString)
	}
}

func TestDoGetOccurencesRequest(t *testing.T) {
	occurrences := occurrencesPage{
		page:        page{"occurrence"},
		occurrences: occurrences{map[string]int{"abc": 2, "def": 1}},
	}

	occurrencesJson, err := json.Marshal(occurrences)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	client := MockClient{
		MockedResponse: http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(occurrencesJson)),
		},
	}

	baseUrl, err := url.Parse("http://localhost/")
	if err != nil {
		t.Fatalf("Url parse error: %v", err)
	}

	response, err := doRequest(client, baseUrl, "occurences")
	if err != nil {
		t.Fatalf("response error: %v", err)
	}

	responseString := response.GetResponse()
	if responseString != "Word\tCount\n----\t-----\n\nabc\t2\ndef\t1\n" {
		t.Fatalf("Response did not match expected output: %s", responseString)
	}
}
