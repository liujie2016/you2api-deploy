# You2API

将 You.com 的 AI 聊天服务转换为 OpenAI 兼容的 API 接口，支持部署到 Vercel。

## 功能特性

- 🚀 **OpenAI 兼容**: 完全兼容 OpenAI Chat Completions API
- 🌊 **流式响应**: 支持流式和非流式响应
- 🤖 **多模型支持**: 支持多种 AI 模型（GPT、Claude、Gemini 等）
- 🔄 **自动模型映射**: 自动将 OpenAI 模型名称映射到 You.com 对应模型
- 🛡️ **CORS 支持**: 内置跨域资源共享支持
- ☁️ **一键部署**: 支持一键部署到 Vercel

## 支持的模型

| OpenAI 模型名称 | You.com 模型 |
|----------------|-------------|
| `deepseek-reasoner` | `deepseek_r1` |
| `deepseek-chat` | `deepseek_v3` |
| `o3-mini-high` | `openai_o3_mini_high` |
| `o3-mini-medium` | `openai_o3_mini_medium` |
| `o1` | `openai_o1` |
| `o1-mini` | `openai_o1_mini` |
| `gpt-4o` | `gpt_4o` |
| `gpt-4o-mini` | `gpt_4o_mini` |
| `claude-3.5-sonnet` | `claude_3_5_sonnet` |
| `gemini-1.5-pro` | `gemini_1_5_pro` |
| 更多... | 更多... |

## 快速部署

### Vercel 部署（推荐）

1. **Fork 此项目**
   ```bash
   git clone https://github.com/yourusername/you2api-deploy.git
   cd you2api-deploy
   ```

2. **部署到 Vercel**
   
   [![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/yourusername/you2api-deploy)
   
   或手动部署：
   ```bash
   npm i -g vercel
   vercel --prod
   ```

3. **配置完成**
   
   部署完成后，你将获得一个 Vercel 域名，如：`https://your-project.vercel.app`

### 本地运行

```bash
# 启动服务
./start.sh

# 停止服务
./stop.sh
```

## API 使用方法

### 基本用法

```bash
curl -X POST https://your-project.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token-here" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {
        "role": "user", 
        "content": "Hello, how are you?"
      }
    ],
    "stream": false
  }'
```

### 流式响应

```bash
curl -X POST https://your-project.vercel.app/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token-here" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {
        "role": "user", 
        "content": "Write a short poem"
      }
    ],
    "stream": true
  }'
```

### Python 示例

```python
import openai

# 配置客户端
client = openai.OpenAI(
    api_key="your-token-here",
    base_url="https://your-project.vercel.app"
)

# 发送请求
response = client.chat.completions.create(
    model="gpt-4o",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

print(response.choices[0].message.content)
```

### Node.js 示例

```javascript
import OpenAI from 'openai';

const openai = new OpenAI({
  apiKey: 'your-token-here',
  baseURL: 'https://your-project.vercel.app',
});

async function main() {
  const completion = await openai.chat.completions.create({
    messages: [{ role: 'user', content: 'Hello!' }],
    model: 'gpt-4o',
  });

  console.log(completion.choices[0].message.content);
}

main();
```

## API 端点

### POST `/v1/chat/completions`

兼容 OpenAI Chat Completions API 的主要端点。

**请求参数:**
- `model`: 模型名称（自动映射到 You.com 对应模型）
- `messages`: 消息数组
- `stream`: 是否使用流式响应（可选，默认 false）

**响应格式:**
- 非流式：标准 OpenAI Chat Completion 响应格式
- 流式：Server-Sent Events 格式

### GET `/`

服务状态检查端点，返回服务运行状态。

## 项目结构

```
you2api-deploy/
├── api/
│   └── main.go          # 主要 API 处理逻辑
├── go.mod               # Go 模块配置
├── vercel.json          # Vercel 部署配置
├── start.sh             # 本地启动脚本
├── stop.sh              # 本地停止脚本
└── README.md            # 项目说明
```

## 注意事项

### 安全性
- 项目使用了 CORS 代理服务来绕过 Cloudflare 保护
- 请确保 API 密钥的安全性
- 建议在生产环境中实施适当的访问控制

### 限制
- 依赖于第三方 CORS 代理服务（proxy.cors.sh）
- You.com 的 API 可能有访问频率限制
- 某些高级功能可能不完全兼容

### 故障排除
1. **部署失败**: 检查 `go.mod` 和 `vercel.json` 配置
2. **请求失败**: 确认 Authorization 头部格式正确
3. **模型不支持**: 检查模型名称是否在支持列表中

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

本项目采用 MIT 许可证。

## 免责声明

本项目仅供学习和研究使用。请遵守相关服务的使用条款和条件。