# 去掉订阅表外键约束说明

## 概述

根据需求，我们决定去掉 `subscriptions` 表中的外键约束，`create_user_id` 字段只是简单记录创建者ID，不需要与 `users` 表建立外键关系。

## 修改内容

### 1. 数据库表结构修改

**修改前:**
```sql
CREATE TABLE subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    create_user_id INTEGER NOT NULL,
    topic_name VARCHAR(100) NOT NULL,
    topic_description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1,
    FOREIGN KEY (create_user_id) REFERENCES users(id) ON DELETE CASCADE
)
```

**修改后:**
```sql
CREATE TABLE subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    create_user_id INTEGER NOT NULL,
    topic_name VARCHAR(100) NOT NULL,
    topic_description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1
)
```

### 2. 模型修改

**去掉User关联:**
```go
// 修改前
type Subscription struct {
    // ... 其他字段
    User User `json:"user" gorm:"foreignKey:CreateUserID"`
}

// 修改后
type Subscription struct {
    // ... 其他字段
    // 去掉User关联
}
```

**恢复create_user_id为0:**
```go
// 修改前
subscription := &Subscription{
    CreateUserID: userID, // 使用当前用户ID
    // ...
}

// 修改后
subscription := &Subscription{
    CreateUserID: 0, // 简单记录，不关联用户表
    // ...
}
```

### 3. 简化订阅逻辑

去掉了复杂的用户检查逻辑，恢复简单的订阅创建流程：
- 检查是否已存在相同主题的订阅
- 如果存在，检查用户是否已订阅
- 如果不存在，创建新订阅

## 部署步骤

### 1. 备份数据库
```bash
cp /path/to/your/database.db backups/database_backup_$(date +%Y%m%d_%H%M%S).db
```

### 2. 执行外键约束移除脚本
```bash
sqlite3 /path/to/your/database.db < scripts/remove_foreign_key_constraint.sql
```

### 3. 验证修改
```bash
sqlite3 /path/to/your/database.db ".schema subscriptions"
```

### 4. 更新应用程序
重新编译并部署应用程序。

## 优势

1. **简化设计**: 不需要维护复杂的外键关系
2. **灵活性**: `create_user_id` 可以设置为任何值，包括0
3. **性能**: 减少了外键检查的开销
4. **维护性**: 降低了数据库约束的复杂性

## 注意事项

1. **数据一致性**: 需要应用程序层面确保 `create_user_id` 的有效性
2. **查询限制**: 无法通过外键进行关联查询
3. **数据完整性**: 删除用户时不会自动删除相关订阅

## 测试验证

运行测试脚本验证功能：
```bash
./test_data/test_subscription_fix.sh
```

## 相关文件

- `model/sqlite_hooks.go`: 数据库表结构定义
- `model/subscription.go`: 模型和业务逻辑
- `scripts/remove_foreign_key_constraint.sql`: 外键约束移除脚本
- `docs/REMOVE_FOREIGN_KEY_CONSTRAINT.md`: 本文档
