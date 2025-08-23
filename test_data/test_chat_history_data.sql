-- 测试聊天历史数据验证脚本
-- 用于验证生成的测试数据是否正确

-- 1. 检查数据总量
SELECT 
    '数据总量检查' as test_name,
    COUNT(*) as total_records,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT session_id) as unique_sessions,
    COUNT(DISTINCT model) as unique_models
FROM chat_histories 
WHERE user_id IN (1, 2, 3);

-- 2. 检查各用户的数据分布
SELECT 
    '用户数据分布' as test_name,
    user_id,
    COUNT(*) as message_count,
    COUNT(DISTINCT session_id) as session_count,
    SUM(tokens) as total_tokens,
    ROUND(SUM(cost), 6) as total_cost,
    MIN(created_time) as first_message_time,
    MAX(created_time) as last_message_time
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY user_id
ORDER BY user_id;

-- 3. 检查各会话的数据
SELECT 
    '会话数据检查' as test_name,
    session_id,
    user_id,
    COUNT(*) as message_count,
    SUM(tokens) as total_tokens,
    ROUND(SUM(cost), 6) as total_cost,
    MIN(created_time) as start_time,
    MAX(created_time) as end_time
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY session_id, user_id
ORDER BY session_id;

-- 4. 检查角色分布
SELECT 
    '角色分布检查' as test_name,
    role,
    COUNT(*) as message_count,
    COUNT(DISTINCT session_id) as session_count,
    AVG(tokens) as avg_tokens,
    ROUND(SUM(cost), 6) as total_cost
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY role
ORDER BY message_count DESC;

-- 5. 检查模型使用情况
SELECT 
    '模型使用情况' as test_name,
    model,
    COUNT(*) as message_count,
    COUNT(DISTINCT user_id) as user_count,
    SUM(tokens) as total_tokens,
    ROUND(SUM(cost), 6) as total_cost,
    ROUND(AVG(tokens), 2) as avg_tokens_per_message
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY model
ORDER BY total_cost DESC;

-- 6. 检查状态分布
SELECT 
    '状态分布检查' as test_name,
    status,
    COUNT(*) as message_count,
    COUNT(DISTINCT session_id) as session_count
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY status
ORDER BY status;

-- 7. 检查时间分布
SELECT 
    '时间分布检查' as test_name,
    FROM_UNIXTIME(MIN(created_time)) as earliest_message,
    FROM_UNIXTIME(MAX(created_time)) as latest_message,
    COUNT(*) as total_messages,
    ROUND((MAX(created_time) - MIN(created_time)) / 3600, 2) as hours_span
FROM chat_histories 
WHERE user_id IN (1, 2, 3);

-- 8. 检查错误记录
SELECT 
    '错误记录检查' as test_name,
    id,
    session_id,
    role,
    LEFT(content, 50) as content_preview,
    error_msg,
    status,
    FROM_UNIXTIME(created_time) as created_time
FROM chat_histories 
WHERE user_id IN (1, 2, 3) AND status = 2
ORDER BY created_time;

-- 9. 检查最活跃的会话
SELECT 
    '最活跃会话' as test_name,
    session_id,
    user_id,
    COUNT(*) as message_count,
    SUM(tokens) as total_tokens,
    ROUND(SUM(cost), 6) as total_cost,
    FROM_UNIXTIME(MIN(created_time)) as start_time,
    FROM_UNIXTIME(MAX(created_time)) as end_time
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY session_id, user_id
ORDER BY message_count DESC
LIMIT 5;

-- 10. 检查费用最高的会话
SELECT 
    '费用最高会话' as test_name,
    session_id,
    user_id,
    COUNT(*) as message_count,
    SUM(tokens) as total_tokens,
    ROUND(SUM(cost), 6) as total_cost,
    ROUND(AVG(cost), 6) as avg_cost_per_message
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY session_id, user_id
ORDER BY total_cost DESC
LIMIT 5;

-- 11. 检查内容长度分布
SELECT 
    '内容长度分布' as test_name,
    CASE 
        WHEN LENGTH(content) < 50 THEN '短消息 (<50字符)'
        WHEN LENGTH(content) < 200 THEN '中等消息 (50-200字符)'
        WHEN LENGTH(content) < 500 THEN '长消息 (200-500字符)'
        ELSE '超长消息 (>500字符)'
    END as content_length_category,
    COUNT(*) as message_count,
    ROUND(AVG(LENGTH(content)), 2) as avg_length,
    MIN(LENGTH(content)) as min_length,
    MAX(LENGTH(content)) as max_length
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY 
    CASE 
        WHEN LENGTH(content) < 50 THEN '短消息 (<50字符)'
        WHEN LENGTH(content) < 200 THEN '中等消息 (50-200字符)'
        WHEN LENGTH(content) < 500 THEN '长消息 (200-500字符)'
        ELSE '超长消息 (>500字符)'
    END
ORDER BY message_count DESC;

-- 12. 检查主题分类（基于会话ID）
SELECT 
    '主题分类检查' as test_name,
    CASE 
        WHEN session_id LIKE '%programming%' THEN '编程学习'
        WHEN session_id LIKE '%math%' THEN '数学问题'
        WHEN session_id LIKE '%english%' THEN '英语学习'
        WHEN session_id LIKE '%error%' THEN '错误处理'
        WHEN session_id LIKE '%writing%' THEN '创意写作'
        WHEN session_id LIKE '%health%' THEN '健康咨询'
        WHEN session_id LIKE '%tech%' THEN '技术问题'
        WHEN session_id LIKE '%travel%' THEN '旅行规划'
        WHEN session_id LIKE '%music%' THEN '音乐学习'
        WHEN session_id LIKE '%philosophy%' THEN '哲学讨论'
        ELSE '其他'
    END as topic_category,
    COUNT(*) as message_count,
    COUNT(DISTINCT session_id) as session_count,
    COUNT(DISTINCT user_id) as user_count,
    SUM(tokens) as total_tokens,
    ROUND(SUM(cost), 6) as total_cost
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY 
    CASE 
        WHEN session_id LIKE '%programming%' THEN '编程学习'
        WHEN session_id LIKE '%math%' THEN '数学问题'
        WHEN session_id LIKE '%english%' THEN '英语学习'
        WHEN session_id LIKE '%error%' THEN '错误处理'
        WHEN session_id LIKE '%writing%' THEN '创意写作'
        WHEN session_id LIKE '%health%' THEN '健康咨询'
        WHEN session_id LIKE '%tech%' THEN '技术问题'
        WHEN session_id LIKE '%travel%' THEN '旅行规划'
        WHEN session_id LIKE '%music%' THEN '音乐学习'
        WHEN session_id LIKE '%philosophy%' THEN '哲学讨论'
        ELSE '其他'
    END
ORDER BY message_count DESC;

