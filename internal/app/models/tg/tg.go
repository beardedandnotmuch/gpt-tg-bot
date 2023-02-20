package tg

type Bot interface {
	GetClientMessage() string
	GetChaId() int
}

type Storage struct {
	ClientMessage string
	ChatId        int
}

func NewStorage(m string, cId int) *Storage {
	return &Storage{
		ClientMessage: m,
		ChatId:        cId,
	}
}

func (s *Storage) GetClientMessage() string {
	return s.ClientMessage
}

func (s *Storage) GetChaId() int {
	return s.ChatId
}
