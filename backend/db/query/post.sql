-- sqlc がこのファイルを読んで型安全な Go コードを自動生成する。
-- コメント行の ":many / :one / :exec / :execrows" が戻り値の型を決める。

-- name: ListPosts :many
-- 全記事を作成日時の降順 (新しい順) で取得する
SELECT id, title, body, created_at, updated_at
FROM posts
ORDER BY created_at DESC;

-- name: GetPost :one
-- 指定 ID の記事を1件取得する。存在しない場合は pgx.ErrNoRows が返る
SELECT id, title, body, created_at, updated_at
FROM posts
WHERE id = $1;

-- name: CreatePost :one
-- 新しい記事を挿入して挿入結果をそのまま返す (RETURNING *)
-- created_at / updated_at は DB 側で NOW() をセットするので渡す必要がない
INSERT INTO posts (id, title, body, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING *;

-- name: UpdatePost :one
-- タイトル・本文を更新して更新後のレコードを返す
-- updated_at も DB 側で NOW() にする
UPDATE posts
SET title      = $1,
    body       = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: DeletePost :execrows
-- 指定 ID の記事を削除する。戻り値は削除された行数 (0 なら対象なし)
DELETE FROM posts
WHERE id = $1;
