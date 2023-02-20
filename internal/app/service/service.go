package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/beardedandnotmuch/gpt-tg-bot/internal/app/models/gpt"
	"github.com/beardedandnotmuch/gpt-tg-bot/internal/app/models/tg"
)

type Service struct {
	tgBot tg.Bot
	gpt   interface{}
}

func New() *Service {
	return &Service{
		&tg.Storage{},
		&gpt.Storage{},
	}
}

func (s *Service) NewBot(m string, cId int) {
	s.tgBot = tg.NewStorage(m, cId)
	s.SendMessage("Processing...")
}

func (s *Service) NewGPT() {
	s.gpt = gpt.NewStorage()
}

func (s *Service) SendMessage(m string) {
	err := s.SendTGApiRequest(m)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) FixGrammar() {
	response, err := s.SendGPTRequest(s.tgBot.GetClientMessage())
	if err != nil {
		log.Fatal(err)
	}

	s.SendMessage(response)
}

func (s *Service) SendGPTRequest(m string) (string, error) {
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
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New(response.Status)
	}

	gptResp := &gpt.Response{}

	resbody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(resbody, &gptResp)

	return gptResp.Choices[0].Text, nil
}

func (s *Service) SendTGApiRequest(m string) error {
	_, err := http.PostForm(fmt.Sprintf("%s%s/sendMessage", os.Getenv("TG_API_URL"), os.Getenv("TG_APITOKEN")), url.Values{
		"chat_id": {fmt.Sprint(s.tgBot.GetChaId())},
		"text":    {m},
	})

	if err != nil {
		return err
	}

	return nil
}
