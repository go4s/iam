package repository

import (
    "github.com/go4s/iam/internal/db"
    "github.com/go4s/iam/internal/model"
)

type UserRepository struct{}

func (r *UserRepository) Create(user *model.User) error {
    _, err := db.Engine.Insert(user)
    return err
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
    user := new(model.User)
    has, err := db.Engine.Where("username = ?", username).Get(user)
    if err != nil {
        return nil, err
    }
    if !has {
        return nil, nil
    }
    return user, nil
}
