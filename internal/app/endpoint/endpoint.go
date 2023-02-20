package endpoint

import (
	"encoding/json"
	"log"
	"net/http"
)

// ------- Received Message structures -----------

type HandledClientMessage struct {
	UpdateID    int           `json:"update_id"`
	Message     ClientMessage `json:"message"`
	ChannelPost ChannelPost   `json:"channel_post"`
}

type ClientMessage struct {
	MessageID int        `json:"message_id"`
	From      From       `json:"from"`
	Chat      Chat       `json:"chat"`
	Date      int        `json:"date"`
	Text      string     `json:"text"`
	Entities  []Entities `json:"entities"`
}

type ChannelPost struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type From struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	UserName     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	ID                          int    `json:"id"`
	FirstName                   string `json:"first_name"`
	UserName                    string `json:"username"`
	Type                        string `json:"type"`
	Title                       string `json:"title"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

type Entities struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type Service interface {
	NewBot(m string, Cid int)
	NewGPT()
	SendMessage(m string)
	FixGrammar()
}

type Endpoint struct {
	s Service
}

func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
	}
}

func (e *Endpoint) HandleClientMessage(w http.ResponseWriter, r *http.Request) {
	message := &HandledClientMessage{}

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		log.Fatal(err)
	}

	e.s.NewBot(message.Message.Text, message.Message.Chat.ID)
	e.s.NewGPT()

	e.s.FixGrammar()
}
