-- 创建聊天历史记录表 (SQLite版本)
-- 用于记录用户与AI的对话内容

CREATE TABLE IF NOT EXISTS chat_histories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_id TEXT NOT NULL,
    message_id TEXT NOT NULL,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    model TEXT NOT NULL,
    tokens INTEGER DEFAULT 0,
    cost REAL DEFAULT 0,
    channel_id INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    error_msg TEXT DEFAULT '',
    created_time INTEGER NOT NULL,
    updated_time INTEGER NOT NULL,
    deleted_at DATETIME DEFAULT NULL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_chat_histories_user_id ON chat_histories(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_histories_session_id ON chat_histories(session_id);
CREATE INDEX IF NOT EXISTS idx_chat_histories_message_id ON chat_histories(message_id);
CREATE INDEX IF NOT EXISTS idx_chat_histories_model ON chat_histories(model);
CREATE INDEX IF NOT EXISTS idx_chat_histories_channel_id ON chat_histories(channel_id);
CREATE INDEX IF NOT EXISTS idx_chat_histories_created_time ON chat_histories(created_time);
CREATE INDEX IF NOT EXISTS idx_chat_histories_deleted_at ON chat_histories(deleted_at);

-- 添加一些示例数据用于测试
INSERT INTO chat_histories (
    user_id, 
    session_id, 
    message_id, 
    role, 
    content, 
    model, 
    tokens, 
    cost, 
    channel_id, 
    status, 
    created_time, 
    updated_time
) VALUES 
-- 用户1的测试对话
(1, 'session_001', 'msg_001', 'user', '你好，请介绍一下你自己', 'gpt-3.5-turbo', 15, 0.0003, 1, 1, strftime('%s', 'now', '-1 day'), strftime('%s', 'now', '-1 day')),
(1, 'session_001', 'msg_002', 'assistant', '你好！我是OpenAI的GPT-3.5-turbo模型，一个大型语言模型。我可以帮助你回答问题、进行对话、编写代码等。有什么我可以帮助你的吗？', 'gpt-3.5-turbo', 45, 0.0009, 1, 1, strftime('%s', 'now', '-1 day', '+10 seconds'), strftime('%s', 'now', '-1 day', '+10 seconds')),
(1, 'session_001', 'msg_003', 'user', '请帮我写一个Python函数来计算斐波那契数列', 'gpt-3.5-turbo', 20, 0.0004, 1, 1, strftime('%s', 'now', '-1 day', '+20 seconds'), strftime('%s', 'now', '-1 day', '+20 seconds')),
(1, 'session_001', 'msg_004', 'assistant', '好的，这里是一个计算斐波那契数列的Python函数：

```python
def fibonacci(n):
    if n <= 0:
        return 0
    elif n == 1:
        return 1
    else:
        return fibonacci(n-1) + fibonacci(n-2)

# 使用示例
for i in range(10):
    print(f"F({i}) = {fibonacci(i)}")
```

这个函数使用递归方式计算斐波那契数列。', 'gpt-3.5-turbo', 120, 0.0024, 1, 1, strftime('%s', 'now', '-1 day', '+30 seconds'), strftime('%s', 'now', '-1 day', '+30 seconds')),

-- 用户2的测试对话
(2, 'session_002', 'msg_005', 'user', '什么是机器学习？', 'gpt-4', 10, 0.0003, 2, 1, strftime('%s', 'now', '-12 hours'), strftime('%s', 'now', '-12 hours')),
(2, 'session_002', 'msg_006', 'assistant', '机器学习是人工智能的一个分支，它使计算机能够在没有明确编程的情况下学习和改进。机器学习算法通过分析数据来识别模式，并使用这些模式来做出预测或决策。', 'gpt-4', 35, 0.0011, 2, 1, strftime('%s', 'now', '-12 hours', '+10 seconds'), strftime('%s', 'now', '-12 hours', '+10 seconds')),

-- 用户3的测试对话（包含错误）
(3, 'session_003', 'msg_007', 'user', '请帮我解决这个数学问题', 'gpt-3.5-turbo', 15, 0.0003, 1, 1, strftime('%s', 'now', '-6 hours'), strftime('%s', 'now', '-6 hours')),
(3, 'session_003', 'msg_008', 'assistant', '', 'gpt-3.5-turbo', 0, 0, 1, 2, strftime('%s', 'now', '-6 hours', '+10 seconds'), strftime('%s', 'now', '-6 hours', '+10 seconds')),
(3, 'session_003', 'msg_009', 'system', '模型响应超时', 'gpt-3.5-turbo', 0, 0, 1, 2, strftime('%s', 'now', '-6 hours', '+20 seconds'), strftime('%s', 'now', '-6 hours', '+20 seconds'));

-- 显示创建的表结构
PRAGMA table_info(chat_histories);

-- 显示插入的测试数据
SELECT 
    id,
    user_id,
    session_id,
    role,
    substr(content, 1, 50) as content_preview,
    model,
    tokens,
    cost,
    status,
    datetime(created_time, 'unixepoch') as created_time
FROM chat_histories 
ORDER BY id;

