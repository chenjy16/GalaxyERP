'use client';

import { useState, useEffect } from 'react';
import { 
  Card, 
  Tabs, 
  Table, 
  Button, 
  Input, 
  Space, 
  Tag,
  Statistic,
  Row,
  Col,
  Progress,
  Modal,
  Form,
  Select,
  DatePicker,
  Typography,
  Avatar,
  Tooltip,
  Dropdown,
  message,
  Badge
} from 'antd';
import { 
  PlusOutlined, 
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  CheckOutlined,
  CloseOutlined,
  EyeOutlined,
  MoreOutlined,
  ShopOutlined,
  PhoneOutlined,
  MailOutlined,
  EnvironmentOutlined,
  FileTextOutlined,
  ShoppingCartOutlined,
  TruckOutlined,
  ExportOutlined,
  UserOutlined,
  ImportOutlined,
  ClockCircleOutlined,
  TeamOutlined
} from '@ant-design/icons';
import { Supplier, PurchaseOrder, PurchaseRequest } from '@/types/api';
import SupplierService from '@/services/supplier';
import PurchaseOrderService from '@/services/purchaseOrder';
import PurchaseRequestService from '@/services/purchaseRequest';
import { withAuth } from '@/contexts/AuthContext';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Option } = Select;

function PurchasePage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'supplier' | 'order' | 'request'>('supplier');
  const [form] = Form.useForm();
  
  // 数据状态
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [purchaseOrders, setPurchaseOrders] = useState<PurchaseOrder[]>([]);
  const [purchaseRequests, setPurchaseRequests] = useState<PurchaseRequest[]>([]);
  
  // 加载状态
  const [suppliersLoading, setSuppliersLoading] = useState(false);
  const [ordersLoading, setOrdersLoading] = useState(false);
  const [requestsLoading, setRequestsLoading] = useState(false);
  
  // 分页状态
  const [suppliersPagination, setSuppliersPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [ordersPagination, setOrdersPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [requestsPagination, setRequestsPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  
  // 搜索状态
  const [suppliersSearch, setSuppliersSearch] = useState('');
  const [ordersSearch, setOrdersSearch] = useState('');
  const [requestsSearch, setRequestsSearch] = useState('');

  // 加载供应商数据
  const loadSuppliers = async (page = 1, search = '') => {
    setSuppliersLoading(true);
    try {
      const response = await SupplierService.getSuppliers({
        page,
        limit: suppliersPagination.pageSize,
        search: search || undefined
      });
      setSuppliers(response.data);
      setSuppliersPagination(prev => ({
        ...prev,
        current: page,
        total: response.total
      }));
    } catch (error) {
      message.error('加载供应商数据失败');
      console.error('Error loading suppliers:', error);
    } finally {
      setSuppliersLoading(false);
    }
  };

  // 加载采购订单数据
  const loadPurchaseOrders = async (page = 1, search = '') => {
    setOrdersLoading(true);
    try {
      const response = await PurchaseOrderService.getPurchaseOrders({
        page,
        limit: ordersPagination.pageSize,
        search: search || undefined
      });
      setPurchaseOrders(response.data);
      setOrdersPagination(prev => ({
        ...prev,
        current: page,
        total: response.total
      }));
    } catch (error) {
      message.error('加载采购订单数据失败');
      console.error('Error loading purchase orders:', error);
    } finally {
      setOrdersLoading(false);
    }
  };

  // 加载采购请求数据
  const loadPurchaseRequests = async (page = 1, search = '') => {
    setRequestsLoading(true);
    try {
      const response = await PurchaseRequestService.getPurchaseRequests({
        page,
        limit: requestsPagination.pageSize,
        search: search || undefined
      });
      setPurchaseRequests(response.data);
      setRequestsPagination(prev => ({
        ...prev,
        current: page,
        total: response.total
      }));
    } catch (error) {
      message.error('加载采购请求数据失败');
      console.error('Error loading purchase requests:', error);
    } finally {
      setRequestsLoading(false);
    }
  };

  // 初始化数据加载
  useEffect(() => {
    loadSuppliers();
    loadPurchaseOrders();
    loadPurchaseRequests();
  }, []);





  const supplierColumns = [
    {
      title: '供应商信息',
      key: 'supplier',
      render: (record: Supplier) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Avatar 
            size={40}
            style={{ backgroundColor: '#52c41a' }}
            icon={<ShopOutlined />}
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
      render: (record: Supplier) => (
        <div>
          <div style={{ marginBottom: 4 }}>
            <UserOutlined style={{ marginRight: 4, color: '#666' }} />
            <Text>{record.contactName}</Text>
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
      key: 'credit',
      render: (record: Supplier) => (
        <div>
          <div style={{ marginBottom: 4 }}>
            <Text type="secondary">信用额度: </Text>
            <Text strong style={{ color: '#52c41a' }}>¥{record.creditLimit.toLocaleString()}</Text>
          </div>
          <div>
            <Text type="secondary">付款条款: </Text>
            <Text>{record.paymentTerms}</Text>
          </div>
        </div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>
          {status === 'active' ? '活跃' : '非活跃'}
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

  const orderColumns = [
    {
      title: '订单编号',
      dataIndex: 'orderNumber',
      key: 'orderNumber',
      render: (orderNumber: string) => <Text strong>{orderNumber}</Text>,
    },
    {
      title: '供应商',
      dataIndex: 'supplierName',
      key: 'supplierName',
    },
    {
      title: '金额',
      dataIndex: 'totalAmount',
      key: 'totalAmount',
      render: (amount: number) => (
        <Text strong style={{ color: '#52c41a' }}>¥{amount.toLocaleString()}</Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusConfig = {
          pending: { color: 'orange', text: '待确认' },
          confirmed: { color: 'blue', text: '已确认' },
          shipped: { color: 'purple', text: '已发货' },
          received: { color: 'green', text: '已收货' },
          cancelled: { color: 'red', text: '已取消' }
        };
        const config = statusConfig[status as keyof typeof statusConfig] || { color: 'default', text: status };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: '订单日期',
      dataIndex: 'orderDate',
      key: 'orderDate',
      render: (date: string) => new Date(date).toLocaleDateString('zh-CN'),
    },
    {
      title: '预期到货',
      dataIndex: 'expectedDeliveryDate',
      key: 'expectedDeliveryDate',
      render: (date: string) => date ? new Date(date).toLocaleDateString('zh-CN') : '-',
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

  const requestColumns = [
    {
      title: '申请编号',
      dataIndex: 'requestNumber',
      key: 'requestNumber',
      render: (requestNumber: string) => <Text strong>{requestNumber}</Text>,
    },
    {
      title: '申请部门',
      dataIndex: 'department',
      key: 'department',
    },
    {
      title: '申请人',
      dataIndex: 'requestedBy',
      key: 'requestedBy',
    },
    {
      title: '预估金额',
      dataIndex: 'totalAmount',
      key: 'totalAmount',
      render: (amount: number) => (
        <Text strong style={{ color: '#52c41a' }}>¥{amount.toLocaleString()}</Text>
      ),
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      key: 'priority',
      render: (priority: string) => {
        const priorityConfig = {
          high: { color: 'red', text: '高' },
          medium: { color: 'orange', text: '中' },
          low: { color: 'green', text: '低' }
        };
        const config = priorityConfig[priority as keyof typeof priorityConfig] || { color: 'default', text: priority };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusConfig = {
          draft: { color: 'default', text: '草稿' },
          submitted: { color: 'blue', text: '已提交' },
          approved: { color: 'green', text: '已批准' },
          rejected: { color: 'red', text: '已拒绝' },
          cancelled: { color: 'red', text: '已取消' }
        };
        const config = statusConfig[status as keyof typeof statusConfig] || { color: 'default', text: status };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: '申请日期',
      dataIndex: 'requestDate',
      key: 'requestDate',
      render: (date: string) => new Date(date).toLocaleDateString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Space>
          <Tooltip title="查看">
            <Button type="text" icon={<EyeOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="审批">
            <Button type="text" icon={<CheckOutlined />} size="small" style={{ color: '#52c41a' }} />
          </Tooltip>
          <Tooltip title="拒绝">
            <Button type="text" icon={<CloseOutlined />} size="small" danger />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const handleModalOk = () => {
    form.validateFields().then(values => {
      console.log('Form values:', values);
      message.success(`${modalType === 'supplier' ? '供应商' : modalType === 'order' ? '采购订单' : '采购申请'}创建成功！`);
      setIsModalVisible(false);
      form.resetFields();
    });
  };

  const showModal = (type: 'supplier' | 'order' | 'request') => {
    setModalType(type);
    setIsModalVisible(true);
  };

  // 计算统计数据
  const totalOrders = ordersPagination.total;
  const pendingOrders = (purchaseOrders || []).filter((o: PurchaseOrder) => o.status === 'pending').length;
  const totalSuppliers = suppliersPagination.total;
  const totalRequests = requestsPagination.total;

  // 定义Tabs的items
  const tabItems = [
    {
      key: 'orders',
      label: '采购订单',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索订单..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
                value={ordersSearch}
                onChange={(e) => setOrdersSearch(e.target.value)}
                onPressEnter={() => loadPurchaseOrders(1, ordersSearch)}
                allowClear
              />
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="待审核">待审核</Option>
                <Option value="已审核">已审核</Option>
                <Option value="已发货">已发货</Option>
                <Option value="已完成">已完成</Option>
                <Option value="已取消">已取消</Option>
              </Select>
              <Select placeholder="供应商筛选" style={{ width: 150 }}>
                <Option value="all">全部</Option>
                <Option value="华为技术有限公司">华为技术有限公司</Option>
                <Option value="小米科技有限公司">小米科技有限公司</Option>
                <Option value="联想集团有限公司">联想集团有限公司</Option>
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
             dataSource={purchaseOrders}
             loading={ordersLoading}
             pagination={{
               current: ordersPagination.current,
               pageSize: ordersPagination.pageSize,
               total: ordersPagination.total,
               showSizeChanger: true,
               showQuickJumper: true,
               showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
               onChange: (page, pageSize) => {
                 setOrdersPagination(prev => ({ ...prev, pageSize: pageSize || 10 }));
                 loadPurchaseOrders(page, ordersSearch);
               }
             }}
             scroll={{ x: 1400 }}
           />
        </>
      )
    },
    {
      key: 'suppliers',
      label: '供应商管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索供应商..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
                value={suppliersSearch}
                onChange={(e) => setSuppliersSearch(e.target.value)}
                onPressEnter={() => loadSuppliers(1, suppliersSearch)}
                allowClear
              />
              <Select placeholder="类型筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="电子产品">电子产品</Option>
                <Option value="办公用品">办公用品</Option>
                <Option value="原材料">原材料</Option>
                <Option value="设备">设备</Option>
              </Select>
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="合作中">合作中</Option>
                <Option value="暂停合作">暂停合作</Option>
                <Option value="黑名单">黑名单</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ImportOutlined />}>导入</Button>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('supplier')}
              >
                新建供应商
              </Button>
            </Space>
          </div>
          <Table 
            columns={supplierColumns} 
            dataSource={suppliers}
            loading={suppliersLoading}
            pagination={{
              current: suppliersPagination.current,
              pageSize: suppliersPagination.pageSize,
              total: suppliersPagination.total,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
              onChange: (page, pageSize) => {
                setSuppliersPagination(prev => ({ ...prev, pageSize: pageSize || 10 }));
                loadSuppliers(page, suppliersSearch);
              }
            }}
            scroll={{ x: 1200 }}
          />
        </>
      )
    },
    {
       key: 'requests',
       label: '采购申请',
       children: (
         <>
           <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
             <Space>
               <Input
                 placeholder="搜索申请..."
                 prefix={<SearchOutlined />}
                 style={{ width: 300 }}
                 value={requestsSearch}
                 onChange={(e) => setRequestsSearch(e.target.value)}
                 onPressEnter={() => loadPurchaseRequests(1, requestsSearch)}
                 allowClear
               />
               <Select placeholder="状态筛选" style={{ width: 120 }}>
                 <Option value="all">全部</Option>
                 <Option value="pending">待审批</Option>
                 <Option value="approved">已批准</Option>
                 <Option value="rejected">已拒绝</Option>
               </Select>
               <Select placeholder="部门筛选" style={{ width: 150 }}>
                 <Option value="all">全部</Option>
                 <Option value="生产部">生产部</Option>
                 <Option value="研发部">研发部</Option>
                 <Option value="销售部">销售部</Option>
               </Select>
               <Select placeholder="紧急程度" style={{ width: 120 }}>
                 <Option value="all">全部</Option>
                 <Option value="high">紧急</Option>
                 <Option value="medium">一般</Option>
                 <Option value="low">不急</Option>
               </Select>
             </Space>
             <Space>
               <Button icon={<ExportOutlined />}>导出</Button>
               <Button 
                 type="primary" 
                 icon={<PlusOutlined />}
                 onClick={() => showModal('request')}
               >
                 新建申请
               </Button>
             </Space>
           </div>
           <Table 
             columns={requestColumns} 
             dataSource={purchaseRequests}
             loading={requestsLoading}
             pagination={{
               current: requestsPagination.current,
               pageSize: requestsPagination.pageSize,
               total: requestsPagination.total,
               showSizeChanger: true,
               showQuickJumper: true,
               showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
               onChange: (page, pageSize) => {
                 setRequestsPagination(prev => ({ ...prev, pageSize: pageSize || 10 }));
                 loadPurchaseRequests(page, requestsSearch);
               }
             }}
             scroll={{ x: 1200 }}
           />
         </>
       )
     }
  ];

  return (
    <div style={{ padding: '0 8px' }}>
      <div style={{ marginBottom: 24 }}>
        <Title level={2} style={{ margin: 0, color: '#1f2937' }}>
          🛒 采购管理
        </Title>
        <Text type="secondary">管理采购订单、供应商和合同</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="采购订单"
              value={totalOrders}
              prefix={<ShoppingCartOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="待审核"
              value={pendingOrders}
              prefix={<ClockCircleOutlined style={{ color: '#faad14' }} />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="供应商"
              value={totalSuppliers}
              prefix={<TeamOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="采购申请"
              value={totalRequests}
              prefix={<FileTextOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <Tabs defaultActiveKey="orders" items={tabItems} />
      </Card>
    </div>
  );
}

export default withAuth(PurchasePage);