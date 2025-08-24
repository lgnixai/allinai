# One-API 自动化测试套件

## 概述

本测试套件使用 Node.js + Jest 框架对 One-API 的所有接口进行自动化测试。

## 测试覆盖范围

### 1. 用户管理模块
- ✅ 发送手机验证码（注册）
- ✅ 用户注册
- ✅ 发送手机验证码（登录）
- ✅ 用户登录
- ✅ 获取用户信息
- ✅ 更新用户信息
- ✅ 用户登出

### 2. 话题管理模块
- ✅ 获取话题列表
- ✅ 创建话题
- ✅ 发送消息
- ✅ 获取话题消息
- ✅ 删除话题

### 3. 订阅管理模块
- ✅ 获取订阅列表
- ✅ 创建订阅
- ✅ 获取订阅文章
- ✅ 取消订阅
- ✅ 重新激活订阅
- ✅ 删除订阅

## 环境要求

- Node.js 16+
- npm 或 yarn

## 安装依赖

```bash
cd tests
npm install
```

## 配置环境

复制环境配置文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件，设置测试环境：

```env
# API配置
API_BASE_URL=http://47.88.91.79:9999

# 测试用户配置
TEST_PHONE=13800138000
TEST_VERIFICATION_CODE=1111

# 测试数据配置
TEST_DISPLAY_NAME=测试用户
TEST_SCHOOL=测试大学
TEST_COLLEGE=计算机学院
TEST_TOPIC_NAME=测试话题
TEST_SUBSCRIPTION_NAME=测试订阅
```

## 运行测试

### 运行所有测试
```bash
npm test
```

### 运行特定模块测试
```bash
# 用户管理测试
npm run test:user

# 话题管理测试
npm run test:topic

# 订阅管理测试
npm run test:subscription
```

### 运行测试并生成报告
```bash
npm run test:report
```

## 测试报告

测试完成后，可以在以下位置查看报告：
- HTML报告：`tests/reports/index.html`
- JSON报告：`tests/reports/results.json`

## 测试数据管理

- 测试数据会在测试过程中自动创建
- 测试完成后会自动清理测试数据
- 支持测试数据隔离，避免影响生产环境

## 持续集成

支持在CI/CD环境中运行：

```yaml
# GitHub Actions 示例
- name: Run API Tests
  run: |
    cd tests
    npm install
    npm test
```

## 故障排除

### 常见问题

1. **连接超时**
   - 检查服务器地址是否正确
   - 确认服务器是否正常运行

2. **认证失败**
   - 检查测试手机号和验证码
   - 确认用户注册和登录流程

3. **测试数据冲突**
   - 清理之前的测试数据
   - 使用不同的测试手机号

### 调试模式

启用详细日志：

```bash
DEBUG=* npm test
```

## 扩展测试

### 添加新接口测试

1. 在对应模块文件夹中创建测试文件
2. 继承基础测试类
3. 实现测试用例
4. 更新测试配置

### 性能测试

```bash
npm run test:performance
```

### 压力测试

```bash
npm run test:stress
```

## 贡献指南

1. 遵循现有的测试结构
2. 添加适当的测试注释
3. 确保测试覆盖率
4. 更新相关文档
