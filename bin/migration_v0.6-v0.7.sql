-- 创建聊天历史记录表
-- 用于记录用户与AI的对话内容

CREATE TABLE IF NOT EXISTS chat_histories (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    user_id INT NOT NULL COMMENT '用户ID',
    session_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    message_id VARCHAR(64) NOT NULL COMMENT '消息ID',
    role VARCHAR(20) NOT NULL COMMENT '角色：user/assistant/system',
    content TEXT NOT NULL COMMENT '消息内容',
    model VARCHAR(50) NOT NULL COMMENT 'AI模型名称',
    tokens INT DEFAULT 0 COMMENT '消耗的token数量',
    cost DECIMAL(10,6) DEFAULT 0 COMMENT '消耗的费用',
    channel_id INT DEFAULT 0 COMMENT '使用的通道ID',
    status INT DEFAULT 1 COMMENT '状态：1-正常，2-错误',
    error_msg VARCHAR(500) DEFAULT '' COMMENT '错误信息',
    created_time BIGINT NOT NULL COMMENT '创建时间',
    updated_time BIGINT NOT NULL COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    
    -- 索引
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_message_id (message_id),
    INDEX idx_model (model),
    INDEX idx_channel_id (channel_id),
    INDEX idx_created_time (created_time),
    INDEX idx_deleted_at (deleted_at),
    
    -- 外键约束（可选，根据实际需要启用）
    -- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    -- FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天历史记录表';

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
(1, 'session_001', 'msg_001', 'user', '你好，请介绍一下你自己', 'gpt-3.5-turbo', 15, 0.0003, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(1, 'session_001', 'msg_002', 'assistant', '你好！我是OpenAI的GPT-3.5-turbo模型，一个大型语言模型。我可以帮助你回答问题、进行对话、编写代码等。有什么我可以帮助你的吗？', 'gpt-3.5-turbo', 45, 0.0009, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(1, 'session_001', 'msg_003', 'user', '请帮我写一个Python函数来计算斐波那契数列', 'gpt-3.5-turbo', 20, 0.0004, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(1, 'session_001', 'msg_004', 'assistant', '好的，这里是一个计算斐波那契数列的Python函数：\n\n```python\ndef fibonacci(n):\n    if n <= 0:\n        return 0\n    elif n == 1:\n        return 1\n    else:\n        return fibonacci(n-1) + fibonacci(n-2)\n\n# 使用示例\nfor i in range(10):\n    print(f"F({i}) = {fibonacci(i)}")\n```\n\n这个函数使用递归方式计算斐波那契数列。', 'gpt-3.5-turbo', 120, 0.0024, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 用户2的测试对话
(2, 'session_002', 'msg_005', 'user', '什么是机器学习？', 'gpt-4', 10, 0.0003, 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 'session_002', 'msg_006', 'assistant', '机器学习是人工智能的一个分支，它使计算机能够在没有明确编程的情况下学习和改进。机器学习算法通过分析数据来识别模式，并使用这些模式来做出预测或决策。', 'gpt-4', 35, 0.0011, 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 用户3的测试对话（包含错误）
(3, 'session_003', 'msg_007', 'user', '请帮我解决这个数学问题', 'gpt-3.5-turbo', 15, 0.0003, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 'session_003', 'msg_008', 'assistant', '', 'gpt-3.5-turbo', 0, 0, 1, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 'session_003', 'msg_009', 'system', '模型响应超时', 'gpt-3.5-turbo', 0, 0, 1, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 显示创建的表结构
DESCRIBE chat_histories;

-- 显示插入的测试数据
SELECT 
    id,
    user_id,
    session_id,
    role,
    LEFT(content, 50) as content_preview,
    model,
    tokens,
    cost,
    status,
    FROM_UNIXTIME(created_time) as created_time
FROM chat_histories 
ORDER BY id;

