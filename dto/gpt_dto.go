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

const (
	IMAGE_SIZE_256      = "256x256"
	IMAGE_SIZE_512      = "512x512"
	IMAGE_SIZE_1024     = "1024x1024"
	IMAGE_FROMAT_URL    = "url"
	IMAGE_FROMAT_BASE64 = "b64_json"
)

type CreateImageReq struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageResp struct {
	Created       int             `json:"created"`
	ImageContents []*ImageContent `json:"data"`
}
type ImageContent struct {
	URL     string `json:"url,omitempty"`
	B64Json string `json:"b64_json,omitempty"`
}

type TextEditReq struct {
	Model       string  `json:"model"`
	Input       string  `json:"input"`
	Instruction string  `json:"instruction,omitempty"`
	N           int     `json:"n,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	TopP        int     `json:"top_p,omitempty"`
}
type TextEditResp struct {
	Object  string            `json:"object"`
	Created int               `json:"created"`
	Choices []*TextEditChoice `json:"choices"`
	Usage   *Usage            `json:"usage"`
}
type TextEditChoice struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}

type ErrorResp struct {
	Error *ErrorContent `json:"error"`
}
type ErrorContent struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Param   interface{} `json:"param"`
	Type    string      `json:"type"`
}
