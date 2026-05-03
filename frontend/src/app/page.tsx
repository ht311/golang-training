/**
 * page.tsx (ルート "/") — 記事一覧ページ
 *
 * async 関数として定義された Server Component。
 * サーバーサイドで apiClient.listPosts() を直接 await できる。
 * ブラウザには HTML として届くため、APIキーやバックエンドURLがクライアントに漏れない。
 */
import Link from "next/link";
import { apiClient } from "@/lib/api";

export default async function HomePage() {
  // サーバー上でバックエンドから記事一覧を取得する
  const { posts, total } = await apiClient.listPosts();

  return (
    <main className="max-w-2xl mx-auto py-10 px-4">
      <div className="flex items-center justify-between mb-8">
        {/* total は記事の総件数 */}
        <h1 className="text-3xl font-bold">Blog Posts ({total})</h1>
        {/* Link は Next.js のクライアントサイドナビゲーション用コンポーネント */}
        <Link
          href="/posts/new"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          New Post
        </Link>
      </div>

      {posts.length === 0 ? (
        <p className="text-gray-500">No posts yet.</p>
      ) : (
        <ul className="space-y-4">
          {posts.map((post) => (
            <li key={post.id} className="border rounded p-4 hover:shadow-sm">
              {/* 記事タイトルをクリックすると詳細ページへ */}
              <Link href={`/posts/${post.id}`}>
                <h2 className="text-xl font-semibold hover:underline">{post.title}</h2>
              </Link>
              {/* line-clamp-2 で本文を2行に切り詰めて表示 */}
              <p className="text-gray-600 mt-1 line-clamp-2">{post.body}</p>
              <p className="text-xs text-gray-400 mt-2">
                {new Date(post.createdAt).toLocaleString()}
              </p>
            </li>
          ))}
        </ul>
      )}
    </main>
  );
}
