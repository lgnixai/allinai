# 用户更新接口修复说明

## 问题描述

在使用 Postman 测试更新用户信息接口 `PUT /api/user/self` 时，遇到以下错误：

```json
{
  "message": "输入不合法 Key: 'User.Phone' Error:Field validation for 'Phone' failed on the 'required' tag",
  "success": false
}
```

## 问题原因

1. **结构体验证问题**: `User` 结构体中的 `Phone` 字段有 `validate:"required,len=11"` 标签，表示该字段是必填的且必须是11位。

2. **更新逻辑问题**: 原来的 `UpdateSelf` 函数会创建一个新的 `User` 结构体，并将请求中的所有字段（包括空的 `Phone` 字段）都赋值给这个结构体，然后进行验证。

3. **部分更新需求**: 用户可能只想更新部分字段（如 `display_name`、`school`、`college`），而不想更新手机号，但原来的逻辑要求所有字段都必须提供。

## 解决方案

### 1. 修改更新逻辑

将 `UpdateSelf` 函数从创建新结构体的方式改为只更新用户提供的字段：

```go
func UpdateSelf(c *gin.Context) {
    var user model.User
    err := json.NewDecoder(c.Request.Body).Decode(&user)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{
            "success": false,
            "message": "无效的参数",
        })
        return
    }

    // 获取当前用户信息
    currentUser, err := model.GetUserById(c.GetInt("id"), false)
    if err != nil {
        common.ApiError(c, err)
        return
    }

    // 只更新用户提供的字段
    updates := make(map[string]interface{})
    
    if user.Username != "" {
        updates["username"] = user.Username
    }
    if user.DisplayName != "" {
        updates["display_name"] = user.DisplayName
    }
    if user.School != "" {
        updates["school"] = user.School
    }
    if user.College != "" {
        updates["college"] = user.College
    }
    if user.Phone != "" {
        // 如果提供了手机号，需要验证格式
        if len(user.Phone) != 11 {
            c.JSON(http.StatusOK, gin.H{
                "success": false,
                "message": "手机号必须是11位数字",
            })
            return
        }
        updates["phone"] = user.Phone
    }

    // 如果没有提供任何要更新的字段
    if len(updates) == 0 {
        c.JSON(http.StatusOK, gin.H{
            "success": false,
            "message": "请提供要更新的字段",
        })
        return
    }

    // 更新用户信息
    if err := model.DB.Model(&currentUser).Updates(updates).Error; err != nil {
        common.ApiError(c, err)
        return
    }

    // 清除缓存，下次获取时会重新加载
    if err := model.InvalidateUserCache(currentUser.Id); err != nil {
        common.SysError("failed to invalidate user cache: " + err.Error())
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "",
    })
    return
}
```

### 2. 添加缓存管理

在 `model/user_cache.go` 中添加了公共的缓存清除函数：

```go
// InvalidateUserCache clears user cache (public function)
func InvalidateUserCache(userId int) error {
    return invalidateUserCache(userId)
}
```

## 修复效果

### 修复前
- ❌ 必须提供所有字段，包括手机号
- ❌ 无法进行部分字段更新
- ❌ 空字段会导致验证失败

### 修复后
- ✅ 支持部分字段更新
- ✅ 不提供手机号也能更新其他字段
- ✅ 提供手机号时会验证格式（11位）
- ✅ 空请求会返回错误提示
- ✅ 缓存管理正确

## 测试用例

### 1. 更新部分字段（不包含手机号）
```json
PUT /api/user/self
{
  "display_name": "新显示名称",
  "school": "新学校",
  "college": "新学院"
}
```
**预期结果**: ✅ 成功

### 2. 更新包含手机号
```json
PUT /api/user/self
{
  "display_name": "新显示名称",
  "phone": "13800138001"
}
```
**预期结果**: ✅ 成功

### 3. 无效手机号
```json
PUT /api/user/self
{
  "display_name": "新显示名称",
  "phone": "123"
}
```
**预期结果**: ❌ 返回"手机号必须是11位数字"

### 4. 空请求
```json
PUT /api/user/self
{}
```
**预期结果**: ❌ 返回"请提供要更新的字段"

## 使用说明

现在用户可以通过以下方式更新信息：

1. **只更新显示名称**:
   ```json
   {
     "display_name": "新名称"
   }
   ```

2. **只更新学校和学院**:
   ```json
   {
     "school": "新学校",
     "college": "新学院"
   }
   ```

3. **更新手机号**:
   ```json
   {
     "phone": "13800138001"
   }
   ```

4. **组合更新**:
   ```json
   {
     "display_name": "新名称",
     "school": "新学校",
     "phone": "13800138001"
   }
   ```

## 注意事项

1. 手机号必须是11位数字
2. 至少需要提供一个要更新的字段
3. 用户名、显示名称、学校、学院字段不能为空字符串
4. 更新后缓存会自动清除，确保数据一致性
