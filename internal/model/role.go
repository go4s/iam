package model

import "time"

type Role struct {
	ID          int64     `xorm:"pk autoincr"`
	Name        string    `xorm:"notnull"`
	Code        string    `xorm:"unique notnull"`
	Description string    `xorm:"text"`
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
}
