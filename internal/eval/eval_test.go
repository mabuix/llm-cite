package eval

import (
	"math"
	"testing"

	"github.com/mabuix/llm-cite/internal/model"
)

func TestRun(t *testing.T) {
	queries := []model.Query{
		{ID: "a", Expect: model.Expectation{ShouldMention: true}},  // 検出true → TP
		{ID: "b", Expect: model.Expectation{ShouldMention: true}},  // 検出false → FN
		{ID: "c", Expect: model.Expectation{ShouldMention: false}}, // 検出true → FP（罠: 否定文脈）
		{ID: "d", Expect: model.Expectation{ShouldMention: false}}, // 検出false → TN
	}
	answers := []model.Answer{
		{QueryID: "a", Text: "iac-guard がおすすめ"},
		{QueryID: "b", Text: "そういうツールは知りません"},
		{QueryID: "c", Text: "iac-guardという別物もありますが扱いません"},
		{QueryID: "d", Text: "gp3とgp2の違いの話"},
	}
	m := Run(queries, answers, []string{"iac-guard"})

	if m.TP != 1 || m.FP != 1 || m.FN != 1 || m.TN != 1 {
		t.Fatalf("TP/FP/FN/TN = %d/%d/%d/%d, want 1/1/1/1", m.TP, m.FP, m.FN, m.TN)
	}
	if !almost(m.Precision, 0.5) || !almost(m.Recall, 0.5) {
		t.Errorf("precision=%.3f recall=%.3f, want 0.5/0.5", m.Precision, m.Recall)
	}
	// 検出true は a,c の2件 / 全4件 = 0.5
	if !almost(m.SoV, 0.5) {
		t.Errorf("SoV=%.3f, want 0.5", m.SoV)
	}
}

func almost(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
