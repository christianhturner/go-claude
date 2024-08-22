package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) CreateMessages(ctx context.Context, body RequestBody) (*ResponseBody, error) {
	reqURL := c.config.BaseURL + c.config.Endpoint
	reqHeaders := map[string]string{
		"X-APi-Key":         c.config.ApiKey,
		"Anthropic-Version": c.config.Version,
		"Content-Type":      contentType,
	}
	if c.config.Beta != "" {
		reqHeaders["Anthropic-Beta"] = c.config.Beta
	}

	jsonBody, err := parseBodyJSON(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	resp, err := c.config.HTTPCLient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result ResponseBody
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}
	if (resp.StatusCode >= 400 && resp.StatusCode <= 500) || resp.StatusCode == 529 {
		var result ResponseError
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("json decode error: %w, status code: %d", err, resp.StatusCode)
		}
		return nil, fmt.Errorf("%s: %s", resp.Status, result.Error.Message)
	}
	return nil, fmt.Errorf("unexpected error: %d", resp.StatusCode)
}

func parseBodyJSON(req RequestBody) ([]byte, error) {
	for i, m := range req.Messages {
		if m.Content != "" {
			req.Messages[i].ContentRaw = m.Content
		}

		if len(m.ContentTypeText) > 0 {
			for j := range m.ContentTypeText {
				m.ContentTypeText[j].Type = "text"
			}
			raw, err := json.Marshal(m.ContentTypeText)
			if err != nil {
				return nil, err
			}
			req.Messages[i].ContentRaw = json.RawMessage(raw)
		}

		// Work on Image processing later
		// if len(m.ContentTypeImage) > 0 {
		// 	for j := range m.ContentTypeImage {
		// 		m.ContentTypeImage[j].Type = "image"
		// 	}
		// 	raw, err := json.Marshal(m.ContentTypeImage)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	req.Messages[i].ContentRaw = json.RawMessage(raw)
		// }
	}
	return json.Marshal(req)
}
