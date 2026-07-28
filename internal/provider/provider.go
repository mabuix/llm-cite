// Package provider は「クエリを投げて回答を得る」プロバイダを抽象化する。
// 既定はオフラインで決定論的な replay。APIキーがあるときだけ実LLMに差し替える。
package provider

import (
	"context"

	"github.com/mabuix/llm-cite/internal/model"
)

// Provider はクエリに対して回答を返す。
type Provider interface {
	Name() string
	Answer(ctx context.Context, q model.Query) (model.Answer, error)
}
