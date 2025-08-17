package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type YouChatResponse struct {
	YouChatToken string `json:"youChatToken"`
}

type OpenAIStreamResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Delta        Delta  `json:"delta"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type Delta struct {
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Model    string    `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
}

type OpenAIChoice struct {
	Message      Message `json:"message"`
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
}

var modelMap = map[string]string{
	"deepseek-reasoner":  "deepseek_r1",
	"deepseek-chat":      "deepseek_v3",
	"o3-mini-high":       "openai_o3_mini_high",
	"o3-mini-medium":     "openai_o3_mini_medium",
	"o1":                 "openai_o1",
	"o1-mini":            "openai_o1_mini",
	"o1-preview":         "openai_o1_preview",
	"gpt-4o":             "gpt_4o",
	"gpt-4o-mini":        "gpt_4o_mini",
	"gpt-4-turbo":        "gpt_4_turbo",
	"gpt-3.5-turbo":      "gpt_3_5",
	"claude-3-opus":      "claude_3_opus",
	"claude-3-sonnet":    "claude_3_sonnet",
	"claude-3.5-sonnet":  "claude_3_5_sonnet",
	"claude-3.5-haiku":   "claude_3_5_haiku",
	"gemini-1.5-pro":     "gemini_1_5_pro",
	"gemini-1.5-flash":   "gemini_1_5_flash",
	"llama-3.2-90b":      "llama3_2_90b",
	"llama-3.1-405b":     "llama3_1_405b",
	"mistral-large-2":    "mistral_large_2",
	"qwen-2.5-72b":       "qwen2p5_72b",
	"qwen-2.5-coder-32b": "qwen2p5_coder_32b",
	"command-r-plus":     "command_r_plus",
}

func getReverseModelMap() map[string]string {
	reverse := make(map[string]string, len(modelMap))
	for k, v := range modelMap {
		reverse[v] = k
	}
	return reverse
}

func mapModelName(openAIModel string) string {
	if mappedModel, exists := modelMap[openAIModel]; exists {
		return mappedModel
	}
	return "deepseek_v3"
}

func reverseMapModelName(youModel string) string {
	reverseMap := getReverseModelMap()
	if mappedModel, exists := reverseMap[youModel]; exists {
		return mappedModel
	}
	return "deepseek-chat"
}

var originalModel string

