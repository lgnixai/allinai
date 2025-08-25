# API 接口文档

## 概述

API 是一个统一的AI接口管理平台，提供用户管理、话题聊天、订阅管理等功能。

**服务器地址**: `http://47.88.91.79:9999`

## 认证方式

所有需要认证的API请求都需要在Header中包含以下两个字段：

```
Authorization: Bearer <your_access_token>
UserID: <your_user_id>
```

**注意**：
- 必须同时提供 `Authorization` 和 `UserID` 两个请求头
- `Authorization` 字段格式为 `Bearer <your_access_token>`
- `UserID` 字段包含用户的ID
- 如果任一字段缺失或无效，将返回 `401 Unauthorized` 错误
- 当认证失效时，前端应跳转到登录页面或显示401状态

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

### 1.1 发送手机验证码

**接口地址**: `GET /api/phone_verification`

**请求参数**:
- `phone`: 手机号（11位数字）
- `purpose`: 用途（register-注册，login-登录）

**请求示例**:
```
GET /api/phone_verification?phone=13800138000&purpose=register
```

**响应示例**:
```json
{
  "success": true,
  "message": "",
  "data": "1111"
}
```

### 1.2 用户注册

**接口地址**: `POST /api/user/register`

**请求参数**:
```json
{
  "phone": "13800138000",
  "phone_verification_code": "1111",
  "display_name": "测试用户",
  "school": "测试大学",
  "college": "计算机学院"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "注册成功"
}
```

### 1.3 用户登录

**接口地址**: `POST /api/user/login`

**请求参数**:
```json
{
  "phone": "13800138000",
  "phone_verification_code": "1111"
}
```

**响应示例**:
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

### 1.4 获取用户信息

**接口地址**: `GET /api/user/self`

**请求头**: 需要认证token和UserID

