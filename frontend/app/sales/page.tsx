'use client';

import { useState, useEffect } from 'react';
import { withAuth } from '@/contexts/AuthContext';
import { CustomerService } from '@/services/customer';
import { QuotationService } from '@/services/quotation';
import { SalesOrderService } from '@/services/salesOrder';
import { Customer, Quotation, SalesOrder } from '@/types/api';
import { 
  Card, 
  Tabs, 
  Table, 
  Button, 
  Input, 
  Tag, 
  Space, 
  Modal, 
  Form, 
  Select, 
  DatePicker, 
  Row, 
  Col, 
  Statistic, 
  Typography,
  Avatar,
  Tooltip,
  Dropdown,
  message
} from 'antd';
import { 
  PlusOutlined, 
  SearchOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  EyeOutlined,
  MoreOutlined,
  UserOutlined,
  PhoneOutlined,
  MailOutlined,
  EnvironmentOutlined,
  DollarOutlined,
  ShoppingCartOutlined,
  FileTextOutlined,
  ExportOutlined,
  ImportOutlined,
  BarChartOutlined,
  ClockCircleOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;
const { Option } = Select;

function SalesPage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'customer' | 'quote' | 'order'>('customer');
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  
  // 数据状态
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [quotations, setQuotations] = useState<Quotation[]>([]);
  const [salesOrders, setSalesOrders] = useState<SalesOrder[]>([]);
  
  // 分页状态
  const [customerPagination, setCustomerPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [quotationPagination, setQuotationPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [orderPagination, setOrderPagination] = useState({ current: 1, pageSize: 10, total: 0 });

  // 加载客户数据
  const loadCustomers = async (page = 1, limit = 10) => {
    try {
      setLoading(true);
      const response = await CustomerService.getCustomers({ page, limit });
      setCustomers(response.data);
      setCustomerPagination({
        current: response.page,
        pageSize: response.limit,
        total: response.total
      });
    } catch (error) {
      message.error('加载客户数据失败');
      console.error('Error loading customers:', error);
    } finally {
      setLoading(false);
    }
  };

  // 加载报价数据
  const loadQuotations = async (page = 1, limit = 10) => {
    try {
      setLoading(true);
      const response = await QuotationService.getQuotations({ page, limit });
      setQuotations(response.data);
      setQuotationPagination({
        current: response.page,
        pageSize: response.limit,
        total: response.total
      });
    } catch (error) {
      message.error('加载报价数据失败');
      console.error('Error loading quotations:', error);
    } finally {
      setLoading(false);
    }
  };

  // 加载销售订单数据
  const loadSalesOrders = async (page = 1, limit = 10) => {
    try {
      setLoading(true);
      const response = await SalesOrderService.getSalesOrders({ page, limit });
      setSalesOrders(response.data);
      setOrderPagination({
        current: response.page,
        pageSize: response.limit,
        total: response.total
      });
    } catch (error) {
      message.error('加载销售订单数据失败');
      console.error('Error loading sales orders:', error);
    } finally {
      setLoading(false);
    }
  };

  // 初始化数据
  useEffect(() => {
    loadCustomers();
    loadQuotations();
    loadSalesOrders();
  }, []);



  const customerColumns = [
    {
      title: '客户信息',
      key: 'customer',
      render: (record: Customer) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Avatar 
            size={40}
            style={{ backgroundColor: '#1890ff' }}
            icon={<UserOutlined />}
          />
          <div>
            <Text strong style={{ display: 'block' }}>{record.name}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>编码: {record.code}</Text>
          </div>
        </div>
      ),
    },
    {
      title: '联系方式',
      key: 'contact',
      render: (record: Customer) => (
        <div>
          <div style={{ marginBottom: 4 }}>
            <UserOutlined style={{ marginRight: 4, color: '#666' }} />
            <Text>{record.contactPerson}</Text>
          </div>
          <div style={{ marginBottom: 4 }}>
            <PhoneOutlined style={{ marginRight: 4, color: '#666' }} />
            <Text>{record.phone}</Text>
          </div>
          <div>
            <MailOutlined style={{ marginRight: 4, color: '#666' }} />
            <Text>{record.email}</Text>
          </div>
        </div>
      ),
    },
    {
      title: '地址',
      dataIndex: 'address',
      key: 'address',
      render: (address: string) => (
        <div>
          <EnvironmentOutlined style={{ marginRight: 4, color: '#666' }} />
          <Text>{address}</Text>
        </div>
      ),
    },
    {
      title: '信用额度',
      key: 'creditLimit',
      render: (record: Customer) => (
        <div>
          <Text type="secondary">信用额度: </Text>
          <Text strong style={{ color: '#52c41a' }}>
            ¥{record.creditLimit ? record.creditLimit.toLocaleString() : '0.00'}
          </Text>
        </div>
      ),
    },
    {
      title: '状态',
      key: 'status',
      render: (record: Customer) => (
        <Tag color={record.status === 'active' ? 'green' : 'red'}>
          {record.status === 'active' ? '活跃' : '非活跃'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Dropdown
          menu={{
            items: [
              {
                key: 'view',
                label: '查看详情',
                icon: <EyeOutlined />,
              },
              {
                key: 'edit',
                label: '编辑',
                icon: <EditOutlined />,
              },
              {
                key: 'delete',
                label: '删除',
                icon: <DeleteOutlined />,
                danger: true,
              },
            ],
          }}
          trigger={['click']}
        >
          <Button type="text" icon={<MoreOutlined />} />
        </Dropdown>
      ),
    },
  ];

  const quoteColumns = [
    {
      title: '报价编号',
      dataIndex: 'quotationNumber',
      key: 'quotationNumber',
      render: (quotationNumber: string) => <Text strong>{quotationNumber}</Text>,
    },
    {
      title: '客户',
      dataIndex: 'customerName',
      key: 'customerName',
      render: (text: any, record: any) => {
        // Fix: Properly handle customer object rendering
        if (typeof text === 'object' && text !== null) {
          return text.name || text.id || '-';
        }
        // Also check if record has a customer object
        if (record.customer && typeof record.customer === 'object' && record.customer !== null) {
          return record.customer.name || record.customer.id || '-';
        }
        return text || '-';
      }
    },
    {
      title: '总金额',
      dataIndex: 'totalAmount',
      key: 'totalAmount',
      render: (amount: number) => (
        <Text strong style={{ color: '#52c41a' }}>
          ¥{amount ? amount.toLocaleString() : '0.00'}
        </Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusMap: { [key: string]: { color: string; text: string } } = {
          'draft': { color: 'default', text: '草稿' },
          'submitted': { color: 'blue', text: '已提交' },
          'accepted': { color: 'green', text: '已接受' },
          'rejected': { color: 'red', text: '已拒绝' },
          'expired': { color: 'orange', text: '已过期' }
        };
        const statusInfo = statusMap[status] || { color: 'default', text: status };
        return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>;
      },
    },
    {
      title: '有效期至',
      dataIndex: 'validUntil',
      key: 'validUntil',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '创建日期',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Space>
          <Tooltip title="查看">
            <Button type="text" icon={<EyeOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="编辑">
            <Button type="text" icon={<EditOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="删除">
            <Button type="text" icon={<DeleteOutlined />} size="small" danger />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const orderColumns = [
    {
      title: '订单编号',
      dataIndex: 'orderNumber',
      key: 'orderNumber',
      render: (orderNumber: string) => <Text strong>{orderNumber}</Text>,
    },
    {
      title: '客户',
      dataIndex: 'customerName',
      key: 'customerName',
      render: (text: any, record: any) => {
        // Fix: Properly handle customer object rendering
        if (typeof text === 'object' && text !== null) {
          return text.name || text.id || '-';
        }
        // Also check if record has a customer object
        if (record.customer && typeof record.customer === 'object' && record.customer !== null) {
          return record.customer.name || record.customer.id || '-';
        }
        return text || '-';
      }
    },
    {
      title: '总金额',
      dataIndex: 'totalAmount',
      key: 'totalAmount',
      render: (amount: number) => (
        <Text strong style={{ color: '#52c41a' }}>
          ¥{amount ? amount.toLocaleString() : '0.00'}
        </Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusMap: { [key: string]: { color: string; text: string } } = {
          'pending': { color: 'orange', text: '待处理' },
          'confirmed': { color: 'blue', text: '已确认' },
          'shipped': { color: 'cyan', text: '已发货' },
          'delivered': { color: 'green', text: '已交付' },
          'cancelled': { color: 'red', text: '已取消' }
        };
        const statusInfo = statusMap[status] || { color: 'default', text: status };
        return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>;
      },
    },
    {
      title: '订单日期',
      dataIndex: 'orderDate',
      key: 'orderDate',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '交付日期',
      dataIndex: 'deliveryDate',
      key: 'deliveryDate',
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Space>
          <Tooltip title="查看">
            <Button type="text" icon={<EyeOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="编辑">
            <Button type="text" icon={<EditOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="删除">
            <Button type="text" icon={<DeleteOutlined />} size="small" danger />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const handleModalOk = () => {
    form.validateFields().then(values => {
      console.log('Form values:', values);
      message.success(`${modalType === 'customer' ? '客户' : modalType === 'quote' ? '报价' : '订单'}创建成功！`);
      setIsModalVisible(false);
      form.resetFields();
    });
  };

  const showModal = (type: 'customer' | 'quote' | 'order') => {
    setModalType(type);
    setIsModalVisible(true);
  };

  // 计算统计数据
  const totalOrders = salesOrders.length;
  const totalCustomers = customers.length;
  const totalRevenue = salesOrders.reduce((sum: number, order: SalesOrder) => sum + order.grandTotal, 0);
  const pendingOrders = salesOrders.filter((order: SalesOrder) => order.status === 'pending').length;

  // 定义Tabs的items
  const tabItems = [
    {
      key: 'orders',
      label: '销售订单',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索订单..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="pending">待处理</Option>
                <Option value="confirmed">已确认</Option>
                <Option value="shipped">已发货</Option>
                <Option value="delivered">已交付</Option>
                <Option value="cancelled">已取消</Option>
              </Select>
              <Select placeholder="客户筛选" style={{ width: 150 }}>
                <Option value="all">全部客户</Option>
                <Option value="阿里巴巴">阿里巴巴</Option>
                <Option value="腾讯科技">腾讯科技</Option>
                <Option value="字节跳动">字节跳动</Option>
              </Select>
              <DatePicker placeholder="订单日期" style={{ width: 150 }} />
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('order')}
              >
                新建订单
              </Button>
            </Space>
          </div>
          <Table 
            columns={orderColumns} 
            dataSource={salesOrders.map(order => ({ ...order, key: order.id }))}
            pagination={{
              current: orderPagination.current,
              pageSize: orderPagination.pageSize,
              total: orderPagination.total,
              showSizeChanger: true,
              onChange: (page, pageSize) => loadSalesOrders(page, pageSize)
            }}
            loading={loading}
            scroll={{ x: 1500 }}
          />
        </>
      )
    },
    {
      key: 'customers',
      label: '客户管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索客户..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="客户类型" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="企业客户">企业客户</Option>
                <Option value="个人客户">个人客户</Option>
                <Option value="代理商">代理商</Option>
              </Select>
              <Select placeholder="客户等级" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="VIP">VIP</Option>
                <Option value="普通">普通</Option>
                <Option value="潜在">潜在</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ImportOutlined />}>导入</Button>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('customer')}
              >
                新建客户
              </Button>
            </Space>
          </div>
          <Table 
            columns={customerColumns} 
            dataSource={customers.map(customer => ({ ...customer, key: customer.id }))}
            pagination={{
              current: customerPagination.current,
              pageSize: customerPagination.pageSize,
              total: customerPagination.total,
              showSizeChanger: true,
              onChange: (page, pageSize) => loadCustomers(page, pageSize)
            }}
            loading={loading}
            scroll={{ x: 1300 }}
          />
        </>
      )
    },
    {
      key: 'quotations',
      label: '报价管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索报价..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="draft">草稿</Option>
                <Option value="submitted">已提交</Option>
                <Option value="accepted">已接受</Option>
                <Option value="rejected">已拒绝</Option>
                <Option value="expired">已过期</Option>
              </Select>
              <Select placeholder="客户筛选" style={{ width: 150 }}>
                <Option value="all">全部客户</Option>
              </Select>
              <DatePicker placeholder="创建日期" style={{ width: 150 }} />
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('quote')}
              >
                新建报价
              </Button>
            </Space>
          </div>
          <Table 
            columns={quoteColumns} 
            dataSource={quotations.map(quotation => ({ ...quotation, key: quotation.id }))}
            pagination={{
              current: quotationPagination.current,
              pageSize: quotationPagination.pageSize,
              total: quotationPagination.total,
              showSizeChanger: true,
              onChange: (page, pageSize) => loadQuotations(page, pageSize)
            }}
            loading={loading}
            scroll={{ x: 1400 }}
          />
        </>
      )
    },
    {
      key: 'reports',
      label: '销售报表',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Select placeholder="报表类型" style={{ width: 150 }}>
                <Option value="daily">日报表</Option>
                <Option value="weekly">周报表</Option>
                <Option value="monthly">月报表</Option>
                <Option value="yearly">年报表</Option>
              </Select>
              <DatePicker placeholder="开始日期" style={{ width: 150 }} />
              <DatePicker placeholder="结束日期" style={{ width: 150 }} />
              <Button type="primary" icon={<BarChartOutlined />}>
                生成报表
              </Button>
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出报表</Button>
            </Space>
          </div>
          <div style={{ textAlign: 'center', padding: '60px 0', color: '#999' }}>
            <BarChartOutlined style={{ fontSize: 48, marginBottom: 16 }} />
            <div>选择报表类型和时间范围生成销售报表</div>
          </div>
        </>
      )
    }
  ];

  return (
    <div style={{ padding: '0 8px' }}>
      <div style={{ marginBottom: 24 }}>
        <Title level={2} style={{ margin: 0, color: '#1f2937' }}>
          💰 销售管理
        </Title>
        <Text type="secondary">管理销售订单、客户信息和销售报表</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="订单总数"
              value={totalOrders}
              prefix={<ShoppingCartOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="客户总数"
              value={totalCustomers}
              prefix={<UserOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="总销售额"
              value={totalRevenue}
              prefix="¥"
              precision={0}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="待处理订单"
              value={pendingOrders}
              prefix={<ClockCircleOutlined style={{ color: '#faad14' }} />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <Tabs defaultActiveKey="orders" items={tabItems} />
      </Card>

      {/* 模态框 */}
      <Modal
        title={modalType === 'customer' ? '新建客户' : modalType === 'quote' ? '新建报价' : '新建订单'}
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={() => {
          setIsModalVisible(false);
          form.resetFields();
        }}
        width={800}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{}}
        >
          {modalType === 'customer' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="name"
                    label="客户名称"
                    rules={[{ required: true, message: '请输入客户名称' }]}
                  >
                    <Input placeholder="请输入客户名称" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="code"
                    label="客户编码"
                    rules={[{ required: true, message: '请输入客户编码' }]}
                  >
                    <Input placeholder="请输入客户编码" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="contactPerson"
                    label="联系人"
                    rules={[{ required: true, message: '请输入联系人' }]}
                  >
                    <Input placeholder="请输入联系人" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="phone"
                    label="联系电话"
                    rules={[{ required: true, message: '请输入联系电话' }]}
                  >
                    <Input placeholder="请输入联系电话" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="email"
                    label="邮箱"
                    rules={[{ type: 'email', message: '请输入有效的邮箱地址' }]}
                  >
                    <Input placeholder="请输入邮箱" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="creditLimit"
                    label="信用额度"
                    rules={[{ required: true, message: '请输入信用额度' }]}
                  >
                    <Input type="number" placeholder="请输入信用额度" />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item
                name="address"
                label="地址"
              >
                <Input.TextArea placeholder="请输入地址" rows={3} />
              </Form.Item>
            </>
          )}
          
          {modalType === 'quote' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="quotationNumber"
                    label="报价编号"
                    rules={[{ required: true, message: '请输入报价编号' }]}
                  >
                    <Input placeholder="请输入报价编号" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="customerId"
                    label="客户"
                    rules={[{ required: true, message: '请选择客户' }]}
                  >
                    <Select placeholder="请选择客户">
                      {customers.map(customer => (
                        <Option key={customer.id} value={customer.id}>
                          {customer.name}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="totalAmount"
                    label="总金额"
                    rules={[{ required: true, message: '请输入总金额' }]}
                  >
                    <Input type="number" placeholder="请输入总金额" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="validUntil"
                    label="有效期至"
                    rules={[{ required: true, message: '请选择有效期' }]}
                  >
                    <DatePicker style={{ width: '100%' }} placeholder="请选择有效期" />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item
                name="description"
                label="备注"
              >
                <Input.TextArea placeholder="请输入备注" rows={3} />
              </Form.Item>
            </>
          )}
          
          {modalType === 'order' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="orderNumber"
                    label="订单编号"
                    rules={[{ required: true, message: '请输入订单编号' }]}
                  >
                    <Input placeholder="请输入订单编号" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="customerId"
                    label="客户"
                    rules={[{ required: true, message: '请选择客户' }]}
                  >
                    <Select placeholder="请选择客户">
                      {customers.map(customer => (
                        <Option key={customer.id} value={customer.id}>
                          {customer.name}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="totalAmount"
                    label="总金额"
                    rules={[{ required: true, message: '请输入总金额' }]}
                  >
                    <Input type="number" placeholder="请输入总金额" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="deliveryDate"
                    label="交付日期"
                  >
                    <DatePicker style={{ width: '100%' }} placeholder="请选择交付日期" />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item
                name="notes"
                label="备注"
              >
                <Input.TextArea placeholder="请输入备注" rows={3} />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>
    </div>
  );
}

export default withAuth(SalesPage);