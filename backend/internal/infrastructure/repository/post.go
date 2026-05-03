// repository パッケージは domain.post.Repository インターフェースの PostgreSQL 実装。
// SQL は sqlc が生成した型安全なメソッド (infrastructure/db) を呼ぶ。
package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/ht311/golang-training/backend/internal/domain/post"
	infradb "github.com/ht311/golang-training/backend/internal/infrastructure/db"
)

// PostRepository は sqlc が生成した *infradb.Queries を使って posts テーブルを操作する。
type PostRepository struct {
	q *infradb.Queries
}

// NewPostRepository はコンストラクタ。
func NewPostRepository(q *infradb.Queries) *PostRepository {
	return &PostRepository{q: q}
}

// toPost は sqlc 生成の infradb.Post をドメインモデルに変換する。
func toPost(p infradb.Post) *post.Post {
	return &post.Post{
		ID:        p.ID,
		Title:     p.Title,
		Body:      p.Body,
		CreatedAt: p.CreatedAt.Time,
		UpdatedAt: p.UpdatedAt.Time,
	}
}

// List は全記事を作成日時の降順で取得する。
func (r *PostRepository) List(ctx context.Context) ([]post.Post, error) {
	rows, err := r.q.ListPosts(ctx)
	if err != nil {
		return nil, err
	}
	posts := make([]post.Post, len(rows))
	for i, row := range rows {
		posts[i] = *toPost(row)
	}
	return posts, nil
}

// Create は新しい記事を INSERT する。ID は UUID v4 をサーバー側で発行する。
func (r *PostRepository) Create(ctx context.Context, title, body string) (*post.Post, error) {
	row, err := r.q.CreatePost(ctx, infradb.CreatePostParams{
		ID:    uuid.New().String(),
		Title: title,
		Body:  body,
	})
	if err != nil {
		return nil, err
	}
	return toPost(row), nil
}

// GetByID は指定 ID の記事を1件取得する。
func (r *PostRepository) GetByID(ctx context.Context, id string) (*post.Post, error) {
	row, err := r.q.GetPost(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, post.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return toPost(row), nil
}

// Update は記事のタイトル・本文を部分更新する。nil のフィールドは現在値を引き継ぐ。
func (r *PostRepository) Update(ctx context.Context, id string, title, body *string) (*post.Post, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	newTitle := current.Title
	newBody := current.Body
	if title != nil {
		newTitle = *title
	}
	if body != nil {
		newBody = *body
	}
	row, err := r.q.UpdatePost(ctx, infradb.UpdatePostParams{
		ID:    id,
		Title: newTitle,
		Body:  newBody,
	})
	if err != nil {
		return nil, err
	}
	return toPost(row), nil
}

// Delete は指定 ID の記事を削除する。削除行数が 0 なら ErrNotFound を返す。
func (r *PostRepository) Delete(ctx context.Context, id string) error {
	affected, err := r.q.DeletePost(ctx, id)
	if err != nil {
		return err
	}
	if affected == 0 {
		return post.ErrNotFound
	}
	return nil
}

// コンパイル時に PostRepository が domain.Repository を完全に実装しているか検証する。
var _ post.Repository = (*PostRepository)(nil)
