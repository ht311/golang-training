// post パッケージはブログ記事のドメイン層。
// DB や HTTP の都合に依存しない純粋なビジネスデータと契約を定義する。
package post

import "time"

// Post はブログ記事を表すエンティティ。
type Post struct {
	ID        string
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
