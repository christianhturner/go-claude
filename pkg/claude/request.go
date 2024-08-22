package claude

type RequestBody struct {
	Model         string                 `json:"model"`
	Messages      []RequestMessages      `json:"messages"`
	System        string                 `json:"system"` // optional
	MaxTokens     int                    `json:"max_tokens"`
	MetaData      map[string]interface{} `json:"metadata"`       // optional
	StopSequences []string               `json:"stop_sequences"` // optional
	Stream        bool                   `json:"stream"`         // optional
	Temparature   float64                `json:"Temparature"`    // optional
	TopP          float64                `json:"top_p"`          // optional
	TopK          float64                `json:"top_k"`
}

type RequestMessages struct {
	Role            string      `json:"role"`
	ContentRaw      interface{} `json:"content"`
	Content         string      `json:"-"`
	ContentTypeText []RequestContentTypeText
	// add option for images
}

const (
	RequestContentTypeTextType = "text"
)

type RequestContentTypeText struct {
	Type string `json:"type"` // always "text"
	Text string `json:"text"`
}

const (
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
)

// Inspirational struct https://github.com/potproject/claude-sdk-go/blob/main/request.go
