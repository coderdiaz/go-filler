package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fatih/color"
)

func makePetition(method, url string, body []byte, token *string) map[string]interface{} {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	if token != nil {
		req.Header.Add("Authorization", *token)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// We need a better handle of this kind of errors
	if res.StatusCode >= 500 {
		fmt.Println(res)
		log.Fatal("Something goes terribly wrong")
	}

	response := make(map[string]interface{})

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode >= 400 {
		data, _ := json.Marshal(response)

		red := color.New(color.FgRed).SprintFunc()
		log.Fatalf("The server has responded with: \"%s\" to the petition: %s on: %s", red(string(data[:])), red(req.Method), red(req.URL))
	}

	return response
}

// The unique difference between this function and the `makePetittion` is the response and what is downloaded
// These functionalities can be made by only one function
func makePetitionResponseArray(method, url string, body []byte, token *string) []map[string]interface{} {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	if token != nil {
		req.Header.Add("Authorization", *token)
	}

	if url == "https://api.culturacolectiva.com/rest-api/apicms/v1/getArticleJson/" {
		req.Header.Add("content-type", "application/x-www-form-urlencoded")
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// We need a better handle of this kind of errors
	if res.StatusCode >= 500 {
		log.Fatal("Something goes terribly wrong")
	}

	bodyResponse, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := make([]map[string]interface{}, 0)

	err = json.Unmarshal(bodyResponse, &response)
	if err != nil {
		log.Fatal(err)
	}

	return response
}
