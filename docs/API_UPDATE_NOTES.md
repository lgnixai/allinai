# API 更新说明

## 更新概述

本次更新主要包含以下内容：
1. 用户授权失效处理优化
2. 新增系统推荐相关接口
3. 完善订阅管理功能
4. 增加文章概要字段
5. 新增管理员功能接口

## 详细更新内容

### 1. 用户授权失效处理

#### 更新内容
- 明确 `Authorization` 头部格式为 `Bearer <your_access_token>`
- 添加401状态码的前端处理说明
- 完善认证相关错误信息分类

#### 前端处理要求
当收到401错误时，前端应：
1. 清除本地存储的token和user_id
2. 跳转到登录页面或显示401状态
3. 提示用户重新登录

#### 相关文件
- `docs/API_Documentation.md`
- `docs/html/API_Documentation.md`

### 2. 系统推荐接口

#### 新增接口
1. **获取系统推荐列表**
   - `GET /api/system-recommendations`
   - 支持分页查询

2. **获取欢迎页面**
   - `GET /api/system-recommendations/welcome`
   - 返回个性化欢迎消息和推荐内容

3. **获取推荐页面**
   - `GET /api/system-recommendations/recommendation`
   - 返回后续访问的推荐内容

4. **搜索系统推荐**
   - `GET /api/system-recommendations/search`
   - 支持关键字搜索

#### 响应格式
```json
{
  "success": true,
  "data": {
    "recommendations": [
      {
        "id": 1,
        "title": "人工智能与机器学习",
        "description": "探索AI和ML的最新发展...",
        "category": "技术",
        "subscription_count": 0,
        "article_count": 0,
        "status": 1,
        "sort_order": 100,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

### 3. 订阅管理功能完善

#### 新增接口
1. **获取所有订阅文章**
   - `GET /api/subscriptions/articles`
   - 获取当前用户所有订阅的文章列表

2. **更新订阅**
   - `PUT /api/subscriptions/{id}`
   - 更新订阅描述

3. **删除订阅**
   - `DELETE /api/subscriptions/{id}`
   - 永久删除订阅

4. **重新激活订阅**
   - `PUT /api/subscriptions/{id}/reactivate`
   - 重新激活已取消的订阅

#### 文章概要字段
在订阅文章响应中新增 `summary` 字段：
```json
{
  "id": 1,
  "subscription_id": 1,
  "title": "文章标题",
  "summary": "文章概要...",
  "content": "文章内容...",
  "author": "作者",
  "published_at": "2024-01-01T00:00:00Z",
  "article_url": "https://example.com/article",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "status": 1
}
```

### 4. 管理员功能接口

#### 新增接口
1. **创建系统推荐**
   - `POST /api/system-recommendations`
   - 需要管理员权限

2. **更新系统推荐**
   - `PUT /api/system-recommendations/{id}`
   - 需要管理员权限

3. **删除系统推荐**
   - `DELETE /api/system-recommendations/{id}`
   - 需要管理员权限

4. **创建订阅文章**
   - `POST /api/subscriptions/{id}/articles`
   - 需要管理员权限

#### 权限要求
- 所有管理员接口都需要管理员权限
- 普通用户调用会返回403错误

### 5. 错误处理优化

#### 错误信息分类
1. **认证相关错误**
   - 未提供访问令牌
   - 访问令牌无效或已过期
   - UserID格式错误
   - UserID与访问令牌不匹配
   - 用户被封禁

2. **权限相关错误**
   - 无权限操作话题
   - 无权限修改订阅
   - 无权限查看文章
   - 无权限创建系统推荐（需要管理员权限）

3. **资源相关错误**
   - 话题不存在
   - 订阅不存在
   - 推荐不存在
   - 已订阅该主题

#### 错误码处理
| 错误码 | 说明 | 前端处理 |
|--------|------|----------|
| 400 | 请求参数错误 | 显示错误信息 |
| 401 | 未授权，认证失败 | 跳转到登录页面或显示401状态 |
| 403 | 权限不足 | 显示权限不足提示 |
| 404 | 资源不存在 | 显示404页面 |
| 500 | 服务器内部错误 | 显示服务器错误提示 |

## 文档更新

### 更新的文档
1. `docs/API_Documentation.md` - 主API文档
2. `docs/html/API_Documentation.md` - HTML版本API文档
3. `docs/Postman_Collection.json` - Postman集合文件

### 新增文档
1. `docs/API_UPDATE_NOTES.md` - 本更新说明文档

## 使用说明

### 认证流程
1. 用户登录后获取 `access_token` 和 `user_id`
2. 在请求头中设置：
   ```
   Authorization: Bearer <access_token>
   UserID: <user_id>
   ```
3. 当收到401错误时，清除本地存储并跳转登录页面

### 权限管理
1. 普通用户只能操作自己的话题和订阅
2. 管理员可以创建和管理系统推荐
3. 管理员可以创建订阅文章

### 测试数据
- 测试环境验证码固定为"1111"
- 可以使用任意11位手机号进行测试

## 兼容性说明

- 所有现有接口保持向后兼容
- 新增字段为可选字段，不影响现有功能
- 认证方式保持不变，只是明确了格式要求

## 后续计划

1. 添加更多系统推荐分类
2. 支持文章标签和分类
3. 增加用户偏好设置
4. 支持文章收藏功能
5. 添加消息推送功能

