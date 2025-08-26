const request = require('supertest');
const { API_BASE_URL } = require('./env');

describe('系统推荐 API 测试', () => {
  let authToken;
  let userId;

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

  describe('获取系统推荐列表', () => {
    test('获取系统推荐列表', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/recommendations')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .query({
          page: 1,
          page_size: 10
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('recommendations');
      expect(response.body.data).toHaveProperty('total');
      expect(Array.isArray(response.body.data.recommendations)).toBe(true);
    });
  });

  describe('获取欢迎页面', () => {
    test('获取欢迎页面', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/welcome')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString());

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('recommendations');
      expect(Array.isArray(response.body.data.recommendations)).toBe(true);
    });
  });

  describe('获取推荐页面', () => {
    test('获取推荐页面', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/recommendations/change')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString());

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('recommendations');
      expect(Array.isArray(response.body.data.recommendations)).toBe(true);
    });
  });

  describe('搜索系统推荐', () => {
    test('搜索系统推荐', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/user/recommendations/search')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          keyword: '技术'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('recommendations');
      expect(Array.isArray(response.body.data.recommendations)).toBe(true);
    });

    test('搜索系统推荐 - 空关键词', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/user/recommendations/search')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          keyword: ''
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('recommendations');
      expect(Array.isArray(response.body.data.recommendations)).toBe(true);
    });
  });

  describe('获取单个系统推荐', () => {
    test('获取单个系统推荐', async () => {
      // 先获取推荐列表，然后获取第一个推荐的详情
      const listResponse = await request(API_BASE_URL)
        .get('/api/user/recommendations')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .query({
          page: 1,
          page_size: 1
        });

      if (listResponse.body.data.recommendations.length > 0) {
        const recommendationId = listResponse.body.data.recommendations[0].id;

        const response = await request(API_BASE_URL)
          .get(`/api/user/recommendations/${recommendationId}`)
          .set('Authorization', `Bearer ${authToken}`)
          .set('UserID', userId.toString());

        expect(response.status).toBe(200);
        expect(response.body.success).toBe(true);
        expect(response.body.data).toHaveProperty('id');
        expect(response.body.data).toHaveProperty('title');
        expect(response.body.data).toHaveProperty('content');
        expect(response.body.data).toHaveProperty('category');
        expect(response.body.data).toHaveProperty('created_at');
        expect(response.body.data).toHaveProperty('updated_at');
      }
    });

    test('获取单个系统推荐 - 不存在', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/recommendations/99999')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString());

      expect(response.status).toBe(404);
      expect(response.body.success).toBe(false);
      expect(response.body.message).toContain('推荐不存在');
    });
  });

  describe('管理员功能', () => {
    let adminToken;
    let adminUserId;

    beforeAll(async () => {
      // 创建管理员用户
      const registerResponse = await request(API_BASE_URL)
        .post('/api/user/register')
        .send({
          phone: '13800138001',
          phone_verification_code: '123456',
          display_name: '管理员用户',
          school: '测试大学',
          college: '计算机学院'
        });

      adminToken = registerResponse.body.data.access_token;
      adminUserId = registerResponse.body.data.user.id;

      // 这里需要手动设置用户为管理员角色，或者使用现有的管理员账户
      // 由于没有直接的管理员设置API，我们假设这个用户已经是管理员
    });

    describe('创建系统推荐', () => {
      test('创建系统推荐', async () => {
        const response = await request(API_BASE_URL)
          .post('/api/admin/recommendations')
          .set('Authorization', `Bearer ${adminToken}`)
          .set('UserID', adminUserId.toString())
          .send({
            title: '测试推荐标题',
            content: '测试推荐内容',
            category: '技术',
            description: '测试推荐描述'
          });

        expect(response.status).toBe(200);
        expect(response.body.success).toBe(true);
        expect(response.body.message).toContain('推荐创建成功');
        expect(response.body.data).toHaveProperty('id');
        expect(response.body.data).toHaveProperty('title');
        expect(response.body.data.title).toBe('测试推荐标题');
      });
    });

    describe('更新系统推荐', () => {
      test('更新系统推荐', async () => {
        // 先创建一个推荐
        const createResponse = await request(API_BASE_URL)
          .post('/api/admin/recommendations')
          .set('Authorization', `Bearer ${adminToken}`)
          .set('UserID', adminUserId.toString())
          .send({
            title: '要更新的推荐',
            content: '要更新的内容',
            category: '技术',
            description: '要更新的描述'
          });

        const recommendationId = createResponse.body.data.id;

        const response = await request(API_BASE_URL)
          .put(`/api/admin/recommendations/${recommendationId}`)
          .set('Authorization', `Bearer ${adminToken}`)
          .set('UserID', adminUserId.toString())
          .send({
            title: '更新后的推荐标题',
            content: '更新后的内容',
            category: '技术',
            description: '更新后的描述'
          });

        expect(response.status).toBe(200);
        expect(response.body.success).toBe(true);
        expect(response.body.message).toContain('推荐更新成功');
        expect(response.body.data.title).toBe('更新后的推荐标题');
      });
    });

    describe('删除系统推荐', () => {
      test('删除系统推荐', async () => {
        // 先创建一个推荐
        const createResponse = await request(API_BASE_URL)
          .post('/api/admin/recommendations')
          .set('Authorization', `Bearer ${adminToken}`)
          .set('UserID', adminUserId.toString())
          .send({
            title: '要删除的推荐',
            content: '要删除的内容',
            category: '技术',
            description: '要删除的描述'
          });

        const recommendationId = createResponse.body.data.id;

        const response = await request(API_BASE_URL)
          .delete(`/api/admin/recommendations/${recommendationId}`)
          .set('Authorization', `Bearer ${adminToken}`)
          .set('UserID', adminUserId.toString());

        expect(response.status).toBe(200);
        expect(response.body.success).toBe(true);
        expect(response.body.message).toContain('推荐删除成功');
      });
    });
  });

  describe('错误处理', () => {
    test('未授权访问', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/recommendations');

      expect(response.status).toBe(401);
      expect(response.body.success).toBe(false);
    });

    test('无效的访问令牌', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/recommendations')
        .set('Authorization', 'Bearer invalid_token')
        .set('UserID', userId.toString());

      expect(response.status).toBe(401);
      expect(response.body.success).toBe(false);
      expect(response.body.message).toContain('无效的访问令牌');
    });
  });
});


