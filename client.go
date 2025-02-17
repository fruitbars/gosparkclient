package gosparkclient

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

type SparkClient struct {
	config    *Config
	transport *http.Transport
}

func NewSparkClient(opts ...ConfigOption) (*SparkClient, error) {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	if err := validateConfig(config); err != nil {
		return nil, newConfigError("invalid configuration", err)
	}

	return &SparkClient{
		config:    config,
		transport: defaultTransport(config.Timeout),
	}, nil
}

// ChatWithCallback initiates a chat session and calls the callback function for each response
func (c *SparkClient) ChatWithCallback(ctx context.Context, req *SparkChatRequest, callback ChatCallback) error {
	dialer := websocket.Dialer{
		HandshakeTimeout: c.config.Timeout,
		NetDialContext:   c.transport.DialContext,
		Proxy:            c.transport.Proxy,
	}

	authURL := c.assembleAuthURL("GET", c.config.HostURL)
	conn, _, err := dialer.DialContext(ctx, authURL, nil)
	if err != nil {
		return newConnectionError("failed to establish WebSocket connection", err)
	}
	defer conn.Close()

	if err := conn.WriteJSON(c.genReqJson(req)); err != nil {
		return newRequestError("failed to send message", err)
	}

	for {
		select {
		case <-ctx.Done():
			return newRequestError("request cancelled", ctx.Err())
		default:
			var response SparkAPIResponse
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return newWebSocketError("failed to read message", err)
			}

			if err := json.Unmarshal(msg, &response); err != nil {
				return newResponseError("failed to parse response", err)
			}

			if response.Header.Code != 0 {
				return newResponseError(response.Header.Message, nil)
			}

			// Call the callback function with the response
			if callback != nil {
				callback(&response)
			}

			if response.Payload.Choices.Status == 2 {
				return nil
			}
		}
	}
}

func (c *SparkClient) Chat(ctx context.Context, req *SparkChatRequest) (*SparkAPIResponse, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: c.config.Timeout,
		NetDialContext:   c.transport.DialContext,
		Proxy:            c.transport.Proxy,
	}

	authURL := c.assembleAuthURL("GET", c.config.HostURL)
	conn, _, err := dialer.DialContext(ctx, authURL, nil)
	if err != nil {
		return nil, newConnectionError("failed to establish WebSocket connection", err)
	}
	defer conn.Close()

	if err := conn.WriteJSON(c.genReqJson(req)); err != nil {
		return nil, newRequestError("failed to send message", err)
	}

	var finalResponse *SparkAPIResponse
	var answer string

	for {
		select {
		case <-ctx.Done():
			return nil, newRequestError("request cancelled", ctx.Err())
		default:
			var response SparkAPIResponse
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return nil, newWebSocketError("failed to read message", err)
			}

			if err := json.Unmarshal(msg, &response); err != nil {
				return nil, newResponseError("failed to parse response", err)
			}

			if response.Header.Code != 0 {
				return nil, newResponseError(response.Header.Message, nil)
			}

			if len(response.Payload.Choices.Text) > 0 {
				answer += response.Payload.Choices.Text[0].Content
			}

			if response.Payload.Choices.Status == 2 {
				if len(response.Payload.Choices.Text) > 0 {
					response.Payload.Choices.Text[0].Content = answer
				}
				finalResponse = &response
				break
			}
		}
		if finalResponse != nil {
			break
		}
	}

	return finalResponse, nil
}

