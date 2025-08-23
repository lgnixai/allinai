# 聊天历史测试数据说明 (SQLite版本)

## 📋 概述

本目录包含了ChatHistory模块的SQLite版本测试数据生成脚本和验证脚本，用于测试聊天历史相关的API功能。

## 📁 文件说明

### 1. `generate_chat_history_test_data_sqlite.sql`
**功能**: 生成聊天历史测试数据 (SQLite版本)
**内容**: 
- 29条聊天记录
- 10个不同的会话
- 3个用户（用户ID: 1, 2, 3）
- 2个AI模型（gpt-3.5-turbo, gpt-4）
- 涵盖多种主题和场景

### 2. `test_chat_history_data_sqlite.sql`
**功能**: 验证测试数据的正确性 (SQLite版本)
**内容**:
- 12个验证查询
- 数据完整性检查
- 统计分析
- 异常数据检测

### 3. `run_chat_history_test_sqlite.sh`
**功能**: 自动化执行脚本 (SQLite版本)
**内容**:
- 依赖检查
- 数据库连接验证
- 自动化数据生成和验证

### 4. `bin/migration_v0.6-v0.7_sqlite.sql`
**功能**: SQLite数据库迁移脚本
**内容**:
- 创建chat_histories表
- 创建必要的索引
- 插入基础测试数据

## 🎯 测试数据特点

### 主题覆盖
| 主题 | 会话ID | 描述 |
|------|--------|------|
| 编程学习 | session_programming_001 | Python编程入门指导 |
| 数学问题 | session_math_001 | 微积分概念解释 |
| 英语学习 | session_english_001 | 英语口语提高方法 |
| 错误处理 | session_error_001 | 模型响应错误示例 |
| 创意写作 | session_writing_001 | 科幻小说创作建议 |
| 健康咨询 | session_health_001 | 疲劳问题咨询 |
| 技术问题 | session_tech_001 | 电脑优化建议 |
| 旅行规划 | session_travel_001 | 日本旅行推荐 |
| 音乐学习 | session_music_001 | 吉他学习指导 |
| 哲学讨论 | session_philosophy_001 | 幸福概念探讨 |

### 数据统计
- **总记录数**: 29条
- **用户数**: 3个
- **会话数**: 10个
- **模型数**: 2个
- **角色类型**: user, assistant, system
- **状态类型**: 正常(1), 错误(2)

### 时间分布
- 数据跨越约24小时
- 模拟真实的聊天时间分布
- 包含不同时间段的对话

## 🚀 使用方法

### 1. 执行数据库迁移
```bash
# 连接到SQLite数据库
sqlite3 one-api.db

# 执行迁移脚本
.read bin/migration_v0.6-v0.7_sqlite.sql
```

### 2. 使用自动化脚本（推荐）
```bash
# 设置数据库文件路径（可选，默认为one-api.db）
export DB_FILE="your_database.db"

# 完整流程（清理+生成+验证）
./test_data/run_chat_history_test_sqlite.sh --full

# 仅生成数据
./test_data/run_chat_history_test_sqlite.sh --generate

# 仅验证数据
./test_data/run_chat_history_test_sqlite.sh --validate
```

### 3. 手动执行SQL
```bash
# 连接到SQLite数据库
sqlite3 one-api.db

# 生成测试数据
.read test_data/generate_chat_history_test_data_sqlite.sql

# 验证数据
.read test_data/test_chat_history_data_sqlite.sql
```

### 4. 清理测试数据（可选）
```sql
-- 删除测试数据
DELETE FROM chat_histories WHERE user_id IN (1, 2, 3);
```

## 📊 预期验证结果

### 数据总量检查
- 总记录数: 29
- 唯一用户数: 3
- 唯一会话数: 10
- 唯一模型数: 2

### 用户数据分布
| 用户ID | 消息数 | 会话数 | 总Token | 总费用 |
|--------|--------|--------|---------|--------|
| 1 | 10 | 4 | ~500 | ~0.01 |
| 2 | 10 | 4 | ~600 | ~0.02 |
| 3 | 9 | 2 | ~400 | ~0.01 |

### 模型使用情况
| 模型 | 消息数 | 用户数 | 总Token | 总费用 |
|------|--------|--------|---------|--------|
| gpt-4 | 12 | 2 | ~700 | ~0.03 |
| gpt-3.5-turbo | 17 | 3 | ~800 | ~0.02 |

