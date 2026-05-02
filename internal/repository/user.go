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

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	user := new(model.User)
	has, err := db.Engine.ID(id).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}

func (r *UserRepository) List(page, size int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	session := db.Engine.Limit(size, (page-1)*size)
	if keyword != "" {
		session = session.Where("username LIKE ?", "%"+keyword+"%")
	}
	total, err := session.FindAndCount(&users)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) Update(user *model.User) error {
	_, err := db.Engine.ID(user.ID).Cols("username", "password_hash", "updated_at").Update(user)
	return err
}
