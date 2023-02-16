package models

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
