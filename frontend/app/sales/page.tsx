'use client';

import { useState, useEffect } from 'react';
import { withAuth } from '@/contexts/AuthContext';
import { CustomerService } from '@/services/customer';
import { QuotationService } from '@/services/quotation';
import { SalesOrderService } from '@/services/salesOrder';
import { QuotationVersionService, QuotationVersion, QuotationVersionHistoryResponse } from '@/services/quotationVersion';
import { Customer, Quotation, SalesOrder } from '@/types/api';
import dayjs from 'dayjs';
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
  Checkbox,
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
  ClockCircleOutlined,
  BranchesOutlined,
  HistoryOutlined,
  SwapOutlined,
  RollbackOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;
const { Option } = Select;

function SalesPage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'customer' | 'quote' | 'order'>('customer');
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [editingRecord, setEditingRecord] = useState<any>(null);
  const [viewModalVisible, setViewModalVisible] = useState(false);
  const [viewingRecord, setViewingRecord] = useState<any>(null);
  
  // 版本管理相关状态
  const [versionModalVisible, setVersionModalVisible] = useState(false);
  const [versionHistoryVisible, setVersionHistoryVisible] = useState(false);
  const [compareModalVisible, setCompareModalVisible] = useState(false);
  const [currentQuotationId, setCurrentQuotationId] = useState<number | null>(null);
  const [versions, setVersions] = useState<QuotationVersionHistoryResponse[]>([]);
  const [selectedVersions, setSelectedVersions] = useState<number[]>([]);
  const [compareResult, setCompareResult] = useState<any>(null);
  const [versionForm] = Form.useForm();
  
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

  // 创建记录
  const handleCreate = async (values: any) => {
    try {
      setLoading(true);
      let response;
      
      switch (modalType) {
        case 'customer':
          response = await CustomerService.createCustomer(values);
          message.success('客户创建成功');
          loadCustomers();
          break;
        case 'quote':
          response = await QuotationService.createQuotation(values);
          message.success('报价单创建成功');
          loadQuotations();
          break;
        case 'order':
          response = await SalesOrderService.createSalesOrder(values);
          message.success('销售订单创建成功');
          loadSalesOrders();
          break;
      }
      
      setIsModalVisible(false);
      form.resetFields();
      setEditingRecord(null);
    } catch (error) {
      message.error(`创建失败: ${error}`);
      console.error('Create error:', error);
    } finally {
      setLoading(false);
    }
  };

  // 更新记录
  const handleUpdate = async (values: any) => {
    try {
      setLoading(true);
      let response;
      
      switch (modalType) {
        case 'customer':
          response = await CustomerService.updateCustomer(editingRecord.id, values);
          message.success('客户更新成功');
          loadCustomers();
          break;
        case 'quote':
          response = await QuotationService.updateQuotation(editingRecord.id, values);
          message.success('报价单更新成功');
          loadQuotations();
          break;
        case 'order':
          response = await SalesOrderService.updateSalesOrder(editingRecord.id, values);
          message.success('销售订单更新成功');
          loadSalesOrders();
          break;
      }
      
      setIsModalVisible(false);
      form.resetFields();
      setEditingRecord(null);
    } catch (error) {
      message.error(`更新失败: ${error}`);
      console.error('Update error:', error);
    } finally {
      setLoading(false);
    }
  };

  // 删除记录
  const handleDelete = async (record: any, type: 'customer' | 'quote' | 'order') => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除这条${type === 'customer' ? '客户' : type === 'quote' ? '报价单' : '销售订单'}记录吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          setLoading(true);
          
          switch (type) {
            case 'customer':
              await CustomerService.deleteCustomer(record.id);
              message.success('客户删除成功');
              loadCustomers();
              break;
            case 'quote':
              await QuotationService.deleteQuotation(record.id);
              message.success('报价单删除成功');
              loadQuotations();
              break;
            case 'order':
              await SalesOrderService.deleteSalesOrder(record.id);
              message.success('销售订单删除成功');
              loadSalesOrders();
              break;
          }
        } catch (error) {
          message.error(`删除失败: ${error}`);
          console.error('Delete error:', error);
        } finally {
          setLoading(false);
        }
      }
    });
  };

  // 查看详情
  const handleView = (record: any) => {
    setViewingRecord(record);
    setViewModalVisible(true);
  };

  // 编辑记录
  const handleEdit = (record: any, type: 'customer' | 'quote' | 'order') => {
    setEditingRecord(record);
    setModalType(type);
    
    // 根据不同类型处理表单数据
    let formData = { ...record };
    if (type === 'order') {
      // 销售订单编辑时，需要处理客户ID和日期格式
      formData = {
        ...record,
        customerId: record.customer?.id || record.customerId,
        orderDate: record.orderDate ? dayjs(record.orderDate) : null,
        deliveryDate: record.deliveryDate ? dayjs(record.deliveryDate) : null,
        totalAmount: record.grandTotal || record.totalAmount
      };
    } else if (type === 'quote') {
      // 报价单编辑时，需要处理客户ID和日期格式
      formData = {
        ...record,
        customerId: record.customer?.id || record.customerId,
        validTill: record.validTill ? dayjs(record.validTill) : null,
        totalAmount: record.grandTotal || record.totalAmount
      };
    }
    
    form.setFieldsValue(formData);
    setIsModalVisible(true);
  };

  // 版本管理相关函数
  const handleVersionManagement = (quotationId: number) => {
    setCurrentQuotationId(quotationId);
    loadVersionHistory(quotationId);
    setVersionHistoryVisible(true);
  };

  const loadVersionHistory = async (quotationId: number) => {
    try {
      setLoading(true);
      const history = await QuotationVersionService.getVersionHistory(quotationId);
      setVersions(history);
    } catch (error) {
      message.error('加载版本历史失败');
      console.error('Error loading version history:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateVersion = async (values: any) => {
    if (!currentQuotationId) return;
    
    try {
      setLoading(true);
      await QuotationVersionService.createVersion({
        quotation_id: currentQuotationId,
        version_name: values.version_name,
        change_reason: values.change_reason
      });
      message.success('版本创建成功');
      setVersionModalVisible(false);
      versionForm.resetFields();
      loadVersionHistory(currentQuotationId);
    } catch (error) {
      message.error('创建版本失败');
      console.error('Error creating version:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSetActiveVersion = async (quotationId: number, versionNumber: number) => {
    try {
      setLoading(true);
      await QuotationVersionService.setActiveVersion(quotationId, versionNumber);
      message.success('版本激活成功');
      loadVersionHistory(quotationId);
    } catch (error) {
      message.error('激活版本失败');
      console.error('Error setting active version:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleRollbackVersion = async (quotationId: number, versionId: number) => {
    Modal.confirm({
      title: '确认回滚',
      content: '确定要回滚到此版本吗？这将覆盖当前版本的数据。',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          setLoading(true);
          await QuotationVersionService.rollbackToVersion({
            quotation_id: quotationId,
            version_id: versionId,
            reason: '手动回滚'
          });
          message.success('版本回滚成功');
          loadVersionHistory(quotationId);
          loadQuotations(); // 刷新报价单列表
        } catch (error) {
          message.error('回滚版本失败');
          console.error('Error rolling back version:', error);
        } finally {
          setLoading(false);
        }
      }
    });
  };

  const handleDeleteVersion = async (versionId: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除此版本吗？此操作不可撤销。',
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          setLoading(true);
          await QuotationVersionService.deleteVersion(versionId);
          message.success('版本删除成功');
          if (currentQuotationId) {
            loadVersionHistory(currentQuotationId);
          }
        } catch (error) {
          message.error('删除版本失败');
          console.error('Error deleting version:', error);
        } finally {
          setLoading(false);
        }
      }
    });
  };

  const handleCompareVersions = async () => {
    if (selectedVersions.length !== 2) {
      message.warning('请选择两个版本进行比较');
      return;
    }

    if (!currentQuotationId) {
      message.error('未找到报价单ID');
      return;
    }

    try {
      setLoading(true);
      const response = await QuotationVersionService.compareVersions({
        quotation_id: currentQuotationId,
        from_version_id: selectedVersions[0],
        to_version_id: selectedVersions[1]
      });
      setCompareResult(response);
      setCompareModalVisible(true);
    } catch (error) {
      message.error('版本比较失败');
      console.error('Error comparing versions:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleVersionSelection = (versionId: number, checked: boolean) => {
    if (checked) {
      if (selectedVersions.length >= 2) {
        message.warning('最多只能选择两个版本进行比较');
        return;
      }
      setSelectedVersions([...selectedVersions, versionId]);
    } else {
      setSelectedVersions(selectedVersions.filter(id => id !== versionId));
    }
  };

  // 处理客户下拉菜单点击
  const handleCustomerMenuClick = (key: string, record: any) => {
    switch (key) {
      case 'view':
        handleView(record);
        break;
      case 'edit':
        handleEdit(record, 'customer');
        break;
      case 'delete':
        handleDelete(record, 'customer');
        break;
    }
  };

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
            onClick: ({ key }) => handleCustomerMenuClick(key, record),
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
      key: 'customer',
      render: (record: Quotation) => {
        // 正确处理客户对象显示
        if (record.customer && typeof record.customer === 'object') {
          return (
            <div>
              <Text strong>{record.customer.name}</Text>
              <br />
              <Text type="secondary" style={{ fontSize: 12 }}>
                {record.customer.code}
              </Text>
            </div>
          );
        }
        return <Text type="secondary">-</Text>;
      }
    },
    {
      title: '总金额',
      dataIndex: 'grandTotal',
      key: 'grandTotal',
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
      dataIndex: 'validTill',
      key: 'validTill',
      render: (date: string) => {
        if (!date) return '-';
        try {
          return new Date(date).toLocaleDateString();
        } catch {
          return '-';
        }
      },
    },
    {
      title: '创建日期',
      dataIndex: 'date',
      key: 'date',
      render: (date: string) => {
        if (!date) return '-';
        try {
          return new Date(date).toLocaleDateString();
        } catch {
          return '-';
        }
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Space>
          <Tooltip title="查看">
            <Button type="text" icon={<EyeOutlined />} size="small" onClick={() => handleView(record)} />
          </Tooltip>
          <Tooltip title="编辑">
            <Button type="text" icon={<EditOutlined />} size="small" onClick={() => handleEdit(record, 'quote')} />
          </Tooltip>
          <Tooltip title="版本管理">
            <Button type="text" icon={<BranchesOutlined />} size="small" onClick={() => handleVersionManagement(record.id)} />
          </Tooltip>
          <Tooltip title="删除">
            <Button type="text" icon={<DeleteOutlined />} size="small" danger onClick={() => handleDelete(record, 'quote')} />
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
      key: 'customer',
      render: (record: SalesOrder) => {
        // 正确处理客户对象显示
        if (record.customer && typeof record.customer === 'object') {
          return (
            <div>
              <Text strong>{record.customer.name}</Text>
              <br />
              <Text type="secondary" style={{ fontSize: 12 }}>
                {record.customer.code}
              </Text>
            </div>
          );
        }
        return <Text type="secondary">-</Text>;
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
            <Button type="text" icon={<EyeOutlined />} size="small" onClick={() => handleView(record)} />
          </Tooltip>
          <Tooltip title="编辑">
            <Button type="text" icon={<EditOutlined />} size="small" onClick={() => handleEdit(record, 'order')} />
          </Tooltip>
          <Tooltip title="删除">
            <Button type="text" icon={<DeleteOutlined />} size="small" danger onClick={() => handleDelete(record, 'order')} />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const handleModalOk = () => {
    form.validateFields().then(values => {
      if (editingRecord) {
        handleUpdate(values);
      } else {
        handleCreate(values);
      }
    }).catch(info => {
      console.log('Validate Failed:', info);
    });
  };

  const showModal = (type: 'customer' | 'quote' | 'order') => {
    setModalType(type);
    setEditingRecord(null);
    form.resetFields();
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
        title={
          editingRecord 
            ? `编辑${modalType === 'customer' ? '客户' : modalType === 'quote' ? '报价' : '订单'}`
            : `新建${modalType === 'customer' ? '客户' : modalType === 'quote' ? '报价' : '订单'}`
        }
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
                    name="validTill"
                    label="有效期至"
                    rules={[{ required: true, message: '请选择有效期' }]}
                  >
                    <DatePicker style={{ width: '100%' }} placeholder="请选择有效期" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="status"
                    label="状态"
                  >
                    <Select placeholder="请选择状态">
                      <Option value="draft">草稿</Option>
                      <Option value="submitted">已提交</Option>
                      <Option value="accepted">已接受</Option>
                      <Option value="rejected">已拒绝</Option>
                      <Option value="expired">已过期</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="subject"
                    label="主题"
                  >
                    <Input placeholder="请输入报价主题" />
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
                    name="orderDate"
                    label="订单日期"
                    rules={[{ required: true, message: '请选择订单日期' }]}
                  >
                    <DatePicker style={{ width: '100%' }} placeholder="请选择订单日期" />
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
                    name="status"
                    label="状态"
                  >
                    <Select placeholder="请选择状态">
                      <Option value="pending">待处理</Option>
                      <Option value="confirmed">已确认</Option>
                      <Option value="shipped">已发货</Option>
                      <Option value="delivered">已交付</Option>
                      <Option value="cancelled">已取消</Option>
                    </Select>
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

      {/* 详情查看模态框 */}
      <Modal
        title="详情信息"
        open={viewModalVisible}
        onCancel={() => setViewModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setViewModalVisible(false)}>
            关闭
          </Button>
        ]}
        width={800}
      >
        {viewingRecord && (
          <div>
            {Object.entries(viewingRecord).map(([key, value]) => (
              <Row key={key} style={{ marginBottom: 8 }}>
                <Col span={6}>
                  <Text strong>{key}:</Text>
                </Col>
                <Col span={18}>
                  <Text>{typeof value === 'object' ? JSON.stringify(value) : String(value)}</Text>
                </Col>
              </Row>
            ))}
          </div>
        )}
      </Modal>

      {/* 版本历史模态框 */}
      <Modal
        title="版本历史管理"
        open={versionHistoryVisible}
        onCancel={() => {
          setVersionHistoryVisible(false);
          setSelectedVersions([]);
        }}
        footer={[
          <Button key="create" type="primary" onClick={() => setVersionModalVisible(true)}>
            创建新版本
          </Button>,
          <Button 
            key="compare" 
            onClick={handleCompareVersions}
            disabled={selectedVersions.length !== 2}
          >
            比较版本 ({selectedVersions.length}/2)
          </Button>,
          <Button key="close" onClick={() => {
            setVersionHistoryVisible(false);
            setSelectedVersions([]);
          }}>
            关闭
          </Button>
        ]}
        width={1000}
      >
        <Table
          dataSource={versions.map(version => ({ ...version, key: version.id }))}
          columns={[
            {
              title: '选择',
              key: 'select',
              width: 60,
              render: (record: any) => (
                <Checkbox
                  checked={selectedVersions.includes(record.id)}
                  onChange={(e) => handleVersionSelection(record.id, e.target.checked)}
                  disabled={!selectedVersions.includes(record.id) && selectedVersions.length >= 2}
                />
              ),
            },
            {
              title: '版本号',
              dataIndex: 'version_number',
              key: 'version_number',
              render: (versionNumber: number, record: any) => (
                <Space>
                  <Text strong>v{versionNumber}</Text>
                  {record.is_active && <Tag color="green">当前版本</Tag>}
                </Space>
              ),
            },
            {
              title: '版本名称',
              dataIndex: 'version_name',
              key: 'version_name',
              render: (name: string) => name || '-',
            },
            {
              title: '变更原因',
              dataIndex: 'change_reason',
              key: 'change_reason',
              render: (reason: string) => reason || '-',
            },
            {
              title: '创建时间',
              dataIndex: 'created_at',
              key: 'created_at',
              render: (date: string) => {
                if (!date) return '-';
                try {
                  const dateObj = new Date(date);
                  return isNaN(dateObj.getTime()) ? '-' : dateObj.toLocaleString();
                } catch (error) {
                  return '-';
                }
              },
            },
            {
              title: '创建者',
              dataIndex: 'creator_name',
              key: 'creator_name',
              render: (name: string) => name || '-',
            },
            {
              title: '操作',
              key: 'action',
              render: (record: any) => (
                <Space>
                  {!record.is_active && (
                    <Tooltip title="设为当前版本">
                      <Button 
                        type="text" 
                        icon={<SwapOutlined />} 
                        size="small" 
                        onClick={() => handleSetActiveVersion(currentQuotationId!, record.version_number)}
                      />
                    </Tooltip>
                  )}
                  <Tooltip title="回滚到此版本">
                    <Button 
                      type="text" 
                      icon={<RollbackOutlined />} 
                      size="small" 
                      onClick={() => handleRollbackVersion(currentQuotationId!, record.id)}
                    />
                  </Tooltip>
                  {!record.is_active && (
                    <Tooltip title="删除版本">
                      <Button 
                        type="text" 
                        icon={<DeleteOutlined />} 
                        size="small" 
                        danger 
                        onClick={() => handleDeleteVersion(record.id)}
                      />
                    </Tooltip>
                  )}
                </Space>
              ),
            },
          ]}
          pagination={false}
          loading={loading}
        />
      </Modal>

      {/* 创建版本模态框 */}
      <Modal
        title="创建新版本"
        open={versionModalVisible}
        onCancel={() => {
          setVersionModalVisible(false);
          versionForm.resetFields();
        }}
        onOk={() => versionForm.submit()}
        confirmLoading={loading}
      >
        <Form
          form={versionForm}
          layout="vertical"
          onFinish={handleCreateVersion}
        >
          <Form.Item
            name="version_name"
            label="版本名称"
            rules={[{ max: 100, message: '版本名称不能超过100个字符' }]}
          >
            <Input placeholder="请输入版本名称（可选）" />
          </Form.Item>
          <Form.Item
            name="change_reason"
            label="变更原因"
            rules={[{ max: 500, message: '变更原因不能超过500个字符' }]}
          >
            <Input.TextArea 
              placeholder="请输入变更原因（可选）" 
              rows={4}
            />
          </Form.Item>
        </Form>
      </Modal>

      {/* 版本比较结果模态框 */}
      <Modal
        title="版本比较结果"
        open={compareModalVisible}
        onCancel={() => {
          setCompareModalVisible(false);
          setCompareResult(null);
          setSelectedVersions([]);
        }}
        footer={[
          <Button key="close" onClick={() => {
            setCompareModalVisible(false);
            setCompareResult(null);
            setSelectedVersions([]);
          }}>
            关闭
          </Button>
        ]}
        width={800}
      >
        {compareResult && (
          <div>
            <Table
              dataSource={compareResult.map((item: any, index: number) => ({ ...item, key: index }))}
              columns={[
                {
                  title: '字段名称',
                  dataIndex: 'field_name',
                  key: 'field_name',
                  width: 150,
                },
                {
                  title: '变更类型',
                  dataIndex: 'change_type',
                  key: 'change_type',
                  width: 100,
                  render: (type: string) => {
                    const colorMap: { [key: string]: string } = {
                      'added': 'green',
                      'modified': 'orange',
                      'deleted': 'red'
                    };
                    const textMap: { [key: string]: string } = {
                      'added': '新增',
                      'modified': '修改',
                      'deleted': '删除'
                    };
                    return <Tag color={colorMap[type]}>{textMap[type] || type}</Tag>;
                  },
                },
                {
                  title: '原值',
                  dataIndex: 'old_value',
                  key: 'old_value',
                  render: (value: any) => (
                    <Text code style={{ wordBreak: 'break-all' }}>
                      {typeof value === 'object' ? JSON.stringify(value) : String(value || '-')}
                    </Text>
                  ),
                },
                {
                  title: '新值',
                  dataIndex: 'new_value',
                  key: 'new_value',
                  render: (value: any) => (
                    <Text code style={{ wordBreak: 'break-all' }}>
                      {typeof value === 'object' ? JSON.stringify(value) : String(value || '-')}
                    </Text>
                  ),
                },
                {
                  title: '说明',
                  dataIndex: 'description',
                  key: 'description',
                  render: (desc: string) => desc || '-',
                },
              ]}
              pagination={false}
              size="small"
            />
          </div>
        )}
      </Modal>
    </div>
  );
}

export default withAuth(SalesPage);