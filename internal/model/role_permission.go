package model

import "time"

type RolePermission struct {
	ID           int64     `xorm:"pk autoincr"`
	RoleID       int64     `xorm:"'role_id' notnull index"`
	PermissionID int64     `xorm:"'permission_id' notnull index"`
	CreatedAt    time.Time `xorm:"created"`
}
