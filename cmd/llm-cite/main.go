// Command llm-cite は、質問群に対する生成AIの回答から、
// 対象（自分のOSS/技術ブログ等）の言及と引用ドメインを計測し、
// 言及判定の品質を precision / recall で評価する開発者向けツール。
//
// 既定はオフラインの replay プロバイダで動く。OPENAI_API_KEY があり
// -provider openai を指定したときだけ実LLMに問い合わせる。
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mabuix/llm-cite/internal/eval"
	"github.com/mabuix/llm-cite/internal/model"
	"github.com/mabuix/llm-cite/internal/provider"
	"github.com/mabuix/llm-cite/internal/report"
)

func main() {
	var (
		queriesPath = flag.String("queries", "testdata/queries.json", "正解ラベル付きクエリ集(JSON)")
		replayPath  = flag.String("replay", "testdata/replay.json", "replay プロバイダの回答(JSON)")
		providerNm  = flag.String("provider", "replay", "プロバイダ: replay | openai")
		targetCSV   = flag.String("target", "iac-guard", "対象のエイリアス(カンマ区切り)")
		modelName   = flag.String("model", "gpt-4o-mini", "openai のモデル名")
		asJSON      = flag.Bool("json", false, "JSON で出力する")
	)
	flag.Parse()

	if err := run(*queriesPath, *replayPath, *providerNm, *targetCSV, *modelName, *asJSON); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(queriesPath, replayPath, providerNm, targetCSV, modelName string, asJSON bool) error {
	queries, err := loadQueries(queriesPath)
	if err != nil {
		return fmt.Errorf("クエリ読み込み: %w", err)
	}

	aliases := splitCSV(targetCSV)
	if len(aliases) == 0 {
		return fmt.Errorf("-target が空です")
	}

	p, err := buildProvider(providerNm, replayPath, modelName)
	if err != nil {
		return err
	}

	// replay は実LLMを叩かない。数字が実計測と誤読されないよう明示する。
	if strings.HasPrefix(p.Name(), "replay") {
		fmt.Fprintln(os.Stderr,
			"注: replay モードです。実LLMには問い合わせず、testdata の想定回答で判定・評価ロジックを自己テストしています。"+
				"実際のGEO露出を計測するには -provider openai（要 OPENAI_API_KEY）を使ってください。")
	}

	ctx := context.Background()
	answers := make([]model.Answer, 0, len(queries))
	for _, q := range queries {
		a, err := p.Answer(ctx, q)
		if err != nil {
			return fmt.Errorf("回答取得(%s): %w", q.ID, err)
		}
		answers = append(answers, a)
	}

	m := eval.Run(queries, answers, aliases)

	if asJSON {
		return report.JSON(os.Stdout, m)
	}
	report.Text(os.Stdout, strings.Join(aliases, "/"), p.Name(), m)
	return nil
}

func buildProvider(name, replayPath, modelName string) (provider.Provider, error) {
	switch name {
	case "replay":
		return provider.NewReplay(replayPath)
	case "openai":
		return provider.NewOpenAI(modelName)
	default:
		return nil, fmt.Errorf("未知のプロバイダ: %s", name)
	}
}

func loadQueries(path string) ([]model.Query, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var qs []model.Query
	if err := json.Unmarshal(b, &qs); err != nil {
		return nil, err
	}
	return qs, nil
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
