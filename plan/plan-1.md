# Go REST API + Next.js フルスタック環境構築計画 (plan-1)

## Context
- **現状**: 空のテンプレートリポジトリ (`/workspaces/golang-training/`)
- **ゴール**: ブログ記事管理APIを題材に、Go+Gin バックエンド / Next.js フロントエンド / Terraform IaC のフルスタック開発環境を構築する
- **API定義フロー**: TypeSpec → OpenAPI 3.0 YAML → oapi-codegen → Go サーバーコード自動生成

---

## 技術スタック

| 層 | 技術 |
|---|---|
| バックエンド | Go 1.22+ / Gin / pgx v5 |
| フロントエンド | Next.js 15 (App Router) / TypeScript |
| データベース | PostgreSQL (Neon serverless 無料枠) |
| API定義 | TypeSpec → OpenAPI 3.0 |
| コード生成 | oapi-codegen (Go server stub) |
| デプロイ (BE) | Fly.io (無料 3台まで) |
| デプロイ (FE) | Vercel (Hobby 無料) |
| IaC | Terraform |

---

## ディレクトリ構成

```
golang-training/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── handler/      # Ginハンドラー実装
│   │   ├── model/        # ドメインモデル
│   │   ├── repository/   # PostgreSQL (pgx v5)
│   │   └── service/      # ビジネスロジック
│   ├── gen/              # oapi-codegen 生成コード (コミット対象)
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── app/          # Next.js App Router
│   │   └── components/
│   ├── package.json
│   └── next.config.ts
├── api/
│   ├── typespec/
│   │   ├── main.tsp      # TypeSpec定義 (source of truth)
│   │   └── tspconfig.yaml
│   └── openapi/
│       └── openapi.yaml  # TypeSpecから生成 (コミット対象)
├── infrastructure/
│   ├── providers.tf      # fly, neon, vercel providers
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── docs/
│   └── architecture.md
└── Makefile              # コマンド集約
```

---

## APIリソース (ブログ記事管理)

```
GET    /posts          記事一覧
POST   /posts          記事作成
GET    /posts/{id}     記事詳細
PUT    /posts/{id}     記事更新
DELETE /posts/{id}     記事削除
```

フィールド: `id` (UUID), `title`, `body`, `createdAt`, `updatedAt`

---

## 実装ステップ

### Step 1: プロジェクト基盤 — ルート構成・Makefile
- [ ] ディレクトリ骨格を作成 (`backend/`, `frontend/`, `api/`, `infrastructure/`, `docs/`)
- [ ] `Makefile` に `generate`, `dev-be`, `dev-fe`, `tf-plan`, `tf-apply` ターゲットを定義

### Step 2: API定義 — TypeSpec → OpenAPI → Go コード生成
- [ ] `api/typespec/` に TypeSpec 環境を初期化 (`npm init` + `@typespec/compiler` + `@typespec/openapi3`)
- [ ] `api/typespec/main.tsp` にブログ記事 CRUD を定義
- [ ] `tsp compile` で `api/openapi/openapi.yaml` を生成
- [ ] `go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest`
- [ ] `oapi-codegen` で `backend/gen/` にサーバーインターフェース・型定義を生成
- [ ] `oapi-codegen.cfg.yaml` で生成設定を固定 (server + types, gin server)

### Step 3: Goバックエンド — モジュール・依存関係
- [ ] `backend/go.mod` を初期化 (`go mod init`)
- [ ] 依存追加: `github.com/gin-gonic/gin`, `github.com/jackc/pgx/v5`, `github.com/oapi-codegen/runtime`
- [ ] `backend/internal/model/post.go` にドメイン型定義
- [ ] `backend/internal/repository/post.go` に PostgreSQL CRUD実装 (pgx v5)
- [ ] `backend/internal/handler/post.go` に生成インターフェースを実装
- [ ] `backend/cmd/server/main.go` でGin + DBコネクション初期化
- [ ] 設定は環境変数から読み込み (`DATABASE_URL` 等)

### Step 4: Dockerfile — Fly.ioデプロイ用
- [ ] `backend/Dockerfile` をマルチステージビルドで作成 (builder → distroless/scratch)
- [ ] ヘルスチェックエンドポイント `GET /health` を追加

### Step 5: Next.js フロントエンド
- [ ] `frontend/` に `create-next-app` でプロジェクト作成 (TypeScript, App Router)
- [ ] OpenAPI生成型を共有 or `openapi-fetch` で型安全なAPIクライアントを設定
- [ ] 記事一覧ページ (`/`) と記事詳細ページ (`/posts/[id]`) を実装

### Step 6: Terraform IaC
- [ ] `infrastructure/providers.tf` に providers 定義

  ```hcl
  terraform {
    required_providers {
      fly    = { source = "fly-apps/fly" }
      neon   = { source = "kislerdm/neon" }
      vercel = { source = "vercel/vercel" }
    }
  }
  ```

- [ ] `infrastructure/main.tf` に以下を定義:
  - `neon_project` + `neon_database` (無料枠)
  - `fly_app` + `fly_machine` (バックエンド)
  - `vercel_project` (フロントエンド)
- [ ] シークレット (DB URL等) は `fly secrets set` 経由 or Terraform の `sensitive = true` 変数で管理

---

## コード生成フロー (Makefile)

```
TypeSpec (main.tsp)
    ↓  tsp compile
OpenAPI (openapi.yaml)
    ↓  oapi-codegen
Go server interface (gen/)
    ↓  手動実装
handler/post.go
```

`make generate` 一発で上記を再生成できるようにする。

---

## テスト方針
- バックエンド: `go test ./...` でハンドラーの単体テスト
- ローカル動作確認: `docker compose up` で Go + PostgreSQL をローカル起動
- E2E確認: `curl` で CRUD 全エンドポイントを叩いて確認

---

## 不明点・確認事項

[Q1] 認証（JWT等）は最初のスコープに含めますか？
[A1] 含めない

[Q2] `docker compose` によるローカル開発環境（DB込み）も作成しますか？
[A2] 作成する → `docker-compose.yml` を Step 4 と合わせて追加