### 主题分类统计
| 主题 | 消息数 | 会话数 | 用户数 |
|------|--------|--------|--------|
| 编程学习 | 4 | 1 | 1 |
| 数学问题 | 4 | 1 | 1 |
| 英语学习 | 4 | 1 | 1 |
| 创意写作 | 4 | 1 | 1 |
| 其他主题 | 13 | 6 | 3 |

## 🔧 SQLite特定注意事项

### 1. 数据类型差异
- SQLite使用 `INTEGER` 而不是 `INT`
- SQLite使用 `REAL` 而不是 `DECIMAL`
- SQLite使用 `TEXT` 而不是 `VARCHAR`

### 2. 时间函数差异
- SQLite使用 `strftime('%s', 'now')` 而不是 `UNIX_TIMESTAMP()`
- SQLite使用 `datetime(timestamp, 'unixepoch')` 而不是 `FROM_UNIXTIME()`

### 3. 字符串函数差异
- SQLite使用 `substr()` 而不是 `LEFT()`
- SQLite使用 `length()` 而不是 `LENGTH()`

### 4. 索引创建
- SQLite使用 `CREATE INDEX IF NOT EXISTS` 语法
- 索引名称需要唯一

## 🔧 API测试建议

### 1. 基础CRUD测试
- 使用用户1的聊天记录测试获取、更新、删除功能
- 验证权限控制（用户只能访问自己的记录）

### 2. 会话管理测试
- 测试获取用户会话列表
- 测试获取指定会话的聊天历史
- 测试删除会话功能

### 3. 统计功能测试
- 测试用户统计信息
- 测试模型使用统计
- 验证费用和token统计

### 4. 搜索功能测试
- 使用关键词"Python"、"数学"、"英语"等测试搜索
- 验证搜索结果的相关性

### 5. 错误处理测试
- 测试错误状态的消息记录
- 验证错误信息的正确显示

## ⚠️ 注意事项

1. **数据隔离**: 测试数据使用用户ID 1, 2, 3，避免与生产数据冲突
2. **时间戳**: 使用相对时间戳，确保数据在不同时间执行时的一致性
3. **费用计算**: 费用数据为模拟值，仅用于测试
4. **内容长度**: 包含不同长度的消息，测试系统处理能力
5. **错误场景**: 包含错误状态的记录，测试错误处理功能
6. **SQLite特性**: 注意SQLite的语法和数据类型差异

## 🧪 测试场景

### 场景1: 用户聊天历史查询
```bash
# 获取用户1的所有聊天历史
curl -X GET "http://localhost:3000/api/chat_history/user/history" \
  -H "Authorization: your_access_token" \
  -H "New-Api-User: 1"
```

### 场景2: 会话历史查询
```bash
# 获取编程学习会话的历史
curl -X GET "http://localhost:3000/api/chat_history/session/session_programming_001" \
  -H "Authorization: your_access_token" \
  -H "New-Api-User: 1"
```

### 场景3: 搜索功能测试
```bash
# 搜索包含"Python"的聊天记录
curl -X GET "http://localhost:3000/api/chat_history/search?keyword=Python" \
  -H "Authorization: your_access_token" \
  -H "New-Api-User: 1"
```

### 场景4: 统计信息查询
```bash
# 获取用户统计信息
curl -X GET "http://localhost:3000/api/chat_history/user/stats" \
  -H "Authorization: your_access_token" \
  -H "New-Api-User: 1"
```

## 📈 性能测试建议

1. **分页测试**: 测试不同page_size的性能表现
2. **搜索性能**: 测试大量数据下的搜索响应时间
3. **并发测试**: 模拟多用户同时访问的场景
4. **数据量测试**: 逐步增加数据量，测试系统承载能力

## 🔄 数据更新

如需更新测试数据，可以：
1. 修改 `generate_chat_history_test_data_sqlite.sql` 文件
2. 重新执行生成脚本
3. 运行验证脚本确认数据正确性

## 📞 技术支持

如有问题，请检查：
1. SQLite数据库文件是否存在
2. 表结构是否正确
3. 用户权限是否足够
4. 时间戳格式是否正确
5. SQLite语法是否正确

## 🔍 故障排除

### 常见问题

1. **数据库文件不存在**
   ```bash
   # 检查数据库文件
   ls -la one-api.db
   ```

2. **表不存在**
   ```sql
   -- 检查表是否存在
   SELECT name FROM sqlite_master WHERE type='table' AND name='chat_histories';
   ```

3. **权限问题**
   ```bash
   # 检查文件权限
   chmod 644 one-api.db
   ```

4. **SQLite版本问题**
   ```bash
   # 检查SQLite版本
   sqlite3 --version
   ```

