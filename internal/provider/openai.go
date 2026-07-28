package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mabuix/llm-cite/internal/model"
)

// OpenAI は OpenAI 互換の Chat Completions API を叩く実プロバイダ。
// OPENAI_API_KEY があるときだけ有効。無ければ NewOpenAI がエラーを返し、
// 呼び出し側は replay にフォールバックする想定。
type OpenAI struct {
	apiKey string
	model  string
	client *http.Client
}

// NewOpenAI は環境変数 OPENAI_API_KEY を読む。未設定ならエラー。
func NewOpenAI(modelName string) (*OpenAI, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}
	return &OpenAI{
		apiKey: key,
		model:  modelName,
		client: &http.Client{Timeout: 60 * time.Second},
	}, nil
}

func (o *OpenAI) Name() string { return "openai:" + o.model }

type chatReq struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResp struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (o *OpenAI) Answer(ctx context.Context, q model.Query) (model.Answer, error) {
	reqBody, err := json.Marshal(chatReq{
		Model: o.model,
		Messages: []chatMessage{
			{Role: "user", Content: q.Text},
		},
	})
	if err != nil {
		return model.Answer{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.openai.com/v1/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return model.Answer{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return model.Answer{}, err
	}
	defer resp.Body.Close()

	var cr chatResp
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return model.Answer{}, err
	}
	if cr.Error != nil {
		return model.Answer{}, fmt.Errorf("openai: %s", cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return model.Answer{}, fmt.Errorf("openai: empty response")
	}

	return model.Answer{
		QueryID:  q.ID,
		Provider: o.Name(),
		Text:     cr.Choices[0].Message.Content,
	}, nil
}
