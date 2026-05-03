# backend

Go + Gin で実装したブログ記事管理 REST API。

## ディレクトリ構成

DDD (ドメイン駆動設計) の4層構造を採用している。依存の向きはドメイン層へ集約する。

```
backend/
├── cmd/server/
│   └── main.go                        # エントリーポイント。依存を組み立てて Gin を起動する
│
├── internal/                          # このモジュール内からしか import できないパッケージ群
│   ├── domain/                        # ドメイン層 (最も内側・外部に依存しない)
│   │   └── post/
│   │       ├── post.go                #   Post エンティティ
│   │       └── repository.go         #   Repository インターフェース + ErrNotFound
│   │
│   ├── usecase/                       # ユースケース層 (ビジネスロジックの置き場)
│   │   └── post/
│   │       └── usecase.go            #   PostUsecase。domain.Repository に依存
│   │
│   ├── handler/                       # インターフェース (アダプター) 層
│   │   └── post.go                   #   HTTP ハンドラ。usecase に依存
│   │
│   └── infrastructure/               # インフラ層 (最も外側・DB などの実装詳細)
│       ├── db/                        #   sqlc が自動生成したコード ← 手で編集しない
│       │   ├── db.go                 #     DBTX インターフェース定義
│       │   ├── models.go             #     DB の行を表す Go 構造体
│       │   └── post.sql.go           #     posts テーブルへの型安全なクエリ関数
│       └── repository/
│           └── post.go               #   domain.Repository の PostgreSQL 実装
│
├── gen/
│   └── blog_api.gen.go               # oapi-codegen が OpenAPI 仕様から自動生成 ← 手で編集しない
│
├── db/query/
│   └── post.sql                      # sqlc に渡す SQL クエリ。ここを編集して sqlc generate を実行する
│
├── migrations/                       # golang-migrate のマイグレーションファイル群
│   ├── 000001_create_posts_table.up.sql    # テーブル作成 (適用)
│   ├── 000001_create_posts_table.down.sql  # テーブル削除 (ロールバック)
│   └── embed.go                      # SQL ファイルをバイナリに埋め込むための go:embed 定義
│
├── sqlc.yaml                         # sqlc の設定 (schema/query/output パスを指定)
├── oapi-codegen.cfg.yaml             # oapi-codegen の設定
├── Dockerfile                        # Fly.io デプロイ用マルチステージビルド
├── go.mod
└── go.sum
```

---

## 処理の流れ

```
HTTP リクエスト
    ↓
Gin ルーター (main.go で登録)
    ↓
handler/post.go        ← OpenAPI 仕様に合わせたリクエスト/レスポンスの変換
    ↓
usecase/post/usecase.go ← ビジネスロジック (現時点は repository への委譲)
    ↓
domain/post/repository.go (interface)
    ↑ 実装
infrastructure/repository/post.go ← domain.Repository の PostgreSQL 実装
    ↓
infrastructure/db/post.sql.go      ← sqlc 生成コード。型安全な SQL を実行
    ↓
PostgreSQL
```

### 依存の向き

```
handler → usecase → domain ← infrastructure/repository
                      ↑
               infrastructure/db
```

インフラ層がドメイン層のインターフェースを実装することで、依存の向きが内側（ドメイン）へ集約される。

---

## コード自動生成

このプロジェクトには「手で書かないファイル」が2種類ある。

### 1. sqlc — SQL から Go コードを生成

`db/query/post.sql` に SQL を書いて以下を実行すると `internal/db/` が更新される。

```bash
sqlc generate
```

**なぜ使うか**: SQL を文字列としてべた書きしないので、タイポや型のズレをコンパイル時に検出できる。

### 2. oapi-codegen — OpenAPI 仕様から Go コードを生成

`api/openapi/openapi.yaml` を変更して以下を実行すると `gen/` が更新される。

```bash
# リポジトリルートで実行
make generate
```

**なぜ使うか**: API の型定義を OpenAPI 仕様と Go コードで二重管理しなくて済む。

---

## マイグレーション

[golang-migrate](https://github.com/golang-migrate/migrate) を使用。Flyway と同じ仕組み。

- `migrations/000001_*.up.sql` — 適用 (テーブル作成など)
- `migrations/000001_*.down.sql` — ロールバック (テーブル削除など)
- サーバー起動時に自動で未適用のマイグレーションだけを実行する
- 適用済みのバージョンは DB の `schema_migrations` テーブルで管理される

新しいマイグレーションを追加するときはファイル名の連番を増やす:

```
migrations/000002_add_author_to_posts.up.sql
migrations/000002_add_author_to_posts.down.sql
```

---

## ローカル起動

```bash
# リポジトリルートで
docker compose up
```

PostgreSQL と API サーバーが起動し、`http://localhost:8080` で動作する。

### 環境変数

| 変数名 | 説明 | 例 |
|---|---|---|
| `DATABASE_URL` | PostgreSQL 接続文字列 | `postgres://user:pass@localhost:5432/blog` |
| `PORT` | リッスンポート (省略時 8080) | `8080` |

---

## API エンドポイント

| メソッド | パス | 説明 |
|---|---|---|
| GET | `/health` | ヘルスチェック |
| GET | `/posts` | 記事一覧 |
| POST | `/posts` | 記事作成 |
| GET | `/posts/:id` | 記事詳細 |
| PUT | `/posts/:id` | 記事更新 |
| DELETE | `/posts/:id` | 記事削除 |

詳細な仕様は `api/openapi/openapi.yaml` を参照。
