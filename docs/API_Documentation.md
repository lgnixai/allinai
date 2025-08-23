# One-API 接口文档

## 概述

One-API 是一个统一的AI接口管理平台，提供用户管理、话题聊天、订阅管理等功能。

**服务器地址**: `http://47.88.91.79:9999`

## 认证方式

所有API请求都需要在Header中包含认证信息：

```
Authorization: Bearer <your_token>
```

## 通用响应格式

```json
{
  "success": true,
  "message": "操作成功",
  "data": {
    // 具体数据
  }
}
```

错误响应格式：

```json
{
  "success": false,
  "message": "错误信息"
}
```

## 1. 用户管理 API

### 1.1 用户登录

**接口地址**: `POST /api/user/login`

**请求参数**:
```json
{
  "username": "your_username",
  "password": "your_password"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin",
      "status": 1
    }
  }
}
```

### 1.2 获取用户信息

**接口地址**: `GET /api/user/info`

**请求头**: 需要认证token

**响应示例**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 1.3 用户注册

**接口地址**: `POST /api/user/register`

**请求参数**:
```json
{
  "username": "new_user",
  "password": "password123",
  "email": "user@example.com"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "注册成功",
  "data": {
    "id": 2,
    "username": "new_user",
    "email": "user@example.com"
  }
}
```

### 1.4 更新用户信息

**接口地址**: `PUT /api/user/self`

**请求头**: 需要认证token

**请求参数**:
```json
{
  "username": "updated_username",
  "display_name": "显示名称",
  "school": "学校名称",
  "college": "学院名称",
  "phone": "手机号码"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "用户信息更新成功"
}
```

**注意事项**:
- 只能更新自己的用户信息
- 用户名不能重复
- 手机号码需要符合格式要求

## 2. 话题管理 API

### 2.1 获取话题列表

**接口地址**: `GET /api/topics`

**请求参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**请求头**: 需要认证token

**响应示例**:
```json
{
  "success": true,
  "data": {
    "topics": [
      {
        "id": 1,
        "user_id": 1,
        "topic_name": "技术讨论",
        "model": "gpt-3.5-turbo",
        "channel_id": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "status": 1,
        "message_count": 5
      }
    ],
    "total": 1
  }
}
```

### 2.2 创建话题

**接口地址**: `POST /api/topics`

**请求头**: 需要认证token

**请求参数**:
```json
{
  "topic_name": "新话题",
  "model": "gpt-3.5-turbo",
  "channel_id": 1
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "话题创建成功",
  "data": {
    "id": 2,
    "user_id": 1,
    "topic_name": "新话题",
    "model": "gpt-3.5-turbo",
    "channel_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "status": 1
  }
}
```

### 2.3 删除话题

**接口地址**: `DELETE /api/topics/{id}`

**请求头**: 需要认证token

**响应示例**:
```json
{
  "success": true,
  "message": "话题已删除"
}
```

### 2.4 获取话题消息

**接口地址**: `GET /api/topics/{id}/messages`

**请求参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**请求头**: 需要认证token

**响应示例**:
```json
{
  "success": true,
  "data": {
    "messages": [
      {
        "id": 1,
        "topic_id": 1,
        "role": "user",
        "content": "你好",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "status": 1
      },
      {
        "id": 2,
        "topic_id": 1,
        "role": "assistant",
        "content": "\"技术讨论\": 你好！有什么可以帮助您的吗？",
        "created_at": "2024-01-01T00:00:01Z",
        "updated_at": "2024-01-01T00:00:01Z",
        "status": 1
      }
    ],
    "total": 2
  }
}
```

### 2.5 发送消息

**接口地址**: `POST /api/topics/{id}/messages`

**请求头**: 需要认证token

**请求参数**:
```json
{
  "content": "用户消息内容",
  "role": "user"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "消息发送成功",
  "data": {
    "user_message": {
      "id": 3,
      "topic_id": 1,
      "role": "user",
      "content": "用户消息内容",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "ai_message": {
      "id": 4,
      "topic_id": 1,
      "role": "assistant",
      "content": "\"技术讨论\": 这是一个很好的问题！让我来为您详细解答。",
      "created_at": "2024-01-01T00:00:01Z"
    }
  }
}
```

## 3. 订阅管理 API

### 3.1 获取订阅列表

**接口地址**: `GET /api/subscriptions`

**请求参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**请求头**: 需要认证token

**响应示例**:
```json
{
  "success": true,
  "data": {
    "subscriptions": [
      {
        "id": 1,
        "user_id": 1,
        "topic_name": "技术新闻",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "article_count": 10
      }
    ],
    "total": 1
  }
}
```

### 3.2 创建订阅

**接口地址**: `POST /api/subscriptions`

**请求头**: 需要认证token

**请求参数**:
```json
{
  "topic_name": "新订阅主题",
  "model": "gpt-3.5-turbo",
  "channel_id": 1
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "订阅创建成功",
  "data": {
    "id": 2,
    "user_id": 1,
    "topic_name": "新订阅主题",
    "model": "gpt-3.5-turbo",
    "channel_id": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "status": 1
  }
}
```

### 3.3 取消订阅

**接口地址**: `PUT /api/subscriptions/{id}/cancel`

**请求头**: 需要认证token

**响应示例**:
```json
{
  "success": true,
  "message": "订阅已取消"
}
```

### 3.4 重新激活订阅

**接口地址**: `PUT /api/subscriptions/{id}/reactivate`

**请求头**: 需要认证token

**响应示例**:
```json
{
  "success": true,
  "message": "订阅重新激活成功"
}
```

### 3.5 删除订阅

**接口地址**: `DELETE /api/subscriptions/{id}`

**请求头**: 需要认证token

**响应示例**:
```json
{
  "success": true,
  "message": "订阅已删除"
}
```

### 3.6 获取订阅文章

**接口地址**: `GET /api/subscriptions/{id}/articles`

**请求参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**请求头**: 需要认证token

**响应示例**:
```json
{
  "success": true,
  "data": {
    "articles": [
      {
        "id": 1,
        "subscription_id": 1,
        "title": "文章标题",
        "content": "文章内容",
        "url": "https://example.com/article",
        "published_at": "2024-01-01T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z",
        "status": 1
      }
    ],
    "total": 1
  }
}
```

## 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或token无效 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 注意事项

1. 所有需要认证的接口都必须在请求头中包含有效的token
2. 分页参数page从1开始计数
3. 时间格式统一使用ISO 8601格式
4. 删除操作通常为软删除，不会真正删除数据
5. 话题和订阅的状态：1表示正常，0表示已删除/取消
