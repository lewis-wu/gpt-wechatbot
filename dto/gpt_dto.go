package dto

type ChatCompleteReq struct {
	Model            string                 `json:"model"`
	Messages         []*Message             `json:"messages"`
	MaxTokens        int                    `json:"max_tokens,omitempty"`
	Temperature      float32                `json:"temperature,omitempty"`
	TopP             int                    `json:"top_p,omitempty"`
	FrequencyPenalty int                    `json:"frequency_penalty,omitempty"`
	PresencePenalty  int                    `json:"presence_penalty,omitempty"`
	N                int                    `json:"n,omitempty"`
	Stop             []*string              `json:"stop,omitempty"`
	LogitBias        map[string]interface{} `json:"logit_bias,omitempty"`
	User             string                 `json:"user,omitempty"`
	Stream           bool                   `json:"stream"`
}
type Message struct {
	Role    string  `json:"role"`
	Content string  `json:"content"`
	Name    *string `json:"name,omitempty"`
}

type ChatCompleteResp struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Choices []*Choice `json:"choices"`
	Usage   *Usage    `json:"usage"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
}
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
