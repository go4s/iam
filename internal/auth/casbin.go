package auth

import (
	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	"github.com/go4s/iam/internal/db"
)

var Enforcer *casbin.Enforcer

func InitCasbin() error {
    // Initialize xorm adapter with the existing engine
    adapter, err := xormadapter.NewAdapterByEngine(db.Engine)
    if err != nil {
        return err
    }
    
    // Load the model from file
    var err2 error
    Enforcer, err2 = casbin.NewEnforcer("configs/model.conf", adapter)
    if err2 != nil {
        return err2
    }
    
    return Enforcer.LoadPolicy()
}
