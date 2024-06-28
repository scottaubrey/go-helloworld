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

func TestDoGetRequest(t *testing.T) {
	words := wordsPage{
		page:  page{"words"},
		words: words{Input: "abc", Words: []string{"abc", "def"}},
	}

	wordsJson, err := json.Marshal(words)
	if err != nil {
		t.Errorf("marshal error: %v", err)
	}

	client := MockClient{
		MockedResponse: http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(wordsJson)),
		},
	}

	baseUrl, err := url.Parse("http://localhost/")
	if err != nil {
		t.Errorf("Url parse error: %v", err)
	}

	response, err := doRequest(client, baseUrl, "words")
	if err != nil {
		t.Errorf("response error: %v", err)
	}

	responseString := response.GetResponse()
	if responseString != "Words\n-----\n\nabc\ndef" {
		t.Errorf("Response did not match expected output: %s", responseString)
	}
}
