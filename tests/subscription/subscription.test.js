const request = require('supertest');
const { API_BASE_URL } = require('../env');

describe('订阅管理 API 测试', () => {
  let authToken;
  let userId;
  let subscriptionId;

  beforeAll(async () => {
    // 先登录获取认证信息
    const loginResponse = await request(API_BASE_URL)
      .post('/api/user/login')
      .send({
        phone: '13800138000',
        phone_verification_code: '123456'
      });

    authToken = loginResponse.body.data.access_token;
    userId = loginResponse.body.data.user.id;
  });

  describe('获取订阅列表', () => {
    test('获取订阅列表', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/subscriptions/')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .query({
          page: 1,
          page_size: 10
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('subscriptions');
      expect(response.body.data).toHaveProperty('total');
      expect(Array.isArray(response.body.data.subscriptions)).toBe(true);
    });
  });

  describe('创建订阅', () => {
    test('创建订阅', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/subscriptions/')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          topic_name: '技术订阅'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('订阅创建成功');
      expect(response.body.data).toHaveProperty('id');

      subscriptionId = response.body.data.id;
    });

    test('创建重复订阅', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/subscriptions/')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          topic_name: '技术订阅'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('订阅创建成功');
    });
  });

  describe('更新订阅', () => {
    test('更新订阅', async () => {
      const response = await request(API_BASE_URL)
        .put(`/api/subscriptions/${subscriptionId}`)
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          topic_name: '更新后的技术订阅',
          description: '更新后的订阅描述'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('订阅更新成功');
      expect(response.body.data).toHaveProperty('id');
      expect(response.body.data).toHaveProperty('topic_name');
      expect(response.body.data.topic_name).toBe('更新后的技术订阅');
    });
  });

  describe('取消订阅', () => {
    test('取消订阅', async () => {
      const response = await request(API_BASE_URL)
        .put(`/api/subscriptions/${subscriptionId}/cancel`)
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString());

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('订阅取消成功');
    });
  });

  describe('重新激活订阅', () => {
    test('重新激活订阅', async () => {
      const response = await request(API_BASE_URL)
        .put(`/api/subscriptions/${subscriptionId}/reactivate`)
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString());

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('订阅重新激活成功');
    });
  });

  describe('获取订阅文章', () => {
    test('获取订阅文章', async () => {
      const response = await request(API_BASE_URL)
        .get(`/api/subscriptions/${subscriptionId}/articles`)
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .query({
          page: 1,
          page_size: 10
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('articles');
      expect(response.body.data).toHaveProperty('total');
      expect(Array.isArray(response.body.data.articles)).toBe(true);
    });
  });

  describe('获取所有订阅文章', () => {
    test('获取所有订阅文章', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/subscriptions/articles')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .query({
          page: 1,
          page_size: 10
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('articles');
      expect(response.body.data).toHaveProperty('total');
      expect(Array.isArray(response.body.data.articles)).toBe(true);
    });
  });

  describe('删除订阅', () => {
    test('删除订阅', async () => {
      // 先创建一个新订阅用于删除测试
      const createResponse = await request(API_BASE_URL)
        .post('/api/subscriptions/')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          topic_name: '要删除的订阅'
        });

      const deleteSubscriptionId = createResponse.body.data.id;

      const response = await request(API_BASE_URL)
        .delete(`/api/subscriptions/${deleteSubscriptionId}`)
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString());

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('订阅删除成功');
    });
  });
});

