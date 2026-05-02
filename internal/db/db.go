package db

import (
    "github.com/go4s/iam/internal/model"
    _ "modernc.org/sqlite"
    "xorm.io/xorm"
)

var Engine *xorm.Engine

func InitDB() error {
    var err error
    // modernc.org/sqlite registers itself as "sqlite"
    Engine, err = xorm.NewEngine("sqlite", "iam.db")
    if err != nil {
        return err
    }
    
    Engine.ShowSQL(true)
    
    // Auto migrate
    if err := Engine.Sync(
        new(model.User),
        new(model.Role),
        new(model.Permission),
        new(model.UserRole),
        new(model.RolePermission),
        new(model.EntityFormat),
    ); err != nil {
        return err
    }
    Engine.SetMaxOpenConns(1)
    return Engine.Ping()
}
