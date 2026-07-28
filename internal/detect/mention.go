// Package detect は、回答テキストからの言及判定と引用ドメイン抽出を行う。
package detect

import "strings"

// normalize は比較用にテキストを正規化する。
// 小文字化し、ハイフン・空白（半角/全角）を除去する。
// これにより "iac-guard" と "iacguard" と "IAC GUARD" を同一視できる。
func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "　", "")
	return s
}

// Mention は、回答テキストに対象のエイリアスのいずれかが含まれるかを返す。
// 正規化した上での部分一致で判定する。
func Mention(answer string, aliases []string) bool {
	na := normalize(answer)
	for _, a := range aliases {
		na2 := normalize(a)
		if na2 == "" {
			continue
		}
		if strings.Contains(na, na2) {
			return true
		}
	}
	return false
}
