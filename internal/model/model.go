// Package model は llm-cite 全体で使う型を定義する。
package model

// Expectation は、あるクエリに対する「期待」（正解ラベル）。
type Expectation struct {
	// ShouldMention は、このクエリの回答で対象（自分のOSS/ブログ等）が言及されるべきかどうか。
	ShouldMention bool `json:"should_mention"`
	// Domains は、引用されることを期待するドメイン（任意）。
	Domains []string `json:"domains,omitempty"`
}

// Query は評価用の1件の質問。
type Query struct {
	ID     string      `json:"id"`
	Text   string      `json:"text"`
	Expect Expectation `json:"expect"`
}

// Answer は、あるプロバイダがクエリに返した回答。
type Answer struct {
	QueryID      string   `json:"query_id"`
	Provider     string   `json:"provider"`
	Text         string   `json:"text"`
	CitedDomains []string `json:"cited_domains,omitempty"`
}

// Metrics は言及判定の評価結果。
type Metrics struct {
	Total     int            `json:"total"`
	TP        int            `json:"tp"`
	FP        int            `json:"fp"`
	FN        int            `json:"fn"`
	TN        int            `json:"tn"`
	Precision float64        `json:"precision"`
	Recall    float64        `json:"recall"`
	F1        float64        `json:"f1"`
	SoV       float64        `json:"sov"` // 言及シェア: 対象が言及された割合
	Domains   map[string]int `json:"domains"`
}
