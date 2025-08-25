const request = require('supertest');
const { API_BASE_URL } = require('../env');

describe('用户管理 API 测试', () => {
  let authToken;
  let userId;

  describe('手机验证码相关', () => {
    test('发送注册验证码', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/phone_verification')
        .query({
          phone: '13800138000',
          purpose: 'register'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('验证码发送成功');
    });

    test('发送登录验证码', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/phone_verification')
        .query({
          phone: '13800138000',
          purpose: 'login'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('验证码发送成功');
    });
  });

  describe('用户注册', () => {
    test('用户注册', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/user/register')
        .send({
          phone: '13800138000',
          phone_verification_code: '123456',
          display_name: '测试用户',
          school: '测试大学',
          college: '计算机学院'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('注册成功');
      expect(response.body.data).toHaveProperty('access_token');
      expect(response.body.data).toHaveProperty('user');
      
      authToken = response.body.data.access_token;
      userId = response.body.data.user.id;
    });

    test('用户注册 - 手机号已存在', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/user/register')
        .send({
          phone: '13800138000',
          phone_verification_code: '123456',
          display_name: '测试用户2',
          school: '测试大学',
          college: '计算机学院'
        });

      expect(response.status).toBe(400);
      expect(response.body.success).toBe(false);
      expect(response.body.message).toContain('手机号已存在');
    });
  });

  describe('用户登录', () => {
    test('用户登录', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/user/login')
        .send({
          phone: '13800138000',
          phone_verification_code: '123456'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('登录成功');
      expect(response.body.data).toHaveProperty('access_token');
      expect(response.body.data).toHaveProperty('user');
      
      authToken = response.body.data.access_token;
      userId = response.body.data.user.id;
    });

    test('用户登录 - 验证码错误', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/user/login')
        .send({
          phone: '13800138000',
          phone_verification_code: '000000'
        });

      expect(response.status).toBe(400);
      expect(response.body.success).toBe(false);
      expect(response.body.message).toContain('验证码错误');
    });

    test('用户登录 - 用户不存在', async () => {
      const response = await request(API_BASE_URL)
        .post('/api/user/login')
        .send({
          phone: '13900139000',
          phone_verification_code: '123456'
        });

      expect(response.status).toBe(404);
      expect(response.body.success).toBe(false);
      expect(response.body.message).toContain('用户不存在');
    });
  });

  describe('用户信息管理', () => {
    test('获取用户信息', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/info')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString());

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.data).toHaveProperty('id');
      expect(response.body.data).toHaveProperty('phone');
      expect(response.body.data).toHaveProperty('display_name');
      expect(response.body.data).toHaveProperty('school');
      expect(response.body.data).toHaveProperty('college');
    });

    test('获取用户信息 - 未授权', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/info');

      expect(response.status).toBe(401);
      expect(response.body.success).toBe(false);
    });

    test('更新用户信息', async () => {
      const response = await request(API_BASE_URL)
        .put('/api/user/info')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', userId.toString())
        .send({
          display_name: '更新后的用户名',
          school: '更新后的大学',
          college: '更新后的学院'
        });

      expect(response.status).toBe(200);
      expect(response.body.success).toBe(true);
      expect(response.body.message).toContain('用户信息更新成功');
      expect(response.body.data.display_name).toBe('更新后的用户名');
      expect(response.body.data.school).toBe('更新后的大学');
      expect(response.body.data.college).toBe('更新后的学院');
    });

    test('更新用户信息 - 未授权', async () => {
      const response = await request(API_BASE_URL)
        .put('/api/user/info')
        .send({
          display_name: '测试用户名'
        });

      expect(response.status).toBe(401);
      expect(response.body.success).toBe(false);
    });
  });

  describe('错误处理', () => {
    test('无效的访问令牌', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/info')
        .set('Authorization', 'Bearer invalid_token')
        .set('UserID', '1');

      expect(response.status).toBe(401);
      expect(response.body.success).toBe(false);
      expect(response.body.message).toContain('无效的访问令牌');
    });

    test('用户ID不匹配', async () => {
      const response = await request(API_BASE_URL)
        .get('/api/user/info')
        .set('Authorization', `Bearer ${authToken}`)
        .set('UserID', '999');

      expect(response.status).toBe(403);
      expect(response.body.success).toBe(false);
      expect(response.body.message).toContain('用户ID不匹配');
    });
  });
});

