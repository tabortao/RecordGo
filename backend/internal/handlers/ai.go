package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"recordgo/internal/common"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AITaskParseReq AI 识别请求
type AITaskParseReq struct {
	Text       string                `form:"text"`
	Image      *multipart.FileHeader `form:"image"`
	AIProvider string                `form:"ai_provider"` // optional
	AIBaseURL  string                `form:"ai_base_url"`
	AIKey      string                `form:"ai_key"`
	AIModel    string                `form:"ai_model"`
	Categories string                `form:"categories"` // Comma separated list of available categories
}

// AITaskItem AI 识别出的任务项
type AITaskItem struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Score       int     `json:"score"`
	PlanMinutes int     `json:"plan_minutes"`
	StartDate   string  `json:"start_date"` // YYYY-MM-DD
	EndDate     *string `json:"end_date"`   // YYYY-MM-DD
	RepeatType  string  `json:"repeat_type"`
	WeeklyDays  []int   `json:"weekly_days"`
	Confidence  string  `json:"confidence"` // High, Medium, Low
}

// ParseTaskByAI 处理 AI 识别请求
func ParseTaskByAI(c *gin.Context) {
	if !hasPermission(c, "tasks", "create") {
		deny(c, "无权限创建任务")
		return
	}

	var req AITaskParseReq
	if err := c.ShouldBind(&req); err != nil {
		common.Error(c, 40001, "参数解析失败")
		return
	}

	// 如果没有 AI 配置，返回错误 (或者回退到简单的 Mock 逻辑，但这里我们强制要求 AI 配置)
	if req.AIBaseURL == "" || req.AIKey == "" {
		// 尝试回退到 Mock 逻辑? 不，直接报错提示用户去设置
		common.Error(c, 40002, "未配置 AI 模型，请在设置页配置 AI")
		return
	}

	var imageBase64 string
	if req.Image != nil {
		file, err := req.Image.Open()
		if err != nil {
			common.Error(c, 40003, "图片读取失败")
			return
		}
		defer file.Close()
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, file); err != nil {
			common.Error(c, 40003, "图片读取失败")
			return
		}
		imageBase64 = base64.StdEncoding.EncodeToString(buf.Bytes())
	}

	tasks, err := callLLM(c, req, imageBase64)
	if err != nil {
		zap.L().Error("AI Call Failed", zap.Error(err))
		common.Error(c, 50001, "AI 服务调用失败: "+err.Error())
		return
	}

	// Rule Engine: Double-check scores
	for i := range tasks {
		tasks[i].Confidence = "Medium" // Default confidence
		refinedScore := refineScoreWithRules(tasks[i].Description)
		if refinedScore > 0 {
			tasks[i].Score = refinedScore
			tasks[i].Confidence = "High" // Rule match increases confidence
		}
	}

	common.Ok(c, gin.H{"tasks": tasks})
}

