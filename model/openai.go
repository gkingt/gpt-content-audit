package model

import (
	"gpt-content-audit/common/config"
	"strings"
)

type OpenAIChatCompletionRequest struct {
	Model    string              `json:"model"`
	Stream   bool                `json:"stream"`
	Messages []OpenAIChatMessage `json:"messages"`
	OpenAIChatCompletionExtraRequest
}

type OpenAIChatCompletionExtraRequest struct {
	ChannelId *string `json:"channelId"`
}

type OpenAIChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type OpenAIErrorResponse struct {
	OpenAIError OpenAIError `json:"error"`
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type OpenAIChatCompletionResponse struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	Choices           []OpenAIChoice `json:"choices"`
	Usage             OpenAIUsage    `json:"usage"`
	SystemFingerprint *string        `json:"system_fingerprint"`
	Suggestions       []string       `json:"suggestions"`
}

type OpenAIChoice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	LogProbs     *string       `json:"logprobs"`
	FinishReason string        `json:"finish_reason"`
	Delta        OpenAIDelta   `json:"delta"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAIDelta struct {
	Content string `json:"content"`
}

type OpenAIImagesGenerationRequest struct {
	OpenAIChatCompletionExtraRequest
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	ResponseFormat string `json:"response_format"`
}

type OpenAIImagesGenerationResponse struct {
	Created     int64                                 `json:"created"`
	DailyLimit  bool                                  `json:"dailyLimit"`
	Data        []*OpenAIImagesGenerationDataResponse `json:"data"`
	Suggestions []string                              `json:"suggestions"`
}

type OpenAIImagesGenerationDataResponse struct {
	URL           string `json:"url"`
	RevisedPrompt string `json:"revised_prompt"`
	B64Json       string `json:"b64_json"`
}

type OpenAIGPT4VImagesReq struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	ImageURL struct {
		URL string `json:"url"`
	} `json:"image_url"`
}

type GetUserContent interface {
	GetUserContent() []string
}

func (r OpenAIChatCompletionRequest) GetUserContent() []string {
	var userContent []string
	if config.AllDialogRecordEnable == 1 {
		foundUser := false
		foundAssistantOrSystem := false
		for i := len(r.Messages) - 1; i >= 0; i-- {
			msg := r.Messages[i]
			if !foundUser && msg.Role == "user" {
				switch contentObj := msg.Content.(type) {
				case string:
					userContent = append(userContent, contentObj)
					foundUser = true
				}
			} else if !foundAssistantOrSystem && (msg.Role == "assistant" || msg.Role == "system") {
				switch contentObj := msg.Content.(type) {
				case string:
					userContent = append(userContent, contentObj)
					foundAssistantOrSystem = true
				}
			}
			if foundUser && foundAssistantOrSystem {
				break
			}
		}
	} else {
		for i := len(r.Messages) - 1; i >= 0; i-- {
			if r.Messages[i].Role == "user" {
				switch contentObj := r.Messages[i].Content.(type) {
				case string:
					userContent = append(userContent, contentObj)
				}
				break
			}
		}
	}
	// 在函数返回前处理 userContent
	if len(userContent) > 0 {
		// 使用换行符连接所有内容
		joinedContent := strings.Join(userContent, "\n")
		// 返回包含单个拼接字符串的切片
		return []string{joinedContent}
	}
	return userContent
}

func (r OpenAIImagesGenerationRequest) GetUserContent() []string {
	return []string{r.Prompt}
}

type OpenAIModerationRequest struct {
	Input string `json:"input"`
}

type OpenAIModerationResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Results []struct {
		Flagged        bool               `json:"flagged"`
		Categories     map[string]bool    `json:"categories"`
		CategoryScores map[string]float64 `json:"category_scores"`
	} `json:"results"`
}
