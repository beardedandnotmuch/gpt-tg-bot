package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	var gpt m.GPT
	var tgBot m.TGBot

	message := &m.TGClientMessage{}

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Println(err)
	}

	tgBot = m.NewTGStorage(message.Message.Text, message.Message.Chat.ID)
	tgBot.SendRequest("Processing...")

	gpt = m.NewGPTStorage()
	gpt.FixGrammar(tgBot.GetClientMessage())

	tgBot.SendRequest(gpt.GetResponse())
}
