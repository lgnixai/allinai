const axios = require('axios');

class ApiClient {
  constructor() {
    this.baseURL = process.env.API_BASE_URL || 'http://47.88.91.79:9999';
    this.token = null;
    this.userId = null;
    
    // 创建axios实例
    this.client = axios.create({
      baseURL: this.baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json'
      }
    });

    // 请求拦截器
    this.client.interceptors.request.use(
      (config) => {
        // 添加认证头
        if (this.token && this.userId) {
          config.headers.Authorization = this.token;
          config.headers.UserID = this.userId;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // 响应拦截器
    this.client.interceptors.response.use(
      (response) => {
        return response;
      },
      (error) => {
        console.error('API请求失败:', {
          url: error.config?.url,
          method: error.config?.method,
          status: error.response?.status,
          message: error.response?.data?.message || error.message
        });
        return Promise.reject(error);
      }
    );
  }

  // 设置认证信息
  setAuth(token, userId) {
    this.token = token;
    this.userId = userId;
  }

  // 清除认证信息
  clearAuth() {
    this.token = null;
    this.userId = null;
  }

  // 发送手机验证码
  async sendVerificationCode(phone, purpose) {
    const response = await this.client.get('/api/phone_verification', {
      params: { phone, purpose }
    });
    return response.data;
  }

  // 用户注册
  async register(userData) {
    const response = await this.client.post('/api/user/register', userData);
    return response.data;
  }

  // 用户登录
  async login(phone, verificationCode) {
    const response = await this.client.post('/api/user/login', {
      phone,
      phone_verification_code: verificationCode
    });
    
    if (response.data.success && response.data.data) {
      this.setAuth(response.data.data.access_token, response.data.data.id);
    }
    
    return response.data;
  }

  // 获取用户信息
  async getUserInfo() {
    const response = await this.client.get('/api/user/self');
    return response.data;
  }

  // 更新用户信息
  async updateUserInfo(userData) {
    const response = await this.client.put('/api/user/self', userData);
    return response.data;
  }

  // 用户登出
  async logout() {
    const response = await this.client.get('/api/user/logout');
    this.clearAuth();
    return response.data;
  }

  // 获取话题列表
  async getTopics(page = 1, size = 10) {
    const response = await this.client.get('/api/topics/', {
      params: { page, size }
    });
    return response.data;
  }

  // 创建话题
  async createTopic(topicData) {
    const response = await this.client.post('/api/topics/', topicData);
    return response.data;
  }

  // 删除话题
  async deleteTopic(topicId) {
    const response = await this.client.delete(`/api/topics/${topicId}`);
    return response.data;
  }

  // 获取话题消息
  async getTopicMessages(topicId, page = 1, size = 20) {
    const response = await this.client.get(`/api/topics/${topicId}/messages`, {
      params: { page, size }
    });
    return response.data;
  }

  // 发送消息
  async sendMessage(topicId, content) {
    const response = await this.client.post(`/api/topics/${topicId}/messages`, {
      content
    });
    return response.data;
  }

  // 获取订阅列表
  async getSubscriptions(page = 1, size = 10) {
    const response = await this.client.get('/api/subscriptions/', {
      params: { page, size }
    });
    return response.data;
  }

  // 创建订阅
  async createSubscription(subscriptionData) {
    const response = await this.client.post('/api/subscriptions/', subscriptionData);
    return response.data;
  }

  // 取消订阅
  async cancelSubscription(subscriptionId) {
    const response = await this.client.put(`/api/subscriptions/${subscriptionId}/cancel`);
    return response.data;
  }

  // 重新激活订阅
  async reactivateSubscription(subscriptionId) {
    const response = await this.client.put(`/api/subscriptions/${subscriptionId}/reactivate`);
    return response.data;
  }

  // 删除订阅
  async deleteSubscription(subscriptionId) {
    const response = await this.client.delete(`/api/subscriptions/${subscriptionId}`);
    return response.data;
  }

  // 获取订阅文章
  async getSubscriptionArticles(subscriptionId, page = 1, size = 10) {
    const response = await this.client.get(`/api/subscriptions/${subscriptionId}/articles`, {
      params: { page, size }
    });
    return response.data;
  }

  // 健康检查
  async healthCheck() {
    try {
      const response = await this.client.get('/api/status');
      return response.data;
    } catch (error) {
      throw new Error(`API服务不可用: ${error.message}`);
    }
  }
}

module.exports = ApiClient;
