package gosparkclient

import "encoding/json"

type SparkAPIRequest struct {
	Header struct {
		AppID string `json:"app_id"`
		UID   string `json:"uid"`
	} `json:"header"`
	Parameter struct {
		Chat struct {
			Domain       string  `json:"domain"`
			Temperature  float64 `json:"temperature,omitempty"`
			MaxTokens    int     `json:"max_tokens,omitempty"`
			TopK         int     `json:"top_k,omitempty"`
			Auditing     string  `json:"auditing"`
			QuestionType string  `json:"question_type,omitempty"`
		} `json:"chat"`
	} `json:"parameter"`
	Payload struct {
		Message struct {
			Text []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"text"`
		} `json:"message"`
	} `json:"payload"`
	Functions json.RawMessage `json:"functions,omitempty"`
}

type SparkAPIResponse struct {
	Header struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Sid     string `json:"sid"`
		Status  int    `json:"status"`
	} `json:"header"`
	Payload struct {
		Choices struct {
			Status int `json:"status"`
			Seq    int `json:"seq"`
			Text   []struct {
				Content      string `json:"content"`
				Role         string `json:"role"`
				ContentType  string `json:"content_type"`
				FunctionCall struct {
					Arguments string `json:"arguments"`
					Name      string `json:"name"`
				} `json:"function_call"`
				Index int `json:"index"`
			} `json:"text"`
		} `json:"choices"`
		Usage struct {
			Text struct {
				QuestionTokens   int `json:"question_tokens"`
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			} `json:"text"`
		} `json:"usage"`

		Plugins struct {
			Text []struct {
				Name        string `json:"name"`
				Content     string `json:"content"`
				ContentType string `json:"content_type"`
				ContentMeTa any    `json:"content_me ta"`
				Role        string `json:"role"`
				Status      string `json:"status"`
				Invoked     struct {
					Namespace  string `json:"namespace"`
					PluginID   string `json:"plugin_id"`
					PluginVer  string `json:"plugin_ver"`
					StatusCode int    `json:"status_code"`
					StatusMsg  string `json:"status_msg"`
					Type       string `json:"type"`
				} `json:"invoked"`
			} `json:"text"`
		} `json:"plugins,omitempty"`
	} `json:"payload"`
}
