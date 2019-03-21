package entity

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

var StatusType = map[byte]string{
	'w': "work",
	'e': "error",
	's': "success",
	'n': "new",
}

type Scheduler struct {
	Channel_id int
	Rss_url    string
	Status     string
	Finish     pq.NullTime
	Start      pq.NullTime
}

func (ch *Channel) GetNewStatus() string {
	return StatusType['n']
}

func (ch *Channel) GetWorkStatus() string {
	return StatusType['w']
}

func (ch *Channel) GetErrorStatus() string {
	return StatusType['e']
}

func (ch *Channel) GetSuccessStatus() string {
	return StatusType['s']
}

func (ch *Channel) SetNewStatus(ent *Channel, status string) {
	ent.Status = ch.GetNewStatus()
}
