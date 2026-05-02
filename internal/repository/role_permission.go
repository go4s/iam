package repository

import (
	"github.com/go4s/iam/internal/db"
	"github.com/go4s/iam/internal/model"
)

type RolePermissionRepository struct{}

func (r *RolePermissionRepository) AddRolePermission(roleID, permissionID int64) error {
	rp := &model.RolePermission{RoleID: roleID, PermissionID: permissionID}
	_, err := db.Engine.Insert(rp)
	return err
}

func (r *RolePermissionRepository) RemoveRolePermission(roleID, permissionID int64) error {
	_, err := db.Engine.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(new(model.RolePermission))
	return err
}

func (r *RolePermissionRepository) GetPermissionIDsByRoleID(roleID int64) ([]int64, error) {
	var rps []model.RolePermission
	if err := db.Engine.Where("role_id = ?", roleID).Find(&rps); err != nil {
		return nil, err
	}
	var permIDs []int64
	for _, rp := range rps {
		permIDs = append(permIDs, rp.PermissionID)
	}
	return permIDs, nil
}

func (r *RolePermissionRepository) GetRoleIDsByPermissionID(permissionID int64) ([]int64, error) {
	var rps []model.RolePermission
	if err := db.Engine.Where("permission_id = ?", permissionID).Find(&rps); err != nil {
		return nil, err
	}
	var roleIDs []int64
	for _, rp := range rps {
		roleIDs = append(roleIDs, rp.RoleID)
	}
	return roleIDs, nil
}

func (r *RolePermissionRepository) DeleteByRoleID(roleID int64) error {
	_, err := db.Engine.Where("role_id = ?", roleID).Delete(new(model.RolePermission))
	return err
}
