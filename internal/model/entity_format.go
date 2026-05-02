package model

import "time"

type EntityFormat struct {
	ID        int64     `xorm:"pk autoincr"`
	Template  string    `xorm:"unique(template_mode) notnull"`
	Mode      string    `xorm:"unique(template_mode) notnull"`
	Fields    string    `xorm:"text notnull"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
