# SparkClient Examples

这个目录包含了 SparkClient 的使用示例。

## 运行示例

1. 首先复制环境配置文件：

```bash
cp .env.example .env
```

2. 编辑 `.env` 文件，填入你的认证信息：
- SPARKAI_APP_ID：你的应用ID
- SPARKAI_API_KEY：你的API密钥
- SPARKAI_API_SECRET：你的API密钥密码
- SPARKAI_URL：API地址（通常不需要修改）
- SPARKAI_DOMAIN：使用的模型版本（通常不需要修改）

3. 运行示例程序：

```bash
go run main.go
```

## 示例说明

示例程序包含了三个主要用例：

1. 简单对话（simpleChat）
    - 展示了最基本的对话用法
    - 使用 ChatSimple 方法直接发送消息

2. 高级对话（advancedChat）
    - 展示了更多参数的使用
    - 包含温度、最大token等设置
    - 展示了token使用统计信息

3. 对话历史（chatHistory）
    - 展示了如何维护对话历史
    - 展示了多轮对话的实现方式

## 注意事项

- 请确保环境变量正确设置
- 示例中的超时时间可以根据需要调整
- 建议在正式环境中适当处理错误情况