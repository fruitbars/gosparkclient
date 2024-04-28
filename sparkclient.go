package gosparkclient

import (
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	SPARK_API_KEY    string
	SPARK_API_SECRET string
	SPARK_API_URL    string
	SPARK_APP_ID     string
	SPARK_DOMAIN     string
)

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
	Prompt       string
	Temperature  float64
	Topk         int
	Maxtokens    int
	System       string
	His          HisContent
	QuestionType string
}

var once sync.Once

func loadEnv(envName string) {
	const (
		AppIdEnvVarName       = "SPARKAI_APP_ID"     //nolint:gosec
		ApiKeyEnvVarName      = "SPARKAI_API_KEY"    //nolint:gosec
		ApiSecretEnvVarName   = "SPARKAI_API_SECRET" //nolint:gosec
		SparkDomainEnvVarName = "SPARKAI_DOMAIN"
		BaseURLEnvVarName     = "SPARKAI_URL" //nolint:gosec
	)

	err := godotenv.Load(envName)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	SPARK_API_KEY = os.Getenv(ApiKeyEnvVarName)
	SPARK_API_SECRET = os.Getenv(ApiSecretEnvVarName)
	SPARK_API_URL = os.Getenv(BaseURLEnvVarName)
	SPARK_APP_ID = os.Getenv(AppIdEnvVarName)
	SPARK_DOMAIN = os.Getenv(SparkDomainEnvVarName)
}

func init() {
	loadEnv(".env")
}

// CallbackFunc 用于回调处理响应
type CallbackFunc func(response SparkAPIResponse)

func NewSparkClient() *SparkClient {
	return &SparkClient{
		AppID:     SPARK_APP_ID,
		ApiSecret: SPARK_API_SECRET,
		ApiKey:    SPARK_API_KEY,
		HostURL:   SPARK_API_URL,
		Domain:    SPARK_DOMAIN,
	}
}

func NewSparkClientWithEnv(envName string) *SparkClient {
	once.Do(func() {
		loadEnv(envName)
	})
	return &SparkClient{
		AppID:     SPARK_APP_ID,
		ApiSecret: SPARK_API_SECRET,
		ApiKey:    SPARK_API_KEY,
		HostURL:   SPARK_API_URL,
		Domain:    SPARK_DOMAIN,
	}
}

func NewSparkClientWithOptions(appid, apikey, apisecret, hostURL, domain string) *SparkClient {
	return &SparkClient{
		AppID:     appid,
		ApiSecret: apikey,
		ApiKey:    apisecret,
		HostURL:   hostURL,
		Domain:    domain,
	}
}

func (client *SparkClient) SparkChatSimple(prompt string) (string, string, error) {
	req := SparkChatRequest{
		Prompt: prompt,
	}
	return client.SparkChatWithCallback(req, nil)
}

// CallSpark 发起与 Spark API 的通信，使用 SparkChatRequest 结构体简化参数
func (client *SparkClient) SparkChatWithCallback(req SparkChatRequest, callback CallbackFunc) (string, string, error) {
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	authURL := client.AssembleAuthURL("GET", client.HostURL)
	conn, resp, err := d.Dial(authURL, nil)
	if err != nil {
		log.Printf("Failed to establish WebSocket connection: %v, %s, %s\n", err, ReadResp(resp), authURL)
		return "", "", err
	}
	defer conn.Close()

	data := client.genReqJson(req)
	if err := conn.WriteJSON(data); err != nil {
		log.Printf("Failed to send message: %v\n", err)
		return "", "", err
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
			return "", response.Header.Sid, errors.New(response.Header.Message)
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
	return answer, response.Header.Sid, err
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
	for _, text := range usrReq.His.Text {
		req.Payload.Message.Text = append(req.Payload.Message.Text, struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{Role: text.Role, Content: text.Content})
	}
	req.Payload.Message.Text = append(req.Payload.Message.Text, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{Role: "user", Content: usrReq.Prompt})

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
