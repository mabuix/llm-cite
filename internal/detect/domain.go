package detect

import (
	"regexp"
	"strings"
)

// urlRe は回答テキスト中のURLからホスト部分を抜き出す。
var urlRe = regexp.MustCompile(`https?://([^/\s)>\]]+)`)

// normalizeDomain はホスト名を正規化する（小文字化、先頭 www. とポートを除去）。
func normalizeDomain(h string) string {
	h = strings.ToLower(strings.TrimSpace(h))
	if i := strings.IndexByte(h, ':'); i >= 0 {
		h = h[:i]
	}
	h = strings.TrimPrefix(h, "www.")
	return h
}

// Domains は、プロバイダが明示した引用ドメインと、回答テキスト中のURLから
// 抽出したドメインをまとめ、重複を除いて返す（出現順を保つ）。
func Domains(text string, provided []string) []string {
	seen := map[string]bool{}
	var out []string
	add := func(d string) {
		d = normalizeDomain(d)
		if d == "" || seen[d] {
			return
		}
		seen[d] = true
		out = append(out, d)
	}
	for _, d := range provided {
		add(d)
	}
	for _, m := range urlRe.FindAllStringSubmatch(text, -1) {
		add(m[1])
	}
	return out
}
