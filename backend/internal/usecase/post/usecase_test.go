package post_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainpost "github.com/ht311/golang-training/backend/internal/domain/post"
	ucpost "github.com/ht311/golang-training/backend/internal/usecase/post"
)

// mockRepo は post.Repository のテスト用モック実装。
type mockRepo struct {
	listFn    func(ctx context.Context) ([]domainpost.Post, error)
	createFn  func(ctx context.Context, title, body string) (*domainpost.Post, error)
	getByIDFn func(ctx context.Context, id string) (*domainpost.Post, error)
	updateFn  func(ctx context.Context, id string, title, body *string) (*domainpost.Post, error)
	deleteFn  func(ctx context.Context, id string) error
}

func (m *mockRepo) List(ctx context.Context) ([]domainpost.Post, error) {
	return m.listFn(ctx)
}
func (m *mockRepo) Create(ctx context.Context, title, body string) (*domainpost.Post, error) {
	return m.createFn(ctx, title, body)
}
func (m *mockRepo) GetByID(ctx context.Context, id string) (*domainpost.Post, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockRepo) Update(ctx context.Context, id string, title, body *string) (*domainpost.Post, error) {
	return m.updateFn(ctx, id, title, body)
}
func (m *mockRepo) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

var fixedTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func samplePost() *domainpost.Post {
	return &domainpost.Post{
		ID:        "test-id",
		Title:     "テストタイトル",
		Body:      "テスト本文",
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}
}

func TestPostUsecase_List(t *testing.T) {
	posts := []domainpost.Post{*samplePost()}
	uc := ucpost.NewPostUsecase(&mockRepo{
		listFn: func(_ context.Context) ([]domainpost.Post, error) { return posts, nil },
	})

	got, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 1 || got[0].ID != "test-id" {
		t.Errorf("List() = %v, want 1 post with ID test-id", got)
	}
}

func TestPostUsecase_List_Error(t *testing.T) {
	want := errors.New("db error")
	uc := ucpost.NewPostUsecase(&mockRepo{
		listFn: func(_ context.Context) ([]domainpost.Post, error) { return nil, want },
	})

	_, err := uc.List(context.Background())
	if !errors.Is(err, want) {
		t.Errorf("List() error = %v, want %v", err, want)
	}
}

func TestPostUsecase_Create(t *testing.T) {
	p := samplePost()
	uc := ucpost.NewPostUsecase(&mockRepo{
		createFn: func(_ context.Context, title, body string) (*domainpost.Post, error) { return p, nil },
	})

	got, err := uc.Create(context.Background(), "テストタイトル", "テスト本文")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if got.ID != p.ID {
		t.Errorf("Create() ID = %v, want %v", got.ID, p.ID)
	}
}

func TestPostUsecase_Create_Error(t *testing.T) {
	want := errors.New("insert failed")
	uc := ucpost.NewPostUsecase(&mockRepo{
		createFn: func(_ context.Context, title, body string) (*domainpost.Post, error) { return nil, want },
	})

	_, err := uc.Create(context.Background(), "t", "b")
	if !errors.Is(err, want) {
		t.Errorf("Create() error = %v, want %v", err, want)
	}
}

func TestPostUsecase_GetByID(t *testing.T) {
	p := samplePost()
	uc := ucpost.NewPostUsecase(&mockRepo{
		getByIDFn: func(_ context.Context, id string) (*domainpost.Post, error) { return p, nil },
	})

	got, err := uc.GetByID(context.Background(), "test-id")
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.ID != "test-id" {
		t.Errorf("GetByID() ID = %v, want test-id", got.ID)
	}
}

func TestPostUsecase_GetByID_NotFound(t *testing.T) {
	uc := ucpost.NewPostUsecase(&mockRepo{
		getByIDFn: func(_ context.Context, id string) (*domainpost.Post, error) {
			return nil, domainpost.ErrNotFound
		},
	})

	_, err := uc.GetByID(context.Background(), "missing")
	if !errors.Is(err, domainpost.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
}

func TestPostUsecase_Update(t *testing.T) {
	title := "新タイトル"
	p := &domainpost.Post{ID: "test-id", Title: title, Body: "本文", CreatedAt: fixedTime, UpdatedAt: fixedTime}
	uc := ucpost.NewPostUsecase(&mockRepo{
		updateFn: func(_ context.Context, id string, t, b *string) (*domainpost.Post, error) { return p, nil },
	})

	got, err := uc.Update(context.Background(), "test-id", &title, nil)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if got.Title != title {
		t.Errorf("Update() Title = %v, want %v", got.Title, title)
	}
}

func TestPostUsecase_Update_NotFound(t *testing.T) {
	uc := ucpost.NewPostUsecase(&mockRepo{
		updateFn: func(_ context.Context, id string, t, b *string) (*domainpost.Post, error) {
			return nil, domainpost.ErrNotFound
		},
	})

	title := "x"
	_, err := uc.Update(context.Background(), "missing", &title, nil)
	if !errors.Is(err, domainpost.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestPostUsecase_Delete(t *testing.T) {
	uc := ucpost.NewPostUsecase(&mockRepo{
		deleteFn: func(_ context.Context, id string) error { return nil },
	})

	if err := uc.Delete(context.Background(), "test-id"); err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}
}

func TestPostUsecase_Delete_NotFound(t *testing.T) {
	uc := ucpost.NewPostUsecase(&mockRepo{
		deleteFn: func(_ context.Context, id string) error { return domainpost.ErrNotFound },
	})

	err := uc.Delete(context.Background(), "missing")
	if !errors.Is(err, domainpost.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
