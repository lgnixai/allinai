# 用户API测试指南 v2.0

## 📋 概述

本文档提供了完整的用户API测试指南，包含Postman集合、测试数据、自动化脚本和详细的使用说明。系统已更新为支持登录后自动生成access_token。

## 📁 文件结构

```
├── postman/
│   ├── User_API_Collection_v2.json    # Postman集合文件
│   └── User_API_Environment_v2.json   # Postman环境文件
├── test_data/
│   ├── generate_test_users_v2.sql     # 测试数据生成脚本
│   └── api_test_script_v2.sh          # 自动化测试脚本
└── docs/
    └── User_API_Documentation_v2.md   # 详细API文档
```

## 🚀 快速开始

### 1. 环境准备

确保以下环境已准备就绪：
- API服务器正在运行（默认端口3000）
- 数据库连接正常
- Turnstile验证已配置

### 2. Postman导入

1. 打开Postman
2. 导入 `postman/User_API_Collection_v2.json`
3. 导入 `postman/User_API_Environment_v2.json`
4. 选择环境：`用户API环境 v2.0`
5. 更新环境变量：
   - `base_url`: 你的API服务器地址
   - `turnstile_token`: 你的Turnstile验证令牌

### 3. 生成测试数据

```bash
# 连接到数据库
mysql -u your_username -p your_database_name

# 执行测试数据脚本
source test_data/generate_test_users_v2.sql
```

### 4. 运行自动化测试

```bash
# 给脚本执行权限
chmod +x test_data/api_test_script_v2.sh

# 运行测试
./test_data/api_test_script_v2.sh
```

## 📝 测试流程

### 基础功能测试

1. **发送验证码**
   - 接口：`GET /api/phone_verification`
   - 目的：获取手机验证码

2. **用户注册**
   - 接口：`POST /api/user/register`
   - 目的：创建新用户账户

3. **用户登录**
   - 接口：`POST /api/user/login`
   - 目的：登录并自动生成access_token

4. **获取用户信息**
   - 接口：`GET /api/user/self`
   - 目的：验证access_token认证

### 高级功能测试

5. **更新用户信息**
   - 接口：`PUT /api/user/self`
   - 目的：修改用户资料

6. **修改密码**
   - 接口：`PUT /api/user/self`
   - 目的：更改用户密码

7. **生成新令牌**
   - 接口：`GET /api/user/token`
   - 目的：获取新的access_token

8. **获取邀请码**
   - 接口：`GET /api/user/aff`
   - 目的：查看邀请相关信息

### 密码重置测试

9. **发送重置验证码**
   - 接口：`GET /api/reset_password`
   - 目的：获取密码重置验证码

10. **验证重置码**
    - 接口：`POST /api/user/verify_reset_code`
    - 目的：验证重置验证码

11. **重置密码**
    - 接口：`POST /api/user/reset_password`
    - 目的：设置新密码

### 清理测试

12. **用户登出**
    - 接口：`GET /api/user/logout`
    - 目的：清除session

## 🔧 测试数据说明

### 预置用户数据

| 手机号 | 用户名 | 角色 | 状态 | 学校 | 学院 | access_token |
|--------|--------|------|------|------|------|--------------|
| 13800138000 | user_8000 | 普通用户 | 正常 | 清华大学 | 计算机学院 | test_access_token_001 |
| 13800138001 | user_8001 | 普通用户 | 正常 | 北京大学 | 信息学院 | test_access_token_002 |
| 13800138002 | user_8002 | 普通用户 | 正常 | 复旦大学 | 软件学院 | test_access_token_003 |
| 13800138003 | admin_8003 | 管理员 | 正常 | 浙江大学 | 管理学院 | admin_access_token_001 |
| 13800138004 | root_8004 | 超级管理员 | 正常 | 上海交通大学 | 电子学院 | root_access_token_001 |
| 13800138005 | disabled_8005 | 普通用户 | 禁用 | 武汉大学 | 物理学院 | disabled_access_token_001 |
| 13800138006 | vip_8006 | 普通用户 | 正常 | 华中科技大学 | 机械学院 | vip_access_token_001 |
| 13800138007 | new_8007 | 普通用户 | 正常 | 西安交通大学 | 化学学院 | NULL |
| 13800138008 | enterprise_8008 | 普通用户 | 正常 | 北京理工大学 | 经济学院 | enterprise_access_token_001 |
| 13800138009 | student_8009 | 普通用户 | 正常 | 南京大学 | 文学院 | student_access_token_001 |

