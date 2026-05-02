-- ============================================
-- 初始化数据脚本
-- 手动运行: sqlite3 iam.db < migrations/init.sql
-- ============================================

-- 1. 插入 admin 角色
INSERT OR IGNORE INTO role (i_d, name, code, description, created_at, updated_at)
VALUES (1, '管理员', 'admin', '系统管理员，拥有所有权限', datetime('now'), datetime('now'));

-- 2. 插入默认权限
INSERT OR IGNORE INTO permission (i_d, name, code, resource, action, description, created_at, updated_at) VALUES
(1, '创建用户', 'user:create', 'user', 'create', '允许创建新用户', datetime('now'), datetime('now')),
(2, '读取用户', 'user:read', 'user', 'read', '允许查看用户信息', datetime('now'), datetime('now')),
(3, '更新用户', 'user:update', 'user', 'update', '允许更新用户信息', datetime('now'), datetime('now')),
(4, '删除用户', 'user:delete', 'user', 'delete', '允许删除用户', datetime('now'), datetime('now')),
(5, '管理角色', 'role:manage', 'role', 'manage', '允许管理角色', datetime('now'), datetime('now')),
(6, '管理权限', 'permission:manage', 'permission', 'manage', '允许管理权限', datetime('now'), datetime('now'));


-- 4. 关联 admin 用户与 admin 角色
INSERT OR IGNORE INTO user_role (user_id, role_id, created_at)
VALUES (1, 1, datetime('now'));

-- 5. 关联 admin 角色与所有权限
INSERT OR IGNORE INTO role_permission (role_id, permission_id, created_at)
SELECT 1, i_d, datetime('now') FROM permission;

-- 6. 插入实体格式配置
-- user summary
INSERT OR IGNORE INTO entity_format (template, mode, fields, created_at, updated_at)
VALUES ('user', 'summary', '[
  {"name":"id","label":"ID","type":"string","visible":true},
  {"name":"username","label":"用户名","type":"string","visible":true}
]', datetime('now'), datetime('now'));

-- user detail
INSERT OR IGNORE INTO entity_format (template, mode, fields, created_at, updated_at)
VALUES ('user', 'detail', '[
  {"name":"id","label":"ID","type":"string","visible":true},
  {"name":"username","label":"用户名","type":"string","visible":true},
  {"name":"roles:","label":"角色","type":"entity","ref":"role","visible":true,"fold":false},
  {"name":"commands","label":"操作","type":"commands","visible":true}
]', datetime('now'), datetime('now'));

-- role summary
INSERT OR IGNORE INTO entity_format (template, mode, fields, created_at, updated_at)
VALUES ('role', 'summary', '[
  {"name":"id","label":"ID","type":"string","visible":true},
  {"name":"name","label":"名称","type":"string","visible":true},
  {"name":"code","label":"编码","type":"string","visible":true}
]', datetime('now'), datetime('now'));

-- role detail
INSERT OR IGNORE INTO entity_format (template, mode, fields, created_at, updated_at)
VALUES ('role', 'detail', '[
  {"name":"id","label":"ID","type":"string","visible":true},
  {"name":"name","label":"名称","type":"string","visible":true},
  {"name":"code","label":"编码","type":"string","visible":true},
  {"name":"description","label":"描述","type":"string","visible":true},
  {"name":"permissions:","label":"权限","type":"entity","ref":"permission","visible":true,"fold":true},
  {"name":"users:","label":"用户","type":"entity","ref":"user","visible":true,"fold":true},
  {"name":"commands","label":"操作","type":"commands","visible":true}
]', datetime('now'), datetime('now'));

-- permission summary
INSERT OR IGNORE INTO entity_format (template, mode, fields, created_at, updated_at)
VALUES ('permission', 'summary', '[
  {"name":"id","label":"ID","type":"string","visible":true},
  {"name":"name","label":"名称","type":"string","visible":true},
  {"name":"code","label":"编码","type":"string","visible":true}
]', datetime('now'), datetime('now'));

-- permission detail
INSERT OR IGNORE INTO entity_format (template, mode, fields, created_at, updated_at)
VALUES ('permission', 'detail', '[
  {"name":"id","label":"ID","type":"string","visible":true},
  {"name":"name","label":"名称","type":"string","visible":true},
  {"name":"code","label":"编码","type":"string","visible":true},
  {"name":"resource","label":"资源","type":"string","visible":true},
  {"name":"action","label":"操作","type":"string","visible":true},
  {"name":"description","label":"描述","type":"string","visible":true},
  {"name":"roles:","label":"角色","type":"entity","ref":"role","visible":true,"fold":true},
  {"name":"commands","label":"操作","type":"commands","visible":true}
]', datetime('now'), datetime('now'));

insert or ignore into user (i_d, username, password_hash, created_at, updated_at)
values (1, 'admin', '$2a$10$pk9Ss9Bpv/0uZ7Uz3aCf8OXkagQlCMq7uk5lZfQekONb2F8NShfY.',   datetime('now'), datetime('now'));