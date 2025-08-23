# One-API 文档中心

## 文档结构

```
docs/
├── README.md                    # 本文档
├── API_Documentation.md         # API接口文档
└── Postman_Usage_Guide.md       # Postman使用指南

postman/
├── One-API_Collection.json      # Postman测试集合
└── One-API_Environment.json     # Postman环境配置
```

## 快速开始

### 1. 阅读API文档

查看 [API_Documentation.md](./API_Documentation.md) 了解所有接口的详细信息。

### 2. 使用Postman测试

1. 导入Postman文件：
   - `One-API_Collection.json` - 测试集合
   - `One-API_Environment.json` - 环境配置

2. 按照 [Postman_Usage_Guide.md](./Postman_Usage_Guide.md) 的说明进行测试

## 主要功能模块

### 1. 用户管理
- 用户登录/注册
- 获取用户信息
- 更新用户信息
- 用户认证

### 2. 话题管理
- 创建/删除话题
- 发送消息
- 查看聊天记录
- AI自动回复

### 3. 订阅管理
- 创建/删除订阅
- 取消/重新激活订阅
- 查看订阅文章

## 服务器信息

- **服务器地址**: `http://47.88.91.79:9999`
- **认证方式**: Bearer Token
- **数据格式**: JSON

## 技术支持

如有问题，请查看：
1. API文档中的错误码说明
2. Postman使用指南中的常见问题
3. 接口响应示例

## 更新日志

- 2024-01-01: 创建文档结构
- 2024-01-01: 添加用户、话题、订阅API文档
- 2024-01-01: 创建Postman测试集合
