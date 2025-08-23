# 用户表字段更新说明

## 概述
本次更新为用户表添加了三个新字段：学校、学院、手机号，以支持更详细的用户信息管理。

## 新增字段

### 1. 学校字段 (school)
- **类型**: VARCHAR(100)
- **默认值**: 空字符串
- **说明**: 存储用户所在学校名称
- **验证**: 最大长度100字符

### 2. 学院字段 (college)
- **类型**: VARCHAR(100)
- **默认值**: 空字符串
- **说明**: 存储用户所在学院名称
- **验证**: 最大长度100字符

### 3. 手机号字段 (phone)
- **类型**: VARCHAR(20)
- **默认值**: 空字符串
- **说明**: 存储用户手机号码
- **验证**: 最大长度20字符
- **索引**: 已添加索引以提高查询性能

## 数据库迁移

### 迁移文件
- 文件路径: `bin/migration_v0.4-v0.5.sql`
- 执行命令: 根据您的数据库类型执行相应的SQL命令

### MySQL示例
```sql
-- 添加用户表新字段：学校、学院、手机号
ALTER TABLE users ADD COLUMN school VARCHAR(100) DEFAULT '' COMMENT '学校';
ALTER TABLE users ADD COLUMN college VARCHAR(100) DEFAULT '' COMMENT '学院';
ALTER TABLE users ADD COLUMN phone VARCHAR(20) DEFAULT '' COMMENT '手机号';

-- 为手机号字段添加索引以提高查询性能
CREATE INDEX idx_users_phone ON users(phone);
```

## 功能更新

### 1. 用户注册
- 支持在注册时填写学校、学院、手机号信息
- 支持手机号验证码验证（可选）
- 手机号验证码通过 `/api/phone_verification` 接口获取

### 2. 用户信息更新
- 用户可以通过 `/api/user/self` 接口更新个人信息
- 支持更新学校、学院、手机号字段

### 3. 手机号验证
- 新增手机号验证功能
- 验证码有效期：10分钟
- 接口路径：`GET /api/phone_verification?phone=手机号`

## API接口更新

### 新增接口
- `GET /api/phone_verification` - 发送手机验证码

### 更新接口
- `POST /api/user/register` - 支持手机号相关字段
- `PUT /api/user/self` - 支持更新新字段
- `GET /api/user/self` - 返回新字段信息

## 前端更新建议

### 注册表单
- 添加学校、学院、手机号输入框
- 添加手机号验证码输入框和发送按钮
- 手机号格式验证（11位数字）

### 个人设置页面
- 添加学校、学院、手机号编辑功能
- 支持手机号验证码绑定

## 测试

### 数据库测试
运行 `test_migration.sql` 文件来验证数据库字段是否正确添加：

```sql
-- 检查表结构
DESCRIBE users;

-- 检查新字段是否存在
SELECT 
    COLUMN_NAME, 
    DATA_TYPE, 
    IS_NULLABLE, 
    COLUMN_DEFAULT, 
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'users' 
    AND COLUMN_NAME IN ('school', 'college', 'phone');
```

### API测试
1. 测试注册接口，包含新字段
2. 测试手机号验证码发送
3. 测试用户信息更新
4. 测试用户信息获取

## 注意事项

1. **向后兼容**: 新字段都有默认值，不会影响现有用户
2. **手机号验证**: 目前验证码会直接返回，生产环境需要集成短信服务
3. **数据验证**: 前端需要添加相应的数据验证
4. **权限控制**: 新字段的访问权限遵循现有规则

## 后续优化建议

1. 集成短信服务提供商（如阿里云、腾讯云等）
2. 添加学校、学院的预设选项
3. 实现手机号登录功能
4. 添加手机号格式的国际支持
