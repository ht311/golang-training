package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ht311/golang-training/backend/gen"
	domainpost "github.com/ht311/golang-training/backend/internal/domain/post"
	"github.com/ht311/golang-training/backend/internal/handler"
	postuc "github.com/ht311/golang-training/backend/internal/usecase/post"
)

// mockRepo は post.Repository のテスト用モック実装。
type mockRepo struct {
	listFn    func(ctx context.Context) ([]domainpost.Post, error)
	createFn  func(ctx context.Context, title, body string) (*domainpost.Post, error)
	getByIDFn func(ctx context.Context, id string) (*domainpost.Post, error)
	updateFn  func(ctx context.Context, id string, title, body *string) (*domainpost.Post, error)
	deleteFn  func(ctx context.Context, id string) error
}

func (m *mockRepo) List(ctx context.Context) ([]domainpost.Post, error) { return m.listFn(ctx) }
func (m *mockRepo) Create(ctx context.Context, title, body string) (*domainpost.Post, error) {
	return m.createFn(ctx, title, body)
}
func (m *mockRepo) GetByID(ctx context.Context, id string) (*domainpost.Post, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockRepo) Update(ctx context.Context, id string, title, body *string) (*domainpost.Post, error) {
	return m.updateFn(ctx, id, title, body)
}
func (m *mockRepo) Delete(ctx context.Context, id string) error { return m.deleteFn(ctx, id) }

var fixedTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func newHandler(repo domainpost.Repository) *handler.PostHandler {
	return handler.NewPostHandler(postuc.NewPostUsecase(repo))
}

func sampleDomainPost() *domainpost.Post {
	return &domainpost.Post{
		ID:        "id-1",
		Title:     "タイトル",
		Body:      "本文",
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}
}

func TestHealthCheck(t *testing.T) {
	h := newHandler(&mockRepo{})
	resp, err := h.HealthCheck(context.Background(), gen.HealthCheckRequestObject{})
	if err != nil {
		t.Fatalf("HealthCheck() error = %v", err)
	}
	got, ok := resp.(gen.HealthCheck200JSONResponse)
	if !ok {
		t.Fatalf("HealthCheck() resp type = %T, want HealthCheck200JSONResponse", resp)
	}
	if got.Status != "ok" {
		t.Errorf("HealthCheck() status = %q, want ok", got.Status)
	}
}

func TestPostsList_Success(t *testing.T) {
	p := sampleDomainPost()
	h := newHandler(&mockRepo{
		listFn: func(_ context.Context) ([]domainpost.Post, error) {
			return []domainpost.Post{*p}, nil
		},
	})

	resp, err := h.PostsList(context.Background(), gen.PostsListRequestObject{})
	if err != nil {
		t.Fatalf("PostsList() error = %v", err)
	}
	body, ok := resp.(gen.PostsList200JSONResponse)
	if !ok {
		t.Fatalf("PostsList() resp type = %T, want PostsList200JSONResponse", resp)
	}
	list, err := body.AsPostList()
	if err != nil {
		t.Fatalf("AsPostList() error = %v", err)
	}
	if list.Total != 1 || list.Posts[0].Id != "id-1" {
		t.Errorf("PostsList() = %+v", list)
	}
}

func TestPostsList_Error(t *testing.T) {
	h := newHandler(&mockRepo{
		listFn: func(_ context.Context) ([]domainpost.Post, error) {
			return nil, errors.New("db error")
		},
	})

	resp, err := h.PostsList(context.Background(), gen.PostsListRequestObject{})
	if err != nil {
		t.Fatalf("PostsList() error = %v", err)
	}
	body, ok := resp.(gen.PostsList200JSONResponse)
	if !ok {
		t.Fatalf("PostsList() resp type = %T, want PostsList200JSONResponse", resp)
	}
	errResp, err := body.AsErrorResponse()
	if err != nil {
		t.Fatalf("AsErrorResponse() error = %v", err)
	}
	if errResp.Message != "db error" {
		t.Errorf("PostsList() error message = %q, want db error", errResp.Message)
	}
}

func TestPostsCreate_Success(t *testing.T) {
	p := sampleDomainPost()
	h := newHandler(&mockRepo{
		createFn: func(_ context.Context, title, body string) (*domainpost.Post, error) {
			return p, nil
		},
	})

	resp, err := h.PostsCreate(context.Background(), gen.PostsCreateRequestObject{
		Body: &gen.CreatePostRequest{Title: "タイトル", Body: "本文"},
	})
	if err != nil {
		t.Fatalf("PostsCreate() error = %v", err)
	}
	got, ok := resp.(gen.PostsCreate201JSONResponse)
	if !ok {
		t.Fatalf("PostsCreate() resp type = %T, want PostsCreate201JSONResponse", resp)
	}
	if got.Id != "id-1" {
		t.Errorf("PostsCreate() id = %v, want id-1", got.Id)
	}
}

