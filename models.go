package gosparkclient

import "encoding/json"

type ChatCallback func(resp *SparkAPIResponse)

// SparkMessage represents a single message in the conversation
type SparkMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SparkHeader represents common header structure
type SparkHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	SID     string `json:"sid"`
	Status  int    `json:"status"`
}

// SparkUsage represents token usage information
type SparkUsage struct {
	QuestionTokens   int `json:"question_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// SparkFunctionCall represents a function call in the response
type SparkFunctionCall struct {
	Arguments string `json:"arguments"`
	Name      string `json:"name"`
}

// SparkChoice represents a single choice in the response
type SparkChoice struct {
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	Role             string            `json:"role"`
	ContentType      string            `json:"content_type"`
	FunctionCall     SparkFunctionCall `json:"function_call"`
	Index            int               `json:"index"`
}

// SparkChatRequest represents a chat request to the Spark API
type SparkChatRequest struct {
	Messages     []SparkMessage  `json:"text"`
	Temperature  float64         `json:"temperature,omitempty"`
	TopK         int             `json:"top_k,omitempty"`
	MaxTokens    int             `json:"max_tokens,omitempty"`
	System       string          `json:"system,omitempty"`
	QuestionType string          `json:"question_type,omitempty"`
	Functions    json.RawMessage `json:"functions,omitempty"`
}

// SparkAPIRequest represents the full API request structure
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
			Text []SparkMessage `json:"text"`
		} `json:"message"`
	} `json:"payload"`
	Functions *struct {
		Text json.RawMessage `json:"text,omitempty"`
	} `json:"functions,omitempty"`
}

// SparkAPIResponse represents the API response structure
type SparkAPIResponse struct {
	Header  SparkHeader `json:"header"`
	Payload struct {
		Choices struct {
			Status int           `json:"status"`
			Seq    int           `json:"seq"`
			Text   []SparkChoice `json:"text"`
		} `json:"choices"`
		Usage struct {
			Text SparkUsage `json:"text"`
		} `json:"usage"`
		Plugins *SparkPlugins `json:"plugins,omitempty"`
	} `json:"payload"`
}

// SparkPlugins represents plugin-related information
type SparkPlugins struct {
	Text []struct {
		Name        string `json:"name"`
		Content     string `json:"content"`
		ContentType string `json:"content_type"`
		ContentMeta any    `json:"content_meta"`
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
}

// SparkFeature represents feature-related information
type SparkFeature struct {
	Encoding string `json:"encoding"`
	Compress string `json:"compress"`
	Format   string `json:"format"`
}

// SparkAPIEmbRequest represents an embedding request
type SparkAPIEmbRequest struct {
	Header struct {
		AppID  string `json:"app_id"`
		UID    string `json:"uid"`
		Status int    `json:"status"`
	} `json:"header"`
	Parameter struct {
		Emb struct {
			Domain  string `json:"domain"`
			Feature struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
			} `json:"feature"`
		} `json:"emb"`
	} `json:"parameter"`
	Payload struct {
		Message struct {
			Encoding string `json:"encoding"`
			Compress string `json:"compress"`
			Format   string `json:"format"`
			Status   int    `json:"status"`
			Text     string `json:"text"`
		} `json:"message"`
	} `json:"payload"`
}

// SparkAPIEmbResponse represents an embedding response
type SparkAPIEmbResponse struct {
	Header  SparkHeader `json:"header"`
	Payload struct {
		Feature SparkFeature `json:"feature"`
	} `json:"payload"`
}
