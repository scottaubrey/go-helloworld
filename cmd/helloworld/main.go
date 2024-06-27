package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/scottaubrey/go-helloworld/pkg/api"
)

func main() {
	var baseUrl string
	var password string

	flag.StringVar(&baseUrl, "url", "", "base Url of API server to make a request to")
	flag.StringVar(&password, "password", "", "password")

	flag.Parse()

	apiClient, err := api.New(api.Options{
		BaseUrl:  baseUrl,
		Password: password,
	})

	if err != nil {
		if requestErr, ok := err.(api.RequestError); ok {
			log.Fatalf(
				"Error with request to %s with status code %d: %s\n\nBody:\n%s",
				requestErr.Url,
				requestErr.HTTPCode,
				requestErr.Err,
				requestErr.Body,
			)
		}
		log.Fatal(err)
	}

	command := flag.Arg(0)
	if command == "getWords" {
		words, err := apiClient.GetWords()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(words.GetResponse())
		os.Exit(0)
	}
	if command == "getOccurences" {
		occurences, err := apiClient.GetOccurences()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(occurences.GetResponse())
		os.Exit(0)
	}
	if command == "addWord" {
		word := flag.Arg(1)
		if word == "" {
			flag.Usage()
			os.Exit(1)
		}
		words, err := apiClient.AddWord(word)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(words.GetResponse())
		os.Exit(0)
	}
	fmt.Println("invalid command or command not given")
	flag.Usage()
	os.Exit(1)
}
