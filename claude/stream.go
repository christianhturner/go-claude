// file: stream.go
package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/tmaxmax/go-sse"
)

const (
	MessagesStreamResponseTypeMessageStart      = "message_start"
	MessagesStreamResponseTypeContentBlockStart = "content_block_start"
	MessagesStreamResponseTypePing              = "ping"
	MessagesStreamResponseTypeContentBlockDelta = "content_block_delta"
	MessagesStreamResponseTypeContentBlockStop  = "content_block_stop"
	MessagesStreamResponseTypeMessageDelta      = "message_delta"
	MessagesStreamResponseTypeMessageStop       = "message_stop"
	MessagesStreamResponseTypeError             = "error"
)

type CreateMessagesStream struct {
	Connection                 *sse.Connection
	Unsubscribe                func()
	Event                      chan sse.Event
	Error                      chan error
	ResponseBodyMessagesStream ResponseBodyStream
}

type ResponseBodyStream struct {
	Id           string                   `json:"id"`
	Type         string                   `json:"type"` // always "message"
	Role         string                   `json:"role"` // always "assistant"
	Content      []ResponseMessagesStream `json:"content"`
	Model        string                   `json:"model"`
	StopReason   string                   `json:"stop_reason"` // "end_turn" or "max_tokens", "stop_sequence", null
	StopSequence string                   `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int64 `json:"input_tokens"`
		OutputTokens int64 `json:"output_tokens"`
	} `json:"usage"`
}

type ResponseMessagesStream struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ResponseMessageStartStream struct {
	Type    string             `json:"type"`
	Message ResponseBodyStream `json:"message"`
}

type ResponseBlockDeltaStream struct {
	Type  string `json:"type"`
	Index int64  `json:"index"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
}

type ResponseMessageDeltaStream struct {
	Type  string `json:"type"`
	Delta struct {
		StopReason   string `json:"stop_reason"`
		StopSequence string `json:"stop_sequence"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int64 `json:"output_tokens"`
	} `json:"usage"`
}

func (c *Client) CreateMessagesStream(ctx context.Context, body RequestBody) (*CreateMessagesStream, error) {
	reqURL := c.config.BaseURL + c.config.Endpoint
	body.Stream = true
	reqHeaders := map[string]string{
		"X-Api-Key":         c.config.ApiKey,
		"Anthropic-Version": c.config.Version,
		"Content-Type":      contentType,
	}
	if c.config.Beta != "" {
		reqHeaders["anthropic-beta"] = c.config.Beta
	}

	jsonBody, err := parseBodyJSON(body)
	if err != nil {
		return nil, err
	}

	client := sse.Client{
		HTTPClient: c.config.HTTPCLient,
		Backoff: sse.Backoff{
			MaxRetries: -1,
		},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}

	conn := client.NewConnection(req)
	chanEvent := make(chan sse.Event)
	connectionError := make(chan error)

	unsubscribe := conn.SubscribeToAll(func(e sse.Event) {
		if e.Type == MessagesStreamResponseTypePing ||
			e.Type == MessagesStreamResponseTypeContentBlockStart ||
			e.Type == MessagesStreamResponseTypeContentBlockStop {
			return
		}
		chanEvent <- e
	})
	go func() {
		err := conn.Connect()
		if !errors.Is(err, io.EOF) && err != nil {
			connectionError <- err
		}
	}()
	return &CreateMessagesStream{
		Connection:                 conn,
		Unsubscribe:                unsubscribe,
		Event:                      chanEvent,
		Error:                      connectionError,
		ResponseBodyMessagesStream: ResponseBodyStream{},
	}, nil
}

func (c *CreateMessagesStream) Close() {
	close(c.Event)
	close(c.Error)
	c.Unsubscribe()
}

func (c *CreateMessagesStream) Recv() (ResponseBodyStream, error) {
	select {
	case e := <-c.Event:
		switch e.Type {
		case MessagesStreamResponseTypeMessageStart:
			d := []byte(e.Data)
			var r ResponseMessageStartStream
			err := json.Unmarshal(d, &r)
			if err != nil {
				return ResponseBodyStream{}, err
			}
			c.ResponseBodyMessagesStream = r.Message
			c.ResponseBodyMessagesStream.Content = []ResponseMessagesStream{
				{
					Type: "text",
					Text: "",
				},
			}
			return c.ResponseBodyMessagesStream, nil
		case MessagesStreamResponseTypeContentBlockDelta:
			d := []byte(e.Data)
			var r ResponseBlockDeltaStream
			err := json.Unmarshal(d, &r)
			if err != nil {
				return ResponseBodyStream{}, err
			}
			c.ResponseBodyMessagesStream.Content = []ResponseMessagesStream{
				{
					Type: "text",
					Text: r.Delta.Text,
				},
			}
			return c.ResponseBodyMessagesStream, nil
		case MessagesStreamResponseTypeMessageDelta:
			d := []byte(e.Data)
			var r ResponseMessageDeltaStream
			err := json.Unmarshal(d, &r)
			if err != nil {
				return ResponseBodyStream{}, err
			}
			c.ResponseBodyMessagesStream.Content = []ResponseMessagesStream{}
			c.ResponseBodyMessagesStream.StopReason = r.Delta.StopReason
			c.ResponseBodyMessagesStream.StopSequence = r.Delta.StopSequence
			c.ResponseBodyMessagesStream.Usage.OutputTokens = r.Usage.OutputTokens
			c.ResponseBodyMessagesStream.Content = []ResponseMessagesStream{
				{
					Type: "text",
					Text: "",
				},
			}
			return c.ResponseBodyMessagesStream, nil

		case MessagesStreamResponseTypeMessageStop:
			return c.ResponseBodyMessagesStream, io.EOF
		case MessagesStreamResponseTypeError:
			d := []byte(e.Data)
			var r ResponseError
			err := json.Unmarshal(d, &r)
			if err != nil {
				return ResponseBodyStream{}, err
			}
			return c.ResponseBodyMessagesStream, errors.New(r.Error.Message)
		}
	case err := <-c.Error:
		return ResponseBodyStream{}, err
	}
	return c.ResponseBodyMessagesStream, nil
}
