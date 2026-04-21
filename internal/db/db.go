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
    if err := Engine.Sync(new(model.User)); err != nil {
        return err
    }
    
    return Engine.Ping()
}
