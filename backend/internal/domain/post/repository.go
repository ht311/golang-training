package post

import (
	"context"
	"errors"
)

// ErrNotFound は指定した ID のレコードが存在しない場合に返すエラー。
var ErrNotFound = errors.New("not found")

// Repository は posts テーブルに対するデータアクセスの契約。
// インフラ層の実装はこのインターフェースを満たす必要がある。
type Repository interface {
	List(ctx context.Context) ([]Post, error)
	Create(ctx context.Context, title, body string) (*Post, error)
	GetByID(ctx context.Context, id string) (*Post, error)
	Update(ctx context.Context, id string, title, body *string) (*Post, error)
	Delete(ctx context.Context, id string) error
}
