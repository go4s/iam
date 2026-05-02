package repository

import (
	"github.com/go4s/iam/internal/db"
	"github.com/go4s/iam/internal/model"
)

type UserRoleRepository struct{}

func (r *UserRoleRepository) AddUserRole(userID, roleID int64) error {
	ur := &model.UserRole{UserID: userID, RoleID: roleID}
	_, err := db.Engine.Insert(ur)
	return err
}

func (r *UserRoleRepository) RemoveUserRole(userID, roleID int64) error {
	_, err := db.Engine.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(new(model.UserRole))
	return err
}

func (r *UserRoleRepository) GetRoleIDsByUserID(userID int64) ([]int64, error) {
	var urs []model.UserRole
	if err := db.Engine.Where("user_id = ?", userID).Find(&urs); err != nil {
		return nil, err
	}
	var roleIDs []int64
	for _, ur := range urs {
		roleIDs = append(roleIDs, ur.RoleID)
	}
	return roleIDs, nil
}

func (r *UserRoleRepository) GetUserIDsByRoleID(roleID int64) ([]int64, error) {
	var urs []model.UserRole
	if err := db.Engine.Where("role_id = ?", roleID).Find(&urs); err != nil {
		return nil, err
	}
	var userIDs []int64
	for _, ur := range urs {
		userIDs = append(userIDs, ur.UserID)
	}
	return userIDs, nil
}

func (r *UserRoleRepository) DeleteByUserID(userID int64) error {
	_, err := db.Engine.Where("user_id = ?", userID).Delete(new(model.UserRole))
	return err
}
