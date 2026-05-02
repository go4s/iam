package model

import "time"

type Permission struct {
	ID          int64     `xorm:"pk autoincr"`
	Name        string    `xorm:"notnull"`
	Code        string    `xorm:"unique notnull"`
	Resource    string    `xorm:"notnull"`
	Action      string    `xorm:"notnull"`
	Description string    `xorm:"text"`
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
}
