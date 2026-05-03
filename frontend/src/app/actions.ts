/**
 * actions.ts — Server Actions
 *
 * "use server" ディレクティブにより、このファイルの関数はサーバー上でのみ実行される。
 * クライアント (ブラウザ) から呼ばれた場合も、実際の処理はサーバー側で走る。
 * フォームの action={createPost} のように JSX から直接渡せるのが特徴。
 */
"use server";

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { apiClient } from "@/lib/api";

/**
 * createPost — 記事作成フォームの送信を処理する Server Action。
 *
 * <form action={createPost}> と組み合わせることで、フォーム送信時にこの関数がサーバーで実行される。
 * 作成成功後は作成した記事の詳細ページへリダイレクトする。
 *
 * @param formData ブラウザが送信した FormData (name 属性でフィールドを取得する)
 */
export async function createPost(formData: FormData) {
  const title = formData.get("title") as string;
  const body = formData.get("body") as string;
  const post = await apiClient.createPost({ title, body });
  // 作成した記事の詳細ページへ転送
  redirect(`/posts/${post.id}`);
}

/**
 * deletePost — 記事削除ボタンを処理する Server Action。
 *
 * deletePost.bind(null, id) で id を束縛してから <form action> に渡す。
 * 削除後は "/" のキャッシュを破棄してから一覧ページへリダイレクトする。
 * revalidatePath("/") を呼ばないと、Next.js のキャッシュが残り削除前の一覧が表示されてしまう。
 *
 * @param id 削除する記事の UUID
 */
export async function deletePost(id: string) {
  await apiClient.deletePost(id);
  // "/" ページのキャッシュを無効化して次回アクセス時に最新データを取得させる
  revalidatePath("/");
  redirect("/");
}
