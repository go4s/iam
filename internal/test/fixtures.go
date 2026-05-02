package test

import (
	"fmt"

	"github.com/go4s/iam/internal/db"
	"github.com/go4s/iam/internal/model"
	"github.com/go4s/iam/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"xorm.io/xorm"
)

// SetupTestDB 初始化测试数据库
func SetupTestDB() error {
	// 使用内存数据库
	var err error
	db.Engine, err = xorm.NewEngine("sqlite", ":memory:")
	if err != nil {
		return err
	}

	// 同步表结构
	if err := db.Engine.Sync(
		new(model.User),
		new(model.Role),
		new(model.Permission),
		new(model.UserRole),
		new(model.RolePermission),
		new(model.EntityFormat),
	); err != nil {
		return err
	}

	return nil
}

// LoadFixtures 加载测试数据
func LoadFixtures() error {
	// 1. 创建角色
	roleRepo := &repository.RoleRepository{}
	adminRole := &model.Role{
		Name:        "管理员",
		Code:        "admin",
		Description: "系统管理员",
	}
	if err := roleRepo.Create(adminRole); err != nil {
		return err
	}

	editorRole := &model.Role{
		Name:        "编辑",
		Code:        "editor",
		Description: "内容编辑",
	}
	if err := roleRepo.Create(editorRole); err != nil {
		return err
	}

	// 2. 创建权限
	permRepo := &repository.PermissionRepository{}
	permissions := []*model.Permission{
		{Name: "创建用户", Code: "user:create", Resource: "user", Action: "create"},
		{Name: "读取用户", Code: "user:read", Resource: "user", Action: "read"},
		{Name: "更新用户", Code: "user:update", Resource: "user", Action: "update"},
		{Name: "删除用户", Code: "user:delete", Resource: "user", Action: "delete"},
		{Name: "管理权限", Code: "permission:manage", Resource: "permission", Action: "manage"},
	}
	for _, p := range permissions {
		if err := permRepo.Create(p); err != nil {
			return err
		}
	}

	// 3. 创建用户（密码: admin123）
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userRepo := &repository.UserRepository{}
	adminUser := &model.User{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
		Status:       "active",
	}
	if err := userRepo.Create(adminUser); err != nil {
		return err
	}

	editorUser := &model.User{
		Username:     "editor",
		PasswordHash: string(hashedPassword),
		Status:       "active",
	}
	if err := userRepo.Create(editorUser); err != nil {
		return err
	}

	// 4. 关联用户角色
	userRoleRepo := &repository.UserRoleRepository{}
	if err := userRoleRepo.AddUserRole(adminUser.ID, adminRole.ID); err != nil {
		return err
	}
	if err := userRoleRepo.AddUserRole(editorUser.ID, editorRole.ID); err != nil {
		return err
	}

	// 5. 关联角色权限
	rolePermRepo := &repository.RolePermissionRepository{}
	for _, p := range permissions {
		if err := rolePermRepo.AddRolePermission(adminRole.ID, p.ID); err != nil {
			return err
		}
	}

	// 6. 插入格式配置
	formats := []model.EntityFormat{
		{
			Template: "user",
			Mode:     "summary",
			Fields:   `[{"name":"id","label":"ID","type":"string","visible":true},{"name":"username","label":"用户名","type":"string","visible":true}]`,
		},
		{
			Template: "user",
			Mode:     "detail",
			Fields:   `[{"name":"id","label":"ID","type":"string","visible":true},{"name":"username","label":"用户名","type":"string","visible":true},{"name":"roles:","label":"角色","type":"entity","ref":"role","visible":true,"fold":false},{"name":"commands","label":"操作","type":"commands","visible":true}]`,
		},
		{
			Template: "role",
			Mode:     "summary",
			Fields:   `[{"name":"id","label":"ID","type":"string","visible":true},{"name":"name","label":"名称","type":"string","visible":true},{"name":"code","label":"编码","type":"string","visible":true}]`,
		},
		{
			Template: "role",
			Mode:     "detail",
			Fields:   `[{"name":"id","label":"ID","type":"string","visible":true},{"name":"name","label":"名称","type":"string","visible":true},{"name":"code","label":"编码","type":"string","visible":true},{"name":"description","label":"描述","type":"string","visible":true},{"name":"permissions:","label":"权限","type":"entity","ref":"permission","visible":true,"fold":true},{"name":"users:","label":"用户","type":"entity","ref":"user","visible":true,"fold":true},{"name":"commands","label":"操作","type":"commands","visible":true}]`,
		},
		{
			Template: "permission",
			Mode:     "summary",
			Fields:   `[{"name":"id","label":"ID","type":"string","visible":true},{"name":"name","label":"名称","type":"string","visible":true},{"name":"code","label":"编码","type":"string","visible":true}]`,
		},
		{
			Template: "permission",
			Mode:     "detail",
			Fields:   `[{"name":"id","label":"ID","type":"string","visible":true},{"name":"name","label":"名称","type":"string","visible":true},{"name":"code","label":"编码","type":"string","visible":true},{"name":"resource","label":"资源","type":"string","visible":true},{"name":"action","label":"操作","type":"string","visible":true},{"name":"description","label":"描述","type":"string","visible":true},{"name":"roles:","label":"角色","type":"entity","ref":"role","visible":true,"fold":true},{"name":"commands","label":"操作","type":"commands","visible":true}]`,
		},
	}
	for _, f := range formats {
		if _, err := db.Engine.Insert(&f); err != nil {
			return err
		}
	}

	fmt.Println("Fixtures loaded successfully")
	return nil
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB() {
	if db.Engine != nil {
		db.Engine.Close()
	}
}

// GetTestDBPath 获取测试数据库路径（用于文件数据库）
func GetTestDBPath() string {
	return ":memory:"
}
