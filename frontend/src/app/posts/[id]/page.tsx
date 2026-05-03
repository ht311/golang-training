/**
 * page.tsx ("/posts/[id]") — 記事詳細ページ
 *
 * [id] はダイナミックルートセグメント。URL の /posts/abc123 の "abc123" が id に入る。
 * Next.js 16 では params が Promise になったため await が必要。
 */
import Link from "next/link";
import { notFound } from "next/navigation";
import { apiClient } from "@/lib/api";
import { deletePost } from "@/app/actions";

/** Next.js App Router がページコンポーネントに渡す props の型 */
interface Props {
  params: Promise<{ id: string }>; // Next.js 16: params は非同期で解決される
}

export default async function PostDetailPage({ params }: Props) {
  // params を await して id を取り出す (Next.js 16 の仕様)
  const { id } = await params;

  let post;
  try {
    post = await apiClient.getPost(id);
  } catch {
    // 404 や取得失敗時は Next.js の notFound() を呼び 404 ページを表示する
    notFound();
  }

  // deletePost は (id: string) => void だが、フォームの action には引数なしの関数が必要。
  // bind(null, id) で id を事前に束縛した新しい関数を作る。
  const deleteWithId = deletePost.bind(null, id);

  return (
    <main className="max-w-2xl mx-auto py-10 px-4">
      <div className="flex items-center gap-4 mb-8">
        <Link href="/" className="text-blue-600 hover:underline">← Back</Link>
      </div>

      <h1 className="text-3xl font-bold mb-2">{post.title}</h1>
      <p className="text-xs text-gray-400 mb-6">
        {new Date(post.createdAt).toLocaleString()}
      </p>
      {/* whitespace-pre-wrap で改行を保持して表示 */}
      <p className="text-gray-800 whitespace-pre-wrap">{post.body}</p>

      {/* 削除フォーム: Server Action で削除してから一覧ページへリダイレクト */}
      <form action={deleteWithId} className="mt-10">
        <button
          type="submit"
          className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
          onClick={(e) => {
            // confirm ダイアログで誤削除を防ぐ。キャンセル時はフォーム送信を止める
            if (!confirm("Delete this post?")) e.preventDefault();
          }}
        >
          Delete Post
        </button>
      </form>
    </main>
  );
}
