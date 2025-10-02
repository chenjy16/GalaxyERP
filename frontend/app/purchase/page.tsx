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
  Badge,
  Upload
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
import { Supplier, PurchaseOrder, PurchaseRequest, Item, CreatePurchaseRequestRequest } from '@/types/api';
import SupplierService from '@/services/supplier';
import PurchaseOrderService from '@/services/purchaseOrder';
import PurchaseRequestService from '@/services/purchaseRequest';
import { ItemService } from '@/services/item';
import { withAuth } from '@/contexts/AuthContext';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Option } = Select;

function PurchasePage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'supplier' | 'order' | 'request' | 'supplierView'>('supplier');
  const [form] = Form.useForm();
  
  // 编辑状态
  const [editingRecord, setEditingRecord] = useState<any>(null);
  const [isEditing, setIsEditing] = useState(false);
  
  // 数据状态
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [purchaseOrders, setPurchaseOrders] = useState<PurchaseOrder[]>([]);
  const [purchaseRequests, setPurchaseRequests] = useState<PurchaseRequest[]>([]);
  const [items, setItems] = useState<Item[]>([]);
  
  // 加载状态
  const [suppliersLoading, setSuppliersLoading] = useState(false);
  const [ordersLoading, setOrdersLoading] = useState(false);
  const [requestsLoading, setRequestsLoading] = useState(false);
  const [itemsLoading, setItemsLoading] = useState(false);
  
  // 分页状态
  const [suppliersPagination, setSuppliersPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [ordersPagination, setOrdersPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [requestsPagination, setRequestsPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  
  // 搜索状态
  const [suppliersSearch, setSuppliersSearch] = useState('');
  const [ordersSearch, setOrdersSearch] = useState('');
  const [requestsSearch, setRequestsSearch] = useState('');
  
  // 筛选状态
  const [supplierTypeFilter, setSupplierTypeFilter] = useState('all');
  const [supplierStatusFilter, setSupplierStatusFilter] = useState('all');

  // 加载供应商数据
  const loadSuppliers = async (page = 1, search = '') => {
    setSuppliersLoading(true);
    try {
      const response = await SupplierService.getSuppliers({
        page,
        limit: suppliersPagination.pageSize,
        search: search || undefined
      });
      
      // 客户端筛选
      let filteredData = response.data;
      
      // 状态筛选
      if (supplierStatusFilter !== 'all') {
        filteredData = filteredData.filter((supplier: Supplier) => {
          if (supplierStatusFilter === 'active') {
            return supplier.isActive;
          } else if (supplierStatusFilter === 'inactive') {
            return !supplier.isActive;
          }
          return true;
        });
      }
      
      setSuppliers(filteredData);
      setSuppliersPagination(prev => ({
        ...prev,
        current: page,
        total: filteredData.length
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

  // 加载物料数据
  const loadItems = async () => {
    setItemsLoading(true);
    try {
      const response = await ItemService.getItems({
        page: 1,
        pageSize: 1000, // 获取所有物料用于选择
      });
      setItems(response.data);
    } catch (error) {
      message.error('加载物料数据失败');
      console.error('Error loading items:', error);
    } finally {
      setItemsLoading(false);
    }
  };

  // 初始化数据加载
  useEffect(() => {
    loadSuppliers();
    loadPurchaseOrders();
    loadPurchaseRequests();
    loadItems();
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
            <Text strong style={{ color: '#52c41a' }}>¥{record.creditLimit ? record.creditLimit.toLocaleString() : '0.00'}</Text>
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
                onClick: () => handleViewSupplier(record),
              },
              {
                key: 'edit',
                label: '编辑',
                icon: <EditOutlined />,
                onClick: () => handleEditSupplier(record),
              },
              {
                key: 'delete',
                label: '删除',
                icon: <DeleteOutlined />,
                danger: true,
                onClick: () => handleDeleteSupplier(record),
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
      dataIndex: 'order_number',
      key: 'order_number',
      render: (orderNumber: string) => <Text strong>{orderNumber || '-'}</Text>,
    },
    {
      title: '供应商',
      dataIndex: 'supplier',
      key: 'supplier',
      render: (supplier: any) => supplier?.name || '-',
    },
    {
      title: '金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (amount: number) => (
        <Text strong style={{ color: '#52c41a' }}>¥{amount ? amount.toLocaleString() : '0.00'}</Text>
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
      dataIndex: 'order_date',
      key: 'order_date',
      render: (date: string) => {
        if (!date) return '-';
        try {
          return new Date(date).toLocaleDateString('zh-CN');
        } catch (error) {
          return '-';
        }
      },
    },
    {
      title: '预期到货',
      dataIndex: 'expected_date',
      key: 'expected_date',
      render: (date: string) => {
        if (!date) return '-';
        try {
          return new Date(date).toLocaleDateString('zh-CN');
        } catch (error) {
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
            <Button 
              type="text" 
              icon={<EyeOutlined />} 
              size="small" 
              onClick={() => handleViewOrder(record)}
            />
          </Tooltip>
          <Tooltip title="编辑">
            <Button 
              type="text" 
              icon={<EditOutlined />} 
              size="small" 
              onClick={() => handleEditOrder(record)}
            />
          </Tooltip>
          <Tooltip title="删除">
            <Button 
              type="text" 
              icon={<DeleteOutlined />} 
              size="small" 
              danger 
              onClick={() => {
                Modal.confirm({
                  title: '确认删除',
                  content: '确定要删除这个采购订单吗？',
                  onOk: () => handleDeleteOrder(record.id),
                });
              }}
            />
          </Tooltip>
          {record.status === 'pending' && (
            <Tooltip title="确认">
              <Button 
                type="text" 
                icon={<CheckOutlined />} 
                size="small" 
                onClick={() => handleConfirmOrder(record.id)}
              />
            </Tooltip>
          )}
          {(record.status === 'pending' || record.status === 'confirmed') && (
            <Tooltip title="取消">
              <Button 
                type="text" 
                icon={<CloseOutlined />} 
                size="small" 
                danger 
                onClick={() => handleCancelOrder(record.id)}
              />
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  const requestColumns = [
    {
      title: '申请编号',
      dataIndex: 'number',
      key: 'number',
      render: (text: string) => <Text strong>{text || '-'}</Text>,
    },
    {
      title: '申请部门',
      dataIndex: 'department',
      key: 'department',
      render: (text: string) => text || '-',
    },
    {
      title: '申请人',
      dataIndex: ['createdBy', 'firstName'],
      key: 'requestedBy',
      render: (text: string, record: any) => {
        if (record.createdBy) {
          return `${record.createdBy.firstName || ''} ${record.createdBy.lastName || ''}`.trim() || '-';
        }
        return '-';
      },
    },
    {
      title: '预估金额',
      dataIndex: 'totalAmount',
      key: 'totalAmount',
      render: (amount: number) => (
        <Text strong style={{ color: '#52c41a' }}>¥{amount ? amount.toLocaleString() : '0.00'}</Text>
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
          <Tooltip title="查看详情">
            <Button 
              type="text" 
              icon={<EyeOutlined />} 
              size="small"
            />
          </Tooltip>
          {record.status === 'draft' && (
            <>
              <Tooltip title="编辑">
                <Button 
                  type="text" 
                  icon={<EditOutlined />} 
                  size="small"
                  onClick={() => showModal('request', record)}
                />
              </Tooltip>
              <Tooltip title="提交">
                <Button 
                  type="text" 
                  icon={<FileTextOutlined />} 
                  size="small"
                  style={{ color: '#1890ff' }}
                  onClick={() => handleSubmitRequest(record.id)}
                />
              </Tooltip>
              <Tooltip title="删除">
                <Button 
                  type="text" 
                  icon={<DeleteOutlined />} 
                  size="small"
                  style={{ color: '#ff4d4f' }}
                  onClick={() => {
                    Modal.confirm({
                      title: '确认删除',
                      content: '确定要删除这个采购申请吗？',
                      onOk: () => handleDeleteRequest(record.id),
                    });
                  }}
                />
              </Tooltip>
            </>
          )}
          {record.status === 'submitted' && (
            <>
              <Tooltip title="审批">
                <Button 
                  type="text" 
                  icon={<CheckOutlined />} 
                  size="small"
                  style={{ color: '#52c41a' }}
                  onClick={() => handleApproveRequest(record.id)}
                />
              </Tooltip>
              <Tooltip title="拒绝">
                <Button 
                  type="text" 
                  icon={<CloseOutlined />} 
                  size="small"
                  style={{ color: '#ff4d4f' }}
                  onClick={() => handleRejectRequest(record.id)}
                />
              </Tooltip>
            </>
          )}
        </Space>
      ),
    },
  ];

  // 显示模态框
  const showModal = (type: 'supplier' | 'order' | 'request', record?: any) => {
    setModalType(type);
    setEditingRecord(record);
    setIsEditing(!!record);
    
    if (record) {
      // 编辑模式，填充表单数据
      if (type === 'request') {
        form.setFieldsValue({
          ...record,
          items: record.items || []
        });
      } else {
        form.setFieldsValue(record);
      }
    } else {
      // 新建模式，重置表单
      form.resetFields();
    }
    
    setIsModalVisible(true);
  };

  // 处理模态框确认
  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      
      if (modalType === 'request') {
        if (isEditing) {
          // 编辑采购申请
          await PurchaseRequestService.updatePurchaseRequest(editingRecord.id, values);
          message.success('采购申请更新成功');
        } else {
          // 创建采购申请
          const createRequest: CreatePurchaseRequestRequest = {
            ...values,
            status: 'draft'
          };
          await PurchaseRequestService.createPurchaseRequest(createRequest);
          message.success('采购申请创建成功');
        }
        loadPurchaseRequests();
      } else if (modalType === 'order') {
        if (isEditing) {
          // 编辑采购订单
          await PurchaseOrderService.updatePurchaseOrder(editingRecord.id, values);
          message.success('采购订单更新成功');
        } else {
          // 创建采购订单
          await PurchaseOrderService.createPurchaseOrder(values);
          message.success('采购订单创建成功');
        }
        loadPurchaseOrders();
      } else if (modalType === 'supplier') {
        if (isEditing) {
          // 编辑供应商
          await SupplierService.updateSupplier(editingRecord.id, values);
          message.success('供应商更新成功');
        } else {
          // 创建供应商
          await SupplierService.createSupplier(values);
          message.success('供应商创建成功');
        }
        loadSuppliers();
      }
      
      setIsModalVisible(false);
      form.resetFields();
      setEditingRecord(null);
      setIsEditing(false);
    } catch (error) {
      console.error('Form validation failed:', error);
      message.error('操作失败，请检查输入信息');
    }
  };

  // 处理模态框取消
  const handleModalCancel = () => {
    setIsModalVisible(false);
    form.resetFields();
    setEditingRecord(null);
    setIsEditing(false);
  };

  // 删除采购申请
  const handleDeleteRequest = async (id: number) => {
    try {
      await PurchaseRequestService.deletePurchaseRequest(id);
      message.success('采购申请删除成功');
      loadPurchaseRequests();
    } catch (error) {
      message.error('删除采购申请失败');
      console.error('Error deleting purchase request:', error);
    }
  };

  // 查看采购订单
  const handleViewOrder = (record: PurchaseOrder) => {
    setEditingRecord(record);
    setModalType('order');
    setIsEditing(true);
    form.setFieldsValue({
      supplierId: record.supplierId,
      orderDate: record.orderDate ? new Date(record.orderDate) : null,
      expectedDate: record.expectedDate ? new Date(record.expectedDate) : null,
      status: record.status,
      notes: record.notes
    });
    setIsModalVisible(true);
  };

  // 编辑采购订单
  const handleEditOrder = (record: PurchaseOrder) => {
    setEditingRecord(record);
    setModalType('order');
    setIsEditing(true);
    form.setFieldsValue({
      supplierId: record.supplierId,
      orderDate: record.orderDate ? new Date(record.orderDate) : null,
      expectedDate: record.expectedDate ? new Date(record.expectedDate) : null,
      status: record.status,
      notes: record.notes
    });
    setIsModalVisible(true);
  };

  // 删除采购订单
  const handleDeleteOrder = async (id: number) => {
    try {
      await PurchaseOrderService.deletePurchaseOrder(id);
      message.success('采购订单删除成功');
      loadPurchaseOrders();
    } catch (error) {
      message.error('删除采购订单失败');
      console.error('Error deleting purchase order:', error);
    }
  };

  // 确认采购订单
  const handleConfirmOrder = async (id: number) => {
    try {
      await PurchaseOrderService.confirmPurchaseOrder(id);
      message.success('采购订单确认成功');
      loadPurchaseOrders();
    } catch (error) {
      message.error('确认采购订单失败');
      console.error('Error confirming purchase order:', error);
    }
  };

  // 取消采购订单
  const handleCancelOrder = async (id: number) => {
    Modal.confirm({
      title: '取消采购订单',
      content: '确定要取消这个采购订单吗？',
      onOk: async () => {
        try {
          await PurchaseOrderService.cancelPurchaseOrder(id);
          message.success('采购订单取消成功');
          loadPurchaseOrders();
        } catch (error) {
          message.error('取消采购订单失败');
          console.error('Error cancelling purchase order:', error);
        }
      }
    });
  };

  // 提交采购申请
  const handleSubmitRequest = async (id: number) => {
    try {
      await PurchaseRequestService.submitPurchaseRequest(id);
      message.success('采购申请提交成功');
      loadPurchaseRequests();
    } catch (error) {
      message.error('提交采购申请失败');
      console.error('Error submitting purchase request:', error);
    }
  };

  // 查看供应商详情
  const handleViewSupplier = (record: Supplier) => {
    setEditingRecord(record);
    setModalType('supplierView');
    setIsModalVisible(true);
  };

  // 编辑供应商
  const handleEditSupplier = (record: Supplier) => {
    setEditingRecord(record);
    setModalType('supplier');
    setIsEditing(true);
    form.setFieldsValue({
      name: record.name,
      contactName: record.contactName,
      phone: record.phone,
      email: record.email,
      address: record.address,
      creditLimit: record.creditLimit,
      paymentTerms: record.paymentTerms,
      isActive: record.isActive
    });
    setIsModalVisible(true);
  };

  // 删除供应商
  const handleDeleteSupplier = async (record: Supplier) => {
    Modal.confirm({
      title: '删除供应商',
      content: `确定要删除供应商 "${record.name}" 吗？此操作不可恢复。`,
      onOk: async () => {
        try {
          await SupplierService.deleteSupplier(record.id);
          message.success('供应商删除成功');
          loadSuppliers();
        } catch (error) {
          message.error('删除供应商失败');
          console.error('Error deleting supplier:', error);
        }
      }
    });
  };

  // 导出供应商数据
  const handleExportSuppliers = () => {
    try {
      // 准备导出数据
      const exportData = (suppliers || []).map((supplier: Supplier) => ({
        '供应商编码': supplier.code,
        '供应商名称': supplier.name,
        '联系人': supplier.contactName,
        '联系电话': supplier.phone,
        '邮箱': supplier.email,
        '税号': supplier.taxNumber,
        '信用额度': supplier.creditLimit,
        '付款条款': supplier.paymentTerms,
        '状态': supplier.isActive ? '启用' : '禁用',
        '地址': supplier.address,
        '创建时间': supplier.createdAt,
        '更新时间': supplier.updatedAt
      }));

      // 转换为CSV格式
      const headers = Object.keys(exportData[0] || {});
      const csvContent = [
        headers.join(','),
        ...exportData.map((row: any) => 
          headers.map(header => `"${row[header] || ''}"`).join(',')
        )
      ].join('\n');

      // 创建下载链接
      const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      const url = URL.createObjectURL(blob);
      link.setAttribute('href', url);
      link.setAttribute('download', `供应商数据_${new Date().toISOString().split('T')[0]}.csv`);
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      message.success('供应商数据导出成功');
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败');
    }
  };

  // 导入供应商数据
  const handleImportSuppliers = async (file: File) => {
    try {
      const fileExtension = file.name.split('.').pop()?.toLowerCase();
      let importData: any[] = [];

      if (fileExtension === 'csv') {
        // 解析CSV文件
        const text = await file.text();
        const lines = text.split('\n');
        const headers = lines[0].split(',').map(h => h.replace(/"/g, '').trim());
        
        for (let i = 1; i < lines.length; i++) {
          if (lines[i].trim()) {
            const values = lines[i].split(',').map(v => v.replace(/"/g, '').trim());
            const row: any = {};
            headers.forEach((header, index) => {
              row[header] = values[index] || '';
            });
            importData.push(row);
          }
        }
      } else if (fileExtension === 'xlsx' || fileExtension === 'xls') {
        // 对于Excel文件，这里需要使用xlsx库
        // 由于没有安装xlsx库，暂时提示用户使用CSV格式
        message.warning('请使用CSV格式文件进行导入');
        return false;
      }

      // 验证和转换数据格式
      const validSuppliers: any[] = [];
      
      for (const row of importData) {
        try {
          const supplier = {
            code: row['供应商编码'] || row['code'] || '',
            name: row['供应商名称'] || row['name'] || '',
            contactName: row['联系人'] || row['contactName'] || '',
            contactPhone: row['联系电话'] || row['contactPhone'] || '',
            email: row['邮箱'] || row['email'] || '',
            taxNumber: row['税号'] || row['taxNumber'] || '',
            creditLimit: parseFloat(row['信用额度'] || row['creditLimit'] || '0'),
            paymentTerms: row['付款条款'] || row['paymentTerms'] || '',
            isActive: (row['状态'] || row['isActive']) === '启用' || (row['状态'] || row['isActive']) === 'true',
            address: row['地址'] || row['address'] || ''
          };
          
          if (supplier.name && supplier.code) {
            validSuppliers.push(supplier);
          }
        } catch (error) {
          console.warn('跳过无效行:', row, error);
        }
      }

      if (validSuppliers.length === 0) {
        message.warning('没有找到有效的数据行');
        return false;
      }

      // 批量创建供应商
      let successCount = 0;
      for (const supplier of validSuppliers) {
        try {
          await SupplierService.createSupplier(supplier);
          successCount++;
        } catch (error) {
          console.error('创建供应商失败:', error);
        }
      }

      message.success(`成功导入 ${successCount} 条供应商数据`);
      await loadSuppliers();
      return true;
    } catch (error) {
      console.error('导入失败:', error);
      message.error('导入失败，请检查文件格式');
      return false;
    }
  };

  // 审批采购申请
  const handleApproveRequest = async (id: number) => {
    try {
      await PurchaseRequestService.approvePurchaseRequest(id);
      message.success('采购申请审批成功');
      loadPurchaseRequests();
    } catch (error) {
      message.error('审批采购申请失败');
      console.error('Error approving purchase request:', error);
    }
  };

  // 拒绝采购申请
  const handleRejectRequest = async (id: number) => {
    try {
      await PurchaseRequestService.rejectPurchaseRequest(id);
      message.success('采购申请已拒绝');
      loadPurchaseRequests();
    } catch (error) {
      message.error('拒绝采购申请失败');
      console.error('Error rejecting purchase request:', error);
    }
  };

  // 导出采购申请数据
  const handleExportRequests = () => {
    try {
      // 准备导出数据
      const exportData = purchaseRequests.map(request => ({
        '申请编号': request.number,
        '申请标题': request.title,
        '申请部门': request.department,
        '申请人': request.createdBy?.firstName || '未知',
        '紧急程度': request.priority === 'high' ? '紧急' : request.priority === 'medium' ? '一般' : '不急',
        '状态': request.status === 'draft' ? '草稿' : 
               request.status === 'submitted' ? '已提交' : 
               request.status === 'approved' ? '已审批' : 
               request.status === 'rejected' ? '已拒绝' : request.status,
        '创建时间': new Date(request.createdAt).toLocaleString(),
        '期望交付日期': request.requiredDate ? new Date(request.requiredDate).toLocaleDateString() : '',
        '申请说明': request.description || ''
      }));

      // 创建CSV内容
      const headers = Object.keys(exportData[0] || {});
      const csvContent = [
        headers.join(','),
        ...exportData.map(row => 
          headers.map(header => `"${(row as any)[header] || ''}"`).join(',')
        )
      ].join('\n');

      // 创建并下载文件
      const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      const url = URL.createObjectURL(blob);
      link.setAttribute('href', url);
      link.setAttribute('download', `采购申请_${new Date().toISOString().split('T')[0]}.csv`);
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      
      message.success('导出成功');
    } catch (error) {
      message.error('导出失败');
      console.error('Export error:', error);
    }
  };

  // 导入采购申请数据
  const handleImportRequests = async (file: File) => {
    try {
      const fileExtension = file.name.split('.').pop()?.toLowerCase();
      let importData: any[] = [];

      if (fileExtension === 'csv') {
        // 解析CSV文件
        const text = await file.text();
        const lines = text.split('\n');
        const headers = lines[0].split(',').map(h => h.replace(/"/g, '').trim());
        
        for (let i = 1; i < lines.length; i++) {
          if (lines[i].trim()) {
            const values = lines[i].split(',').map(v => v.replace(/"/g, '').trim());
            const row: any = {};
            headers.forEach((header, index) => {
              row[header] = values[index] || '';
            });
            importData.push(row);
          }
        }
      } else if (fileExtension === 'xlsx' || fileExtension === 'xls') {
        // 对于Excel文件，这里需要使用xlsx库
        // 由于没有安装xlsx库，暂时提示用户使用CSV格式
        message.warning('请使用CSV格式文件进行导入');
        return false;
      }

      // 验证和转换数据格式
      const validRequests: CreatePurchaseRequestRequest[] = [];
      
      for (const row of importData) {
        try {
          const request: CreatePurchaseRequestRequest = {
             title: row['申请标题'] || row['title'] || '',
             description: row['申请说明'] || row['description'] || '',
             priority: row['紧急程度'] === '紧急' ? 'high' : 
                      row['紧急程度'] === '一般' ? 'medium' : 'low',
             requiredDate: row['期望交付日期'] || row['requiredDate'] || '',
             items: [] // 暂时为空，可以后续扩展支持物料导入
           };
          
          if (request.title) {
            validRequests.push(request);
          }
        } catch (error) {
          console.warn('跳过无效行:', row, error);
        }
      }

      if (validRequests.length === 0) {
        message.warning('没有找到有效的数据行');
        return false;
      }

      // 批量创建采购申请
      let successCount = 0;
      for (const request of validRequests) {
        try {
          await PurchaseRequestService.createPurchaseRequest(request);
          successCount++;
        } catch (error) {
          console.error('创建采购申请失败:', error);
        }
      }
      
      message.success(`成功导入 ${successCount} 条采购申请`);
      loadPurchaseRequests();
    } catch (error) {
      console.error('导入失败:', error);
      message.error('导入失败，请检查文件格式');
    }
    return false; // 阻止默认上传行为
  };

  // 导出采购订单数据
  const handleExportOrders = () => {
    try {
      // 准备导出数据
      const exportData = (purchaseOrders || []).map((order: PurchaseOrder) => ({
        '订单编号': order.orderNumber,
        '供应商': suppliers.find(s => s.id === order.supplierId)?.name || '',
        '订单日期': order.orderDate,
        '期望交付日期': order.expectedDate,
        '状态': order.status === 'confirmed' ? '已确认' :
               order.status === 'shipped' ? '已发货' :
               order.status === 'delivered' ? '已交付' :
               order.status === 'cancelled' ? '已取消' : '待确认',
        '总金额': order.totalAmount || 0,
        '创建时间': order.createdAt,
        '更新时间': order.updatedAt
      }));

      // 转换为CSV格式
      const headers = Object.keys(exportData[0] || {});
      const csvContent = [
        headers.join(','),
        ...exportData.map((row: any) => 
          headers.map(header => `"${row[header] || ''}"`).join(',')
        )
      ].join('\n');

      // 创建下载链接
      const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      const url = URL.createObjectURL(blob);
      link.setAttribute('href', url);
      link.setAttribute('download', `采购订单_${new Date().toISOString().split('T')[0]}.csv`);
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      message.success('采购订单数据导出成功');
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败');
    }
  };

  // 导入采购订单数据
  const handleImportOrders = async (file: File) => {
    try {
      const fileExtension = file.name.split('.').pop()?.toLowerCase();
      let importData: any[] = [];

      if (fileExtension === 'csv') {
        // 解析CSV文件
        const text = await file.text();
        const lines = text.split('\n');
        const headers = lines[0].split(',').map(h => h.replace(/"/g, '').trim());
        
        for (let i = 1; i < lines.length; i++) {
          if (lines[i].trim()) {
            const values = lines[i].split(',').map(v => v.replace(/"/g, '').trim());
            const row: any = {};
            headers.forEach((header, index) => {
              row[header] = values[index] || '';
            });
            importData.push(row);
          }
        }
      } else if (fileExtension === 'xlsx' || fileExtension === 'xls') {
        // 对于Excel文件，这里需要使用xlsx库
        // 由于没有安装xlsx库，暂时提示用户使用CSV格式
        message.warning('请使用CSV格式文件进行导入');
        return false;
      }

      // 验证和转换数据格式
      const validOrders: any[] = [];
      
      for (const row of importData) {
        try {
          // 查找供应商ID
          const supplierName = row['供应商'] || row['supplier'];
          const supplier = suppliers.find(s => s.name === supplierName);
          
          if (!supplier) {
            console.warn('未找到供应商:', supplierName);
            continue;
          }

          const order = {
            supplierId: supplier.id,
            orderDate: row['订单日期'] || row['orderDate'] || new Date().toISOString().split('T')[0],
            expectedDate: row['期望交付日期'] || row['expectedDate'] || '',
            items: [] // 暂时为空，可以后续扩展支持物料导入
          };
          
          if (order.supplierId) {
            validOrders.push(order);
          }
        } catch (error) {
          console.warn('跳过无效行:', row, error);
        }
      }

      if (validOrders.length === 0) {
        message.warning('没有找到有效的数据行');
        return false;
      }

      // 批量创建采购订单
      let successCount = 0;
      for (const order of validOrders) {
        try {
          await PurchaseOrderService.createPurchaseOrder(order);
          successCount++;
        } catch (error) {
          console.error('创建采购订单失败:', error);
        }
      }

      message.success(`成功导入 ${successCount} 条采购订单`);
      await loadPurchaseOrders();
      return true;
    } catch (error) {
      console.error('导入失败:', error);
      message.error('导入失败，请检查文件格式');
      return false;
    }
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
              <Select 
                placeholder="类型筛选" 
                style={{ width: 120 }}
                value={supplierTypeFilter}
                onChange={(value) => {
                  setSupplierTypeFilter(value);
                  loadSuppliers(1, suppliersSearch);
                }}
              >
                <Option value="all">全部</Option>
                <Option value="电子产品">电子产品</Option>
                <Option value="办公用品">办公用品</Option>
                <Option value="原材料">原材料</Option>
                <Option value="设备">设备</Option>
              </Select>
              <Select 
                placeholder="状态筛选" 
                style={{ width: 120 }}
                value={supplierStatusFilter}
                onChange={(value) => {
                  setSupplierStatusFilter(value);
                  loadSuppliers(1, suppliersSearch);
                }}
              >
                <Option value="all">全部</Option>
                <Option value="active">启用</Option>
                <Option value="inactive">禁用</Option>
              </Select>
            </Space>
            <Space>
              <Upload
                accept=".csv,.xlsx,.xls"
                showUploadList={false}
                beforeUpload={(file) => {
                  handleImportSuppliers(file);
                  return false;
                }}
              >
                <Button icon={<ImportOutlined />}>导入</Button>
              </Upload>
              <Button icon={<ExportOutlined />} onClick={handleExportSuppliers}>导出</Button>
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
               <Button icon={<ExportOutlined />} onClick={handleExportRequests}>导出</Button>
               <Upload
                 accept=".csv,.xlsx,.xls"
                 showUploadList={false}
                 beforeUpload={(file) => {
                   handleImportRequests(file);
                   return false; // 阻止自动上传
                 }}
               >
                 <Button icon={<ImportOutlined />}>
                   导入
                 </Button>
               </Upload>
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

      {/* 模态框 */}
      <Modal
        title={
          modalType === 'request' 
            ? (isEditing ? '编辑采购申请' : '新建采购申请')
            : modalType === 'supplier' 
            ? (isEditing ? '编辑供应商' : '新建供应商')
            : (isEditing ? '编辑采购订单' : '新建采购订单')
        }
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        width={800}
        destroyOnHidden
      >
        {modalType === 'request' && (
          <Form
            form={form}
            layout="vertical"
            initialValues={{
              priority: 'medium',
              status: 'draft',
              items: [{}]
            }}
          >
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="title"
                  label="申请标题"
                  rules={[{ required: true, message: '请输入申请标题' }]}
                >
                  <Input placeholder="请输入申请标题" />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item
                  name="department"
                  label="申请部门"
                  rules={[{ required: true, message: '请输入申请部门' }]}
                >
                  <Input placeholder="请输入申请部门" />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="priority"
                  label="紧急程度"
                  rules={[{ required: true, message: '请选择紧急程度' }]}
                >
                  <Select placeholder="请选择紧急程度">
                    <Option value="high">紧急</Option>
                    <Option value="medium">一般</Option>
                    <Option value="low">不急</Option>
                  </Select>
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item
                  name="expectedDate"
                  label="期望交付日期"
                  rules={[{ required: true, message: '请选择期望交付日期' }]}
                >
                  <DatePicker 
                    style={{ width: '100%' }} 
                    placeholder="请选择期望交付日期"
                  />
                </Form.Item>
              </Col>
            </Row>

            <Form.Item
              name="description"
              label="申请说明"
            >
              <Input.TextArea 
                rows={3} 
                placeholder="请输入申请说明"
              />
            </Form.Item>

            <Form.List name="items">
              {(fields, { add, remove }) => (
                <>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                    <Text strong>申请物料</Text>
                    <Button 
                      type="dashed" 
                      onClick={() => add()} 
                      icon={<PlusOutlined />}
                    >
                      添加物料
                    </Button>
                  </div>
                  {fields.map(({ key, name, ...restField }) => (
                    <Card 
                      key={key} 
                      size="small" 
                      style={{ marginBottom: 16 }}
                      extra={
                        fields.length > 1 && (
                          <Button 
                            type="text" 
                            danger 
                            icon={<DeleteOutlined />}
                            onClick={() => remove(name)}
                          />
                        )
                      }
                    >
                      <Row gutter={16}>
                        <Col span={8}>
                          <Form.Item
                            {...restField}
                            name={[name, 'itemId']}
                            label="物料"
                            rules={[{ required: true, message: '请选择物料' }]}
                          >
                            <Select 
                              placeholder="请选择物料"
                              showSearch
                              filterOption={(input, option) =>
                                (option?.children as unknown as string)
                                  ?.toLowerCase()
                                  ?.includes(input.toLowerCase())
                              }
                              loading={itemsLoading}
                            >
                              {items.map(item => (
                                <Option key={item.id} value={item.id}>
                                  {item.name} ({item.code})
                                </Option>
                              ))}
                            </Select>
                          </Form.Item>
                        </Col>
                        <Col span={6}>
                          <Form.Item
                            {...restField}
                            name={[name, 'quantity']}
                            label="数量"
                            rules={[
                              { required: true, message: '请输入数量' },
                              { type: 'number', min: 1, message: '数量必须大于0' }
                            ]}
                          >
                            <Input type="number" placeholder="请输入数量" />
                          </Form.Item>
                        </Col>
                        <Col span={6}>
                          <Form.Item
                            {...restField}
                            name={[name, 'estimatedPrice']}
                            label="预估单价"
                          >
                            <Input type="number" placeholder="预估单价" />
                          </Form.Item>
                        </Col>
                        <Col span={4}>
                          <Form.Item
                            {...restField}
                            name={[name, 'unit']}
                            label="单位"
                          >
                            <Input placeholder="单位" />
                          </Form.Item>
                        </Col>
                      </Row>
                      <Form.Item
                        {...restField}
                        name={[name, 'notes']}
                        label="备注"
                      >
                        <Input.TextArea rows={2} placeholder="物料备注" />
                      </Form.Item>
                    </Card>
                  ))}
                </>
              )}
            </Form.List>
          </Form>
        )}

        {modalType === 'supplier' && (
          <Form form={form} layout="vertical">
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="name"
                  label="供应商名称"
                  rules={[{ required: true, message: '请输入供应商名称' }]}
                >
                  <Input placeholder="请输入供应商名称" />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item
                  name="code"
                  label="供应商编码"
                  rules={[{ required: true, message: '请输入供应商编码' }]}
                >
                  <Input placeholder="请输入供应商编码" />
                </Form.Item>
              </Col>
            </Row>
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="contactName"
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
                  name="taxNumber"
                  label="税号"
                >
                  <Input placeholder="请输入税号" />
                </Form.Item>
              </Col>
            </Row>
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="creditLimit"
                  label="信用额度"
                >
                  <Input type="number" placeholder="请输入信用额度" />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item
                  name="paymentTerms"
                  label="付款条款"
                >
                  <Input placeholder="请输入付款条款" />
                </Form.Item>
              </Col>
            </Row>
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="isActive"
                  label="状态"
                  valuePropName="checked"
                >
                  <Select placeholder="请选择状态">
                    <Option value={true}>活跃</Option>
                    <Option value={false}>非活跃</Option>
                  </Select>
                </Form.Item>
              </Col>
            </Row>
            <Form.Item
              name="address"
              label="地址"
            >
              <Input.TextArea rows={3} placeholder="请输入地址" />
            </Form.Item>
          </Form>
        )}

        {modalType === 'supplierView' && editingRecord && (
          <div>
            <Row gutter={16}>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">供应商名称:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Text strong>{(editingRecord as Supplier).name}</Text>
                  </div>
                </div>
              </Col>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">供应商编码:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Text strong>{(editingRecord as Supplier).code}</Text>
                  </div>
                </div>
              </Col>
            </Row>
            <Row gutter={16}>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">联系人:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Text strong>{(editingRecord as Supplier).contactName || '-'}</Text>
                  </div>
                </div>
              </Col>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">联系电话:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Text strong>{(editingRecord as Supplier).phone || '-'}</Text>
                  </div>
                </div>
              </Col>
            </Row>
            <Row gutter={16}>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">邮箱:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Text strong>{(editingRecord as Supplier).email || '-'}</Text>
                  </div>
                </div>
              </Col>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">税号:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Text strong>{(editingRecord as Supplier).taxNumber || '-'}</Text>
                  </div>
                </div>
              </Col>
            </Row>
            <Row gutter={16}>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">信用额度:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Text strong style={{ color: '#52c41a' }}>
                      ¥{(editingRecord as Supplier).creditLimit ? (editingRecord as Supplier).creditLimit.toLocaleString() : '0.00'}
                    </Text>
                  </div>
                </div>
              </Col>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">付款条款:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Text strong>{(editingRecord as Supplier).paymentTerms || '-'}</Text>
                  </div>
                </div>
              </Col>
            </Row>
            <Row gutter={16}>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">状态:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Tag color={(editingRecord as Supplier).isActive ? 'green' : 'red'}>
                      {(editingRecord as Supplier).isActive ? '活跃' : '非活跃'}
                    </Tag>
                  </div>
                </div>
              </Col>
              <Col span={12}>
                <div style={{ marginBottom: 16 }}>
                  <Text type="secondary">创建时间:</Text>
                  <div style={{ marginTop: 4 }}>
                    <Text strong>{new Date((editingRecord as Supplier).createdAt).toLocaleString()}</Text>
                  </div>
                </div>
              </Col>
            </Row>
            <div style={{ marginBottom: 16 }}>
              <Text type="secondary">地址:</Text>
              <div style={{ marginTop: 4 }}>
                <Text strong>{(editingRecord as Supplier).address || '-'}</Text>
              </div>
            </div>
          </div>
        )}

        {modalType === 'order' && (
          <Form form={form} layout="vertical">
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="supplierId"
                  label="供应商"
                  rules={[{ required: true, message: '请选择供应商' }]}
                >
                  <Select placeholder="请选择供应商">
                    {suppliers.map(supplier => (
                      <Option key={supplier.id} value={supplier.id}>
                        {supplier.name}
                      </Option>
                    ))}
                  </Select>
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item
                  name="orderDate"
                  label="订单日期"
                  rules={[{ required: true, message: '请选择订单日期' }]}
                >
                  <DatePicker style={{ width: '100%' }} />
                </Form.Item>
              </Col>
            </Row>
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  name="expectedDate"
                  label="期望交付日期"
                  rules={[{ required: true, message: '请选择期望交付日期' }]}
                >
                  <DatePicker style={{ width: '100%' }} />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item
                  name="status"
                  label="状态"
                  rules={[{ required: true, message: '请选择状态' }]}
                >
                  <Select placeholder="请选择状态">
                    <Option value="pending">待确认</Option>
                    <Option value="confirmed">已确认</Option>
                    <Option value="shipped">已发货</Option>
                    <Option value="delivered">已交付</Option>
                    <Option value="cancelled">已取消</Option>
                  </Select>
                </Form.Item>
              </Col>
            </Row>
            
            {/* 物料项目 */}
            <Form.Item label="物料项目">
              <Form.List name="items">
                {(fields, { add, remove }) => (
                  <>
                    {fields.map(({ key, name, ...restField }) => (
                      <Row key={key} gutter={16} style={{ marginBottom: 16 }}>
                        <Col span={8}>
                          <Form.Item
                            {...restField}
                            name={[name, 'itemId']}
                            label="物料"
                            rules={[{ required: true, message: '请选择物料' }]}
                          >
                            <Select placeholder="请选择物料" showSearch>
                              {items.map(item => (
                                <Option key={item.id} value={item.id}>
                                  {item.name} ({item.code})
                                </Option>
                              ))}
                            </Select>
                          </Form.Item>
                        </Col>
                        <Col span={4}>
                          <Form.Item
                            {...restField}
                            name={[name, 'quantity']}
                            label="数量"
                            rules={[
                              { required: true, message: '请输入数量' },
                              { type: 'number', min: 1, message: '数量必须大于0' }
                            ]}
                          >
                            <Input type="number" placeholder="数量" />
                          </Form.Item>
                        </Col>
                        <Col span={4}>
                          <Form.Item
                            {...restField}
                            name={[name, 'unitPrice']}
                            label="单价"
                            rules={[{ required: true, message: '请输入单价' }]}
                          >
                            <Input type="number" placeholder="单价" />
                          </Form.Item>
                        </Col>
                        <Col span={4}>
                          <Form.Item
                            {...restField}
                            name={[name, 'taxRate']}
                            label="税率(%)"
                          >
                            <Input type="number" placeholder="税率" />
                          </Form.Item>
                        </Col>
                        <Col span={3}>
                          <Form.Item
                            {...restField}
                            name={[name, 'notes']}
                            label="备注"
                          >
                            <Input placeholder="备注" />
                          </Form.Item>
                        </Col>
                        <Col span={1}>
                          <Form.Item label=" ">
                            <Button 
                              type="text" 
                              danger 
                              icon={<DeleteOutlined />} 
                              onClick={() => remove(name)}
                            />
                          </Form.Item>
                        </Col>
                      </Row>
                    ))}
                    <Form.Item>
                      <Button 
                        type="dashed" 
                        onClick={() => add()} 
                        block 
                        icon={<PlusOutlined />}
                      >
                        添加物料
                      </Button>
                    </Form.Item>
                  </>
                )}
              </Form.List>
            </Form.Item>

            <Form.Item
              name="notes"
              label="备注"
            >
              <Input.TextArea rows={3} placeholder="请输入备注" />
            </Form.Item>
          </Form>
        )}
      </Modal>
    </div>
  );
}

export default withAuth(PurchasePage);