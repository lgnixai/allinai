-- 生成聊天历史测试数据 (SQLite版本)
-- 包含29条聊天记录，涵盖不同的主题和场景

-- 清理现有测试数据（可选）
-- DELETE FROM chat_histories WHERE user_id IN (1, 2, 3);

-- 插入测试数据
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
-- 会话1：编程学习对话
(1, 'session_programming_001', 'msg_001', 'user', '你好，我想学习Python编程，应该从哪里开始？', 'gpt-3.5-turbo', 25, 0.0005, 1, 1, strftime('%s', 'now', '-1 day'), strftime('%s', 'now', '-1 day')),
(1, 'session_programming_001', 'msg_002', 'assistant', '欢迎学习Python！我建议你从以下几个方面开始：

1. 安装Python环境
2. 学习基础语法（变量、数据类型、控制流）
3. 练习简单的程序
4. 学习函数和模块
5. 了解面向对象编程

你想先从哪个方面开始？', 'gpt-3.5-turbo', 85, 0.0017, 1, 1, strftime('%s', 'now', '-1 day', '+10 seconds'), strftime('%s', 'now', '-1 day', '+10 seconds')),
(1, 'session_programming_001', 'msg_003', 'user', '我想先学习基础语法，能给我一个简单的例子吗？', 'gpt-3.5-turbo', 20, 0.0004, 1, 1, strftime('%s', 'now', '-1 day', '+20 seconds'), strftime('%s', 'now', '-1 day', '+20 seconds')),
(1, 'session_programming_001', 'msg_004', 'assistant', '当然！这里是一个简单的Python基础语法示例：

```python
# 变量定义
name = "小明"
age = 18
height = 175.5

# 打印输出
print(f"我叫{name}，今年{age}岁，身高{height}cm")

# 条件判断
if age >= 18:
    print("你已经成年了")
else:
    print("你还未成年")

# 循环
for i in range(5):
    print(f"这是第{i+1}次循环")
```

你可以先试试运行这个例子！', 'gpt-3.5-turbo', 120, 0.0024, 1, 1, strftime('%s', 'now', '-1 day', '+30 seconds'), strftime('%s', 'now', '-1 day', '+30 seconds')),

-- 会话2：数学问题讨论
(2, 'session_math_001', 'msg_005', 'user', '请帮我解释一下什么是微积分？', 'gpt-4', 15, 0.0006, 2, 1, strftime('%s', 'now', '-20 hours'), strftime('%s', 'now', '-20 hours')),
(2, 'session_math_001', 'msg_006', 'assistant', '微积分是数学的一个重要分支，主要包含两个部分：

**微分学**：研究函数的变化率，即导数。比如速度是位移对时间的导数。

**积分学**：研究函数的累积效果，即积分。比如位移是速度对时间的积分。

微积分在物理学、工程学、经济学等领域都有广泛应用。你想了解哪个具体方面？', 'gpt-4', 95, 0.0038, 2, 1, strftime('%s', 'now', '-20 hours', '+10 seconds'), strftime('%s', 'now', '-20 hours', '+10 seconds')),
(2, 'session_math_001', 'msg_007', 'user', '能给我一个简单的导数例子吗？', 'gpt-4', 15, 0.0006, 2, 1, strftime('%s', 'now', '-20 hours', '+20 seconds'), strftime('%s', 'now', '-20 hours', '+20 seconds')),
(2, 'session_math_001', 'msg_008', 'assistant', '好的！这里是一个简单的导数例子：

假设函数 f(x) = x²

要计算 f(x) 的导数，我们使用幂法则：

f''(x) = 2x

这意味着：
- 当 x = 1 时，f''(1) = 2
- 当 x = 2 时，f''(2) = 4
- 当 x = 3 时，f''(3) = 6

导数表示函数在该点的瞬时变化率。', 'gpt-4', 110, 0.0044, 2, 1, strftime('%s', 'now', '-20 hours', '+30 seconds'), strftime('%s', 'now', '-20 hours', '+30 seconds')),

-- 会话3：英语学习对话
(3, 'session_english_001', 'msg_009', 'user', '我想提高英语口语，有什么好的方法吗？', 'gpt-3.5-turbo', 20, 0.0004, 1, 1, strftime('%s', 'now', '-10 hours'), strftime('%s', 'now', '-10 hours')),
(3, 'session_english_001', 'msg_010', 'assistant', '提高英语口语有很多有效的方法：

1. **每天练习**：坚持每天说英语，哪怕只有10分钟
2. **模仿发音**：跟着英语视频或音频模仿发音
3. **找语伴**：和英语母语者或学习者交流
4. **录音自评**：录下自己的发音，找出问题
5. **学习常用表达**：掌握日常对话中的常用短语

你目前英语水平如何？我可以给你更具体的建议。', 'gpt-3.5-turbo', 130, 0.0026, 1, 1, strftime('%s', 'now', '-10 hours', '+10 seconds'), strftime('%s', 'now', '-10 hours', '+10 seconds')),
(3, 'session_english_001', 'msg_011', 'user', '我是初学者，能教我一些简单的日常对话吗？', 'gpt-3.5-turbo', 20, 0.0004, 1, 1, strftime('%s', 'now', '-10 hours', '+20 seconds'), strftime('%s', 'now', '-10 hours', '+20 seconds')),
(3, 'session_english_001', 'msg_012', 'assistant', '当然可以！这里是一些简单的日常对话：

**问候**：
- Hello! / Hi! (你好！)
- How are you? (你好吗？)
- I''m fine, thank you. (我很好，谢谢。)

**自我介绍**：
- My name is [你的名字]. (我的名字是...)
- Nice to meet you. (很高兴认识你。)

**简单交流**：
- What do you do? (你做什么工作？)
- I am a student. (我是学生。)

你想练习哪个对话？', 'gpt-3.5-turbo', 140, 0.0028, 1, 1, strftime('%s', 'now', '-10 hours', '+30 seconds'), strftime('%s', 'now', '-10 hours', '+30 seconds')),

-- 会话4：错误示例（用于测试错误处理）
(1, 'session_error_001', 'msg_013', 'user', '请帮我解决一个复杂的问题', 'gpt-4', 15, 0.0006, 2, 1, strftime('%s', 'now', '-5 hours'), strftime('%s', 'now', '-5 hours')),
(1, 'session_error_001', 'msg_014', 'assistant', '', 'gpt-4', 0, 0, 2, 2, strftime('%s', 'now', '-5 hours', '+10 seconds'), strftime('%s', 'now', '-5 hours', '+10 seconds')),
(1, 'session_error_001', 'msg_015', 'system', '模型响应超时，请重试', 'gpt-4', 0, 0, 2, 2, strftime('%s', 'now', '-5 hours', '+20 seconds'), strftime('%s', 'now', '-5 hours', '+20 seconds')),

-- 会话5：创意写作
(2, 'session_writing_001', 'msg_016', 'user', '我想写一个科幻小说，能给我一些创意建议吗？', 'gpt-4', 20, 0.0008, 2, 1, strftime('%s', 'now', '-2.5 hours'), strftime('%s', 'now', '-2.5 hours')),
(2, 'session_writing_001', 'msg_017', 'assistant', '科幻小说创作可以从以下几个角度入手：

**世界观设定**：
- 未来科技发展（AI、太空旅行、时间旅行）
- 平行宇宙或多元宇宙
- 外星文明接触

**主题探索**：
- 人性与科技的关系
- 人工智能的伦理问题
- 环境变化与人类适应

**角色设计**：
- 科学家、探险家、普通人
- 人工智能、外星生物

你有什么具体的想法吗？', 'gpt-4', 150, 0.0060, 2, 1, strftime('%s', 'now', '-2.5 hours', '+20 seconds'), strftime('%s', 'now', '-2.5 hours', '+20 seconds')),
(2, 'session_writing_001', 'msg_018', 'user', '我想写关于时间旅行的故事，有什么好的情节建议？', 'gpt-4', 20, 0.0008, 2, 1, strftime('%s', 'now', '-2.5 hours', '+30 seconds'), strftime('%s', 'now', '-2.5 hours', '+30 seconds')),
(2, 'session_writing_001', 'msg_019', 'assistant', '时间旅行故事有很多有趣的情节方向：

**经典情节**：
- 蝴蝶效应：小改变导致大后果
- 祖父悖论：回到过去杀死祖父
- 时间循环：重复同一段时间

**创新角度**：
- 时间旅行者的身份危机
- 不同时间线的冲突
- 时间旅行技术的副作用

**具体建议**：
- 设定明确的时间旅行规则
- 考虑时间旅行的社会影响
- 平衡科幻元素和人性探索

你想从哪个角度开始？', 'gpt-4', 160, 0.0064, 2, 1, strftime('%s', 'now', '-2.5 hours', '+40 seconds'), strftime('%s', 'now', '-2.5 hours', '+40 seconds')),

-- 会话6：健康咨询
(3, 'session_health_001', 'msg_020', 'user', '我最近总是感觉疲劳，有什么建议吗？', 'gpt-3.5-turbo', 15, 0.0003, 1, 1, strftime('%s', 'now', '-1.25 hours'), strftime('%s', 'now', '-1.25 hours')),
(3, 'session_health_001', 'msg_021', 'assistant', '疲劳可能由多种原因引起，以下是一些建议：

**生活习惯调整**：
- 保证充足的睡眠（7-9小时）
- 规律作息，避免熬夜
- 适量运动，增强体质

**饮食建议**：
- 均衡营养，多吃蔬果
- 适量补充维生素
- 避免过度依赖咖啡因

**其他建议**：
- 减少压力，学会放松
- 定期体检
- 如果持续疲劳，建议咨询医生

你的作息时间规律吗？', 'gpt-3.5-turbo', 140, 0.0028, 1, 1, strftime('%s', 'now', '-1.25 hours', '+10 seconds'), strftime('%s', 'now', '-1.25 hours', '+10 seconds')),

-- 会话7：技术问题
(1, 'session_tech_001', 'msg_022', 'user', '我的电脑运行很慢，怎么优化？', 'gpt-3.5-turbo', 15, 0.0003, 1, 1, strftime('%s', 'now', '-50 minutes'), strftime('%s', 'now', '-50 minutes')),
(1, 'session_tech_001', 'msg_023', 'assistant', '电脑运行慢可以从以下几个方面优化：

**系统清理**：
- 清理临时文件和垃圾文件
- 卸载不常用的软件
- 清理浏览器缓存

**硬件检查**：
- 检查硬盘空间是否充足
- 考虑升级内存或SSD
- 清理风扇和散热器

**软件优化**：
- 关闭开机自启动程序
- 更新驱动和系统
- 使用优化软件

**使用习惯**：
- 避免同时运行过多程序
- 定期重启电脑

你想先从哪个方面开始？', 'gpt-3.5-turbo', 150, 0.0030, 1, 1, strftime('%s', 'now', '-50 minutes', '+10 seconds'), strftime('%s', 'now', '-50 minutes', '+10 seconds')),

-- 会话8：旅行规划
(2, 'session_travel_001', 'msg_024', 'user', '我想去日本旅行，有什么推荐的地方吗？', 'gpt-4', 20, 0.0008, 2, 1, strftime('%s', 'now', '-25 minutes'), strftime('%s', 'now', '-25 minutes')),
(2, 'session_travel_001', 'msg_025', 'assistant', '日本有很多值得游览的地方，以下是一些推荐：

**经典城市**：
- 东京：现代化都市，购物美食天堂
- 京都：古都文化，寺庙神社众多
- 大阪：美食之都，关西文化中心
- 奈良：古都，有著名的鹿公园

**自然景观**：
- 富士山：日本象征
- 箱根：温泉胜地
- 北海道：自然风光优美

**季节建议**：
- 春季：樱花季（3-4月）
- 秋季：红叶季（11月）
- 冬季：滑雪和温泉

你计划什么时候去？预算大概多少？', 'gpt-4', 140, 0.0056, 2, 1, strftime('%s', 'now', '-25 minutes', '+10 seconds'), strftime('%s', 'now', '-25 minutes', '+10 seconds')),

-- 会话9：音乐学习
(3, 'session_music_001', 'msg_026', 'user', '我想学习吉他，应该从哪里开始？', 'gpt-3.5-turbo', 15, 0.0003, 1, 1, strftime('%s', 'now', '-13 minutes'), strftime('%s', 'now', '-13 minutes')),
(3, 'session_music_001', 'msg_027', 'assistant', '学习吉他可以从以下几个方面开始：

**基础准备**：
- 选择适合的吉他（古典、民谣、电吉他）
- 学习正确的持琴姿势
- 了解吉他的基本构造

**基础技能**：
- 学习基本和弦（C、G、Am、F等）
- 练习右手拨弦技巧
- 学习简单的节奏型

**学习方法**：
- 跟着视频教程学习
- 练习简单的歌曲
- 坚持每天练习

**推荐资源**：
- 在线教程网站
- 吉他教学APP
- 找老师指导

你有吉他吗？想学哪种类型的吉他？', 'gpt-3.5-turbo', 160, 0.0032, 1, 1, strftime('%s', 'now', '-13 minutes', '+10 seconds'), strftime('%s', 'now', '-13 minutes', '+10 seconds')),

-- 会话10：哲学讨论
(1, 'session_philosophy_001', 'msg_028', 'user', '什么是幸福？', 'gpt-4', 10, 0.0004, 2, 1, strftime('%s', 'now', '-5 minutes'), strftime('%s', 'now', '-5 minutes')),
(1, 'session_philosophy_001', 'msg_029', 'assistant', '幸福是一个复杂而深刻的哲学问题，不同文化和学派有不同的理解：

**古希腊哲学**：
- 亚里士多德认为幸福是人生的最高目标
- 伊壁鸠鲁认为快乐是幸福的基础

**东方哲学**：
- 儒家强调道德修养和人际关系
- 道家追求内心的平静和自然

**现代观点**：
- 积极心理学关注主观幸福感
- 社会学家强调社会关系的重要性

**个人理解**：
幸福可能包括：
- 内心的满足感
- 良好的人际关系
- 有意义的目标
- 健康的身体

你认为什么让你感到幸福？', 'gpt-4', 180, 0.0072, 2, 1, strftime('%s', 'now', '-5 minutes', '+10 seconds'), strftime('%s', 'now', '-5 minutes', '+10 seconds'));

-- 显示插入的测试数据统计
SELECT 
    '数据统计' as info,
    COUNT(*) as total_records,
    COUNT(DISTINCT user_id) as total_users,
    COUNT(DISTINCT session_id) as total_sessions,
    COUNT(DISTINCT model) as total_models,
    SUM(tokens) as total_tokens,
    ROUND(SUM(cost), 6) as total_cost
FROM chat_histories 
WHERE user_id IN (1, 2, 3);

-- 显示各用户的统计信息
SELECT 
    user_id,
    COUNT(*) as message_count,
    COUNT(DISTINCT session_id) as session_count,
    SUM(tokens) as total_tokens,
    ROUND(SUM(cost), 6) as total_cost,
    MAX(created_time) as last_message_time
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY user_id
ORDER BY user_id;

-- 显示各模型的统计信息
SELECT 
    model,
    COUNT(*) as message_count,
    SUM(tokens) as total_tokens,
    ROUND(SUM(cost), 6) as total_cost,
    ROUND(AVG(tokens), 2) as avg_tokens_per_message
FROM chat_histories 
WHERE user_id IN (1, 2, 3)
GROUP BY model
ORDER BY total_cost DESC;

