// Jest 测试设置文件
require('dotenv').config();

// 设置全局测试超时
jest.setTimeout(parseInt(process.env.TEST_TIMEOUT) || 30000);

// 全局测试数据
global.testData = {
  user: {
    phone: process.env.TEST_PHONE || '13800138000',
    verificationCode: process.env.TEST_VERIFICATION_CODE || '1111',
    displayName: process.env.TEST_DISPLAY_NAME || '测试用户',
    school: process.env.TEST_SCHOOL || '测试大学',
    college: process.env.TEST_COLLEGE || '计算机学院'
  },
  topic: {
    name: process.env.TEST_TOPIC_NAME || '测试话题',
    model: 'gpt-3.5-turbo',
    channelId: 1
  },
  subscription: {
    name: process.env.TEST_SUBSCRIPTION_NAME || '测试订阅',
    description: '测试订阅描述'
  }
};

// 全局测试工具函数
global.testUtils = {
  // 生成随机手机号
  generateRandomPhone: () => {
    return '138' + Math.random().toString().slice(2, 11);
  },
  
  // 生成随机字符串
  generateRandomString: (length = 10) => {
    return Math.random().toString(36).substring(2, length + 2);
  },
  
  // 等待函数
  wait: (ms) => new Promise(resolve => setTimeout(resolve, ms)),
  
  // 重试函数
  retry: async (fn, maxRetries = 3, delay = 1000) => {
    for (let i = 0; i < maxRetries; i++) {
      try {
        return await fn();
      } catch (error) {
        if (i === maxRetries - 1) throw error;
        await global.testUtils.wait(delay);
      }
    }
  }
};

// 测试环境检查
beforeAll(() => {
  console.log('🚀 开始运行 One-API 测试套件');
  console.log(`📡 API地址: ${process.env.API_BASE_URL}`);
  console.log(`📱 测试手机号: ${global.testData.user.phone}`);
});

// 测试完成清理
afterAll(() => {
  console.log('✅ One-API 测试套件运行完成');
});
