package entity

import (
	"time"
)

type ChannelContent struct {
	Channel_id  int
	Link        string
	Title       string
	Author      string
	Category    string
	Description string
	Pub_date    time.Time
}
