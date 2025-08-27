# 认证架构重构说明

## 问题描述

当前的认证架构存在一个设计问题：`authHelper` 函数同时支持session认证和token认证，这导致API调用时可能依赖session，这在架构上是不合理的。

### 问题分析

1. **混合认证方式**：
   - `authHelper` 首先尝试从session获取用户信息
   - 如果session为空，才检查access token
   - 这种设计导致API调用可能依赖session

2. **架构不合理**：
   - API应该是无状态的，不应该依赖session
   - session认证适合Web界面，token认证适合API
   - 混合使用会增加复杂性和安全风险

3. **当前路由使用情况**：
   - `/api/*` 路由使用 `UserAuth()` 和 `AdminAuth()`（基于session）
   - `/v1/*` 路由使用 `TokenAuth()`（基于token）
   - 这种分离是正确的，但实现上还有问题

## 解决方案

### 1. 分离认证方式

创建专门的认证中间件：

```go
// WebAuth 专门用于Web界面的session认证
func WebAuth() func(c *gin.Context)

// APIAuth 专门用于API的token认证  
func APIAuth() func(c *gin.Context)
```

### 2. 认证方式选择

| 使用场景 | 认证方式 | 中间件 | 说明 |
|---------|---------|--------|------|
| Web界面 | Session | `WebAuth()` | 用户通过浏览器登录 |
| API调用 | Token | `APIAuth()` | 程序通过API调用 |
| AI服务 | Token | `TokenAuth()` | 专门的AI API认证 |

### 3. 迁移计划

#### 第一阶段：添加新的认证中间件
- [x] 创建 `WebAuth()` 和 `APIAuth()` 函数
- [x] 保留原有的 `authHelper` 用于向后兼容

#### 第二阶段：逐步迁移路由
- [ ] 将Web界面路由迁移到 `WebAuth()`
- [ ] 将API路由迁移到 `APIAuth()`
- [ ] 保持AI服务路由使用 `TokenAuth()`

#### 第三阶段：清理旧代码
- [ ] 移除 `authHelper` 函数
- [ ] 移除 `UserAuth()` 和 `AdminAuth()` 函数
- [ ] 更新相关文档

### 4. 路由迁移示例

#### 当前路由配置：
```go
// Web界面路由（应该使用WebAuth）
apiRouter.GET("/models", middleware.UserAuth(), controller.DashboardListModels)

// API路由（应该使用APIAuth）
selfRoute.Use(middleware.UserAuth())
{
    selfRoute.GET("/self", controller.GetSelf)
    selfRoute.PUT("/self", controller.UpdateSelf)
}
```

#### 迁移后的路由配置：
```go
// Web界面路由
apiRouter.GET("/models", middleware.WebAuth(), controller.DashboardListModels)

// API路由
selfRoute.Use(middleware.APIAuth())
{
    selfRoute.GET("/self", controller.GetSelf)
    selfRoute.PUT("/self", controller.UpdateSelf)
}
```

## 优势

1. **架构清晰**：Web认证和API认证完全分离
2. **安全性提升**：API不再依赖session，减少CSRF风险
3. **性能优化**：API调用不需要检查session
4. **维护性**：代码结构更清晰，易于维护
5. **扩展性**：可以独立优化Web和API的认证逻辑

## 注意事项

1. **向后兼容**：在迁移期间保持向后兼容
2. **测试充分**：确保迁移后功能正常
3. **文档更新**：更新API文档和开发文档
4. **客户端适配**：确保客户端正确使用新的认证方式

## 相关文件

- `middleware/auth.go` - 认证中间件实现
- `router/api-router.go` - API路由配置
- `router/relay-router.go` - AI服务路由配置
- `controller/user.go` - 用户认证相关控制器
