const ApiClient = require('../utils/apiClient');

describe('订阅管理模块测试', () => {
  let apiClient;
  let testUserData;
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

    // 注册并登录用户
    await apiClient.sendVerificationCode(testUserData.phone, 'register');
    await apiClient.register(testUserData);
    await apiClient.sendVerificationCode(testUserData.phone, 'login');
    await apiClient.login(testUserData.phone, testUserData.phone_verification_code);
  });

  afterAll(async () => {
    // 清理测试数据
    try {
      if (createdSubscriptionId) {
        await apiClient.deleteSubscription(createdSubscriptionId);
      }
      await apiClient.logout();
    } catch (error) {
      console.log('清理测试数据时出错:', error.message);
    }
  });

  describe('获取订阅列表', () => {
    test('应该能够获取用户的订阅列表', async () => {
      const result = await apiClient.getSubscriptions();
      
      expect(result.success).toBe(true);
      expect(result.data).toBeDefined();
      expect(result.data.subscriptions).toBeDefined();
      expect(Array.isArray(result.data.subscriptions)).toBe(true);
      expect(result.data.total).toBeDefined();
      expect(result.data.page).toBeDefined();
      expect(result.data.size).toBeDefined();
    });

    test('应该支持分页参数', async () => {
      const result = await apiClient.getSubscriptions(1, 5);
      
      expect(result.success).toBe(true);
      expect(result.data.page).toBe(1);
      expect(result.data.size).toBe(5);
    });

    test('应该返回空列表当用户没有订阅时', async () => {
      const result = await apiClient.getSubscriptions();
      
      expect(result.success).toBe(true);
      expect(result.data.subscriptions).toEqual([]);
      expect(result.data.total).toBe(0);
    });
  });

  describe('创建订阅', () => {
    test('应该能够成功创建订阅', async () => {
      const subscriptionData = {
        topic_name: global.testData.subscription.name,
        topic_description: global.testData.subscription.description
      };
      
      const result = await apiClient.createSubscription(subscriptionData);
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('成功');
      expect(result.data).toBeDefined();
      expect(result.data.id).toBeDefined();
      expect(result.data.topic_name).toBe(subscriptionData.topic_name);
      expect(result.data.topic_description).toBe(subscriptionData.topic_description);
      expect(result.data.status).toBe(1); // 活跃状态
      
      // 保存创建的订阅ID用于后续测试
      createdSubscriptionId = result.data.id;
    });

    test('应该拒绝创建重复名称的订阅', async () => {
      const subscriptionData = {
        topic_name: global.testData.subscription.name,
        topic_description: global.testData.subscription.description
      };
      
      await expect(
        apiClient.createSubscription(subscriptionData)
      ).rejects.toThrow();
    });

    test('应该使用默认值创建订阅', async () => {
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10)
      };
      
      const result = await apiClient.createSubscription(subscriptionData);
      
      expect(result.success).toBe(true);
      expect(result.data.topic_name).toBe(subscriptionData.topic_name);
      expect(result.data.status).toBeDefined();
    });
  });

  describe('获取订阅文章', () => {
    test('应该能够获取订阅的文章列表', async () => {
      const result = await apiClient.getSubscriptionArticles(createdSubscriptionId);
      
      expect(result.success).toBe(true);
      expect(result.data).toBeDefined();
      expect(result.data.articles).toBeDefined();
      expect(Array.isArray(result.data.articles)).toBe(true);
      expect(result.data.total).toBeDefined();
      expect(result.data.page).toBeDefined();
      expect(result.data.size).toBeDefined();
    });

    test('应该支持分页参数', async () => {
      const result = await apiClient.getSubscriptionArticles(createdSubscriptionId, 1, 5);
      
      expect(result.success).toBe(true);
      expect(result.data.page).toBe(1);
      expect(result.data.size).toBe(5);
    });

    test('文章应该包含正确的字段', async () => {
      const result = await apiClient.getSubscriptionArticles(createdSubscriptionId);
      
      if (result.data.articles.length > 0) {
        const article = result.data.articles[0];
        expect(article.id).toBeDefined();
        expect(article.title).toBeDefined();
        expect(article.content).toBeDefined();
        expect(article.author).toBeDefined();
        expect(article.published_at).toBeDefined();
        expect(article.article_url).toBeDefined();
      }
    });

    test('应该拒绝获取不存在订阅的文章', async () => {
      await expect(
        apiClient.getSubscriptionArticles(99999)
      ).rejects.toThrow();
    });
  });

  describe('取消订阅', () => {
    test('应该能够取消订阅', async () => {
      // 先创建一个新订阅用于取消测试
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10),
        topic_description: '用于取消测试的订阅'
      };
      
      const createResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createResult.data.id;
      
      const result = await apiClient.cancelSubscription(subscriptionId);
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('取消');
      
      // 验证订阅状态已变为取消
      const listResult = await apiClient.getSubscriptions();
      const cancelledSubscription = listResult.data.subscriptions.find(s => s.id === subscriptionId);
      expect(cancelledSubscription.status).toBe(0); // 取消状态
    });

    test('应该拒绝取消不存在的订阅', async () => {
      await expect(
        apiClient.cancelSubscription(99999)
      ).rejects.toThrow();
    });

    test('应该拒绝重复取消已取消的订阅', async () => {
      // 创建一个订阅并取消
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10),
        topic_description: '用于重复取消测试的订阅'
      };
      
      const createResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createResult.data.id;
      
      await apiClient.cancelSubscription(subscriptionId);
      
      // 尝试再次取消
      await expect(
        apiClient.cancelSubscription(subscriptionId)
      ).rejects.toThrow();
    });
  });

  describe('重新激活订阅', () => {
    test('应该能够重新激活已取消的订阅', async () => {
      // 先创建一个订阅并取消
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10),
        topic_description: '用于重新激活测试的订阅'
      };
      
      const createResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createResult.data.id;
      
      await apiClient.cancelSubscription(subscriptionId);
      
      const result = await apiClient.reactivateSubscription(subscriptionId);
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('激活');
      
      // 验证订阅状态已恢复为活跃
      const listResult = await apiClient.getSubscriptions();
      const reactivatedSubscription = listResult.data.subscriptions.find(s => s.id === subscriptionId);
      expect(reactivatedSubscription.status).toBe(1); // 活跃状态
    });

    test('应该拒绝重新激活不存在的订阅', async () => {
      await expect(
        apiClient.reactivateSubscription(99999)
      ).rejects.toThrow();
    });

    test('应该拒绝重新激活活跃状态的订阅', async () => {
      // 创建一个活跃状态的订阅
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10),
        topic_description: '用于重复激活测试的订阅'
      };
      
      const createResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createResult.data.id;
      
      // 尝试重新激活活跃状态的订阅
      await expect(
        apiClient.reactivateSubscription(subscriptionId)
      ).rejects.toThrow();
    });
  });

  describe('删除订阅', () => {
    test('应该能够删除订阅', async () => {
      // 先创建一个新订阅用于删除测试
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10),
        topic_description: '用于删除测试的订阅'
      };
      
      const createResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createResult.data.id;
      
      const result = await apiClient.deleteSubscription(subscriptionId);
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('删除');
    });

    test('应该拒绝删除不存在的订阅', async () => {
      await expect(
        apiClient.deleteSubscription(99999)
      ).rejects.toThrow();
    });

    test('删除订阅后应该无法获取文章', async () => {
      // 创建一个订阅并删除
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10),
        topic_description: '用于删除后测试的订阅'
      };
      
      const createResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createResult.data.id;
      
      await apiClient.deleteSubscription(subscriptionId);
      
      // 尝试获取已删除订阅的文章
      await expect(
        apiClient.getSubscriptionArticles(subscriptionId)
      ).rejects.toThrow();
    });
  });

  describe('订阅数据完整性', () => {
    test('创建订阅后应该出现在列表中', async () => {
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10),
        topic_description: '用于完整性测试的订阅'
      };
      
      const createResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createResult.data.id;
      
      const listResult = await apiClient.getSubscriptions();
      const subscriptionInList = listResult.data.subscriptions.find(s => s.id === subscriptionId);
      
      expect(subscriptionInList).toBeDefined();
      expect(subscriptionInList.topic_name).toBe(subscriptionData.topic_name);
      
      // 清理
      await apiClient.deleteSubscription(subscriptionId);
    });

    test('订阅状态变化应该正确反映', async () => {
      // 创建一个订阅
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10),
        topic_description: '用于状态测试的订阅'
      };
      
      const createResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createResult.data.id;
      
      // 检查初始状态为活跃
      let listResult = await apiClient.getSubscriptions();
      let subscription = listResult.data.subscriptions.find(s => s.id === subscriptionId);
      expect(subscription.status).toBe(1);
      
      // 取消订阅
      await apiClient.cancelSubscription(subscriptionId);
      
      // 检查状态变为取消
      listResult = await apiClient.getSubscriptions();
      subscription = listResult.data.subscriptions.find(s => s.id === subscriptionId);
      expect(subscription.status).toBe(0);
      
      // 重新激活
      await apiClient.reactivateSubscription(subscriptionId);
      
      // 检查状态恢复为活跃
      listResult = await apiClient.getSubscriptions();
      subscription = listResult.data.subscriptions.find(s => s.id === subscriptionId);
      expect(subscription.status).toBe(1);
      
      // 清理
      await apiClient.deleteSubscription(subscriptionId);
    });
  });

  describe('权限控制', () => {
    test('未认证用户应该无法访问订阅接口', async () => {
      const unauthenticatedClient = new ApiClient();
      
      await expect(
        unauthenticatedClient.getSubscriptions()
      ).rejects.toThrow();
      
      await expect(
        unauthenticatedClient.createSubscription({ topic_name: 'test' })
      ).rejects.toThrow();
    });

    test('用户应该只能访问自己的订阅', async () => {
      const result = await apiClient.getSubscriptions();
      
      // 所有订阅都应该属于当前用户
      for (const subscription of result.data.subscriptions) {
        // 这里假设订阅数据中包含用户ID，如果没有则跳过此检查
        if (subscription.user_id) {
          expect(subscription.user_id).toBe(apiClient.userId);
        }
      }
    });
  });

  describe('错误处理', () => {
    test('应该正确处理无效的订阅ID', async () => {
      await expect(
        apiClient.getSubscriptionArticles('invalid_id')
      ).rejects.toThrow();
    });

    test('应该正确处理网络错误', async () => {
      const invalidClient = new ApiClient();
      invalidClient.baseURL = 'http://invalid-url:9999';
      
      await expect(
        invalidClient.getSubscriptions()
      ).rejects.toThrow();
    });
  });

  describe('订阅生命周期', () => {
    test('完整的订阅生命周期测试', async () => {
      // 1. 创建订阅
      const subscriptionData = {
        topic_name: global.testUtils.generateRandomString(10),
        topic_description: '生命周期测试订阅'
      };
      
      const createResult = await apiClient.createSubscription(subscriptionData);
      const subscriptionId = createResult.data.id;
      
      // 2. 获取文章
      const articlesResult = await apiClient.getSubscriptionArticles(subscriptionId);
      expect(articlesResult.success).toBe(true);
      
      // 3. 取消订阅
      await apiClient.cancelSubscription(subscriptionId);
      
      // 4. 重新激活
      await apiClient.reactivateSubscription(subscriptionId);
      
      // 5. 删除订阅
      await apiClient.deleteSubscription(subscriptionId);
      
      // 6. 验证删除后无法访问
      await expect(
        apiClient.getSubscriptionArticles(subscriptionId)
      ).rejects.toThrow();
    });
  });
});

