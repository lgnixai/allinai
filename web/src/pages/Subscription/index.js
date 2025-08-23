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
} from '@douyinfe/semi-ui';
import {
  IconPlus,
  IconDelete,
  IconBookmark,
  IconTextStroked,
  IconRefresh,
} from '@douyinfe/semi-icons';

const { Content, Sider } = Layout;
const { Title, Text } = Typography;

const Subscription = () => {
  const { t } = useTranslation();
  const [subscriptions, setSubscriptions] = useState([]);
  const [selectedSubscription, setSelectedSubscription] = useState(null);
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [articlesLoading, setArticlesLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [formData, setFormData] = useState({
    topic_name: '',
    topic_description: ''
  });
  const [total, setTotal] = useState(0);
  const [articlesTotal, setArticlesTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [articlesPage, setArticlesPage] = useState(1);
  const [pageSize] = useState(10);
  const [articlesPageSize] = useState(10);

  // 获取订阅列表
  const fetchSubscriptions = async (page = 1) => {
    setLoading(true);
    try {
      const response = await API.get(`/api/subscriptions/?page=${page}&page_size=${pageSize}`);
      if (response.data.success) {
        setSubscriptions(response.data.data.subscriptions || []);
        setTotal(response.data.data.total || 0);
        setCurrentPage(page);
      } else {
        showError(response.data.message || '获取订阅列表失败');
      }
    } catch (error) {
      showError('获取订阅列表失败: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  // 重新激活订阅
  const reactivateSubscription = async (id) => {
    try {
      const response = await API.put(`/api/subscriptions/${id}/reactivate`);
      if (response.data.success) {
        showSuccess('订阅重新激活成功');
        fetchSubscriptions(currentPage);
      } else {
        showError(response.data.message || '重新激活订阅失败');
      }
    } catch (error) {
      showError('重新激活订阅失败: ' + error.message);
    }
  };

  // 获取订阅文章列表
  const fetchArticles = async (subscriptionId, page = 1) => {
    if (!subscriptionId) return;
    
    setArticlesLoading(true);
    try {
      const response = await API.get(`/api/subscriptions/${subscriptionId}/articles?page=${page}&page_size=${articlesPageSize}`);
      if (response.data.success) {
        setArticles(response.data.data.articles || []);
        setArticlesTotal(response.data.data.total || 0);
        setArticlesPage(page);
      } else {
        showError(response.data.message || '获取文章列表失败');
      }
    } catch (error) {
      showError('获取文章列表失败: ' + error.message);
    } finally {
      setArticlesLoading(false);
    }
  };

  // 创建订阅
  const createSubscription = async () => {
    try {
      const response = await API.post('/api/subscriptions/', formData);
      if (response.data.success) {
        showSuccess('订阅创建成功');
        setModalVisible(false);
        setFormData({ topic_name: '', topic_description: '' });
        fetchSubscriptions(currentPage);
      } else {
        showError(response.data.message || '创建订阅失败');
      }
    } catch (error) {
      showError('创建订阅失败: ' + error.message);
    }
  };

  // 取消订阅
  const cancelSubscription = async (id) => {
    try {
      const response = await API.put(`/api/subscriptions/${id}/cancel`);
      if (response.data.success) {
        showSuccess('订阅已取消');
        fetchSubscriptions(currentPage);
        if (selectedSubscription && selectedSubscription.id === id) {
          setSelectedSubscription(null);
          setArticles([]);
        }
      } else {
        showError(response.data.message || '取消订阅失败');
      }
    } catch (error) {
      showError('取消订阅失败: ' + error.message);
    }
  };

  // 删除订阅
  const deleteSubscription = async (id) => {
    try {
      const response = await API.delete(`/api/subscriptions/${id}`);
      if (response.data.success) {
        showSuccess('订阅已删除');
        fetchSubscriptions(currentPage);
        if (selectedSubscription && selectedSubscription.id === id) {
          setSelectedSubscription(null);
          setArticles([]);
        }
      } else {
        showError(response.data.message || '删除订阅失败');
      }
    } catch (error) {
      showError('删除订阅失败: ' + error.message);
    }
  };

  // 选择订阅
  const handleSelectSubscription = (subscription) => {
    setSelectedSubscription(subscription);
    fetchArticles(subscription.id, 1);
  };

  // 处理分页
  const handlePageChange = (page) => {
    fetchSubscriptions(page);
  };

  const handleArticlesPageChange = (page) => {
    if (selectedSubscription) {
      fetchArticles(selectedSubscription.id, page);
    }
  };

  // 初始化
  useEffect(() => {
    fetchSubscriptions(1);
  }, []);

  // 文章表格列定义
  const articleColumns = [
    {
      title: t('标题'),
      dataIndex: 'title',
      key: 'title',
      render: (text) => (
        <Text ellipsis={{ showTooltip: true }} style={{ maxWidth: 200 }}>
          {text}
        </Text>
      ),
    },
    {
      title: t('作者'),
      dataIndex: 'author',
      key: 'author',
      render: (text) => text || '-',
    },
    {
      title: t('发布时间'),
      dataIndex: 'published_at',
      key: 'published_at',
      render: (text) => {
        if (!text) return '-';
        return new Date(text).toLocaleString();
      },
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
              <IconBookmark />
              <Title level={4}>{t('我的订阅')}</Title>
            </Space>
          }
          style={{ height: '100%' }}
        >
          {/* 头部添加订阅按钮 */}
          <div style={{ marginBottom: '16px', textAlign: 'center' }}>
            <Button
              type="primary"
              icon={<IconPlus />}
              onClick={() => setModalVisible(true)}
              size="large"
              style={{ width: '100%' }}
            >
              {t('添加订阅')}
            </Button>
          </div>
          <Spin spinning={loading}>
            {subscriptions.length === 0 ? (
              <div style={{ textAlign: 'center', padding: '20px' }}>
                <Empty description={t('暂无订阅')} />
              </div>
            ) : (
              <List
                dataSource={subscriptions}
                renderItem={(item) => (
                  <List.Item
                    key={item.id}
                    style={{
                      cursor: 'pointer',
                      backgroundColor: selectedSubscription?.id === item.id ? 'var(--semi-color-fill-0)' : 'transparent',
                      borderRadius: '6px',
                      padding: '8px',
                      marginBottom: '8px',
                    }}
                    onClick={() => handleSelectSubscription(item)}
                  >
                    <div style={{ width: '100%' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <Text strong>{item.topic_name}</Text>
                        <Space>
                          <Tag color={item.status === 1 ? "blue" : "grey"} size="small">
                            {item.article_count} {t('篇文章')}
                          </Tag>
                          {item.status === 1 ? (
                            // 活跃状态：显示取消和删除按钮
                            <>
                              <Popconfirm
                                title={t('确定要取消订阅吗？')}
                                onConfirm={() => cancelSubscription(item.id)}
                              >
                                <Button
                                  type="warning"
                                  size="small"
                                  onClick={(e) => e.stopPropagation()}
                                >
                                  {t('取消')}
                                </Button>
                              </Popconfirm>

                            </>
                          ) : (
                            // 取消状态：显示重新订阅和删除按钮
                            <>
                              {/*<Button*/}
                              {/*  type="primary"*/}
                              {/*  size="small"*/}
                              {/*  onClick={(e) => {*/}
                              {/*    e.stopPropagation();*/}
                              {/*    reactivateSubscription(item.id);*/}
                              {/*  }}*/}
                              {/*>{item.status}*/}
                              {/*  {t('重新订阅')}*/}
                              {/*</Button>*/}
                              {/*<Popconfirm*/}
                              {/*  title={t('确定要删除订阅吗？删除后无法恢复！')}*/}
                              {/*  onConfirm={() => deleteSubscription(item.id)}*/}
                              {/*>*/}
                              {/*  <Button*/}
                              {/*    type="danger"*/}
                              {/*    size="small"*/}
                              {/*    icon={<IconDelete />}*/}
                              {/*    onClick={(e) => e.stopPropagation()}*/}
                              {/*  >*/}
                              {/*    {t('删除')}*/}
                              {/*  </Button>*/}
                              {/*</Popconfirm>*/}
                            </>
                          )}
                        </Space>
                      </div>
                      {item.topic_description && (
                        <Text type="secondary" size="small" style={{ display: 'block', marginTop: '4px' }}>
                          {item.topic_description}
                        </Text>
                      )}
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
              <IconTextStroked />
              <Title level={4}>
                {selectedSubscription ? `${selectedSubscription.topic_name} - ${t('文章列表')}` : t('文章列表')}
              </Title>
            </Space>
          }
          extra={
            selectedSubscription && (
              <Button
                icon={<IconRefresh />}
                onClick={() => fetchArticles(selectedSubscription.id, articlesPage)}
              >
                {t('刷新')}
              </Button>
            )
          }
          style={{ height: '100%' }}
        >
          {!selectedSubscription ? (
            <Empty description={t('请选择一个订阅查看文章')} style={{ marginTop: 60 }} />
          ) : (
            <>
              <Spin spinning={articlesLoading}>
                <Table
                  columns={articleColumns}
                  dataSource={articles}
                  pagination={false}
                  rowKey="id"
                  empty={<Empty description={t('暂无文章')} />}
                />
              </Spin>
              
              {articlesTotal > articlesPageSize && (
                <div style={{ marginTop: '20px', textAlign: 'center' }}>
                  <Pagination
                    currentPage={articlesPage}
                    pageSize={articlesPageSize}
                    total={articlesTotal}
                    onPageChange={handleArticlesPageChange}
                    showSizeChanger={false}
                  />
                </div>
              )}
            </>
          )}
        </Card>
      </Content>

      {/* 添加订阅模态框 */}
      <Modal
        title={t('创建新订阅')}
        visible={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          setFormData({ topic_name: '', topic_description: '' });
        }}
        footer={null}
        width={500}
      >
        <div style={{ padding: '20px 0' }}>
          <div style={{ marginBottom: '16px' }}>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: '500' }}>
              {t('主题名称')} *
            </label>
            <Input
              value={formData.topic_name}
              onChange={(value) => setFormData({ ...formData, topic_name: value })}
              placeholder={t('请输入主题名称')}
              maxLength={100}
            />
          </div>
          
          <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
            <Button
              onClick={() => {
                setModalVisible(false);
                setFormData({ topic_name: '', topic_description: '' });
              }}
            >
              {t('取消')}
            </Button>
            <Button
              type="primary"
              onClick={createSubscription}
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

export default Subscription;
