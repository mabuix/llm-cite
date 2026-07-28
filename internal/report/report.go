// Package report は評価結果を人間向け／JSON で出力する。
package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/mabuix/llm-cite/internal/model"
)

// Text は指標を読みやすいテキストで書き出す。
func Text(w io.Writer, target string, provider string, m model.Metrics) {
	fmt.Fprintf(w, "対象: %s（プロバイダ: %s）\n", target, provider)
	fmt.Fprintf(w, "クエリ数: %d\n", m.Total)
	fmt.Fprintf(w, "言及シェア(SoV): %.1f%%\n", m.SoV*100)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "言及判定の品質（正解ラベルとの比較）")
	fmt.Fprintf(w, "  TP=%d FP=%d FN=%d TN=%d\n", m.TP, m.FP, m.FN, m.TN)
	fmt.Fprintf(w, "  precision=%.3f  recall=%.3f  f1=%.3f\n", m.Precision, m.Recall, m.F1)

	if len(m.Domains) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "引用ドメイン（多い順）")
		for _, dc := range sortDomains(m.Domains) {
			fmt.Fprintf(w, "  %3d  %s\n", dc.count, dc.domain)
		}
	}
}

// JSON は指標を JSON で書き出す。
func JSON(w io.Writer, m model.Metrics) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}

type domainCount struct {
	domain string
	count  int
}

// sortDomains は件数の降順、同数はドメイン名の昇順で安定に並べる。
func sortDomains(d map[string]int) []domainCount {
	out := make([]domainCount, 0, len(d))
	for k, v := range d {
		out = append(out, domainCount{domain: k, count: v})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].count != out[j].count {
			return out[i].count > out[j].count
		}
		return out[i].domain < out[j].domain
	})
	return out
}