**响应示例**:
```json
{
  "success": true,
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

### 1.5 更新用户信息

**接口地址**: `PUT /api/user/self`

**请求头**: 需要认证token和UserID

**请求参数**:
```json
{
  "display_name": "更新后的显示名称",
  "school": "更新后的学校",
  "college": "更新后的学院"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "更新成功"
}
```

### 1.6 用户登出

**接口地址**: `GET /api/user/logout`

**请求头**: 需要认证token和UserID

**响应示例**:
```json
{
  "success": true,
  "message": "登出成功"
}
```

## 2. 话题管理 API

### 2.1 获取话题列表

**接口地址**: `GET /api/topics/`

**请求头**: 需要认证token和UserID

**请求参数**:
- `page`: 页码（默认1）
- `size`: 每页数量（默认10）

**响应示例**:
```json
{
  "success": true,
  "data": {
    "topics": [
      {
        "id": 1,
        "topic_name": "测试话题",
        "model": "gpt-3.5-turbo",
        "channel_id": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "message_count": 5
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

### 2.2 创建话题

**接口地址**: `POST /api/topics/`

**请求头**: 需要认证token和UserID

**请求参数**:
```json
{
  "topic_name": "测试话题",
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
    "id": 1,
    "topic_name": "测试话题",
    "model": "gpt-3.5-turbo",
    "channel_id": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 2.3 删除话题

**接口地址**: `DELETE /api/topics/{id}`

**请求头**: 需要认证token和UserID

**响应示例**:
```json
{
  "success": true,
  "message": "话题已删除"
}
```

### 2.4 获取话题消息

**接口地址**: `GET /api/topics/{id}/messages`

**请求头**: 需要认证token和UserID

**请求参数**:
- `page`: 页码（默认1）
- `size`: 每页数量（默认20）

**响应示例**:
```json
{
  "success": true,
  "data": {
    "messages": [
      {
        "id": 1,
        "role": "user",
        "content": "你好",
        "created_at": "2024-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "role": "assistant",
        "content": "\"测试话题\": 你好！我是AI助手，有什么可以帮助你的吗？",
        "created_at": "2024-01-01T00:01:00Z"
      }
    ],
    "total": 2,
    "page": 1,
    "size": 20
  }
}
```

### 2.5 发送消息

**接口地址**: `POST /api/topics/{id}/messages`

**请求头**: 需要认证token和UserID

**说明**: 
- 当 `{id}` 为 0 时，系统会自动创建新话题，话题标题为用户输入内容的前10个字符
- 当 `{id}` 为有效话题ID时，在指定话题中发送消息

**请求参数**:
```json
{
  "content": "你好，这是一个测试消息"
}
```

**响应示例** (发送到现有话题):
```json
{
  "success": true,
  "message": "消息发送成功",
  "data": {
    "user_message": {
      "id": 3,
      "role": "user",
      "content": "你好，这是一个测试消息",
      "created_at": "2024-01-01T00:02:00Z"
    },
    "ai_message": {
      "id": 4,
      "role": "assistant",
      "content": "\"测试话题\": 你好！我收到了你的测试消息，有什么可以帮助你的吗？",
      "created_at": "2024-01-01T00:02:01Z"
    }
  }
}
```

**响应示例** (自动创建话题):
```json
{
  "success": true,
  "message": "消息发送成功",
  "data": {
    "user_message": {
      "id": 3,
      "role": "user",
      "content": "你好，这是一个测试消息",
      "created_at": "2024-01-01T00:02:00Z"
    },
    "ai_message": {
      "id": 4,
      "role": "assistant",
      "content": "\"你好，这是一个\": 你好！我收到了你的测试消息，有什么可以帮助你的吗？",
      "created_at": "2024-01-01T00:02:01Z"
    },
    "topic": {
      "id": 5,
      "user_id": 1,
      "topic_name": "你好，这是一个",
      "model": "gpt-3.5-turbo",
      "channel_id": 1,
      "status": 1,
      "created_at": "2024-01-01T00:02:00Z",
      "updated_at": "2024-01-01T00:02:00Z"
    }
  }
}
```

## 3. 订阅管理 API

### 3.1 获取订阅列表

**接口地址**: `GET /api/subscriptions/`

**请求头**: 需要认证token和UserID

**请求参数**:
- `page`: 页码（默认1）
- `size`: 每页数量（默认10）

**响应示例**:
```json
{
  "success": true,
  "data": {
    "subscriptions": [
      {
        "id": 1,
        "topic_name": "技术订阅",
        "topic_description": "技术相关文章订阅",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

### 3.2 创建订阅

**接口地址**: `POST /api/subscriptions/`

**请求头**: 需要认证token和UserID

**请求参数**:
```json
{
  "topic_name": "技术订阅",
  "topic_description": "技术相关文章订阅"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "订阅创建成功",
  "data": {
    "id": 1,
    "topic_name": "技术订阅",
    "topic_description": "技术相关文章订阅",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3.3 取消订阅

**接口地址**: `PUT /api/subscriptions/{id}/cancel`

**请求头**: 需要认证token和UserID

**响应示例**:
```json
{
  "success": true,
  "message": "订阅已取消"
}
```
 

### 3.6 获取订阅文章

**接口地址**: `GET /api/subscriptions/{id}/articles`

**请求头**: 需要认证token和UserID

**请求参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**响应示例**:
```json
{
  "success": true,
  "data": {
    "articles": [
      {
        "id": 1,
        "subscription_id": 1,
        "title": "示例文章标题",
        "summary": "文章概要...",
        "content": "文章内容...",
        "author": "作者",
        "published_at": "2024-01-01T00:00:00Z",
        "article_url": "https://example.com/article",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "status": 1
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

### 3.7 获取所有订阅文章

**接口地址**: `GET /api/subscriptions/articles`

**请求头**: 需要认证token和UserID

**请求参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**说明**: 获取当前用户所有订阅的文章列表

**响应示例**:
```json
{
  "success": true,
  "data": {
    "articles": [
      {
        "id": 1,
        "subscription_id": 1,
        "title": "示例文章标题",
        "summary": "文章概要...",
        "content": "文章内容...",
        "author": "作者",
        "published_at": "2024-01-01T00:00:00Z",
        "article_url": "https://example.com/article",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z",
        "status": 1
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

### 3.8 更新订阅

**接口地址**: `PUT /api/subscriptions/{id}`

**请求头**: 需要认证token和UserID

**请求参数**:
```json
{
  "topic_description": "更新后的订阅描述"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "订阅更新成功"
}
```

### 3.9 删除订阅

**接口地址**: `DELETE /api/subscriptions/{id}`

**请求头**: 需要认证token和UserID

**响应示例**:
```json
{
  "success": true,
  "message": "订阅已删除"
}
```

### 3.10 重新激活订阅

**接口地址**: `PUT /api/subscriptions/{id}/reactivate`

**请求头**: 需要认证token和UserID

**响应示例**:
```json
{
  "success": true,
  "message": "订阅重新激活成功"
}
```

## 4. 系统推荐 API

### 4.1 获取系统推荐列表

**接口地址**: `GET /api/system-recommendations`

**请求头**: 需要认证token和UserID

**请求参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**响应示例**:
```json
{
  "success": true,
  "data": {
    "recommendations": [
      {
        "id": 1,
        "title": "人工智能与机器学习",
        "description": "探索AI和ML的最新发展，包括深度学习、自然语言处理、计算机视觉等前沿技术",
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

### 4.2 获取欢迎页面

**接口地址**: `GET /api/system-recommendations/welcome`

**请求头**: 需要认证token和UserID

**说明**: 获取首次访问的欢迎页面，包含个性化欢迎消息和推荐内容

**响应示例**:
```json
{
  "success": true,
  "data": {
    "welcome_message": "Hi, 张三,我是 Moyo 安排给你的科研合伙人, 我叫IU。今天是咱们俩第一次见面,为了可以更好的开展后面的工作,给你初步介绍下我现在可以做的事情。因为还不知道你想让我做什么,我根据你的专业帮你选择了几个可能感兴趣的话题。",
    "recommendations": [
      {
        "id": 1,
        "title": "人工智能与机器学习",
        "description": "探索AI和ML的最新发展，包括深度学习、自然语言处理、计算机视觉等前沿技术",
        "category": "技术",
        "subscription_count": 0,
        "article_count": 0,
        "status": 1,
        "sort_order": 100,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 4.3 获取推荐页面

**接口地址**: `GET /api/system-recommendations/recommendation`

**请求头**: 需要认证token和UserID

**说明**: 获取后续访问的推荐页面

**响应示例**:
```json
{
  "success": true,
  "data": {
    "recommendations": [
      {
        "id": 1,
        "title": "人工智能与机器学习",
        "description": "探索AI和ML的最新发展，包括深度学习、自然语言处理、计算机视觉等前沿技术",
        "category": "技术",
        "subscription_count": 0,
        "article_count": 0,
        "status": 1,
        "sort_order": 100,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 4.4 搜索系统推荐

**接口地址**: `GET /api/system-recommendations/search`

**请求头**: 需要认证token和UserID

**请求参数**:
- `keyword`: 搜索关键字（必填）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**响应示例**:
```json
{
  "success": true,
  "data": {
    "recommendations": [
      {
        "id": 1,
        "title": "人工智能与机器学习",
        "description": "探索AI和ML的最新发展，包括深度学习、自然语言处理、计算机视觉等前沿技术",
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
  },
  "keyword": "人工智能"
}
```

## 5. 管理员功能 API

### 5.1 创建系统推荐（管理员）

**接口地址**: `POST /api/system-recommendations`

**请求头**: 需要认证token和UserID（管理员权限）

**请求参数**:
```json
{
  "title": "新推荐标题",
  "description": "推荐描述",
  "category": "技术",
  "sort_order": 100
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "推荐创建成功",
  "data": {
    "id": 1
  }
}
```

### 5.2 更新系统推荐（管理员）

**接口地址**: `PUT /api/system-recommendations/{id}`

**请求头**: 需要认证token和UserID（管理员权限）

**请求参数**:
```json
{
  "title": "更新后的标题",
  "description": "更新后的描述",
  "category": "技术",
  "sort_order": 100
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "推荐更新成功"
}
```

### 5.3 删除系统推荐（管理员）

**接口地址**: `DELETE /api/system-recommendations/{id}`

**请求头**: 需要认证token和UserID（管理员权限）

**响应示例**:
```json
{
  "success": true,
  "message": "推荐删除成功"
}
```

### 5.4 创建订阅文章（管理员）

**接口地址**: `POST /api/subscriptions/{id}/articles`

**请求头**: 需要认证token和UserID（管理员权限）

**请求参数**:
```json
{
  "title": "文章标题",
  "summary": "文章概要",
  "content": "文章内容",
  "author": "作者",
  "published_at": "2024-01-01T00:00:00Z",
  "article_url": "https://example.com/article"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "文章创建成功",
  "data": {
    "id": 1
  }
}
```

## 错误码说明

| 错误码 | 说明 | 前端处理 |
|--------|------|----------|
| 400 | 请求参数错误 | 显示错误信息 |
| 401 | 未授权，认证失败 | 跳转到登录页面或显示401状态 |
| 403 | 权限不足 | 显示权限不足提示 |
| 404 | 资源不存在 | 显示404页面 |
| 500 | 服务器内部错误 | 显示服务器错误提示 |

## 常见错误信息

### 认证相关错误
- `"无权进行此操作，未登录且未提供 access token"`: 未提供访问令牌
- `"无权进行此操作，access token 无效"`: 访问令牌无效或已过期
- `"无权进行此操作，未提供 UserID"`: 缺少UserID请求头
- `"无权进行此操作，UserID 格式错误"`: UserID格式不正确
- `"无权进行此操作，UserID 与登录用户不匹配"`: UserID与访问令牌不匹配
- `"用户已被封禁"`: 用户账户被禁用

### 用户管理错误
- `"手机号格式错误"`: 手机号必须是11位数字
- `"手机验证码错误或已过期"`: 验证码不正确或已过期
- `"手机号未注册"`: 登录时使用的手机号未注册
- `"手机号已被占用"`: 注册时使用的手机号已存在

### 权限相关错误
- `"无权限操作此话题"`: 用户无权操作该话题
- `"无权限修改此订阅"`: 用户无权修改该订阅
- `"无权限查看此订阅的文章"`: 用户无权查看该订阅的文章
- `"无权限创建文章"`: 用户无权创建文章（需要管理员权限）
- `"无权限创建系统推荐"`: 用户无权创建系统推荐（需要管理员权限）
- `"无权限更新系统推荐"`: 用户无权更新系统推荐（需要管理员权限）
- `"无权限删除系统推荐"`: 用户无权删除系统推荐（需要管理员权限）

### 资源相关错误
- `"话题不存在"`: 指定的话题ID不存在
- `"订阅不存在"`: 指定的订阅ID不存在
- `"推荐不存在"`: 指定的推荐ID不存在
- `"您已经订阅了该主题"`: 用户已订阅该主题

## 使用说明

1. **注册流程**：
   - 先调用发送验证码接口获取验证码
   - 使用验证码进行注册

2. **登录流程**：
   - 先调用发送验证码接口获取验证码
   - 使用验证码进行登录
   - 保存返回的access_token和user_id

3. **API调用**：
   - 在请求头中添加Authorization和UserID
   - Authorization格式：`Bearer <your_access_token>`
   - 使用保存的access_token和user_id

4. **错误处理**：
   - 当收到401错误时，表示认证失效
   - 前端应清除本地存储的token和user_id
   - 跳转到登录页面或显示401状态
   - 用户需要重新登录获取新的token

5. **权限管理**：
   - 普通用户只能操作自己的话题和订阅
   - 管理员可以创建和管理系统推荐
   - 管理员可以创建订阅文章

6. **测试数据**：
   - 测试环境验证码固定为"1111"
   - 可以使用任意11位手机号进行测试
