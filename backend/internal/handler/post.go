// handler パッケージは HTTP リクエストを受け取り、usecase を呼び出してレスポンスを返す層。
// oapi-codegen が生成した StrictServerInterface を実装することで、
// OpenAPI 仕様と実装の乖離をコンパイル時に検出できる。
package handler

import (
	"context"
	"errors"

	"github.com/ht311/golang-training/backend/gen"
	"github.com/ht311/golang-training/backend/internal/domain/post"
	postuc "github.com/ht311/golang-training/backend/internal/usecase/post"
)

// PostHandler は全エンドポイントのハンドラをまとめた構造体。
// uc フィールド経由でユースケースを呼び出し、HTTPの詳細(ステータスコード等)は gen 側に任せる。
type PostHandler struct {
	uc *postuc.PostUsecase
}

// NewPostHandler はコンストラクタ。
func NewPostHandler(uc *postuc.PostUsecase) *PostHandler {
	return &PostHandler{uc: uc}
}

// toGenPost はドメインモデル (post.Post) を OpenAPI 生成型 (gen.Post) に変換するヘルパー。
func toGenPost(p *post.Post) gen.Post {
	return gen.Post{
		Id:        p.ID,
		Title:     p.Title,
		Body:      p.Body,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// HealthCheck は GET /health のハンドラ。
func (h *PostHandler) HealthCheck(_ context.Context, _ gen.HealthCheckRequestObject) (gen.HealthCheckResponseObject, error) {
	return gen.HealthCheck200JSONResponse{Status: "ok"}, nil
}

// PostsList は GET /posts のハンドラ。全記事一覧を返す。
func (h *PostHandler) PostsList(ctx context.Context, _ gen.PostsListRequestObject) (gen.PostsListResponseObject, error) {
	posts, err := h.uc.List(ctx)
	if err != nil {
		var resp gen.PostsList200JSONResponse
		_ = resp.FromErrorResponse(gen.ErrorResponse{Message: err.Error()})
		return resp, nil
	}
	genPosts := make([]gen.Post, len(posts))
	for i, p := range posts {
		genPosts[i] = toGenPost(&p)
	}
	var resp gen.PostsList200JSONResponse
	_ = resp.FromPostList(gen.PostList{Posts: genPosts, Total: int32(len(genPosts))})
	return resp, nil
}

// PostsCreate は POST /posts のハンドラ。新しい記事を作成し 201 を返す。
func (h *PostHandler) PostsCreate(ctx context.Context, req gen.PostsCreateRequestObject) (gen.PostsCreateResponseObject, error) {
	p, err := h.uc.Create(ctx, req.Body.Title, req.Body.Body)
	if err != nil {
		return gen.PostsCreate200JSONResponse{Message: err.Error()}, nil
	}
	return gen.PostsCreate201JSONResponse(toGenPost(p)), nil
}

// PostsGet は GET /posts/:id のハンドラ。指定 ID の記事を返す。
func (h *PostHandler) PostsGet(ctx context.Context, req gen.PostsGetRequestObject) (gen.PostsGetResponseObject, error) {
	p, err := h.uc.GetByID(ctx, req.Id)
	if errors.Is(err, post.ErrNotFound) {
		return gen.PostsGet404JSONResponse{Message: "post not found"}, nil
	}
	if err != nil {
		return nil, err
	}
	return gen.PostsGet200JSONResponse(toGenPost(p)), nil
}

// PostsUpdate は PUT /posts/:id のハンドラ。タイトル・本文を部分更新する。
func (h *PostHandler) PostsUpdate(ctx context.Context, req gen.PostsUpdateRequestObject) (gen.PostsUpdateResponseObject, error) {
	p, err := h.uc.Update(ctx, req.Id, req.Body.Title, req.Body.Body)
	if errors.Is(err, post.ErrNotFound) {
		return gen.PostsUpdate404JSONResponse{Message: "post not found"}, nil
	}
	if err != nil {
		return nil, err
	}
	return gen.PostsUpdate200JSONResponse(toGenPost(p)), nil
}

// PostsDelete は DELETE /posts/:id のハンドラ。成功時は 204 No Content を返す。
func (h *PostHandler) PostsDelete(ctx context.Context, req gen.PostsDeleteRequestObject) (gen.PostsDeleteResponseObject, error) {
	err := h.uc.Delete(ctx, req.Id)
	if errors.Is(err, post.ErrNotFound) {
		return gen.PostsDelete404JSONResponse{Message: "post not found"}, nil
	}
	if err != nil {
		return nil, err
	}
	return gen.PostsDelete204Response{}, nil
}

// コンパイル時に PostHandler が StrictServerInterface を完全に実装しているか検証する。
var _ gen.StrictServerInterface = (*PostHandler)(nil)
