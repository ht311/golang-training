# アーキテクチャ概要

## 技術スタック

| 層 | 技術 | 備考 |
|---|---|---|
| バックエンド | Go 1.22+ / Gin | oapi-codegen でサーバースタブ自動生成 |
| フロントエンド | Next.js 16 / TypeScript | App Router, Server Components |
| データベース | PostgreSQL (Neon) | サーバーレス, 無料枠 |
| デプロイ (BE) | Fly.io | 無料枠 3台まで |
| デプロイ (FE) | Vercel | Hobby プラン無料 |
| IaC | Terraform | fly / neon / vercel providers |
| API定義 | TypeSpec → OpenAPI 3.0 | single source of truth |

## コード生成フロー

```
api/typespec/main.tsp
    ↓  tsp compile (make generate-openapi)
api/openapi/openapi.yaml
    ↓  oapi-codegen (make generate-go)
backend/gen/blog_api.gen.go
    ↓  手動実装
backend/internal/handler/post.go
```

## ローカル開発

```bash
# Docker Compose で全サービス起動
make up

# バックエンドのみ
export DATABASE_URL=postgres://blog:blog@localhost:5432/blog
make dev-be

# フロントエンドのみ
make dev-fe
```

## デプロイ手順

1. 各サービスのAPIトークンを環境変数に設定 (`.env.example` 参照)
2. `make tf-init` → `make tf-plan` → `make tf-apply`
3. バックエンドイメージを `flyctl deploy --app <app-name>` でプッシュ
4. Vercel はGitHubプッシュで自動デプロイ
