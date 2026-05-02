package model

import "time"

type UserRole struct {
	ID        int64     `xorm:"pk autoincr"`
	UserID    int64     `xorm:"'user_id' notnull index"`
	RoleID    int64     `xorm:"'role_id' notnull index"`
	CreatedAt time.Time `xorm:"created"`
}
