package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ddg-community/go-ddg-api"
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
	"gpt-3.5-turbo":      "gpt_3.5",
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

	var openAIReq OpenAIRequest
	if err := json.NewDecoder(r.Body).Decode(&openAIReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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

	youReq, _ := http.NewRequest("GET", "https://you.com/api/streamingSearch", nil)

	q := youReq.URL.Query()
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
	youReq.URL.RawQuery = q.Encode()

	youReq.Header = http.Header{
		"sec-ch-ua-platform":         {"Windows"},
		"Cache-Control":              {"no-cache"},
		"sec-ch-ua":                  {`"Not(A:Brand";v="99", "Microsoft Edge";v="133", "Chromium";v="133"`},
		"sec-ch-ua-bitness":          {"64"},
		"sec-ch-ua-model":            {"