/**
 * page.tsx ("/posts/new") — 記事作成ページ
 *
 * Server Action を使ったフォーム送信の例。
 * <form action={createPost}> と書くだけで、Submit 時に Server Action が呼ばれる。
 * JavaScript が無効なブラウザでも動作する (Progressive Enhancement)。
 */
import Link from "next/link";
import { createPost } from "@/app/actions";

export default function NewPostPage() {
  return (
    <main className="max-w-2xl mx-auto py-10 px-4">
      <div className="flex items-center gap-4 mb-8">
        <Link href="/" className="text-blue-600 hover:underline">← Back</Link>
        <h1 className="text-3xl font-bold">New Post</h1>
      </div>

      {/*
        action={createPost} — フォーム送信時に Server Action (actions.ts) が実行される。
        input の name 属性が FormData のキーになる (createPost 側で formData.get("title") で取得)。
      */}
      <form action={createPost} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Title</label>
          <input
            name="title"
            required
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400"
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Body</label>
          <textarea
            name="body"
            required
            rows={8}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400"
          />
        </div>
        <button
          type="submit"
          className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700"
        >
          Create
        </button>
      </form>
    </main>
  );
}
