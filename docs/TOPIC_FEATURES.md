# 话题管理功能说明

## 功能概述

话题管理模块提供了完整的用户话题创建、管理和交互功能，支持AI自动回复和智能话题管理。

## 核心功能

### 1. 话题列表管理
- **获取话题列表**: `GET /api/topics/`
- **创建话题**: `POST /api/topics/` (无需参数，自动创建默认话题)
- **更新话题名称**: `PUT /api/topics/{id}`
- **删除话题**: `DELETE /api/topics/{id}`

### 2. 消息管理
- **获取话题消息**: `GET /api/topics/{id}/messages`
- **发送消息**: `POST /api/topics/{id}/messages`
- **自动话题创建**: 当话题ID为0时，自动创建新话题

### 3. 智能特性
- **权限验证**: 确保用户只能操作自己的话题
- **AI自动回复**: 发送消息后自动生成AI回复
- **中文支持**: 正确处理中文字符的话题标题生成

## 技术特性

### 数据模型
- **Topic**: 话题表，包含用户ID、话题名称、状态等字段
- **Message**: 消息表，包含话题ID、角色、内容等字段
- **软删除**: 使用status字段标记删除状态

### 权限管理
- 用户隔离：每个用户只能操作自己的话题
- 权限验证：所有操作都验证用户身份
- 错误处理：完善的权限错误提示

## 使用示例

### 创建话题并发送消息
```bash
# 自动创建话题并发送消息
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -H "UserID: <user_id>" \
  -H "Content-Type: application/json" \
  -d '{"content": "你好，这是一个测试消息"}' \
  "http://localhost:9999/api/topics/0"
```

### 更新话题名称
```bash
curl -X PUT \
  -H "Authorization: Bearer <token>" \
  -H "UserID: <user_id>" \
  -H "Content-Type: application/json" \
  -d '{"topic_name": "我的新话题"}' \
  "http://localhost:9999/api/topics/1"
```

## 错误处理

### 常见错误码
- **400**: 请求参数错误
- **401**: 未授权，认证失败
- **403**: 无权限操作此话题
- **404**: 话题不存在

## 最佳实践

1. **话题命名**: 使用简洁明了的话题名称
2. **消息管理**: 定期清理不需要的话题
3. **权限管理**: 确保用户只能操作自己的话题

## 扩展功能

### 计划中的功能
1. 话题分类和标签
2. 话题内容搜索
3. 话题分享功能
4. 消息历史导出
5. 集成真实AI API
