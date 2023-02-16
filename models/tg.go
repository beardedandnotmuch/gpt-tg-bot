package models

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// ------- Received Message structures -----------

type TGClientMessage struct {
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

// ------- TGBot -----------

type TGBot interface {
	SendRequest(m string)
	GetClientMessage() string
	GetChaId() int
}

type TGStorage struct {
	ClientMessage string
	ChatId        int
}

func NewTGStorage(m string, cId int) *TGStorage {
	return &TGStorage{
		ClientMessage: m,
		ChatId:        cId,
	}
}

func (s *TGStorage) GetClientMessage() string {
	return s.ClientMessage
}

func (s *TGStorage) GetChaId() int {
	return s.ChatId
}

func (s *TGStorage) SendRequest(m string) {
	_, err := http.PostForm(fmt.Sprintf("%s%s/sendMessage", os.Getenv("TG_API_URL"), os.Getenv("TG_APITOKEN")), url.Values{
		"chat_id": {fmt.Sprint(s.GetChaId())},
		"text":    {m},
	})

	if err != nil {
		log.Fatal(err)
	}
}
