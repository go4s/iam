package model

import "time"

type User struct {
	ID           int64     `xorm:"pk autoincr"`
	Username     string    `xorm:"unique notnull"`
	PasswordHash string    `xorm:"notnull"`
	Role         string    `xorm:"notnull"`
	CreatedAt    time.Time `xorm:"created"`
	UpdatedAt    time.Time `xorm:"updated"`
}
