package detect

import (
	"reflect"
	"testing"
)

func TestMention(t *testing.T) {
	aliases := []string{"iac-guard"}
	cases := []struct {
		name   string
		answer string
		want   bool
	}{
		{"ハイフンあり", "iac-guard がおすすめ", true},
		{"ハイフンなし表記揺れ", "iacguard というツール", true},
		{"大文字", "IAC-GUARD is a tool", true},
		{"全角空白入り", "iac　guard を使う", true},
		{"言及なし", "go-task が定番です", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Mention(c.answer, aliases); got != c.want {
				t.Errorf("Mention(%q) = %v, want %v", c.answer, got, c.want)
			}
		})
	}
}

func TestDomains(t *testing.T) {
	text := "詳しくは https://github.com/mabuix/iac-guard と https://www.example.com:8080/x を参照。"
	got := Domains(text, []string{"Extra.com"})
	want := []string{"extra.com", "github.com", "example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Domains() = %v, want %v", got, want)
	}
}
