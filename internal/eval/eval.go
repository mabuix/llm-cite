// Package eval は言及判定の品質を precision / recall で評価する。
package eval

import (
	"github.com/mabuix/llm-cite/internal/detect"
	"github.com/mabuix/llm-cite/internal/model"
)

// Run は、正解ラベル付きのクエリ群と回答群から評価指標を計算する。
//
//   - ShouldMention=true かつ 検出=true  → TP
//   - ShouldMention=false かつ 検出=true → FP（誤検知。同名の別物や否定文脈など）
//   - ShouldMention=true かつ 検出=false → FN（見逃し）
//   - それ以外                           → TN
func Run(queries []model.Query, answers []model.Answer, aliases []string) model.Metrics {
	byID := make(map[string]model.Answer, len(answers))
	for _, a := range answers {
		byID[a.QueryID] = a
	}

	m := model.Metrics{Domains: map[string]int{}}
	mentioned := 0
	for _, q := range queries {
		a := byID[q.ID]
		found := detect.Mention(a.Text, aliases)
		for _, d := range detect.Domains(a.Text, a.CitedDomains) {
			m.Domains[d]++
		}
		if found {
			mentioned++
		}
		switch {
		case q.Expect.ShouldMention && found:
			m.TP++
		case !q.Expect.ShouldMention && found:
			m.FP++
		case q.Expect.ShouldMention && !found:
			m.FN++
		default:
			m.TN++
		}
	}

	m.Total = len(queries)
	m.Precision = ratio(m.TP, m.TP+m.FP)
	m.Recall = ratio(m.TP, m.TP+m.FN)
	if m.Precision+m.Recall > 0 {
		m.F1 = 2 * m.Precision * m.Recall / (m.Precision + m.Recall)
	}
	if m.Total > 0 {
		m.SoV = float64(mentioned) / float64(m.Total)
	}
	return m
}

func ratio(num, den int) float64 {
	if den == 0 {
		return 0
	}
	return float64(num) / float64(den)
}
