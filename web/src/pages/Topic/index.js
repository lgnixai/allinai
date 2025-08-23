import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { API, showError, showSuccess } from '../../helpers/index.js';
import {
  Layout,
  Card,
  Button,
  Typography,
  Modal,
  Input,
  Space,
  List,
  Tag,
  Empty,
  Spin,
  Table,
  Pagination,
  Popconfirm,
  TextArea,
} from '@douyinfe/semi-ui';
import {
  IconPlus,
  IconDelete,
  IconComment,
  IconSend,
} from '@douyinfe/semi-icons';

const { Content, Sider } = Layout;
const { Title, Text } = Typography;

const Topic = () => {
  const { t } = useTranslation();
  const [topics, setTopics] = useState([]);
  const [selectedTopic, setSelectedTopic] = useState(null);
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [messagesLoading, setMessagesLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [formData, setFormData] = useState({
    topic_name: '',
    model: 'gpt-3.5-turbo',
    channel_id: 1
  });
  const [total, setTotal] = useState(0);
  const [messagesTotal, setMessagesTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [messagesPage, setMessagesPage] = useState(1);
  const [pageSize] = useState(10);
  const [messagesPageSize] = useState(10);
  const [chatInput, setChatInput] = useState('');
  const [chatLoading, setChatLoading] = useState(false);

  // 获取话题列表
  const fetchTopics = async (page = 1) => {
    setLoading(true);
    try {
      const response = await API.get(`/api/topics/?page=${page}&page_size=${pageSize}`);
      if (response.data.success) {
        setTopics(response.data.data.topics || []);
        setTotal(response.data.data.total || 0);
        setCurrentPage(page);
      } else {
        showError(response.data.message || '获取话题列表失败');
      }
    } catch (error) {
      showError('获取话题列表失败: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  // 获取话题下的消息
  const fetchMessages = async (topicId, page = 1) => {
    if (!topicId) return;
    
    setMessagesLoading(true);
    try {
      const response = await API.get(`/api/topics/${topicId}/messages?page=${page}&page_size=${messagesPageSize}`);
      if (response.data.success) {
        setMessages(response.data.data.messages || []);
        setMessagesTotal(response.data.data.total || 0);
        setMessagesPage(page);
      } else {
        showError(response.data.message || '获取消息列表失败');
      }
    } catch (error) {
      showError('获取消息列表失败: ' + error.message);
    } finally {
      setMessagesLoading(false);
    }
  };

  // 创建话题
  const createTopic = async () => {
    try {
      const response = await API.post('/api/topics/', formData);
      if (response.data.success) {
        showSuccess('话题创建成功');
        setModalVisible(false);
        setFormData({ topic_name: '', model: 'gpt-3.5-turbo', channel_id: 1 });
        fetchTopics(currentPage);
      } else {
        showError(response.data.message || '创建话题失败');
      }
    } catch (error) {
      showError('创建话题失败: ' + error.message);
    }
  };

  // 删除话题
  const deleteTopic = async (id) => {
    try {
      const response = await API.delete(`/api/topics/${id}`);
      if (response.data.success) {
        showSuccess('话题已删除');
        fetchTopics(currentPage);
        if (selectedTopic && selectedTopic.id === id) {
          setSelectedTopic(null);
          setMessages([]);
        }
      } else {
        showError(response.data.message || '删除话题失败');
      }
    } catch (error) {
      showError('删除话题失败: ' + error.message);
    }
  };

  // 发送消息
  const sendMessage = async () => {
    if (!selectedTopic || !chatInput.trim()) return;

    setChatLoading(true);
    try {
      const response = await API.post(`/api/topics/${selectedTopic.id}/messages`, {
        content: chatInput,
        role: 'user'
      });
      
      if (response.data.success) {
        setChatInput('');
        // 重新获取消息列表
        fetchMessages(selectedTopic.id, messagesPage);
      } else {
        showError(response.data.message || '发送消息失败');
      }
    } catch (error) {
      showError('发送消息失败: ' + error.message);
    } finally {
      setChatLoading(false);
    }
  };

  // 选择话题
  const handleSelectTopic = (topic) => {
    setSelectedTopic(topic);
    fetchMessages(topic.id, 1);
  };

  // 处理分页
  const handlePageChange = (page) => {
    fetchTopics(page);
  };

  const handleMessagesPageChange = (page) => {
    if (selectedTopic) {
      fetchMessages(selectedTopic.id, page);
    }
  };

  // 初始化
  useEffect(() => {
    fetchTopics(1);
  }, []);

  // 消息表格列定义
  const messageColumns = [
    {
      title: t('角色'),
      dataIndex: 'role',
      key: 'role',
      render: (text) => (
        <Tag color={text === 'user' ? 'blue' : 'green'} size="small">
          {text === 'user' ? t('用户') : t('AI助手')}
        </Tag>
      ),
    },
    {
      title: t('内容'),
      dataIndex: 'content',
      key: 'content',
      render: (text) => (
        <Text ellipsis={{ showTooltip: true }} style={{ maxWidth: 300 }}>
          {text}
        </Text>
      ),
    },
    {
      title: t('创建时间'),
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text) => new Date(text).toLocaleString(),
    },
  ];

  return (
    <Layout style={{ height: 'calc(100vh - 120px)', padding: '20px' }}>
      <Sider width={400} style={{ backgroundColor: 'transparent', marginRight: '20px' }}>
        <Card
          title={
            <Space>
              <IconComment />
              <Title level={4}>{t('我的话题')}</Title>
            </Space>
          }
          style={{ height: '100%' }}
        >
          {/* 头部添加话题按钮 */}
          <div style={{ marginBottom: '16px', textAlign: 'center' }}>
            <Button
              type="primary"
              icon={<IconPlus />}
              onClick={() => setModalVisible(true)}
              size="large"
              style={{ width: '100%' }}
            >
              {t('添加话题')}
            </Button>
          </div>

          <Spin spinning={loading}>
            {topics.length === 0 ? (
              <div style={{ textAlign: 'center', padding: '20px' }}>
                <Empty description={t('暂无话题')} />
              </div>
            ) : (
              <List
                dataSource={topics}
                renderItem={(item) => (
                  <List.Item
                    key={item.id}
                    style={{
                      cursor: 'pointer',
                      backgroundColor: selectedTopic?.id === item.id ? 'var(--semi-color-fill-0)' : 'transparent',
                      borderRadius: '6px',
                      padding: '8px',
                      marginBottom: '8px',
                    }}
                    onClick={() => handleSelectTopic(item)}
                  >
                    <div style={{ width: '100%' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <Text strong>{item.topic_name}</Text>
                        <Space>
                          <Tag color="blue" size="small">
                            {item.message_count || 0} {t('条消息')}
                          </Tag>
                          <Popconfirm
                            title={t('确定要删除话题吗？删除后无法恢复！')}
                            onConfirm={() => deleteTopic(item.id)}
                          >
                            <Button
                              type="danger"
                              size="small"
                              icon={<IconDelete />}
                              onClick={(e) => e.stopPropagation()}
                            >
                              {t('删除')}
                            </Button>
                          </Popconfirm>
                        </Space>
                      </div>
                      <Text type="secondary" size="small" style={{ display: 'block', marginTop: '4px' }}>
                        {t('模型')}: {item.model}
                      </Text>
                      <Text type="tertiary" size="small" style={{ display: 'block', marginTop: '4px' }}>
                        {new Date(item.created_at).toLocaleDateString()}
                      </Text>
                    </div>
                  </List.Item>
                )}
              />
            )}
          </Spin>
          
          {total > pageSize && (
            <div style={{ marginTop: '20px', textAlign: 'center' }}>
              <Pagination
                currentPage={currentPage}
                pageSize={pageSize}
                total={total}
                onPageChange={handlePageChange}
                showSizeChanger={false}
              />
            </div>
          )}
        </Card>
      </Sider>

      <Content>
        <Card
          title={
            <Space>
              <IconComment />
              <Title level={4}>
                {selectedTopic ? `${selectedTopic.topic_name} - ${t('聊天记录')}` : t('聊天记录')}
              </Title>
            </Space>
          }
          style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
        >
          {!selectedTopic ? (
            <Empty description={t('请选择一个话题查看聊天记录')} style={{ marginTop: 60 }} />
          ) : (
            <>
              {/* 消息列表 */}
              <div style={{ flex: 1, overflow: 'hidden' }}>
                <Spin spinning={messagesLoading}>
                  <Table
                    columns={messageColumns}
                    dataSource={messages}
                    pagination={false}
                    rowKey="id"
                    empty={<Empty description={t('暂无消息')} />}
                    scroll={{ y: 300 }}
                  />
                </Spin>
                
                {messagesTotal > messagesPageSize && (
                  <div style={{ marginTop: '20px', textAlign: 'center' }}>
                    <Pagination
                      currentPage={messagesPage}
                      pageSize={messagesPageSize}
                      total={messagesTotal}
                      onPageChange={handleMessagesPageChange}
                      showSizeChanger={false}
                    />
                  </div>
                )}
              </div>

              {/* 聊天输入框 */}
              <div style={{ marginTop: '20px', borderTop: '1px solid var(--semi-color-border)', paddingTop: '20px' }}>
                <div style={{ display: 'flex', gap: '10px' }}>
                  <TextArea
                    value={chatInput}
                    onChange={setChatInput}
                    placeholder={t('请输入消息内容...')}
                    rows={3}
                    style={{ flex: 1 }}
                    onKeyPress={(e) => {
                      if (e.key === 'Enter' && !e.shiftKey) {
                        e.preventDefault();
                        sendMessage();
                      }
                    }}
                  />
                  <Button
                    type="primary"
                    icon={<IconSend />}
                    onClick={sendMessage}
                    loading={chatLoading}
                    disabled={!chatInput.trim()}
                    style={{ alignSelf: 'flex-end' }}
                  >
                    {t('发送')}
                  </Button>
                </div>
              </div>
            </>
          )}
        </Card>
      </Content>

      {/* 添加话题模态框 */}
      <Modal
        title={t('创建新话题')}
        visible={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          setFormData({ topic_name: '', model: 'gpt-3.5-turbo', channel_id: 1 });
        }}
        footer={null}
        width={500}
      >
        <div style={{ padding: '20px 0' }}>
          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: '500' }}>
              {t('话题名称')} *
            </label>
            <Input
              value={formData.topic_name}
              onChange={(value) => setFormData({ ...formData, topic_name: value })}
              placeholder={t('请输入话题名称')}
              maxLength={100}
            />
          </div>
          
          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: '500' }}>
              {t('AI模型')}
            </label>
            <Input
              value={formData.model}
              onChange={(value) => setFormData({ ...formData, model: value })}
              placeholder={t('请输入AI模型')}
            />
          </div>

          <div style={{ marginBottom: '20px' }}>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: '500' }}>
              {t('渠道ID')}
            </label>
            <Input
              value={formData.channel_id.toString()}
              onChange={(value) => setFormData({ ...formData, channel_id: parseInt(value) || 1 })}
              placeholder={t('请输入渠道ID')}
              type="number"
            />
          </div>
          
          <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
            <Button
              onClick={() => {
                setModalVisible(false);
                setFormData({ topic_name: '', model: 'gpt-3.5-turbo', channel_id: 1 });
              }}
            >
              {t('取消')}
            </Button>
            <Button
              type="primary"
              onClick={createTopic}
              disabled={!formData.topic_name.trim()}
            >
              {t('确定')}
            </Button>
          </Space>
        </div>
      </Modal>
    </Layout>
  );
};

export default Topic;
