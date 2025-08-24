const ApiClient = require('../utils/apiClient');

describe('话题管理模块测试', () => {
  let apiClient;
  let testUserData;
  let createdTopicId;

  beforeAll(async () => {
    apiClient = new ApiClient();
    testUserData = {
      phone: global.testUtils.generateRandomPhone(),
      phone_verification_code: global.testData.user.verificationCode,
      display_name: global.testData.user.displayName,
      school: global.testData.user.school,
      college: global.testData.user.college
    };

    // 注册并登录用户
    await apiClient.sendVerificationCode(testUserData.phone, 'register');
    await apiClient.register(testUserData);
    await apiClient.sendVerificationCode(testUserData.phone, 'login');
    await apiClient.login(testUserData.phone, testUserData.phone_verification_code);
  });

  afterAll(async () => {
    // 清理测试数据
    try {
      if (createdTopicId) {
        await apiClient.deleteTopic(createdTopicId);
      }
      await apiClient.logout();
    } catch (error) {
      console.log('清理测试数据时出错:', error.message);
    }
  });

  describe('获取话题列表', () => {
    test('应该能够获取用户的话题列表', async () => {
      const result = await apiClient.getTopics();
      
      expect(result.success).toBe(true);
      expect(result.data).toBeDefined();
      expect(result.data.topics).toBeDefined();
      expect(Array.isArray(result.data.topics)).toBe(true);
      expect(result.data.total).toBeDefined();
      expect(result.data.page).toBeDefined();
      expect(result.data.size).toBeDefined();
    });

    test('应该支持分页参数', async () => {
      const result = await apiClient.getTopics(1, 5);
      
      expect(result.success).toBe(true);
      expect(result.data.page).toBe(1);
      expect(result.data.size).toBe(5);
    });

    test('应该返回空列表当用户没有话题时', async () => {
      const result = await apiClient.getTopics();
      
      expect(result.success).toBe(true);
      expect(result.data.topics).toEqual([]);
      expect(result.data.total).toBe(0);
    });
  });

  describe('创建话题', () => {
    test('应该能够成功创建话题', async () => {
      const topicData = {
        topic_name: global.testData.topic.name,
        model: global.testData.topic.model,
        channel_id: global.testData.topic.channelId
      };
      
      const result = await apiClient.createTopic(topicData);
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('成功');
      expect(result.data).toBeDefined();
      expect(result.data.id).toBeDefined();
      expect(result.data.topic_name).toBe(topicData.topic_name);
      expect(result.data.model).toBe(topicData.model);
      expect(result.data.channel_id).toBe(topicData.channel_id);
      
      // 保存创建的话题ID用于后续测试
      createdTopicId = result.data.id;
    });

    test('应该拒绝创建重复名称的话题', async () => {
      const topicData = {
        topic_name: global.testData.topic.name,
        model: global.testData.topic.model,
        channel_id: global.testData.topic.channelId
      };
      
      await expect(
        apiClient.createTopic(topicData)
      ).rejects.toThrow();
    });

    test('应该使用默认值创建话题', async () => {
      const topicData = {
        topic_name: global.testUtils.generateRandomString(10)
      };
      
      const result = await apiClient.createTopic(topicData);
      
      expect(result.success).toBe(true);
      expect(result.data.topic_name).toBe(topicData.topic_name);
      expect(result.data.model).toBeDefined();
      expect(result.data.channel_id).toBeDefined();
    });
  });

  describe('发送消息', () => {
    test('应该能够发送消息到话题', async () => {
      const messageContent = '你好，这是一个测试消息';
      const result = await apiClient.sendMessage(createdTopicId, messageContent);
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('成功');
      expect(result.data).toBeDefined();
      expect(result.data.user_message).toBeDefined();
      expect(result.data.ai_message).toBeDefined();
      expect(result.data.user_message.content).toBe(messageContent);
      expect(result.data.ai_message.content).toContain(global.testData.topic.name);
    });

    test('应该能够发送多条消息', async () => {
      const messages = [
        '第一条测试消息',
        '第二条测试消息',
        '第三条测试消息'
      ];
      
      for (const message of messages) {
        const result = await apiClient.sendMessage(createdTopicId, message);
        expect(result.success).toBe(true);
        expect(result.data.user_message.content).toBe(message);
      }
    });

    test('应该拒绝发送空消息', async () => {
      await expect(
        apiClient.sendMessage(createdTopicId, '')
      ).rejects.toThrow();
    });

    test('应该拒绝向不存在的话题发送消息', async () => {
      await expect(
        apiClient.sendMessage(99999, '测试消息')
      ).rejects.toThrow();
    });
  });

  describe('获取话题消息', () => {
    test('应该能够获取话题的消息列表', async () => {
      const result = await apiClient.getTopicMessages(createdTopicId);
      
      expect(result.success).toBe(true);
      expect(result.data).toBeDefined();
      expect(result.data.messages).toBeDefined();
      expect(Array.isArray(result.data.messages)).toBe(true);
      expect(result.data.total).toBeDefined();
      expect(result.data.page).toBeDefined();
      expect(result.data.size).toBeDefined();
    });

    test('应该支持分页参数', async () => {
      const result = await apiClient.getTopicMessages(createdTopicId, 1, 5);
      
      expect(result.success).toBe(true);
      expect(result.data.page).toBe(1);
      expect(result.data.size).toBe(5);
    });

    test('消息应该包含正确的字段', async () => {
      const result = await apiClient.getTopicMessages(createdTopicId);
      
      if (result.data.messages.length > 0) {
        const message = result.data.messages[0];
        expect(message.id).toBeDefined();
        expect(message.role).toBeDefined();
        expect(message.content).toBeDefined();
        expect(message.created_at).toBeDefined();
      }
    });

    test('AI回复应该包含话题名称前缀', async () => {
      const result = await apiClient.getTopicMessages(createdTopicId);
      
      const aiMessages = result.data.messages.filter(msg => msg.role === 'assistant');
      if (aiMessages.length > 0) {
        const aiMessage = aiMessages[0];
        expect(aiMessage.content).toContain(`"${global.testData.topic.name}"`);
      }
    });

    test('应该拒绝获取不存在话题的消息', async () => {
      await expect(
        apiClient.getTopicMessages(99999)
      ).rejects.toThrow();
    });
  });

  describe('删除话题', () => {
    test('应该能够删除话题', async () => {
      // 先创建一个新话题用于删除测试
      const topicData = {
        topic_name: global.testUtils.generateRandomString(10),
        model: global.testData.topic.model,
        channel_id: global.testData.topic.channelId
      };
      
      const createResult = await apiClient.createTopic(topicData);
      const topicId = createResult.data.id;
      
      const result = await apiClient.deleteTopic(topicId);
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('删除');
    });

    test('应该拒绝删除不存在的话题', async () => {
      await expect(
        apiClient.deleteTopic(99999)
      ).rejects.toThrow();
    });

    test('删除话题后应该无法获取消息', async () => {
      // 创建一个话题并删除
      const topicData = {
        topic_name: global.testUtils.generateRandomString(10),
        model: global.testData.topic.model,
        channel_id: global.testData.topic.channelId
      };
      
      const createResult = await apiClient.createTopic(topicData);
      const topicId = createResult.data.id;
      
      await apiClient.deleteTopic(topicId);
      
      // 尝试获取已删除话题的消息
      await expect(
        apiClient.getTopicMessages(topicId)
      ).rejects.toThrow();
    });
  });

  describe('话题数据完整性', () => {
    test('创建话题后应该出现在列表中', async () => {
      const topicData = {
        topic_name: global.testUtils.generateRandomString(10),
        model: global.testData.topic.model,
        channel_id: global.testData.topic.channelId
      };
      
      const createResult = await apiClient.createTopic(topicData);
      const topicId = createResult.data.id;
      
      const listResult = await apiClient.getTopics();
      const topicInList = listResult.data.topics.find(t => t.id === topicId);
      
      expect(topicInList).toBeDefined();
      expect(topicInList.topic_name).toBe(topicData.topic_name);
      
      // 清理
      await apiClient.deleteTopic(topicId);
    });

    test('话题消息计数应该正确', async () => {
      const listResult = await apiClient.getTopics();
      const topic = listResult.data.topics.find(t => t.id === createdTopicId);
      
      if (topic) {
        const messagesResult = await apiClient.getTopicMessages(createdTopicId);
        expect(topic.message_count).toBe(messagesResult.data.total);
      }
    });
  });

  describe('权限控制', () => {
    test('未认证用户应该无法访问话题接口', async () => {
      const unauthenticatedClient = new ApiClient();
      
      await expect(
        unauthenticatedClient.getTopics()
      ).rejects.toThrow();
      
      await expect(
        unauthenticatedClient.createTopic({ topic_name: 'test' })
      ).rejects.toThrow();
    });

    test('用户应该只能访问自己的话题', async () => {
      const result = await apiClient.getTopics();
      
      // 所有话题都应该属于当前用户
      for (const topic of result.data.topics) {
        // 这里假设话题数据中包含用户ID，如果没有则跳过此检查
        if (topic.user_id) {
          expect(topic.user_id).toBe(apiClient.userId);
        }
      }
    });
  });

  describe('错误处理', () => {
    test('应该正确处理无效的话题ID', async () => {
      await expect(
        apiClient.getTopicMessages('invalid_id')
      ).rejects.toThrow();
    });

    test('应该正确处理网络错误', async () => {
      const invalidClient = new ApiClient();
      invalidClient.baseURL = 'http://invalid-url:9999';
      
      await expect(
        invalidClient.getTopics()
      ).rejects.toThrow();
    });
  });
});