func Handler(w http.ResponseWriter, r *http.Request) {
	// 检查是否应该使用备用处理器
	if os.Getenv("USE_FALLBACK") == "true" || os.Getenv("FALLBACK_MODE") == "true" {
		FallbackHandler(w, r)
		return
	}

	if r.URL.Path != "/v1/chat/completions" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "You2Api Service Running...",
			"message": "MoLoveSze...",
		})
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid authorization header", http.StatusUnauthorized)
		return
	}
	dsToken := strings.TrimPrefix(authHeader, "Bearer ")
	_ = dsToken // Token for future use

	var openAIReq OpenAIRequest
	if err := json.NewDecoder(r.Body).Decode(&openAIReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(openAIReq.Messages) == 0 {
		http.Error(w, "Messages array cannot be empty", http.StatusBadRequest)
		return
	}

	originalModel = openAIReq.Model
	lastMessage := openAIReq.Messages[len(openAIReq.Messages)-1].Content
	var chatHistory []map[string]interface{}
	for _, msg := range openAIReq.Messages {
		chatMsg := map[string]interface{}{
			"question": msg.Content,
			"answer":   "",
		}
		if msg.Role == "assistant" {
			chatMsg["question"] = ""
			chatMsg["answer"] = msg.Content
		}
		chatHistory = append(chatHistory, chatMsg)
	}

	chatHistoryJSON, _ := json.Marshal(chatHistory)

	// Construct the original URL
	originalURL := "https://you.com/api/streamingSearch"
	q := url.Values{}
	q.Add("q", lastMessage)
	q.Add("page", "1")
	q.Add("count", "10")
	q.Add("safeSearch", "Moderate")
	q.Add("mkt", "zh-HK")
	q.Add("enable_worklow_generation_ux", "true")
	q.Add("domain", "youchat")
	q.Add("use_personalization_extraction", "true")
	q.Add("pastChatLength", fmt.Sprintf("%d", len(chatHistory)-1))
	q.Add("selectedChatMode", "custom")
	q.Add("selectedAiModel", mapModelName(openAIReq.Model))
	q.Add("enable_agent_clarification_questions", "true")
	q.Add("use_nested_youchat_updates", "true")
	q.Add("chat", string(chatHistoryJSON))

	// Use a proxy to bypass Cloudflare
	proxyURL := "https://proxy.cors.sh/" + originalURL + "?" + q.Encode()

	youReq, err := http.NewRequest("GET", proxyURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	youReq.Header = http.Header{
		"x-cors-api-key":     {"live_a48b9b66e68b4b0bb41a3df6de21e59b4a28cfc55b1343a0b0b0f5b5c2e8e8c7"},
		"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"Accept":             {"text/event-stream"},
		"Accept-Language":    {"en-US,en;q=0.9"},
		"Accept-Encoding":    {"gzip, deflate, br"},
		"Referer":            {"https://you.com/"},
		"Origin":             {"https://you.com"},
		"DNT":                {"1"},
		"Connection":         {"keep-alive"},
		"Sec-Fetch-Dest":     {"empty"},
		"Sec-Fetch-Mode":     {"cors"},
		"Sec-Fetch-Site":     {"same-origin"},
		"sec-ch-ua":          {"\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"Windows\""},
	}

	client := &http.Client{
		Timeout: 300 * time.Second,
	}

	resp, err := client.Do(youReq)
	if err != nil {
		log.Printf("Request error: %v", err)
		http.Error(w, "Failed to send request to You.com API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("You.com response status: %d", resp.StatusCode)

	if resp.StatusCode != 200 {
		log.Printf("You.com API returned non-200 status: %d", resp.StatusCode)
		http.Error(w, "You.com API request failed", http.StatusBadGateway)
		return
	}

	if openAIReq.Stream {
		content := handleStreamResponse(w, resp, originalModel)
		// 如果主方法失败，尝试备用方法
		if content == "" {
			log.Printf("Primary method failed, trying fallback...")
			fallbackContent := tryMultipleMethods(lastMessage, originalModel)
			if fallbackContent != "" {
				sendStreamResponse(w, fallbackContent, originalModel)
			}
		}
	} else {
		content := handleNonStreamResponse(w, resp, originalModel)
		// 如果主方法失败，尝试备用方法
		if content == "" {
			log.Printf("Primary method failed, trying fallback...")
			fallbackContent := tryMultipleMethods(lastMessage, originalModel)
			if fallbackContent == "" {
				fallbackContent = generateFallbackResponse(lastMessage)
			}
			sendNormalResponse(w, fallbackContent, originalModel)
		}
	}
}

func handleStreamResponse(w http.ResponseWriter, resp *http.Response, model string) string {
	var totalContent strings.Builder
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	scanner := bufio.NewScanner(resp.Body)
	responseID := "chatcmpl-" + strconv.FormatInt(time.Now().Unix(), 10)
	created := time.Now().Unix()
	isDebug := os.Getenv("DEBUG") == "true"

	for scanner.Scan() {
		line := scanner.Text()
		if isDebug {
			log.Printf("Raw line: %s", line)
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if isDebug {
				log.Printf("Data content: %s", data)
			}

			if data == "[DONE]" {
				// Send final chunk
				finalChunk := OpenAIStreamResponse{
					ID:      responseID,
					Object:  "chat.completion.chunk",
					Created: created,
					Model:   model,
					Choices: []Choice{{
						Delta:        Delta{Content: ""},
						Index:        0,
						FinishReason: "stop",
					}},
				}
				if jsonData, err := json.Marshal(finalChunk); err == nil {
					fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
				}
				fmt.Fprint(w, "data: [DONE]\n\n")
				break
			} else if data != "" && data != "{}" {
				var youResp map[string]interface{}
				if err := json.Unmarshal([]byte(data), &youResp); err == nil {
					if isDebug {
						log.Printf("Parsed response: %+v", youResp)
					}

					// Try multiple possible field names for the content
					var content string
					if youChatToken, exists := youResp["youChatToken"].(string); exists && youChatToken != "" {
						content = youChatToken
					} else if text, exists := youResp["text"].(string); exists && text != "" {
						content = text
					} else if message, exists := youResp["message"].(string); exists && message != "" {
						content = message
					} else if delta, exists := youResp["delta"].(map[string]interface{}); exists {
						if deltaContent, exists := delta["content"].(string); exists && deltaContent != "" {
							content = deltaContent
						}
					}

					if content != "" {
						totalContent.WriteString(content)
						chunk := OpenAIStreamResponse{
							ID:      responseID,
							Object:  "chat.completion.chunk",
							Created: created,
							Model:   model,
							Choices: []Choice{{
								Delta: Delta{Content: content},
								Index: 0,
							}},
						}
						if jsonData, err := json.Marshal(chunk); err == nil {
							fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
						}
					} else if isDebug {
						log.Printf("No content found in response: %+v", youResp)
					}
				} else if isDebug {
					log.Printf("Failed to parse JSON: %v, data: %s", err, data)
				}
			}
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
	return totalContent.String()
}

func handleNonStreamResponse(w http.ResponseWriter, resp *http.Response, model string) string {
	w.Header().Set("Content-Type", "application/json")

	var fullResponse strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	isDebug := os.Getenv("DEBUG") == "true"
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if isDebug {
			log.Printf("Non-stream line %d: %s", lineCount, line)
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if isDebug {
				log.Printf("Non-stream data: %s", data)
			}

			if data != "[DONE]" && data != "" && data != "{}" {
				var youResp map[string]interface{}
				if err := json.Unmarshal([]byte(data), &youResp); err == nil {
					if isDebug {
						log.Printf("Non-stream parsed: %+v", youResp)
					}

					// Try multiple possible field names for the content
					var content string
					if youChatToken, exists := youResp["youChatToken"].(string); exists && youChatToken != "" {
						content = youChatToken
					} else if text, exists := youResp["text"].(string); exists && text != "" {
						content = text
					} else if message, exists := youResp["message"].(string); exists && message != "" {
						content = message
					} else if delta, exists := youResp["delta"].(map[string]interface{}); exists {
						if deltaContent, exists := delta["content"].(string); exists && deltaContent != "" {
							content = deltaContent
						}
					}

					if content != "" {
						fullResponse.WriteString(content)
					} else if isDebug {
						log.Printf("No content found in non-stream response: %+v", youResp)
					}
				} else if isDebug {
					log.Printf("Failed to parse non-stream JSON: %v, data: %s", err, data)
				}
			}
		}
	}

	finalContent := fullResponse.String()
	if isDebug {
		log.Printf("Final response content length: %d, content: %s", len(finalContent), finalContent)
	}

	// If no content was found, provide a fallback message
	if finalContent == "" {
		finalContent = "I apologize, but I couldn't process your request at this time. Please try again."
		if isDebug {
			log.Printf("Using fallback content")
		}
	}

	openAIResp := OpenAIResponse{
		ID:      "chatcmpl-" + strconv.FormatInt(time.Now().Unix(), 10),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []OpenAIChoice{{
			Message: Message{
				Role:    "assistant",
				Content: finalContent,
			},
			Index:        0,
			FinishReason: "stop",
		}},
	}

	if err := json.NewEncoder(w).Encode(openAIResp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return ""
	}
	return finalContent
}
