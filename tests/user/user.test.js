const ApiClient = require('../utils/apiClient');

describe('用户管理模块测试', () => {
  let apiClient;
  let testUserData;

  beforeAll(() => {
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
      await apiClient.logout();
    } catch (error) {
      console.log('清理测试数据时出错:', error.message);
    }
  });

  describe('API健康检查', () => {
    test('应该能够连接到API服务器', async () => {
      const result = await apiClient.healthCheck();
      expect(result).toBeDefined();
    });
  });

  describe('手机验证码功能', () => {
    test('应该能够发送注册验证码', async () => {
      const result = await apiClient.sendVerificationCode(
        testUserData.phone,
        'register'
      );
      
      expect(result.success).toBe(true);
      expect(result.data).toBe(global.testData.user.verificationCode);
    });

    test('应该能够发送登录验证码', async () => {
      const result = await apiClient.sendVerificationCode(
        testUserData.phone,
        'login'
      );
      
      expect(result.success).toBe(true);
      expect(result.data).toBe(global.testData.user.verificationCode);
    });

    test('应该拒绝无效的手机号格式', async () => {
      await expect(
        apiClient.sendVerificationCode('123', 'register')
      ).rejects.toThrow();
    });
  });

  describe('用户注册功能', () => {
    test('应该能够成功注册新用户', async () => {
      const result = await apiClient.register(testUserData);
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('成功');
    });

    test('应该拒绝重复注册同一手机号', async () => {
      await expect(
        apiClient.register(testUserData)
      ).rejects.toThrow();
    });

    test('应该拒绝无效的验证码', async () => {
      const invalidData = {
        ...testUserData,
        phone: global.testUtils.generateRandomPhone(),
        phone_verification_code: '9999'
      };
      
      await expect(
        apiClient.register(invalidData)
      ).rejects.toThrow();
    });
  });

  describe('用户登录功能', () => {
    test('应该能够成功登录已注册用户', async () => {
      const result = await apiClient.login(
        testUserData.phone,
        testUserData.phone_verification_code
      );
      
      expect(result.success).toBe(true);
      expect(result.data).toBeDefined();
      expect(result.data.access_token).toBeDefined();
      expect(result.data.id).toBeDefined();
      expect(result.data.phone).toBe(testUserData.phone);
    });

    test('应该拒绝未注册用户的登录', async () => {
      const unregisteredPhone = global.testUtils.generateRandomPhone();
      
      await expect(
        apiClient.login(unregisteredPhone, testUserData.phone_verification_code)
      ).rejects.toThrow();
    });

    test('应该拒绝错误的验证码', async () => {
      await expect(
        apiClient.login(testUserData.phone, '9999')
      ).rejects.toThrow();
    });
  });

  describe('用户信息管理', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.login(
        testUserData.phone,
        testUserData.phone_verification_code
      );
    });

    test('应该能够获取用户信息', async () => {
      const result = await apiClient.getUserInfo();
      
      expect(result.success).toBe(true);
      expect(result.data).toBeDefined();
      expect(result.data.phone).toBe(testUserData.phone);
      expect(result.data.display_name).toBe(testUserData.display_name);
    });

    test('应该能够更新用户信息', async () => {
      const updateData = {
        display_name: '更新后的显示名称',
        school: '更新后的学校',
        college: '更新后的学院'
      };
      
      const result = await apiClient.updateUserInfo(updateData);
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('成功');
      
      // 验证更新是否生效
      const userInfo = await apiClient.getUserInfo();
      expect(userInfo.data.display_name).toBe(updateData.display_name);
      expect(userInfo.data.school).toBe(updateData.school);
      expect(userInfo.data.college).toBe(updateData.college);
    });

    test('应该能够部分更新用户信息', async () => {
      const partialUpdate = {
        display_name: '部分更新名称'
      };
      
      const result = await apiClient.updateUserInfo(partialUpdate);
      
      expect(result.success).toBe(true);
      
      // 验证部分更新是否生效
      const userInfo = await apiClient.getUserInfo();
      expect(userInfo.data.display_name).toBe(partialUpdate.display_name);
    });
  });

  describe('用户登出功能', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.login(
        testUserData.phone,
        testUserData.phone_verification_code
      );
    });

    test('应该能够成功登出', async () => {
      const result = await apiClient.logout();
      
      expect(result.success).toBe(true);
      expect(result.message).toContain('成功');
    });

    test('登出后应该无法访问需要认证的接口', async () => {
      await apiClient.logout();
      
      await expect(
        apiClient.getUserInfo()
      ).rejects.toThrow();
    });
  });

  describe('认证机制', () => {
    beforeEach(async () => {
      // 确保用户已登录
      await apiClient.login(
        testUserData.phone,
        testUserData.phone_verification_code
      );
    });

    test('应该正确设置认证头', async () => {
      const userInfo = await apiClient.getUserInfo();
      expect(userInfo.success).toBe(true);
    });

    test('应该拒绝无效的认证信息', async () => {
      apiClient.setAuth('invalid_token', '999');
      
      await expect(
        apiClient.getUserInfo()
      ).rejects.toThrow();
    });
  });

  describe('错误处理', () => {
    test('应该正确处理网络错误', async () => {
      const invalidClient = new ApiClient();
      invalidClient.baseURL = 'http://invalid-url:9999';
      
      await expect(
        invalidClient.healthCheck()
      ).rejects.toThrow();
    });

    test('应该正确处理服务器错误', async () => {
      await expect(
        apiClient.client.get('/api/nonexistent-endpoint')
      ).rejects.toThrow();
    });
  });
});
