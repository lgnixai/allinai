# One API 首次登录 is_first_use 字段测试

## 问题描述

在首次登录时，API 响应中 `is_first_use` 字段值为 0，但应该为 1。

## 问题原因

在 `controller/user.go` 的 `setupLogin` 函数中，`cleanUser` 结构体没有包含 `IsFirstUse` 字段，导致即使数据库中有正确的值，返回给前端的 JSON 中也不会包含这个字段。

**根本原因：**
在 `setupLogin` 函数中，返回给前端的 `cleanUser` 结构体缺少了 `IsFirstUse` 字段，导致前端无法获取到用户的首次使用状态。

## 修复方案

已在 `controller/user.go` 的 `setupLogin` 函数中添加了 `IsFirstUse` 字段：

```go
cleanUser := model.User{
    Id:          user.Id,
    Username:    user.Username,
    DisplayName: user.DisplayName,
    Role:        user.Role,
    Status:      user.Status,
    Group:       user.Group,
    School:      user.School,
    College:     user.College,
    Phone:       user.Phone,
    AccessToken: user.AccessToken,
    IsFirstUse:  user.IsFirstUse,  // 新增：包含首次使用标识
}
```

## 测试脚本

### 1. 安装依赖

```bash
npm install
```

### 2. 运行测试

#### 快速测试
```bash
npm run test-simple
# 或者
node test_quick.js
```

#### 完整测试（包含数据库检查）
```bash
npm test
# 或者
node test_first_use_complete.js
```

### 3. 测试脚本说明

- `test_quick.js`: 快速测试修复效果
- `test_first_use_login.js`: 基础登录测试
- `test_first_use_complete.js`: 完整测试，包含数据库检查和修复

## 配置说明

在运行测试前，请确保：

1. **修改 BASE_URL**: 根据实际情况修改测试脚本中的 `BASE_URL`
2. **修改测试数据**: 根据实际情况修改 `TEST_PHONE` 和 `TEST_VERIFICATION_CODE`
3. **确保服务运行**: 确保 One API 服务正在运行
4. **数据库路径**: 确保 `test_first_use_complete.js` 中的数据库路径正确

## 预期结果

修复后，登录 API 响应应该包含：

```json
{
  "success": true,
  "message": "",
  "data": {
    "id": 1,
    "username": "testuser",
    "phone": "17629726688",
    "is_first_use": 1,  // 应该是 1 而不是 0
    "role": 1,
    "status": 1,
    "school": "测试大学",
    "college": "计算机学院"
  }
}
```

## 验证步骤

1. 运行快速测试确认修复效果
2. 检查数据库中用户的 `is_first_use` 值
3. 验证 API 响应中的 `is_first_use` 字段
4. 确认数据库值和 API 返回值一致

## 测试结果

✅ **修复成功！**

- 原有用户 (17629726688): `is_first_use` 从 0 修复为 1
- 新注册用户 (13900139000): `is_first_use` 正确返回 1
- API 响应现在包含 `is_first_use` 字段
- 数据库中的值正确 (1)
- 前端可以正确获取用户的首次使用状态

## 技术细节

1. **问题定位**: 通过调试发现 `setupLogin` 函数中的 `cleanUser` 结构体缺少 `IsFirstUse` 字段
2. **修复方法**: 在 `cleanUser` 结构体中添加 `IsFirstUse: user.IsFirstUse` 字段
3. **验证方法**: 使用 curl 命令测试 API 响应，确认字段存在且值正确
4. **影响范围**: 所有用户登录 API 都会正确返回 `is_first_use` 字段
