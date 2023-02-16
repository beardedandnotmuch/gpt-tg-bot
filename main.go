package main

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

	m "banm-gpt-tg-bot/models"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/api/v1/update", handleClientMessage)

	fmt.Println("Listenning on port", os.Getenv("PORT"), ".")
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func handleClientMessage(w http.ResponseWriter, r *http.Request) {

	message := &m.ReceiveMessage{}

	chatID := 0
	msgText := ""

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(message.Message.Chat.ID, message.Message.Text)
	chatID = message.Message.Chat.ID
	msgText = message.Message.Text

	tgResponse("Processing...", chatID)

	gptResponseText, err := gptRequest(msgText)

	if err != nil {
		fmt.Println(err)
	}

	tgResponse(gptResponseText, chatID)
}

func tgResponse(message string, chatID int) {
	response, err := http.PostForm(fmt.Sprintf("%s%s/sendMessage", os.Getenv("TG_API_URL"), os.Getenv("TG_APITOKEN")), url.Values{
		"chat_id": {fmt.Sprint(chatID)},
		"text":    {message},
	})

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(body))
}

func gptRequest(text string) (string, error) {
	var gptResp m.GPTResponse
	var err error

	data := []byte(`{
		"model":             "text-davinci-003",
		"prompt":            "Correct this to standard English:\n\n` + text + `",
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

	// TODO: create error handler for statuses
	if response.StatusCode != 200 {
		err = errors.New(response.Status)
	}

	resbody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(resbody, &gptResp)

	return gptResp.GPTChoices[0].Text, nil
}
