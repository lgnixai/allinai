const request = require('supertest');
const { API_BASE_URL } = require('../env');

describe('话题管理 API 测试', () => {
  let authToken;
  let userId;
  let topicId;

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

  describe('获取话题列表', () => {
    test('获取话题列表', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/topics/')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .query({
          page: 1,
          page_size: 10
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('topics');
      expect(response.body.data).toHaveProperty('total');
      expect(Array.isArray(response.body.data.topics)).toBe(true);
    });
  });

  describe('创建话题', () => {
    test('创建话题', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/topics/')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          topic_name: '测试话题'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('话题创建成功');
      expect(response.body.data).toHaveProperty('id');
      expect(response.body.data).toHaveProperty('topic_name');

      topicId = response.body.data.id;
    });
  });

  describe('更新话题名称', () => {
    test('更新话题名称', async () => {
      const response = await request(API_BASE_URL)
        .put(`/api/topics/${topicId}`)
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          topic_name: '更新后的话题名称'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('话题名称更新成功');
      expect(response.body.data.topic_name).toBe('更新后的话题名称');
    });
  });

  describe('获取话题消息', () => {
    test('获取话题消息', async () => {
      const response = await request(API_BASE_URL)
        .get(`/api/topics/${topicId}/messages`)
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .query({
          page: 1,
          page_size: 20
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('messages');
      expect(response.body.data).toHaveProperty('total');
      expect(Array.isArray(response.body.data.messages)).toBe(true);
    });
  });

  describe('发送消息', () => {
    test('发送消息', async () => {
      const response = await request(API_BASE_URL)
        .post(`/api/topics/${topicId}/messages`)
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          content: '这是一条测试消息'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('消息发送成功');
      expect(response.body.data).toHaveProperty('id');
      expect(response.body.data).toHaveProperty('content');
      expect(response.body.data).toHaveProperty('topic_id');
      expect(response.body.data.content).toBe('这是一条测试消息');
    });
  });

  describe('自动创建话题并发送消息', () => {
    test('自动创建话题并发送消息', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/topics/0/messages')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          content: '这是一条会自动创建话题的消息，内容很长，应该会被截取作为话题标题'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('消息发送成功');
      expect(response.body.data).toHaveProperty('id');
      expect(response.body.data).toHaveProperty('content');
      expect(response.body.data).toHaveProperty('topic_id');
    });
  });

  describe('删除话题', () => {
    test('删除话题', async () => {
      // 先创建一个新话题用于删除测试
      const createResponse = await request(API_BASE_URL)
        .post('/api/topics/')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          topic_name: '要删除的话题'
        });

      const deleteTopicId = createResponse.body.data.id;

      const response = await request(API_BASE_URL)
        .delete(`/api/topics/${deleteTopicId}`)
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString());

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('话题删除成功');
    });
  });
});

