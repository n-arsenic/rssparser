package entity

import (
	"github.com/lib/pq"
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

func (sh *Scheduler) GetNewStatus() string {
	return StatusType['n']
}

func (sh *Scheduler) GetWorkStatus() string {
	return StatusType['w']
}

func (sh *Scheduler) GetErrorStatus() string {
	return StatusType['e']
}

func (sh *Scheduler) GetSuccessStatus() string {
	return StatusType['s']
}

func (sh *Scheduler) SetNewStatus(ent *Scheduler, status string) {
	ent.Status = sh.GetNewStatus()
}
