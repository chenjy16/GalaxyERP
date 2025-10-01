'use client';

import { Card, Row, Col, Statistic, Progress, List, Avatar, Typography, Space, Divider } from 'antd';
import { withAuth } from '@/contexts/AuthContext';
import { 
  DollarOutlined, 
  ShoppingCartOutlined, 
  UserOutlined, 
  BoxPlotOutlined,
  RiseOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;

function Page() {
  // 模拟数据
  const recentActivities = [
    {
      id: 1,
      title: '新订单 #SO-2025-001',
      description: '来自北京科技有限公司的采购订单',
      time: '2分钟前',
      type: 'order'
    },
    {
      id: 2,
      title: '库存预警',
      description: '产品 iPhone 15 库存不足',
      time: '15分钟前',
      type: 'warning'
    },
    {
      id: 3,
      title: '付款确认',
      description: '订单 #SO-2025-002 已收到付款',
      time: '1小时前',
      type: 'payment'
    },
    {
      id: 4,
      title: '新客户注册',
      description: '深圳创新科技有限公司',
      time: '2小时前',
      type: 'customer'
    }
  ];

  const getActivityIcon = (type: string) => {
    switch (type) {
      case 'order': return <ShoppingCartOutlined style={{ color: '#1890ff' }} />;
      case 'warning': return <ExclamationCircleOutlined style={{ color: '#faad14' }} />;
      case 'payment': return <CheckCircleOutlined style={{ color: '#52c41a' }} />;
      case 'customer': return <UserOutlined style={{ color: '#722ed1' }} />;
      default: return <ClockCircleOutlined />;
    }
  };

  return (
    <div style={{ padding: '0 8px' }}>
      <div style={{ marginBottom: 24 }}>
        <Title level={2} style={{ margin: 0, color: '#1f2937' }}>
          📊 仪表盘
        </Title>
        <Text type="secondary">欢迎回来！这里是您的业务概览</Text>
      </div>
      
      {/* 核心指标卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card 
            style={{ 
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              border: 'none',
              borderRadius: 12
            }}
          >
            <Statistic
              title={<span style={{ color: 'rgba(255,255,255,0.8)' }}>总销售额</span>}
              value={1128.5}
              precision={2}
              valueStyle={{ color: '#fff', fontSize: 28, fontWeight: 'bold' }}
              prefix={<DollarOutlined style={{ color: '#fff' }} />}
              suffix={<span style={{ color: 'rgba(255,255,255,0.8)' }}>万</span>}
            />
            <div style={{ marginTop: 8 }}>
              <Text style={{ color: 'rgba(255,255,255,0.8)', fontSize: 12 }}>
                <RiseOutlined /> 较上月增长 12.5%
              </Text>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card 
            style={{ 
              background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
              border: 'none',
              borderRadius: 12
            }}
          >
            <Statistic
              title={<span style={{ color: 'rgba(255,255,255,0.8)' }}>订单数量</span>}
              value={93}
              valueStyle={{ color: '#fff', fontSize: 28, fontWeight: 'bold' }}
              prefix={<ShoppingCartOutlined style={{ color: '#fff' }} />}
            />
            <div style={{ marginTop: 8 }}>
              <Text style={{ color: 'rgba(255,255,255,0.8)', fontSize: 12 }}>
                <RiseOutlined /> 较上月增长 8.2%
              </Text>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card 
            style={{ 
              background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
              border: 'none',
              borderRadius: 12
            }}
          >
            <Statistic
              title={<span style={{ color: 'rgba(255,255,255,0.8)' }}>客户数量</span>}
              value={1128}
              valueStyle={{ color: '#fff', fontSize: 28, fontWeight: 'bold' }}
              prefix={<UserOutlined style={{ color: '#fff' }} />}
            />
            <div style={{ marginTop: 8 }}>
              <Text style={{ color: 'rgba(255,255,255,0.8)', fontSize: 12 }}>
                <RiseOutlined /> 较上月增长 15.3%
              </Text>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card 
            style={{ 
              background: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
              border: 'none',
              borderRadius: 12
            }}
          >
            <Statistic
              title={<span style={{ color: 'rgba(255,255,255,0.8)' }}>库存价值</span>}
              value={856.2}
              precision={1}
              valueStyle={{ color: '#fff', fontSize: 28, fontWeight: 'bold' }}
              prefix={<BoxPlotOutlined style={{ color: '#fff' }} />}
              suffix={<span style={{ color: 'rgba(255,255,255,0.8)' }}>万</span>}
            />
            <div style={{ marginTop: 8 }}>
              <Text style={{ color: 'rgba(255,255,255,0.8)', fontSize: 12 }}>
                <RiseOutlined /> 较上月增长 5.7%
              </Text>
            </div>
          </Card>
        </Col>
      </Row>

      {/* 详细信息区域 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card 
            title={<Title level={4} style={{ margin: 0 }}>📈 业务进度</Title>}
            style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}
          >
            <Row gutter={[16, 24]}>
              <Col span={12}>
                <div>
                  <Text strong>月度销售目标</Text>
                  <Progress 
                    percent={75} 
                    strokeColor="#667eea"
                    style={{ marginTop: 8 }}
                  />
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    已完成 75% (目标: 1500万)
                  </Text>
                </div>
              </Col>
              <Col span={12}>
                <div>
                  <Text strong>客户满意度</Text>
                  <Progress 
                    percent={92} 
                    strokeColor="#52c41a"
                    style={{ marginTop: 8 }}
                  />
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    客户评分: 4.6/5.0
                  </Text>
                </div>
              </Col>
              <Col span={12}>
                <div>
                  <Text strong>库存周转率</Text>
                  <Progress 
                    percent={68} 
                    strokeColor="#faad14"
                    style={{ marginTop: 8 }}
                  />
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    平均周转: 45天
                  </Text>
                </div>
              </Col>
              <Col span={12}>
                <div>
                  <Text strong>订单及时率</Text>
                  <Progress 
                    percent={88} 
                    strokeColor="#1890ff"
                    style={{ marginTop: 8 }}
                  />
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    按时交付: 88%
                  </Text>
                </div>
              </Col>
            </Row>
          </Card>
        </Col>
        
        <Col xs={24} lg={8}>
          <Card 
            title={<Title level={4} style={{ margin: 0 }}>🔔 最近活动</Title>}
            style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}
          >
            <List
              itemLayout="horizontal"
              dataSource={recentActivities}
              renderItem={(item) => (
                <List.Item style={{ padding: '12px 0', border: 'none' }}>
                  <List.Item.Meta
                    avatar={
                      <Avatar 
                        icon={getActivityIcon(item.type)}
                        style={{ backgroundColor: '#f5f5f5' }}
                      />
                    }
                    title={
                      <Text strong style={{ fontSize: 14 }}>
                        {item.title}
                      </Text>
                    }
                    description={
                      <div>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          {item.description}
                        </Text>
                        <br />
                        <Text type="secondary" style={{ fontSize: 11 }}>
                          <ClockCircleOutlined /> {item.time}
                        </Text>
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>

      {/* 快速操作区域 */}
      <Card 
        title={<Title level={4} style={{ margin: 0 }}>⚡ 快速操作</Title>}
        style={{ marginTop: 16, borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}
      >
        <Row gutter={[16, 16]}>
          <Col xs={12} sm={8} md={6}>
            <Card 
              hoverable
              style={{ textAlign: 'center', borderRadius: 8 }}
              bodyStyle={{ padding: 16 }}
            >
              <ShoppingCartOutlined style={{ fontSize: 24, color: '#1890ff', marginBottom: 8 }} />
              <div>
                <Text strong>创建订单</Text>
              </div>
            </Card>
          </Col>
          <Col xs={12} sm={8} md={6}>
            <Card 
              hoverable
              style={{ textAlign: 'center', borderRadius: 8 }}
              bodyStyle={{ padding: 16 }}
            >
              <UserOutlined style={{ fontSize: 24, color: '#722ed1', marginBottom: 8 }} />
              <div>
                <Text strong>添加客户</Text>
              </div>
            </Card>
          </Col>
          <Col xs={12} sm={8} md={6}>
            <Card 
              hoverable
              style={{ textAlign: 'center', borderRadius: 8 }}
              bodyStyle={{ padding: 16 }}
            >
              <BoxPlotOutlined style={{ fontSize: 24, color: '#52c41a', marginBottom: 8 }} />
              <div>
                <Text strong>库存管理</Text>
              </div>
            </Card>
          </Col>
          <Col xs={12} sm={8} md={6}>
            <Card 
              hoverable
              style={{ textAlign: 'center', borderRadius: 8 }}
              bodyStyle={{ padding: 16 }}
            >
              <DollarOutlined style={{ fontSize: 24, color: '#faad14', marginBottom: 8 }} />
              <div>
                <Text strong>财务报表</Text>
              </div>
            </Card>
          </Col>
        </Row>
      </Card>
    </div>
  );
}

export default withAuth(Page);