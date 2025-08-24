const axios = require('axios');

// 测试配置
const BASE_URL = 'http://localhost:4000';

// 创建axios实例
const api = axios.create({
  baseURL: BASE_URL,
  timeout: 10000,
});

// 测试文档列表API
async function testDocsList() {
  try {
    console.log('📋 测试文档列表API...');
    const response = await api.get('/api/docs/list');
    
    if (response.data.success) {
      console.log('✅ 文档列表获取成功');
      console.log('📚 可用文档:');
      response.data.data.forEach(doc => {
        console.log(`  - ${doc.title}: ${doc.url}`);
      });
    } else {
      console.log('❌ 文档列表获取失败:', response.data.message);
    }
  } catch (error) {
    console.log('❌ 文档列表请求失败:', error.message);
  }
}

// 测试文档首页
async function testDocsIndex() {
  try {
    console.log('\n🏠 测试文档首页...');
    const response = await api.get('/api/docs/');
    
    if (response.status === 200) {
      console.log('✅ 文档首页访问成功');
      console.log('📄 响应类型:', response.headers['content-type']);
      console.log('📏 内容长度:', response.data.length, '字符');
    } else {
      console.log('❌ 文档首页访问失败:', response.status);
    }
  } catch (error) {
    console.log('❌ 文档首页请求失败:', error.message);
  }
}

// 测试具体文档
async function testSpecificDoc(docType) {
  try {
    console.log(`\n📖 测试${docType}文档...`);
    const response = await api.get(`/api/docs/${docType}`);
    
    if (response.status === 200) {
      console.log(`✅ ${docType}文档访问成功`);
      console.log('📄 响应类型:', response.headers['content-type']);
      console.log('📏 内容长度:', response.data.length, '字符');
      
      // 检查是否是Markdown内容
      if (response.data.includes('#')) {
        console.log('✅ 确认是Markdown格式文档');
      }
    } else {
      console.log(`❌ ${docType}文档访问失败:`, response.status);
    }
  } catch (error) {
    console.log(`❌ ${docType}文档请求失败:`, error.message);
  }
}

// 测试不存在的文档
async function testNonExistentDoc() {
  try {
    console.log('\n🚫 测试不存在的文档...');
    const response = await api.get('/api/docs/nonexistent');
    
    if (response.status === 404) {
      console.log('✅ 正确处理不存在的文档');
    } else {
      console.log('❌ 未正确处理不存在的文档:', response.status);
    }
  } catch (error) {
    if (error.response && error.response.status === 404) {
      console.log('✅ 正确处理不存在的文档');
    } else {
      console.log('❌ 处理不存在的文档时出错:', error.message);
    }
  }
}

// 主测试函数
async function runTests() {
  console.log('🚀 开始测试One-API文档功能\n');
  
  await testDocsList();
  await testDocsIndex();
  await testSpecificDoc('api');
  await testSpecificDoc('postman');
  await testSpecificDoc('deployment');
  await testSpecificDoc('auth');
  await testNonExistentDoc();
  
  console.log('\n✨ 测试完成！');
  console.log('\n📝 使用说明:');
  console.log('  - 文档首页: http://localhost:4000/api/docs/');
  console.log('  - API文档: http://localhost:4000/api/docs/api');
  console.log('  - 文档列表: http://localhost:4000/api/docs/list');
}

// 运行测试
runTests().catch(console.error);
