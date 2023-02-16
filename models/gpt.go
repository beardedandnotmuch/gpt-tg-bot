package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type GPTResponse struct {
	Id         string       `json:"id"`
	Object     string       `json:"object"`
	Created    int          `json:"created"`
	Model      string       `json:"model"`
	GPTChoices []GPTChoices `json:"choices"`
	GPTUsage   GPTUsage     `json:"usage"`
}

type GPTChoices struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     string `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type GPTUsage struct {
	PromprTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GPT interface {
	FixGrammar(m string)
	SetResponse(t string)
	GetResponse() string
}

type GPTStorage struct {
	Response string
}

func NewGPTStorage() *GPTStorage {
	return &GPTStorage{
		Response: "",
	}
}

func (s *GPTStorage) SetResponse(t string) {
	s.Response = t
}
func (s *GPTStorage) GetResponse() string {
	return s.Response
}

func (s *GPTStorage) FixGrammar(m string) {
	data := []byte(`{
		"model":             "text-davinci-003",
		"prompt":            "Correct this to standard English:\n\n` + m + `",
		"temperature":       0,
		"max_tokens":        60,
		"top_p":             1.0,
		"frequency_penalty": 0.0,
		"presence_penalty":  0.0
	}`)

	req, err := http.NewRequest("POST", os.Getenv("GPT_API_URL"), bytes.NewBuffer(data))

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("GPT_API_TOKEN"))

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println(response.Status)
	}

	gptResp := &GPTResponse{}

	resbody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(resbody, &gptResp)

	s.SetResponse(gptResp.GPTChoices[0].Text)
}
