# API 文档更新日志

## 2025-08-26 - 添加 is_first_use 字段

### 更新内容

在用户登录和获取用户信息 API 的响应中添加了 `is_first_use` 字段。

### 更新的 API

1. **用户登录** (`POST /api/user/login`)
   - 在响应数据中添加了 `is_first_use` 字段
   - 添加了详细的响应字段说明

2. **获取用户信息** (`GET /api/user/self`)
   - 在响应数据中添加了 `is_first_use` 字段

### 字段说明

- `is_first_use`: 首次使用标识
  - `1`: 首次使用
  - `0`: 非首次使用

### 更新前的响应示例

```json
{
  "success": true,
  "message": "登录成功",
  "data": {
    "id": 1,
    "username": "user_8000",
    "display_name": "测试用户",
    "role": 1,
    "status": 1,
    "group": "default",
    "school": "测试大学",
    "college": "计算机学院",
    "phone": "13800138000",
    "access_token": "your_access_token_here"
  }
}
```

### 更新后的响应示例

```json
{
  "success": true,
  "message": "登录成功",
  "data": {
    "id": 1,
    "username": "user_8000",
    "display_name": "测试用户",
    "role": 1,
    "status": 1,
    "group": "default",
    "school": "测试大学",
    "college": "计算机学院",
    "phone": "13800138000",
    "access_token": "your_access_token_here",
    "is_first_use": 1
  }
}
```

### 影响范围

- 前端可以根据 `is_first_use` 字段判断用户是否为首次使用
- 可以用于显示新手引导、欢迎页面等功能
- 向后兼容，不影响现有功能

### 技术背景

此更新是为了修复首次登录时 `is_first_use` 字段缺失的问题。在 `controller/user.go` 的 `setupLogin` 函数中添加了 `IsFirstUse: user.IsFirstUse` 字段，确保 API 响应包含用户的首次使用状态。