// Helper to refine score from text using Regex
func refineScoreWithRules(text string) int {
	if text == "" {
		return 0
	}
	// Match pattern: (number or chinese_number) + (space) + (coin/score keywords)
	// Supports: "20金币", "20 金币", "二十金币", "10积分"
	re := regexp.MustCompile(`([0-9]+|[零一二三四五六七八九十百]+)\s*(?:金币|积分|points|coins|进步数量)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 2 {
		numStr := matches[1]
		return parseNumber(numStr)
	}
	return 0
}

func parseNumber(s string) int {
	// Try Arabic number first
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	// Try Chinese number
	return parseChineseNumber(s)
}

func parseChineseNumber(s string) int {
	cnMap := map[rune]int{
		'零': 0, '一': 1, '二': 2, '三': 3, '四': 4,
		'五': 5, '六': 6, '七': 7, '八': 8, '九': 9,
		'十': 10, '百': 100,
	}

	val := 0
	temp := 0
	for _, r := range s {
		if v, ok := cnMap[r]; ok {
			if v == 100 {
				if temp == 0 {
					temp = 1
				}
				val += temp * 100
				temp = 0
			} else if v == 10 {
				if temp == 0 {
					temp = 1
				}
				val += temp * 10
				temp = 0
			} else {
				temp = v
			}
		}
	}
	val += temp
	return val
}

// LLMRequest OpenAI Compatible Request
type LLMRequest struct {
	Model    string       `json:"model"`
	Messages []LLMMessage `json:"messages"`
	Stream   bool         `json:"stream"`
}

type LLMMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // string or []ContentPart
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type LLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func callLLM(ctx *gin.Context, req AITaskParseReq, imageBase64 string) ([]AITaskItem, error) {
	nowStr := time.Now().Format("2006-01-02 15:04:05 Monday")

	categoryPrompt := ""
	if req.Categories != "" {
		categoryPrompt = fmt.Sprintf("5. Category Constraint: The 'category' MUST be one of the following: [%s]. If no match, choose the most relevant one or the first one.", req.Categories)
	} else {
		categoryPrompt = "5. Category: Choose a reasonable category (e.g. Learning, Sport, Life, Work)."
	}

	systemPrompt := fmt.Sprintf(`You are a smart task extraction assistant for a Task Management App.
Current Time: %s
Your goal is to extract structured tasks from the user's input (text and/or image).

Output Format: Return ONLY a JSON array of objects. Do not include Markdown formatting like `+"```json"+`.
JSON Structure:
[
  {
    "name": "Concise Title (Verb + Object, NO time/score/duration)",
    "description": "Full original details",
    "category": "Must be one of the provided categories",
    "score": 10,
    "plan_minutes": 30,
    "start_date": "YYYY-MM-DD",
    "end_date": "YYYY-MM-DD" or null,
    "repeat_type": "none" | "daily" | "weekly" | "monthly" | "weekdays",
    "weekly_days": [1, 2] (1=Monday, 7=Sunday)
  }
]

Rules:
1. Title Summarization: The "name" MUST be concise. Remove any mention of time, duration, coins, points, or scores. 
   - Example: "每天早上7点背单词30分钟 10积分" -> name: "背单词", description: "每天早上7点背单词30分钟 10积分"
   - Example: "本周末下午跑步5公里 1小时 20金币" -> name: "跑步", description: "本周末下午跑步5公里 1小时 20金币"
2. Time Recognition: 
   - Identify relative dates (e.g. "tomorrow", "next Monday") based on "Current Time".
   - "start_date" is the date of the FIRST occurrence.
   - "This weekend" ("本周末"): usually means the upcoming Saturday or Sunday. If today is Monday-Friday, it's the next Saturday.
3. Repetition & Dates:
   - "Every day" -> repeat_type: "daily"
   - "Every week" -> repeat_type: "weekly"
   - "Weekdays" -> repeat_type: "weekdays"
   - "Tomorrow 7am" (one time) -> repeat_type: "none", start_date="YYYY-MM-DD" (tomorrow), end_date=null
   - "Every day for 1 month" -> repeat_type: "daily", start_date="tomorrow" (or specified), end_date="1 month later"
   - If it's a repetitive task with a specific end date (e.g., "until next month"), set "end_date". Otherwise, "end_date" is null.
4. Score/Coins: 
   - STRICTLY extract the number immediately preceding "金币", "积分", "points", "coins", "金币数量".
   - Ignore numbers related to distance (km, m), duration (minutes, hours), or time unless they are explicitly labeled as score.
   - Example: "20金币" -> score: 20.
   - Example: "10积分" -> score: 10.
   - Example: "跑步5公里 1小时 20金币" -> score: 20 (NOT 5, NOT 1).
   - If NO score keywords are found, ONLY THEN use default score: 1.
%s
6. Language: Output tasks in the same language as input (Chinese/English).
`, nowStr, categoryPrompt)

	userContent := []ContentPart{}
	if req.Text != "" {
		userContent = append(userContent, ContentPart{Type: "text", Text: req.Text})
	} else if imageBase64 == "" {
		userContent = append(userContent, ContentPart{Type: "text", Text: "Extract tasks from this."})
	}

	if imageBase64 != "" {
		userContent = append(userContent, ContentPart{
			Type: "image_url",
			ImageURL: &ImageURL{
				URL: "data:image/jpeg;base64," + imageBase64,
			},
		})
	}

	llmReq := LLMRequest{
		Model: req.AIModel,
		Messages: []LLMMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userContent},
		},
		Stream: false,
	}

	reqBody, _ := json.Marshal(llmReq)

	// Create Request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", strings.TrimRight(req.AIBaseURL, "/")+"/chat/completions", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+req.AIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API Error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var llmResp LLMResponse
	if err := json.Unmarshal(bodyBytes, &llmResp); err != nil {
		return nil, fmt.Errorf("JSON Parse Error: %v", err)
	}

	if llmResp.Error != nil {
		return nil, fmt.Errorf("API Error: %s", llmResp.Error.Message)
	}

	if len(llmResp.Choices) == 0 {
		return nil, fmt.Errorf("No response choices")
	}

	content := llmResp.Choices[0].Message.Content
	// Clean markdown code blocks if present
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var tasks []AITaskItem
	if err := json.Unmarshal([]byte(content), &tasks); err != nil {
		// Try to find JSON array in the text if simple unmarshal fails
		start := strings.Index(content, "[")
		end := strings.LastIndex(content, "]")
		if start != -1 && end != -1 && end > start {
			jsonPart := content[start : end+1]
			if err2 := json.Unmarshal([]byte(jsonPart), &tasks); err2 != nil {
				return nil, fmt.Errorf("Result Parse Error: %v, Content: %s", err, content)
			}
		} else {
			return nil, fmt.Errorf("Result Parse Error: %v, Content: %s", err, content)
		}
	}

	return tasks, nil
}
