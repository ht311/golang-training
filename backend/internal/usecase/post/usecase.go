// post パッケージはブログ記事に関するユースケース層。
// ビジネスロジックをここに集約し、handler と infrastructure の中間に位置する。
package post

import (
	"context"

	"github.com/ht311/golang-training/backend/internal/domain/post"
)

// PostUsecase はブログ記事に関するユースケースをまとめた構造体。
// domain.Repository インターフェース経由で DB アクセスするため、実装に非依存。
type PostUsecase struct {
	repo post.Repository
}

// NewPostUsecase はコンストラクタ。
func NewPostUsecase(repo post.Repository) *PostUsecase {
	return &PostUsecase{repo: repo}
}

// List は全記事一覧を返す。
func (uc *PostUsecase) List(ctx context.Context) ([]post.Post, error) {
	return uc.repo.List(ctx)
}

// Create は新しい記事を作成する。
func (uc *PostUsecase) Create(ctx context.Context, title, body string) (*post.Post, error) {
	return uc.repo.Create(ctx, title, body)
}

// GetByID は指定 ID の記事を取得する。
func (uc *PostUsecase) GetByID(ctx context.Context, id string) (*post.Post, error) {
	return uc.repo.GetByID(ctx, id)
}

// Update は記事のタイトル・本文を部分更新する。
func (uc *PostUsecase) Update(ctx context.Context, id string, title, body *string) (*post.Post, error) {
	return uc.repo.Update(ctx, id, title, body)
}

// Delete は指定 ID の記事を削除する。
func (uc *PostUsecase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
