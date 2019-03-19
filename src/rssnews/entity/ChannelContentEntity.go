package entity

import (
	"time"
)

type ChannelContent struct {
	Channel_id  int
	Link        string
	Title       string
	Description string
	Pub_date    time.Time
}
