package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Phone    string
	Password string
	NickName string
	Birthday time.Time
	AboutMe  string
	Status   uint8
	Ctime    time.Time
}
