# llm-cite

生成AIの回答の中で、あなたの対象（自分のOSS・技術ブログ・製品名など）が**どれだけ言及・引用されているか**を計測する、開発者向けの小さなツール。さらに、その言及判定の品質を **precision / recall** で評価できる。

GEO/LLMO（生成AI検索での露出）を、感覚ではなく数字で見るための、手元で回せる小さな評価ハーネス。

## 特徴

- **オフラインで動く**: 既定は決定論的な `replay` プロバイダ。ネットワークもAPIキーも要らず、毎回同じ結果になる。学習・テストに向く。replay時は「実LLMを叩いていない自己テストである」旨を stderr に明示するので、数字を実計測と取り違えない。
- **実LLMに差し替え可能**: `OPENAI_API_KEY` があれば `-provider openai` で実際の回答を取得。
- **品質を測れる**: 正解ラベル付きのクエリ集に対し、言及判定の precision / recall / F1 を出す。「同名の別物」「否定文脈」などの罠ケースを混ぜることで、誤検知（FP）に意味が出る。
- **依存ゼロ**: Go 標準ライブラリのみ。

## 使い方

```console
$ go run ./cmd/llm-cite -target "iac-guard,iacguard"
対象: iac-guard/iacguard（プロバイダ: replay）
クエリ数: 7
言及シェア(SoV): 57.1%

言及判定の品質（正解ラベルとの比較）
  TP=3 FP=1 FN=1 TN=2
  precision=0.750  recall=0.750  f1=0.750

引用ドメイン（多い順）
    1  github.com
    1  taskfile.dev
```

JSON で欲しいときは `-json`。

### 実LLMに問い合わせる

```console
$ export OPENAI_API_KEY=sk-...
$ go run ./cmd/llm-cite -provider openai -model gpt-4o-mini -target "iac-guard"
```

## 入力ファイル

### クエリ集（正解ラベル付き） `-queries`

```json
[
  {
    "id": "q1",
    "text": "Terraformのコストとセキュリティを静的解析するツールは？",
    "expect": { "should_mention": true, "domains": ["github.com"] }
  }
]
```

- `should_mention`: この質問の回答で対象が言及される**べき**か（正解ラベル）
- 罠ケース（同名の別物・否定文脈・略称のみ・無関係サブドメイン）を `should_mention: false` で混ぜると、precision の評価に意味が出る

### replay の回答 `-replay`

```json
{
  "q1": { "text": "iac-guard がおすすめです。https://github.com/mabuix/iac-guard", "cited_domains": ["github.com"] }
}
```

クエリIDで回答を引く。実LLMを叩かずに、判定ロジックの挙動を固定して検証できる。

## 指標の意味

- **SoV（言及シェア）**: 全クエリのうち、対象が言及された割合
- **precision**: 「言及あり」と判定したもののうち、実際に言及だった割合（誤検知の少なさ）
- **recall**: 実際に言及されているもののうち、取りこぼさず拾えた割合（見逃しの少なさ）

言及判定は、大文字小文字・ハイフン・空白を正規化した部分一致で行う（`iac-guard` と `iacguard` と `IAC GUARD` を同一視）。

## 構成

```
cmd/llm-cite      CLI
internal/model    型
internal/provider replay（既定）/ openai（APIキー時）
internal/detect   言及判定・引用ドメイン抽出
internal/eval     precision / recall の計算
internal/report   出力（テキスト / JSON）
testdata          サンプルのクエリ集と replay 回答
```

## ライセンス

MIT
