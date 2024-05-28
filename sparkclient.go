package gosparkclient

import (
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	AppIdEnvVarName       = "SPARKAI_APP_ID"     //nolint:gosec
	ApiKeyEnvVarName      = "SPARKAI_API_KEY"    //nolint:gosec
	ApiSecretEnvVarName   = "SPARKAI_API_SECRET" //nolint:gosec
	SparkDomainEnvVarName = "SPARKAI_DOMAIN"
	BaseURLEnvVarName     = "SPARKAI_URL" //nolint:gosec
)

const defaultEnvName = ".env"

// SparkClient 包含与 API 交互所需的配置信息
type SparkClient struct {
	AppID     string
	ApiSecret string
	ApiKey    string
	HostURL   string
	Domain    string
}

// SparkChatRequest 封装了调用 Spark API 所需的所有参数
type SparkChatRequest struct {
	Message []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"text"`
	Temperature  float64
	Topk         int
	Maxtokens    int
	System       string
	QuestionType string
}

var (
	loadEnvLock   sync.Mutex
	loadedEnvs    map[string]bool
	clientConfigs map[string]*SparkClient // 存储每个环境的客户端配置
)

// CallbackFunc 用于回调处理响应
type CallbackFunc func(response SparkAPIResponse)

func init() {
	loadedEnvs = make(map[string]bool)
	clientConfigs = make(map[string]*SparkClient)
	loadEnvIfNeeded(defaultEnvName) // 使用改进后的函数直接加载默认环境
}

// loadEnvIfNeeded 检查并加载指定环境配置，如果尚未加载
func loadEnvIfNeeded(envName string) *SparkClient {
	loadEnvLock.Lock()
	defer loadEnvLock.Unlock()

	// 检查配置是否已加载
	if client, exists := clientConfigs[envName]; exists {
		return client
	}

	// 加载环境配置文件
	env, err := godotenv.Read(envName)
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	// 读取环境变量并创建新的SparkClient实例
	client := &SparkClient{
		AppID:     env[AppIdEnvVarName],
		ApiSecret: env[ApiSecretEnvVarName],
		ApiKey:    env[ApiKeyEnvVarName],
		HostURL:   env[BaseURLEnvVarName],
		Domain:    env[SparkDomainEnvVarName],
	}

	// 保存到全局配置存储中
	clientConfigs[envName] = client
	loadedEnvs[envName] = true

	return client
}

func NewSparkClient() *SparkClient {
	return loadEnvIfNeeded(defaultEnvName)
}

func NewSparkClientWithEnv(envName string) *SparkClient {
	return loadEnvIfNeeded(envName)
}

func NewSparkClientWithOptions(appid, apikey, apisecret, hostURL, domain string) *SparkClient {
	return &SparkClient{
		AppID:     appid,
		ApiSecret: apisecret,
		ApiKey:    apikey,
		HostURL:   hostURL,
		Domain:    domain,
	}
}

func (client *SparkClient) SparkChatSimple(prompt string) (*SparkAPIResponse, error) {
	req := SparkChatRequest{}
	newMessage := struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		Role:    "user",
		Content: prompt,
	}
	req.Message = append(req.Message, newMessage)

	return client.SparkChatWithCallback(req, nil)
}

func (client *SparkClient) SparkChatWithCallback(req SparkChatRequest, callback CallbackFunc) (*SparkAPIResponse, error) {
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	authURL := client.AssembleAuthURL("GET", client.HostURL)
	conn, resp, err := d.Dial(authURL, nil)
	if err != nil {
		log.Printf("Failed to establish WebSocket connection: %v, %s, %s\n", err, ReadResp(resp), authURL)
		return nil, err
	}
	defer conn.Close()

	data := client.genReqJson(req)
	if err := conn.WriteJSON(data); err != nil {
		log.Printf("Failed to send message: %v\n", err)
		return nil, err
	}
	var response SparkAPIResponse
	var answer string
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read message error:", err)
			break
		}
		if err := json.Unmarshal(msg, &response); err != nil {
			log.Println("Error parsing JSON:", err)
			break
		}
		if response.Header.Code != 0 {
			return nil, errors.New(response.Header.Message)
		}
		if len(response.Payload.Choices.Text) > 0 {
			answer += response.Payload.Choices.Text[0].Content
			if callback != nil {
				callback(response)
			}
		}
		if response.Payload.Choices.Status == 2 {
			break
		}
	}

	if len(response.Payload.Choices.Text) > 0 {
		response.Payload.Choices.Text[0].Content = answer
	}

	return &response, err
}

// genReqJson 生成请求 JSON
func (client *SparkClient) genReqJson(usrReq SparkChatRequest) *SparkAPIRequest {
	var req SparkAPIRequest
	if usrReq.System != "" {
		req.Payload.Message.Text = append(req.Payload.Message.Text, struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{Role: "system", Content: usrReq.System})
	}

	// 填充Payload.Message.Text
	for _, msg := range usrReq.Message {
		req.Payload.Message.Text = append(req.Payload.Message.Text, struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	req.Header.AppID = client.AppID
	req.Header.UID = "12345"
	req.Parameter.Chat.Domain = client.Domain
	req.Parameter.Chat.Temperature = usrReq.Temperature
	req.Parameter.Chat.TopK = usrReq.Topk
	req.Parameter.Chat.MaxTokens = usrReq.Maxtokens
	req.Parameter.Chat.Auditing = "default"
	if usrReq.QuestionType != "" {
		req.Parameter.Chat.QuestionType = usrReq.QuestionType
	}
	return &req
}
