package entity

import (
	"time"
)

type User struct {
	Id         int
	Name       string
	Password   string
	Created_at time.Time
	//Last_login time.Time
}
