# バックエンドのDDDパッケージ構成への移行計画 (plan-2)

## Context

- **現状**: `internal/` 以下が `model/`, `repository/`, `handler/`, `db/` のフラットな層構造
- **問題**:
  1. `handler.PostHandler` が `*repository.PostRepository`（具象型）に直接依存 → テスト時にモックを差し込めない（DIP違反）
  2. ユースケース層がない → ビジネスロジックが handler または repository に混入するリスク
  3. リポジトリのインターフェースが存在しない
- **ゴール**: DDD の4層（domain / usecase / handler / infrastructure）に再編し、依存の向きをドメイン層へ集約する

---

## 設計方針

- `interface` は Go のキーワードなので `internal/interface/` は使えない → handler はそのまま `internal/handler/` に残す
- usecase はこの規模ではインターフェース化しない（具象型依存で十分）
- sqlc 生成コードは `internal/infrastructure/db/` に移動し、sqlc.yaml を更新する

### 依存の向き（変更後）

```
handler → usecase → domain（interface）← infrastructure/repository（実装）
                              ↑
                        infrastructure/db（sqlc生成）
```

---

## アーキテクチャ / ディレクトリ構成

```
backend/internal/
├── domain/
│   └── post/
│       ├── post.go         # Post エンティティ（model/post.go を移動）
│       └── repository.go   # PostRepository インターフェース（新規）
├── usecase/
│   └── post/
│       └── usecase.go      # PostUsecase（新規）
├── handler/
│   └── post.go             # HTTPハンドラ（usecase 依存に変更）
└── infrastructure/
    ├── db/                 # sqlc 生成コード（internal/db/ から移動）
    │   ├── db.go
    │   ├── models.go
    │   └── post.sql.go
    └── repository/
        └── post.go         # domain.PostRepository 実装（repository/post.go から移動）
```

---

## 実装ステップ

### Step 1: `internal/domain/post/post.go` を作成
- [ ] `internal/model/post.go` の `Post` 構造体を `package post` として移植
- インポートパス: `github.com/ht311/golang-training/backend/internal/domain/post`

### Step 2: `internal/domain/post/repository.go` を作成（新規）
- [ ] `PostRepository` インターフェースを定義
  ```go
  type Repository interface {
      List(ctx context.Context) ([]Post, error)
      Create(ctx context.Context, title, body string) (*Post, error)
      GetByID(ctx context.Context, id string) (*Post, error)
      Update(ctx context.Context, id string, title, body *string) (*Post, error)
      Delete(ctx context.Context, id string) error
  }
  ```
- [ ] `ErrNotFound` もここに移動（現在は `internal/repository/post.go` に定義）

### Step 3: `internal/infrastructure/db/` に sqlc コードを移動
- [ ] `internal/db/db.go`, `models.go`, `post.sql.go` を `internal/infrastructure/db/` にコピー
- [ ] `sqlc.yaml` の `queries`/`out` パスを `internal/infrastructure/db` に更新
- [ ] 旧 `internal/db/` を削除

### Step 4: `internal/infrastructure/repository/post.go` を作成
- [ ] `internal/repository/post.go` の実装を `package repository` として移植
- [ ] `*db.Queries` は `internal/infrastructure/db` の型を使う
- [ ] `domain.PostRepository` インターフェースを実装（コンパイル時検証: `var _ domain.Repository = (*PostRepository)(nil)`）
- [ ] `ErrNotFound` の参照を `domain.ErrNotFound` に変更
- [ ] 旧 `internal/repository/` を削除

### Step 5: `internal/usecase/post/usecase.go` を作成（新規）
- [ ] `PostUsecase` struct（`domain.Repository` インターフェースに依存）
- [ ] List, Create, GetByID, Update, Delete メソッドを実装
  - 現時点ではビジネスロジックなし → repository への委譲のみ
  - 将来のビジネスロジック追加箇所として機能する層

### Step 6: `internal/handler/post.go` を更新
- [ ] `*repository.PostRepository` への直接依存を `*usecase.PostUsecase` に変更
- [ ] `toGenPost` ヘルパーはそのまま維持
- [ ] `ErrNotFound` の参照を `domain.ErrNotFound` に変更

### Step 7: `cmd/server/main.go` を更新
- [ ] import パスを更新（`infrastructure/db`, `infrastructure/repository`, `usecase/post`）
- [ ] 依存組み立て順:
  ```go
  queries := infradb.New(pool)
  repo    := infrarepo.NewPostRepository(queries)
  uc      := postuc.NewPostUsecase(repo)
  h       := handler.NewPostHandler(uc)
  ```

### Step 8: 旧ファイルを削除
- [ ] `internal/model/post.go`（domain/post/post.go に移行済み）
- [ ] `internal/repository/post.go`（infrastructure/repository に移行済み）

### Step 9: `backend/README.md` のディレクトリ構成を更新

---

## テスト方針

```bash
# コンパイルと静的解析
cd backend && go build ./...
go vet ./...

# サーバー起動確認
docker compose up
curl http://localhost:8080/health
curl http://localhost:8080/posts
```

---

## 不明点・確認事項

なし（ユースケース層のインターフェース化は今回スコープ外とし、次回の拡張時に検討）
