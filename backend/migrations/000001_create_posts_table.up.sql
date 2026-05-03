-- posts テーブルを作成する
-- golang-migrate がこのファイルを一度だけ実行し、schema_migrations テーブルに記録する
CREATE TABLE IF NOT EXISTS posts (
    id         TEXT        PRIMARY KEY,
    title      TEXT        NOT NULL,
    body       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
