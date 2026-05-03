// main パッケージはサーバーのエントリーポイント。
// 依存関係 (DB → マイグレーション → repository → usecase → handler) を組み立てて Gin を起動する。
package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	gmigrate "github.com/golang-migrate/migrate/v4"
	// postgres マイグレーションドライバを副作用でロード (init() が "postgres" ドライバを登録する)
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ht311/golang-training/backend/gen"
	"github.com/ht311/golang-training/backend/internal/handler"
	infradb "github.com/ht311/golang-training/backend/internal/infrastructure/db"
	infrarepo "github.com/ht311/golang-training/backend/internal/infrastructure/repository"
	"github.com/ht311/golang-training/backend/migrations"
	postuc "github.com/ht311/golang-training/backend/internal/usecase/post"
)

func main() {
	// DATABASE_URL は必須。未設定のままでは起動させない
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// golang-migrate でバージョン管理されたマイグレーションを実行する。
	// migrations.FS はバイナリに埋め込まれた SQL ファイル群 (go:embed)。
	srcDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("create migration source: %v", err)
	}
	m, err := gmigrate.NewWithSourceInstance("iofs", srcDriver, dbURL)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, gmigrate.ErrNoChange) {
		// ErrNoChange は "すでに最新" を意味するので正常系として無視する
		log.Fatalf("run migrations: %v", err)
	}
	log.Println("migrations applied")

	ctx := context.Background()

	// pgxpool.New でコネクションプールを作成する。
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	// 依存を domain → infrastructure → usecase → handler の順に組み立てる
	queries := infradb.New(pool)
	repo := infrarepo.NewPostRepository(queries)
	uc := postuc.NewPostUsecase(repo)
	h := handler.NewPostHandler(uc)

	// gin.Default() はログ出力とパニックリカバリのミドルウェアを内包した Gin ルーターを返す
	r := gin.Default()

	// CORS ミドルウェア: フロントエンド (localhost:3000) からの fetch を許可する。
	// OPTIONS はブラウザが本リクエスト前に送るプリフライトリクエストで、204 で即返す。
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// oapi-codegen が生成した RegisterHandlersWithOptions で全ルートを一括登録する。
	gen.RegisterHandlersWithOptions(r, gen.NewStrictHandler(h, nil), gen.GinServerOptions{})

	// PORT 環境変数が設定されていればそれを使い、なければ 8080 をデフォルトにする
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
