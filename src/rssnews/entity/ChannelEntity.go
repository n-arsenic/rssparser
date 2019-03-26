package entity

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type Channel struct {
	Id          int
	Rss_url     string
	Description sql.NullString
	Pub_date    pq.NullTime
	Created_at  time.Time
}

/*
func (ch *Channel) GetCreatedAt() time.Time {
	return time.Now()
}
*/
