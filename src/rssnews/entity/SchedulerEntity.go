package entity

import (
	"github.com/lib/pq"
	"time"
)

var StatusType = map[byte]string{
	'w': "work",
	'e': "error",
	's': "success",
}

type Scheduler struct {
	Channel_id int
	Rss_url    string
	Status     string
	Finish     pq.NullTime
	Start      pq.NullTime
	Plan_start pq.NullTime
	Message    string
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

func (sh *Scheduler) SetWorkStatus() {
	sh.Status = sh.GetWorkStatus()
}

func (sh *Scheduler) SetErrorStatus() {
	sh.Status = sh.GetErrorStatus()
}

func (sh *Scheduler) SetSuccessStatus() {
	sh.Status = sh.GetSuccessStatus()
}

func (sh *Scheduler) SetError(mess string) {
	sh.Status = sh.GetErrorStatus()
	sh.Message = mess
	sh.SetFinish()
}

func (sh *Scheduler) SetFinish() {
	sh.Finish = pq.NullTime{Time: time.Now(), Valid: true}
}

func (sh *Scheduler) SetPlanStart(period time.Duration) {
	date := time.Now().Add(time.Duration(period))
	sh.Plan_start = pq.NullTime{Time: date, Valid: true}
}
