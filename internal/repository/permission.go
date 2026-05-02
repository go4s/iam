package repository

import (
	"github.com/go4s/iam/internal/db"
	"github.com/go4s/iam/internal/model"
)

type PermissionRepository struct{}

func (r *PermissionRepository) Create(perm *model.Permission) error {
	_, err := db.Engine.Insert(perm)
	return err
}

func (r *PermissionRepository) GetByID(id int64) (*model.Permission, error) {
	perm := new(model.Permission)
	has, err := db.Engine.ID(id).Get(perm)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return perm, nil
}

func (r *PermissionRepository) GetByCode(code string) (*model.Permission, error) {
	perm := new(model.Permission)
	has, err := db.Engine.Where("code = ?", code).Get(perm)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return perm, nil
}

func (r *PermissionRepository) List(page, size int) ([]model.Permission, int64, error) {
	var perms []model.Permission
	total, err := db.Engine.Limit(size, (page-1)*size).FindAndCount(&perms)
	if err != nil {
		return nil, 0, err
	}
	return perms, total, nil
}

func (r *PermissionRepository) Update(perm *model.Permission) error {
	_, err := db.Engine.ID(perm.ID).Update(perm)
	return err
}

func (r *PermissionRepository) Delete(id int64) error {
	_, err := db.Engine.ID(id).Delete(new(model.Permission))
	return err
}
