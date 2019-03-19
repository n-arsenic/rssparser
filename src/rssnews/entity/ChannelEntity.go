package entity

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

var StatusType = map[byte]string{
	'w': "wait", // in progress
	'e': "error",
	's': "success",
	'n': "new",
}

type Channel struct {
	Id          int
	Rss_url     string
	Description sql.NullString
	Pub_date    pq.NullTime
	Parsed_at   pq.NullTime
	Status      string
	Created_at  time.Time
	Start_parse pq.NullTime
}

func (ch *Channel) GetNewStatus() string {
	return StatusType['n']
}

func (ch *Channel) GetWaitStatus() string {
	return StatusType['w']
}

func (ch *Channel) GetErrorStatus() string {
	return StatusType['e']
}

func (ch *Channel) GetSuccessStatus() string {
	return StatusType['s']
}

func (ch *Channel) GetCreatedAt() time.Time {
	return time.Now()
}

func (ch *Channel) SetNewStatus(ent *Channel, status string) {
	ent.Status = ch.GetNewStatus()
}