### 密码信息

所有测试用户的密码都是：`12345678`

### 验证码

测试环境中的验证码已硬编码为：`1111`

## 🔐 认证机制

### Session认证（Web端）
- 用于浏览器环境
- 通过Cookie维护登录状态
- 适用于Web界面操作

### Token认证（API客户端）
- 用于程序调用
- 需要两个请求头：
  - `Authorization`: access_token
  - `New-Api-User`: user_id
- 适用于API集成

### 自动生成机制
- 用户登录成功后自动生成access_token
- 如果用户已有access_token，不会重新生成
- access_token在登录响应中返回

## ⚠️ 重要更新

### v2.0 新特性

1. **登录自动生成access_token**
   - 登录成功后自动生成访问令牌
   - 无需单独调用token生成接口
   - 简化了API客户端集成流程

2. **完整的字段支持**
   - 新增学校、学院、手机号字段
   - 所有字段都有中文备注
   - 支持字段验证和约束

3. **改进的密码重置流程**
   - 两步验证流程
   - 用户可自定义新密码
   - 更安全的验证机制

## 🧪 自动化测试

### 脚本功能

`api_test_script_v2.sh` 提供了完整的自动化测试：

1. **自动提取认证信息**
   - 从登录响应中提取access_token和user_id
   - 自动用于后续API调用

2. **完整的测试覆盖**
   - 包含所有用户相关接口
   - 测试正常和异常情况

3. **详细的输出信息**
   - 彩色输出便于阅读
   - 显示请求和响应详情
   - 成功/失败状态明确

### 运行方式

```bash
# 基本运行
./api_test_script_v2.sh

# 查看帮助
./api_test_script_v2.sh --help

# 自定义配置
# 编辑脚本中的 BASE_URL 和 TURNSTILE_TOKEN 变量
```

## 📊 测试结果验证

### 成功指标

1. **HTTP状态码**: 200 OK
2. **响应格式**: JSON格式正确
3. **success字段**: true
4. **数据完整性**: 所有必要字段都存在

### 常见问题排查

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| 429错误 | 请求频率过高 | 降低请求频率，等待限制解除 |
| 401错误 | 认证失败 | 检查access_token和user_id |
| 400错误 | 参数错误 | 检查请求参数格式和必填字段 |
| 验证码错误 | 验证码过期或错误 | 重新发送验证码 |

## 🔧 配置说明

### 环境变量

| 变量名 | 描述 | 默认值 | 说明 |
|--------|------|--------|------|
| `base_url` | API基础地址 | http://localhost:3000 | 根据实际部署调整 |
| `turnstile_token` | Turnstile令牌 | your_turnstile_token_here | 需要配置有效的令牌 |

### 数据库配置

确保数据库连接正常，并执行了必要的迁移脚本：
- `bin/migration_v0.4-v0.5.sql` - 添加新字段
- `bin/migration_v0.5-v0.6.sql` - 更新约束

## 📞 技术支持

### 常见问题

1. **Q: 登录后没有access_token？**
   A: 检查用户是否已有access_token，系统不会重复生成

2. **Q: 验证码总是错误？**
   A: 确认使用的是硬编码的"1111"，或检查验证码发送是否成功

3. **Q: 429错误频繁出现？**
   A: 检查频率限制配置，或增加限制阈值

### 联系支持

如遇到问题，请：
1. 查看系统日志获取详细错误信息
2. 检查数据库连接和配置
3. 确认所有依赖服务正常运行
4. 联系技术支持团队

## 📝 更新日志

### v2.0 (2024-01-01)
- ✅ 登录后自动生成access_token
- ✅ 新增学校、学院、手机号字段
- ✅ 改进密码重置流程
- ✅ 更新Postman集合和环境
- ✅ 完善测试数据和脚本
- ✅ 详细的中文文档

### v1.0 (2023-12-01)
- ✅ 基础用户API功能
- ✅ 手机号认证系统
- ✅ 验证码验证机制
- ✅ 基础测试套件