func TestPostsCreate_Error(t *testing.T) {
	h := newHandler(&mockRepo{
		createFn: func(_ context.Context, title, body string) (*domainpost.Post, error) {
			return nil, errors.New("create failed")
		},
	})

	resp, err := h.PostsCreate(context.Background(), gen.PostsCreateRequestObject{
		Body: &gen.CreatePostRequest{Title: "t", Body: "b"},
	})
	if err != nil {
		t.Fatalf("PostsCreate() error = %v", err)
	}
	if _, ok := resp.(gen.PostsCreate200JSONResponse); !ok {
		t.Errorf("PostsCreate() resp type = %T, want PostsCreate200JSONResponse", resp)
	}
}

func TestPostsGet_Success(t *testing.T) {
	p := sampleDomainPost()
	h := newHandler(&mockRepo{
		getByIDFn: func(_ context.Context, id string) (*domainpost.Post, error) { return p, nil },
	})

	resp, err := h.PostsGet(context.Background(), gen.PostsGetRequestObject{Id: "id-1"})
	if err != nil {
		t.Fatalf("PostsGet() error = %v", err)
	}
	got, ok := resp.(gen.PostsGet200JSONResponse)
	if !ok {
		t.Fatalf("PostsGet() resp type = %T, want PostsGet200JSONResponse", resp)
	}
	if got.Id != "id-1" {
		t.Errorf("PostsGet() id = %v, want id-1", got.Id)
	}
}

func TestPostsGet_NotFound(t *testing.T) {
	h := newHandler(&mockRepo{
		getByIDFn: func(_ context.Context, id string) (*domainpost.Post, error) {
			return nil, domainpost.ErrNotFound
		},
	})

	resp, err := h.PostsGet(context.Background(), gen.PostsGetRequestObject{Id: "missing"})
	if err != nil {
		t.Fatalf("PostsGet() error = %v", err)
	}
	if _, ok := resp.(gen.PostsGet404JSONResponse); !ok {
		t.Errorf("PostsGet() resp type = %T, want PostsGet404JSONResponse", resp)
	}
}

func TestPostsUpdate_Success(t *testing.T) {
	title := "新タイトル"
	p := &domainpost.Post{ID: "id-1", Title: title, Body: "本文", CreatedAt: fixedTime, UpdatedAt: fixedTime}
	h := newHandler(&mockRepo{
		updateFn: func(_ context.Context, id string, t, b *string) (*domainpost.Post, error) {
			return p, nil
		},
	})

	resp, err := h.PostsUpdate(context.Background(), gen.PostsUpdateRequestObject{
		Id:   "id-1",
		Body: &gen.UpdatePostRequest{Title: &title},
	})
	if err != nil {
		t.Fatalf("PostsUpdate() error = %v", err)
	}
	got, ok := resp.(gen.PostsUpdate200JSONResponse)
	if !ok {
		t.Fatalf("PostsUpdate() resp type = %T, want PostsUpdate200JSONResponse", resp)
	}
	if got.Title != title {
		t.Errorf("PostsUpdate() title = %v, want %v", got.Title, title)
	}
}

func TestPostsUpdate_NotFound(t *testing.T) {
	h := newHandler(&mockRepo{
		updateFn: func(_ context.Context, id string, t, b *string) (*domainpost.Post, error) {
			return nil, domainpost.ErrNotFound
		},
	})

	title := "x"
	resp, err := h.PostsUpdate(context.Background(), gen.PostsUpdateRequestObject{
		Id:   "missing",
		Body: &gen.UpdatePostRequest{Title: &title},
	})
	if err != nil {
		t.Fatalf("PostsUpdate() error = %v", err)
	}
	if _, ok := resp.(gen.PostsUpdate404JSONResponse); !ok {
		t.Errorf("PostsUpdate() resp type = %T, want PostsUpdate404JSONResponse", resp)
	}
}

func TestPostsDelete_Success(t *testing.T) {
	h := newHandler(&mockRepo{
		deleteFn: func(_ context.Context, id string) error { return nil },
	})

	resp, err := h.PostsDelete(context.Background(), gen.PostsDeleteRequestObject{Id: "id-1"})
	if err != nil {
		t.Fatalf("PostsDelete() error = %v", err)
	}
	if _, ok := resp.(gen.PostsDelete204Response); !ok {
		t.Errorf("PostsDelete() resp type = %T, want PostsDelete204Response", resp)
	}
}

func TestPostsDelete_NotFound(t *testing.T) {
	h := newHandler(&mockRepo{
		deleteFn: func(_ context.Context, id string) error { return domainpost.ErrNotFound },
	})

	resp, err := h.PostsDelete(context.Background(), gen.PostsDeleteRequestObject{Id: "missing"})
	if err != nil {
		t.Fatalf("PostsDelete() error = %v", err)
	}
	if _, ok := resp.(gen.PostsDelete404JSONResponse); !ok {
		t.Errorf("PostsDelete() resp type = %T, want PostsDelete404JSONResponse", resp)
	}
}
