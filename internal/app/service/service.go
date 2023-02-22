package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/beardedandnotmuch/gpt-tg-bot/internal/app/cache"
	"github.com/beardedandnotmuch/gpt-tg-bot/internal/app/models/gpt"
	"github.com/beardedandnotmuch/gpt-tg-bot/internal/app/models/tg"
)

type Service struct {
	tgBot tg.Bot
	gpt   interface{}
	cache cache.GPTMessageCache
}

func New() *Service {
	return &Service{
		&tg.Storage{},
		&gpt.Storage{},
		&cache.RedisCache{},
	}
}

func (s *Service) NewBot(m string, cId int) {
	s.tgBot = tg.NewStorage(ProcessingTGMessage(m), cId)
}

func (s *Service) NewGPT() {
	s.gpt = gpt.NewStorage()
	s.cache = cache.NewRedisCache("redis-db:6379", 0, 10)
}

func (s *Service) SendMessage(m string) {
	err := s.SendTGApiRequest(m)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Service) FixGrammar() {
	cache := s.cache.Get(s.tgBot.GetClientMessage())

	if cache != "" {
		s.SendMessage(cache)
		return
	}

	s.SendMessage("Processing...")

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

	fmt.Println(bytes.NewBuffer(data))

	req, err := http.NewRequest("POST", os.Getenv("GPT_API_URL"), bytes.NewBuffer(data))

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("GPT_API_TOKEN"))

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New(response.Status)
	}

	gptResp := &gpt.Response{}

	resbody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(resbody, &gptResp)

	GPTMessage := gptResp.Choices[0].Text

	s.cache.Set(s.tgBot.GetClientMessage(), GPTMessage)

	return GPTMessage, nil
}

func (s *Service) SendTGApiRequest(m string) error {
	response, err := http.PostForm(fmt.Sprintf("%s%s/sendMessage", os.Getenv("TG_API_URL"), os.Getenv("TG_APITOKEN")), url.Values{
		"chat_id": {fmt.Sprint(s.tgBot.GetChaId())},
		"text":    {m},
	})

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}

func ProcessingTGMessage(s string) string {
	var lines string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			lines += "\\n"
		} else {
			lines += line + "\\n"
		}
	}
	return lines
}
