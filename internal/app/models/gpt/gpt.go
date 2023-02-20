package gpt

type Response struct {
	Id      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
	Usage   Usage     `json:"usage"`
}

type Choices struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     string `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	PromprTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Storage struct {
}

func NewStorage() *Storage {
	return &Storage{}
}
