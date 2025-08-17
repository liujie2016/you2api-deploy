# You2API 测试命令

## 基本测试命令

### 1. 服务状态检查
```bash
curl https://you2api-deploy.vercel.app/
```

### 2. 调试端点测试
```bash
curl -X POST https://you2api-deploy.vercel.app/test \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Hello, how are you?"
      }
    ]
  }'
```

### 3. 基本API测试（英文）
```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {
        "role": "user",
        "content": "Hello! How can you help me today?"
      }
    ],
    "stream": false
  }'
```

### 4. 中文测试
```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {
        "role": "user",
        "content": "你好，请用中文介绍一下你自己"
      }
    ],
    "stream": false
  }'
```

### 5. 流式响应测试
```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {
        "role": "user",
        "content": "Write a short poem about artificial intelligence"
      }
    ],
    "stream": true
  }'
```

## 不同模型测试

### DeepSeek 模型
```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "model": "deepseek-chat",
    "messages": [
      {
        "role": "user",
        "content": "解释一下什么是深度学习"
      }
    ],
    "stream": false
  }'
```

### Claude 模型
```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "model": "claude-3.5-sonnet",
    "messages": [
      {
        "role": "user",
        "content": "What are the benefits of renewable energy?"
      }
    ],
    "stream": false
  }'
```

### Gemini 模型
```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "model": "gemini-1.5-pro",
    "messages": [
      {
        "role": "user",
        "content": "Explain quantum computing in simple terms"
      }
    ],
    "stream": false
  }'
```

## 多轮对话测试

```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {
        "role": "user",
        "content": "What is the capital of France?"
      },
      {
        "role": "assistant",
        "content": "The capital of France is Paris."
      },
      {
        "role": "user",
        "content": "What is the population of that city?"
      }
    ],
    "stream": false
  }'
```

## 错误处理测试

### 无效模型测试
```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "model": "invalid-model",
    "messages": [
      {
        "role": "user",
        "content": "Hello"
      }
    ],
    "stream": false
  }'
```

### 空消息测试
```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "model": "gpt-4o",
    "messages": [],
    "stream": false
  }'
```

### 缺少Authorization测试
```bash
curl -X POST https://you2api-deploy.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {
        "role": "user",
        "content": "Hello"
      }
    ],
    "stream": false
  }'
```

## 使用Python测试

```python
import requests
import json

# 基本请求
def test_basic_request():
    url = "https://you2api-deploy.vercel.app/v1/chat/completions"
    headers = {
        "Content-Type": "application/json",
        "Authorization": "Bearer test-token"
    }
    data = {
        "model": "gpt-4o",
        "messages": [
            {
                "role": "user",
                "content": "Hello! Can you help me with Python?"
            }
        ],
        "stream": False
    }
    
    response = requests.post(url, headers=headers, json=data)
    print(f"Status: {response.status_code}")
    print(f"Response: {response.json()}")

# 流式请求
def test_stream_request():
    url = "https://you2api-deploy.vercel.app/v1/chat/completions"
    headers = {
        "Content-Type": "application/json",
        "Authorization": "Bearer test-token"
    }
    data = {
        "model": "gpt-4o",
        "messages": [
            {
                "role": "user",
                "content": "Write a short story about a robot"
            }
        ],
        "stream": True
    }
    
    response = requests.post(url, headers=headers, json=data, stream=True)
    
    for line in response.iter_lines():
        if line:
            print(line.decode('utf-8'))

if __name__ == "__main__":
    test_basic_request()
    # test_stream_request()  # 取消注释以测试流式响应
```

## 使用Node.js测试

```javascript
// 使用 fetch API
async function testBasicRequest() {
    const response = await fetch('https://you2api-deploy.vercel.app/v1/chat/completions', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer test-token'
        },
        body: JSON.stringify({
            model: 'gpt-4o',
            messages: [
                {
                    role: 'user',
                    content: 'Hello! Can you help me with JavaScript?'
                }
            ],
            stream: false
        })
    });

    const data = await response.json();
    console.log('Status:', response.status);
    console.log('Response:', data);
}

// 使用 OpenAI SDK
import OpenAI from 'openai';

const openai = new OpenAI({
    apiKey: 'test-token',
    baseURL: 'https://you2api-deploy.vercel.app'
});

async function testWithSDK() {
    try {
        const completion = await openai.chat.completions.create({
            messages: [{ role: 'user', content: 'Hello from Node.js!' }],
            model: 'gpt-4o',
        });

        console.log(completion.choices[0].message.content);
    } catch (error) {
        console.error('Error:', error);
    }
}

testBasicRequest();
// testWithSDK();  // 需要安装 openai 包
```

## 测试建议

1. **按顺序测试**：
   - 先测试服务状态检查
   - 然后测试调试端点
   - 最后测试实际API

2. **查看日志**：
   - 在Vercel Dashboard中查看Function Logs
   - 如果需要详细日志，设置环境变量 `DEBUG=true`

3. **如果遇到问题**：
   - 设置环境变量 `USE_FALLBACK=true` 启用备用模式
   - 检查网络连接
   - 确认API密钥格式正确

4. **性能测试**：
   - 测试不同长度的输入
   - 测试并发请求
   - 测试流式vs非流式响应性能