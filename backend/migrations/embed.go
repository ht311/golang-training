// migrations パッケージは SQL マイグレーションファイルをバイナリに埋め込む。
// //go:embed で *.sql を埋め込むことで、デプロイ時にファイルを別途配置せずに済む。
package migrations

import "embed"

// FS はビルド時に *.sql ファイルを埋め込んだ仮想ファイルシステム。
// golang-migrate の iofs ドライバがこれを読み込んでマイグレーションを実行する。
//
//go:embed *.sql
var FS embed.FS
