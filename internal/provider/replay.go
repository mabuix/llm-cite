package provider

import (
	"context"
	"encoding/json"
	"os"

	"github.com/mabuix/llm-cite/internal/model"
)

// fixture は replay 用の1件の回答。
type fixture struct {
	Text         string   `json:"text"`
	CitedDomains []string `json:"cited_domains,omitempty"`
}

// Replay は、あらかじめ用意した回答をクエリIDで引く決定論的プロバイダ。
// ネットワークもAPIキーも要らず、毎回同じ結果になる。テストと学習の既定。
type Replay struct {
	fx map[string]fixture
}

// NewReplay は JSON ファイル（{queryID: {text, cited_domains}}）から Replay を作る。
func NewReplay(path string) (*Replay, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fx map[string]fixture
	if err := json.Unmarshal(b, &fx); err != nil {
		return nil, err
	}
	return &Replay{fx: fx}, nil
}

func (r *Replay) Name() string { return "replay" }

func (r *Replay) Answer(_ context.Context, q model.Query) (model.Answer, error) {
	a := model.Answer{QueryID: q.ID, Provider: r.Name()}
	if f, ok := r.fx[q.ID]; ok {
		a.Text = f.Text
		a.CitedDomains = f.CitedDomains
	}
	return a, nil
}