func (c *SparkClient) ChatSimple(ctx context.Context, prompt string) (*SparkAPIResponse, error) {
	req := &SparkChatRequest{
		Messages: []SparkMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	return c.Chat(ctx, req)
}

func (c *SparkClient) Embedding(ctx context.Context, query, domain string) (*SparkAPIEmbResponse, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: c.config.Timeout,
		NetDialContext:   c.transport.DialContext,
		Proxy:            c.transport.Proxy,
	}

	authURL := c.assembleAuthURL("GET", c.config.EMBURL)
	conn, _, err := dialer.DialContext(ctx, authURL, nil)
	if err != nil {
		return nil, newConnectionError("failed to establish WebSocket connection", err)
	}
	defer conn.Close()

	req := c.getEmbeddingRequest(query, domain)
	if err := conn.WriteJSON(req); err != nil {
		return nil, newRequestError("failed to send embedding request", err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, newWebSocketError("failed to read message", err)
	}

	var response SparkAPIEmbResponse
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, newResponseError("failed to parse response", err)
	}

	if response.Header.Code != 0 {
		return nil, newResponseError(response.Header.Message, nil)
	}

	return &response, nil
}

func (c *SparkClient) WithNewConfig(opts ...ConfigOption) (*SparkClient, error) {
	newConfig := *c.config
	for _, opt := range opts {
		opt(&newConfig)
	}

	if err := validateConfig(&newConfig); err != nil {
		return nil, newConfigError("invalid configuration", err)
	}

	return &SparkClient{
		config:    &newConfig,
		transport: defaultTransport(newConfig.Timeout),
	}, nil
}

func (c *SparkClient) genReqJson(req *SparkChatRequest) *SparkAPIRequest {
	apiReq := &SparkAPIRequest{}
	apiReq.Header.AppID = c.config.AppID
	apiReq.Header.UID = c.config.UID
	apiReq.Parameter.Chat.Domain = c.config.Domain
	apiReq.Parameter.Chat.Temperature = req.Temperature
	apiReq.Parameter.Chat.TopK = req.TopK
	apiReq.Parameter.Chat.MaxTokens = req.MaxTokens
	apiReq.Parameter.Chat.Auditing = c.config.Auditing
	apiReq.Parameter.Chat.QuestionType = req.QuestionType

	if req.System != "" {
		apiReq.Payload.Message.Text = append(apiReq.Payload.Message.Text, SparkMessage{
			Role:    "system",
			Content: req.System,
		})
	}

	apiReq.Payload.Message.Text = append(apiReq.Payload.Message.Text, req.Messages...)

	if req.Functions != nil && len(req.Functions) > 0 {
		apiReq.Functions = &struct {
			Text json.RawMessage `json:"text,omitempty"`
		}{
			Text: req.Functions,
		}
	}

	return apiReq
}

func (c *SparkClient) getEmbeddingRequest(query, domain string) *SparkAPIEmbRequest {
	return &SparkAPIEmbRequest{
		Header: struct {
			AppID  string `json:"app_id"`
			UID    string `json:"uid"`
			Status int    `json:"status"`
		}{
			AppID:  c.config.AppID,
			UID:    c.config.UID,
			Status: 3,
		},
		Parameter: struct {
			Emb struct {
				Domain  string `json:"domain"`
				Feature struct {
					Encoding string `json:"encoding"`
					Compress string `json:"compress"`
					Format   string `json:"format"`
				} `json:"feature"`
			} `json:"emb"`
		}{
			Emb: struct {
				Domain  string `json:"domain"`
				Feature struct {
					Encoding string `json:"encoding"`
					Compress string `json:"compress"`
					Format   string `json:"format"`
				} `json:"feature"`
			}{
				Domain: domain,
				Feature: struct {
					Encoding string `json:"encoding"`
					Compress string `json:"compress"`
					Format   string `json:"format"`
				}{
					Encoding: "utf8",
					Compress: "raw",
					Format:   "plain",
				},
			},
		},
		Payload: struct {
			Message struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
				Status   int    `json:"status"`
				Text     string `json:"text"`
			} `json:"message"`
		}{
			Message: struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
				Status   int    `json:"status"`
				Text     string `json:"text"`
			}{
				Encoding: "utf8",
				Compress: "raw",
				Format:   "json",
				Status:   3,
				Text:     query,
			},
		},
	}
}
