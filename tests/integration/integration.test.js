const ApiClient = require('../utils/apiClient');

describe('One-API 集成测试', () => {
  let apiClient;
  let testUserData;
  let createdTopicId;
  let createdSubscriptionId;

  beforeAll(async () => {
    apiClient = new ApiClient();
    testUserData = {
      phone: global.testUtils.generateRandomPhone(),
      phone_verification_code: global.testData.user.verificationCode,
      display_name: global.testData.user.displayName,
      school: global.testData.user.school,
      college: global.testData.user.college
    };
  });

  afterAll(async () => {
    // 清理测试数据
    try {
      if (createdTopicId) {
        await apiClient.deleteTopic(createdTopicId);
      }
      if (createdSubscriptionId) {
        await apiClient.deleteSubscription(createdSubscriptionId);
      }
      await apiClient.logout();
    } catch (error) {
      console.log('清理测试数据时出错:', error.message);
    }
  });

  describe('完整的用户注册和登录流程', () => {
    test('应该能够完成完整的用户注册和登录流程', async () => {
      // 1. 发送注册验证码
      const registerCodeResult = await apiClient.sendVerificationCode(
        testUserData.phone,
        'register'
      );
      expect(registerCodeResult.success).toBe(true);

      // 2. 用户注册
      const registerResult = await apiClient.register(testUserData);
      expect(registerResult.success).toBe(true);

      // 3. 发送登录验证码
      const loginCodeResult = await apiClient.sendVerificationCode(
        testUserData.phone,
        'login'
      );
      expect(loginCodeResult.success).toBe(true);

      // 4. 用户登录
      const loginResult = await apiClient.login(
        testUserData.phone,
        testUserData.phone_verification_code
      );
      expect(loginResult.success).toBe(true);
      expect(loginResult.data.access_token).toBeDefined();
      expect(loginResult.data.id).toBeDefined();
    });
  });

  describe('用户信息管理流程', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.sendVerificationCode(testUserData.phone, 'login');
      await apiClient.login(testUserData.phone, testUserData.phone_verification_code);
    });

    test('应该能够完整管理用户信息', async () => {
      // 1. 获取用户信息
      const userInfoResult = await apiClient.getUserInfo();
      expect(userInfoResult.success).toBe(true);
      expect(userInfoResult.data.phone).toBe(testUserData.phone);

      // 2. 更新用户信息
      const updateData = {
        display_name: '集成测试用户',
        school: '集成测试大学',
        college: '集成测试学院'
      };
      const updateResult = await apiClient.updateUserInfo(updateData);
      expect(updateResult.success).toBe(true);

      // 3. 验证更新结果
      const updatedUserInfo = await apiClient.getUserInfo();
      expect(updatedUserInfo.data.display_name).toBe(updateData.display_name);
      expect(updatedUserInfo.data.school).toBe(updateData.school);
      expect(updatedUserInfo.data.college).toBe(updateData.college);
    });
  });

  describe('话题管理完整流程', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.sendVerificationCode(testUserData.phone, 'login');
      await apiClient.login(testUserData.phone, testUserData.phone_verification_code);
    });

    test('应该能够完整管理话题', async () => {
      // 1. 获取初始话题列表
      const initialTopics = await apiClient.getTopics();
      expect(initialTopics.success).toBe(true);

      // 2. 创建话题
      const topicData = {
        topic_name: '集成测试话题',
        model: 'gpt-3.5-turbo',
        channel_id: 1
      };
      const createTopicResult = await apiClient.createTopic(topicData);
      expect(createTopicResult.success).toBe(true);
      createdTopicId = createTopicResult.data.id;

      // 3. 验证话题出现在列表中
      const topicsAfterCreate = await apiClient.getTopics();
      const createdTopic = topicsAfterCreate.data.topics.find(t => t.id === createdTopicId);
      expect(createdTopic).toBeDefined();
      expect(createdTopic.topic_name).toBe(topicData.topic_name);

      // 4. 发送消息
      const messageContent = '集成测试消息';
      const sendMessageResult = await apiClient.sendMessage(createdTopicId, messageContent);
      expect(sendMessageResult.success).toBe(true);
      expect(sendMessageResult.data.user_message.content).toBe(messageContent);
      expect(sendMessageResult.data.ai_message.content).toContain(topicData.topic_name);

      // 5. 获取话题消息
      const messagesResult = await apiClient.getTopicMessages(createdTopicId);
      expect(messagesResult.success).toBe(true);
      expect(messagesResult.data.messages.length).toBeGreaterThan(0);

      // 6. 发送多条消息
      const additionalMessages = [
        '第二条测试消息',
        '第三条测试消息',
        '第四条测试消息'
      ];
      for (const message of additionalMessages) {
        const result = await apiClient.sendMessage(createdTopicId, message);
        expect(result.success).toBe(true);
      }

      // 7. 验证消息计数
      const finalMessages = await apiClient.getTopicMessages(createdTopicId);
      expect(finalMessages.data.total).toBe(5); // 1条初始消息 + 4条测试消息
    });
  });

  describe('订阅管理完整流程', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.sendVerificationCode(testUserData.phone, 'login');
      await apiClient.login(testUserData.phone, testUserData.phone_verification_code);
    });

    test('应该能够完整管理订阅', async () => {
      // 1. 获取初始订阅列表
      const initialSubscriptions = await apiClient.getSubscriptions();
      expect(initialSubscriptions.success).toBe(true);

      // 2. 创建订阅
      const subscriptionData = {
        topic_name: '集成测试订阅',
        topic_description: '用于集成测试的订阅'
      };
      const createSubscriptionResult = await apiClient.createSubscription(subscriptionData);
      expect(createSubscriptionResult.success).toBe(true);
      createdSubscriptionId = createSubscriptionResult.data.id;

      // 3. 验证订阅出现在列表中
      const subscriptionsAfterCreate = await apiClient.getSubscriptions();
      const createdSubscription = subscriptionsAfterCreate.data.subscriptions.find(s => s.id === createdSubscriptionId);
      expect(createdSubscription).toBeDefined();
      expect(createdSubscription.topic_name).toBe(subscriptionData.topic_name);
      expect(createdSubscription.status).toBe(1); // 活跃状态

      // 4. 获取订阅文章
      const articlesResult = await apiClient.getSubscriptionArticles(createdSubscriptionId);
      expect(articlesResult.success).toBe(true);

      // 5. 取消订阅
      const cancelResult = await apiClient.cancelSubscription(createdSubscriptionId);
      expect(cancelResult.success).toBe(true);

      // 6. 验证订阅状态变为取消
      const subscriptionsAfterCancel = await apiClient.getSubscriptions();
      const cancelledSubscription = subscriptionsAfterCancel.data.subscriptions.find(s => s.id === createdSubscriptionId);
      expect(cancelledSubscription.status).toBe(0); // 取消状态

      // 7. 重新激活订阅
      const reactivateResult = await apiClient.reactivateSubscription(createdSubscriptionId);
      expect(reactivateResult.success).toBe(true);

      // 8. 验证订阅状态恢复为活跃
      const subscriptionsAfterReactivate = await apiClient.getSubscriptions();
      const reactivatedSubscription = subscriptionsAfterReactivate.data.subscriptions.find(s => s.id === createdSubscriptionId);
      expect(reactivatedSubscription.status).toBe(1); // 活跃状态
    });
  });

  describe('跨模块功能测试', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.sendVerificationCode(testUserData.phone, 'login');
      await apiClient.login(testUserData.phone, testUserData.phone_verification_code);
    });

    test('应该能够在话题和订阅之间切换', async () => {
      // 1. 创建话题
      const topicData = {
        topic_name: '跨模块测试话题',
        model: 'gpt-3.5-turbo',
        channel_id: 1
      };
      const createTopicResult = await apiClient.createTopic(topicData);
      const topicId = createTopicResult.data.id;

      // 2. 在话题中发送消息
      await apiClient.sendMessage(topicId, '跨模块测试消息');

      // 3. 创建订阅
      const subscriptionData = {
        topic_name: '跨模块测试订阅',
        topic_description: '跨模块测试订阅描述'
      };
      const createSubscriptionResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createSubscriptionResult.data.id;

      // 4. 获取订阅文章
      await apiClient.getSubscriptionArticles(subscriptionId);

      // 5. 再次回到话题发送消息
      await apiClient.sendMessage(topicId, '第二条跨模块测试消息');

      // 6. 验证两个模块的数据都正确
      const topicsResult = await apiClient.getTopics();
      const subscriptionsResult = await apiClient.getSubscriptions();

      expect(topicsResult.data.topics.find(t => t.id === topicId)).toBeDefined();
      expect(subscriptionsResult.data.subscriptions.find(s => s.id === subscriptionId)).toBeDefined();

      // 清理
      await apiClient.deleteTopic(topicId);
      await apiClient.deleteSubscription(subscriptionId);
    });

    test('应该能够处理并发操作', async () => {
      // 并发创建多个话题
      const topicPromises = [];
      for (let i = 0; i < 3; i++) {
        const topicData = {
          topic_name: `并发测试话题${i}`,
          model: 'gpt-3.5-turbo',
          channel_id: 1
        };
        topicPromises.push(apiClient.createTopic(topicData));
      }

      const topicResults = await Promise.all(topicPromises);
      expect(topicResults.length).toBe(3);
      topicResults.forEach(result => {
        expect(result.success).toBe(true);
      });

      // 并发创建多个订阅
      const subscriptionPromises = [];
      for (let i = 0; i < 3; i++) {
        const subscriptionData = {
          topic_name: `并发测试订阅${i}`,
          topic_description: `并发测试订阅描述${i}`
        };
        subscriptionPromises.push(apiClient.createSubscription(subscriptionData));
      }

      const subscriptionResults = await Promise.all(subscriptionPromises);
      expect(subscriptionResults.length).toBe(3);
      subscriptionResults.forEach(result => {
        expect(result.success).toBe(true);
      });

      // 清理
      const topicIds = topicResults.map(r => r.data.id);
      const subscriptionIds = subscriptionResults.map(r => r.data.id);

      await Promise.all([
        ...topicIds.map(id => apiClient.deleteTopic(id)),
        ...subscriptionIds.map(id => apiClient.deleteSubscription(id))
      ]);
    });
  });

  describe('错误恢复测试', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.sendVerificationCode(testUserData.phone, 'login');
      await apiClient.login(testUserData.phone, testUserData.phone_verification_code);
    });

    test('应该能够从认证失败中恢复', async () => {
      // 1. 清除认证信息
      apiClient.clearAuth();

      // 2. 尝试访问需要认证的接口（应该失败）
      await expect(apiClient.getUserInfo()).rejects.toThrow();

      // 3. 重新登录
      await apiClient.sendVerificationCode(testUserData.phone, 'login');
      await apiClient.login(testUserData.phone, testUserData.phone_verification_code);

      // 4. 验证可以正常访问
      const userInfo = await apiClient.getUserInfo();
      expect(userInfo.success).toBe(true);
    });

    test('应该能够处理网络中断和恢复', async () => {
      // 1. 正常操作
      const topics = await apiClient.getTopics();
      expect(topics.success).toBe(true);

      // 2. 模拟网络问题（通过设置无效的baseURL）
      const originalBaseURL = apiClient.baseURL;
      apiClient.baseURL = 'http://invalid-url:9999';

      // 3. 尝试操作（应该失败）
      await expect(apiClient.getTopics()).rejects.toThrow();

      // 4. 恢复网络连接
      apiClient.baseURL = originalBaseURL;

      // 5. 验证可以正常操作
      const topicsAfterRecovery = await apiClient.getTopics();
      expect(topicsAfterRecovery.success).toBe(true);
    });
  });

  describe('性能测试', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.sendVerificationCode(testUserData.phone, 'login');
      await apiClient.login(testUserData.phone, testUserData.phone_verification_code);
    });

    test('应该能够快速响应基本操作', async () => {
      const startTime = Date.now();

      // 执行一系列基本操作
      await apiClient.getUserInfo();
      await apiClient.getTopics();
      await apiClient.getSubscriptions();

      const endTime = Date.now();
      const duration = endTime - startTime;

      // 验证响应时间在合理范围内（5秒内）
      expect(duration).toBeLessThan(5000);
    });

    test('应该能够处理大量数据', async () => {
      // 创建多个话题
      const topicIds = [];
      for (let i = 0; i < 5; i++) {
        const topicData = {
          topic_name: `性能测试话题${i}`,
          model: 'gpt-3.5-turbo',
          channel_id: 1
        };
        const result = await apiClient.createTopic(topicData);
        topicIds.push(result.data.id);
      }

      // 在每个话题中发送消息
      for (const topicId of topicIds) {
        for (let i = 0; i < 3; i++) {
          await apiClient.sendMessage(topicId, `性能测试消息${i}`);
        }
      }

      // 获取所有话题的消息
      const messagePromises = topicIds.map(id => apiClient.getTopicMessages(id));
      const messageResults = await Promise.all(messagePromises);

      // 验证所有操作都成功
      messageResults.forEach(result => {
        expect(result.success).toBe(true);
        expect(result.data.messages.length).toBeGreaterThan(0);
      });

      // 清理
      await Promise.all(topicIds.map(id => apiClient.deleteTopic(id)));
    });
  });

  describe('数据一致性测试', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.sendVerificationCode(testUserData.phone, 'login');
      await apiClient.login(testUserData.phone, testUserData.phone_verification_code);
    });

    test('应该保持数据一致性', async () => {
      // 1. 创建话题
      const topicData = {
        topic_name: '一致性测试话题',
        model: 'gpt-3.5-turbo',
        channel_id: 1
      };
      const createTopicResult = await apiClient.createTopic(topicData);
      const topicId = createTopicResult.data.id;

      // 2. 发送消息
      await apiClient.sendMessage(topicId, '一致性测试消息');

      // 3. 多次获取数据并验证一致性
      const results = [];
      for (let i = 0; i < 3; i++) {
        const topics = await apiClient.getTopics();
        const messages = await apiClient.getTopicMessages(topicId);
        results.push({ topics, messages });
      }

      // 4. 验证数据一致性
      const firstResult = results[0];
      results.forEach(result => {
        expect(result.topics.data.total).toBe(firstResult.topics.data.total);
        expect(result.messages.data.total).toBe(firstResult.messages.data.total);
      });

      // 清理
      await apiClient.deleteTopic(topicId);
    });
  });
});

