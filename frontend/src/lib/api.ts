/**
 * api.ts — バックエンド REST API への型安全なクライアント
 *
 * ここで定義した型は OpenAPI 仕様 (api/openapi/openapi.yaml) と手動で同期している。
 * 仕様変更時はこのファイルの型も合わせて更新すること。
 * baseUrl は環境変数 NEXT_PUBLIC_API_URL で切り替え可能 (ローカル/本番)。
 */

/** バックエンドの Post モデルに対応する型 */
export interface Post {
  id: string;
  title: string;
  body: string;
  createdAt: string; // ISO 8601 文字列 (例: "2024-01-01T00:00:00Z")
  updatedAt: string;
}

/** 一覧取得レスポンス: 記事の配列と総件数を持つ */
export interface PostList {
  posts: Post[];
  total: number;
}

/** 記事作成リクエストのボディ */
export interface CreatePostRequest {
  title: string;
  body: string;
}

/**
 * 記事更新リクエストのボディ。
 * フィールドは全て省略可能 (undefined の場合はサーバー側で変更しない)。
 */
export interface UpdatePostRequest {
  title?: string;
  body?: string;
}

/**
 * NEXT_PUBLIC_ プレフィックスがつく環境変数はブラウザにも公開される。
 * ローカル開発では .env.local に NEXT_PUBLIC_API_URL=http://localhost:8080 を設定する。
 */
const baseUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

/**
 * apiClient — バックエンドの各エンドポイントを呼ぶ関数をまとめたオブジェクト。
 * Server Component からサーバーサイドで呼ばれる場合と、
 * Server Action から呼ばれる場合の両方に対応する。
 */
export const apiClient = {
  /** GET /posts — 全記事を取得する */
  async listPosts(): Promise<PostList> {
    const res = await fetch(`${baseUrl}/posts`);
    if (!res.ok) throw new Error("Failed to fetch posts");
    return res.json();
  },

  /** GET /posts/:id — 指定 ID の記事を1件取得する */
  async getPost(id: string): Promise<Post> {
    const res = await fetch(`${baseUrl}/posts/${id}`);
    if (!res.ok) throw new Error("Post not found");
    return res.json();
  },

  /** POST /posts — 新しい記事を作成する */
  async createPost(data: CreatePostRequest): Promise<Post> {
    const res = await fetch(`${baseUrl}/posts`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to create post");
    return res.json();
  },

  /** PUT /posts/:id — 既存の記事を更新する */
  async updatePost(id: string, data: UpdatePostRequest): Promise<Post> {
    const res = await fetch(`${baseUrl}/posts/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error("Failed to update post");
    return res.json();
  },

  /** DELETE /posts/:id — 記事を削除する (レスポンスボディなし) */
  async deletePost(id: string): Promise<void> {
    const res = await fetch(`${baseUrl}/posts/${id}`, { method: "DELETE" });
    if (!res.ok) throw new Error("Failed to delete post");
  },
};
