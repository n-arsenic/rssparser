package entity

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type Channel struct {
	Id          int
	Rss_url     string
	Link        sql.NullString
	Title       sql.NullString
	Description sql.NullString
	Pub_date    pq.NullTime
	Created_at  time.Time
}
