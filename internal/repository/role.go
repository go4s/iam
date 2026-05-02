package repository

import (
	"github.com/go4s/iam/internal/db"
	"github.com/go4s/iam/internal/model"
)

type RoleRepository struct{}

func (r *RoleRepository) Create(role *model.Role) error {
	_, err := db.Engine.Insert(role)
	return err
}

func (r *RoleRepository) GetByID(id int64) (*model.Role, error) {
	role := new(model.Role)
	has, err := db.Engine.ID(id).Get(role)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return role, nil
}

func (r *RoleRepository) GetByCode(code string) (*model.Role, error) {
	role := new(model.Role)
	has, err := db.Engine.Where("code = ?", code).Get(role)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return role, nil
}

func (r *RoleRepository) List(page, size int) ([]model.Role, int64, error) {
	var roles []model.Role
	total, err := db.Engine.Limit(size, (page-1)*size).FindAndCount(&roles)
	if err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (r *RoleRepository) Update(role *model.Role) error {
	_, err := db.Engine.ID(role.ID).Update(role)
	return err
}

func (r *RoleRepository) Delete(id int64) error {
	_, err := db.Engine.ID(id).Delete(new(model.Role))
	return err
}
