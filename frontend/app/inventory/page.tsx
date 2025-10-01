'use client';

import { useState, useEffect } from 'react';
import { ProductService } from '@/services/product';
import { InventoryService } from '@/services/inventory';
import { WarehouseService } from '@/services/warehouse';
import { withAuth } from '@/contexts/AuthContext';
import { 
  Product, 
  Stock, 
  StockMovement, 
  Warehouse, 
  CreateProductRequest,
  CreateStockRequest,
  CreateStockMovementRequest,
  CreateWarehouseRequest,
  PaginatedResponse 
} from '@/types/api';
import { 
  Card, 
  Tabs, 
  Table, 
  Button, 
  Input, 
  Space, 
  Tag,
  Progress,
  Statistic,
  Row,
  Col,
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
  Alert,
  Divider
} from 'antd';
import { 
  PlusOutlined, 
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  WarningOutlined,
  EyeOutlined,
  MoreOutlined,
  BoxPlotOutlined,
  InboxOutlined,
  SendOutlined,
  AlertOutlined,
  BarChartOutlined,
  ExportOutlined,
  ImportOutlined,
  ShopOutlined,
  HomeOutlined,
  RiseOutlined,
  FallOutlined,
  ExclamationCircleOutlined,
  SwapOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Option } = Select;

function InventoryPage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'product' | 'inbound' | 'outbound' | 'warehouse' | 'movement'>('product');
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  // API数据状态
  const [products, setProducts] = useState<Product[]>([]);
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [stockMovements, setStockMovements] = useState<StockMovement[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  
  // 分页状态
  const [productsPagination, setProductsPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0
  });
  const [stocksPagination, setStocksPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0
  });
  const [movementsPagination, setMovementsPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0
  });
  const [warehousesPagination, setWarehousesPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0
  });

  // 搜索状态
  const [productsSearch, setProductsSearch] = useState('');
  const [stocksSearch, setStocksSearch] = useState('');
  const [movementsSearch, setMovementsSearch] = useState('');
  const [warehousesSearch, setWarehousesSearch] = useState('');

  // 数据加载函数
  const loadProducts = async (page = 1, search = '', pageSize?: number) => {
    try {
      setLoading(true);
      const limit = pageSize || productsPagination.pageSize;
      const response = await ProductService.getProducts({
        page,
        limit,
        search: search || undefined
      });
      setProducts(response.data || []);
      setProductsPagination({
        current: response.page,
        pageSize: response.limit,
        total: response.total
      });
    } catch (error) {
      console.error('加载产品数据失败:', error);
      message.error('加载产品数据失败');
      setProducts([]);
    } finally {
      setLoading(false);
    }
  };

  const loadStocks = async (page = 1, search = '', pageSize?: number) => {
    try {
      setLoading(true);
      const limit = pageSize || stocksPagination.pageSize;
      const response = await InventoryService.getStocks({
        page,
        limit,
        search: search || undefined
      });
      setStocks(response.data || []);
      setStocksPagination({
        current: response.page,
        pageSize: response.limit,
        total: response.total
      });
    } catch (error) {
      console.error('加载库存数据失败:', error);
      message.error('加载库存数据失败');
      setStocks([]);
    } finally {
      setLoading(false);
    }
  };

  const loadStockMovements = async (page = 1, search = '', pageSize?: number) => {
    try {
      setLoading(true);
      const limit = pageSize || movementsPagination.pageSize;
      const response = await InventoryService.getStockMovements({
        page,
        limit,
        search: search || undefined
      });
      setStockMovements(response.data || []);
      setMovementsPagination({
        current: response.page,
        pageSize: response.limit,
        total: response.total
      });
    } catch (error) {
      console.error('加载库存移动数据失败:', error);
      message.error('加载库存移动数据失败');
      setStockMovements([]);
    } finally {
      setLoading(false);
    }
  };

  const loadWarehouses = async (page = 1, search = '', pageSize?: number) => {
    try {
      setLoading(true);
      const limit = pageSize || warehousesPagination.pageSize;
      const response = await WarehouseService.getWarehouses({
        page,
        limit,
        search: search || undefined
      });
      setWarehouses(response.data || []);
      setWarehousesPagination({
        current: response.page,
        pageSize: response.limit,
        total: response.total
      });
    } catch (error) {
      console.error('加载仓库数据失败:', error);
      message.error('加载仓库数据失败');
      setWarehouses([]);
    } finally {
      setLoading(false);
    }
  };

  // 初始化数据加载
  useEffect(() => {
    loadProducts();
    loadStocks();
    loadStockMovements();
    loadWarehouses();
  }, []);



  const getStockStatus = (stock: number, minStock: number, maxStock: number) => {
    if (stock <= minStock * 0.5) return 'critical';
    if (stock <= minStock) return 'low';
    if (stock >= maxStock * 0.9) return 'high';
    return 'normal';
  };

  const getStockStatusColor = (status: string) => {
    const colors: { [key: string]: string } = {
      critical: 'red',
      low: 'orange',
      normal: 'green',
      high: 'blue'
    };
    return colors[status] || 'default';
  };

  const getStockStatusText = (status: string) => {
    const texts: { [key: string]: string } = {
      critical: '严重不足',
      low: '库存不足',
      normal: '正常',
      high: '库存充足'
    };
    return texts[status] || status;
  };

  const getTypeColor = (type: string) => {
    return type === 'inbound' ? 'green' : 'red';
  };

  const getTypeText = (type: string) => {
    return type === 'inbound' ? '入库' : '出库';
  };

  const productColumns = [
    {
      title: '商品信息',
      key: 'product',
      render: (record: Product) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Avatar 
            size={40}
            style={{ backgroundColor: '#1890ff' }}
            icon={<InboxOutlined />}
          />
          <div>
            <Text strong style={{ display: 'block' }}>{record.name}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {record.code} | {record.category}
            </Text>
          </div>
        </div>
      ),
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      render: (category: string) => <Tag color="blue">{category}</Tag>,
    },
    {
      title: '单位',
      dataIndex: 'unit',
      key: 'unit',
      render: (unit: string) => <Text>{unit}</Text>,
    },
    {
      title: '价格信息',
      key: 'price',
      render: (record: Product) => (
        <div>
          <div style={{ marginBottom: 4 }}>
            <Text type="secondary">售价: </Text>
            <Text strong style={{ color: '#52c41a' }}>¥{record.price.toLocaleString()}</Text>
          </div>
          <div>
            <Text type="secondary">成本: </Text>
            <Text>¥{record.cost.toLocaleString()}</Text>
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
          {status === 'active' ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: Product) => (
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
                key: 'inbound',
                label: '入库',
                icon: <InboxOutlined />,
              },
              {
                key: 'outbound',
                label: '出库',
                icon: <SendOutlined />,
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

  const stockColumns = [
    {
      title: '单据编号',
      dataIndex: 'id',
      key: 'id',
      render: (id: string) => <Text strong>{id}</Text>,
    },
    {
      title: '商品名称',
      dataIndex: 'product',
      key: 'product',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={getTypeColor(type)} icon={type === 'inbound' ? <InboxOutlined /> : <SendOutlined />}>
          {getTypeText(type)}
        </Tag>
      ),
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      key: 'quantity',
      render: (quantity: number, record: any) => (
        <Text strong style={{ color: record.type === 'inbound' ? '#52c41a' : '#ff4d4f' }}>
          {record.type === 'inbound' ? '+' : '-'}{quantity}
        </Text>
      ),
    },
    {
      title: '仓库',
      dataIndex: 'warehouse',
      key: 'warehouse',
    },
    {
      title: '操作人',
      dataIndex: 'operator',
      key: 'operator',
    },
    {
      title: '操作时间',
      dataIndex: 'time',
      key: 'time',
    },
    {
      title: '原因',
      dataIndex: 'reason',
      key: 'reason',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'completed' ? 'green' : 'orange'}>
          {status === 'completed' ? '已完成' : '待处理'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Space>
          <Tooltip title="查看">
            <Button type="text" icon={<EyeOutlined />} size="small" />
          </Tooltip>
          {record.status === 'pending' && (
            <Tooltip title="审核">
              <Button type="text" icon={<EditOutlined />} size="small" style={{ color: '#52c41a' }} />
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  const warehouseColumns = [
    {
      title: '仓库信息',
      key: 'warehouse',
      render: (record: Warehouse) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Avatar 
            size={40}
            style={{ backgroundColor: '#1890ff' }}
            icon={<BoxPlotOutlined />}
          />
          <div>
            <Text strong style={{ display: 'block' }}>{record.name}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {record.code}
            </Text>
          </div>
        </div>
      ),
    },
    {
      title: '地址',
      dataIndex: 'address',
      key: 'address',
    },
    {
      title: '负责人',
      dataIndex: 'manager',
      key: 'manager',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>
          {status === 'active' ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: Warehouse) => (
        <Space>
          <Tooltip title="查看">
            <Button type="text" icon={<EyeOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="编辑">
            <Button type="text" icon={<EditOutlined />} size="small" />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const movementColumns = [
    {
      title: '变动信息',
      key: 'movement',
      render: (record: StockMovement) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Avatar 
            size={40}
            style={{ backgroundColor: record.movementType === 'IN' ? '#52c41a' : record.movementType === 'OUT' ? '#ff4d4f' : '#1890ff' }}
            icon={<SwapOutlined />}
          />
          <div>
            <Text strong style={{ display: 'block' }}>{record.itemCode}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {record.itemName}
            </Text>
          </div>
        </div>
      ),
    },
    {
      title: '变动类型',
      dataIndex: 'movementType',
      key: 'movementType',
      render: (type: string) => {
        const typeMap: { [key: string]: { text: string; color: string } } = {
          'IN': { text: '入库', color: 'green' },
          'OUT': { text: '出库', color: 'red' },
          'TRANSFER': { text: '调拨', color: 'blue' },
          'ADJUSTMENT': { text: '调整', color: 'orange' }
        };
        const config = typeMap[type] || { text: type, color: 'default' };
        return (
          <Tag color={config.color}>
            {config.text}
          </Tag>
        );
      },
    },
    {
      title: '数量变动',
      dataIndex: 'quantity',
      key: 'quantity',
      render: (quantity: number, record: StockMovement) => (
        <div>
          <Text strong style={{ color: record.movementType === 'IN' ? '#52c41a' : '#ff4d4f' }}>
            {record.movementType === 'IN' ? '+' : '-'}{Math.abs(quantity)}
          </Text>
        </div>
      ),
    },
    {
      title: '仓库',
      dataIndex: 'warehouseName',
      key: 'warehouseName',
    },
    {
      title: '单价',
      dataIndex: 'unitCost',
      key: 'unitCost',
      render: (cost: number) => `¥${cost.toFixed(2)}`,
    },
    {
      title: '总价值',
      dataIndex: 'totalValue',
      key: 'totalValue',
      render: (value: number) => `¥${value.toFixed(2)}`,
    },
    {
      title: '变动时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (date: string) => new Date(date).toLocaleString(),
    },
    {
      title: '备注',
      dataIndex: 'notes',
      key: 'notes',
      render: (notes: string) => (
        <Tooltip title={notes}>
          <Text ellipsis style={{ maxWidth: 100 }}>
            {notes || '-'}
          </Text>
        </Tooltip>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: StockMovement) => (
        <Space>
          <Tooltip title="查看详情">
            <Button type="text" icon={<EyeOutlined />} size="small" />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const handleModalOk = () => {
    form.validateFields().then(values => {
      console.log('Form values:', values);
      message.success(`${modalType === 'product' ? '商品' : modalType === 'inbound' ? '入库单' : modalType === 'outbound' ? '出库单' : '仓库'}创建成功！`);
      setIsModalVisible(false);
      form.resetFields();
    });
  };

  const showModal = (type: 'product' | 'inbound' | 'outbound' | 'warehouse' | 'movement') => {
    setModalType(type);
    setIsModalVisible(true);
  };

  // 计算统计数据
  const totalProducts = products?.length || 0;
  const lowStockProducts = stocks?.filter(s => s.quantity < s.reorderLevel)?.length || 0;
  const totalWarehouse = warehouses?.length || 0;
  const totalMovements = stockMovements?.length || 0;

  // 定义Tabs的items
  const tabItems = [
    {
      key: 'products',
      label: '商品管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索商品..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
                value={productsSearch}
                onChange={(e) => setProductsSearch(e.target.value)}
                onPressEnter={() => loadProducts()}
              />
              <Select placeholder="分类筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="电子产品">电子产品</Option>
                <Option value="办公用品">办公用品</Option>
                <Option value="原材料">原材料</Option>
                <Option value="成品">成品</Option>
              </Select>
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="正常">正常</Option>
                <Option value="缺货">缺货</Option>
                <Option value="停产">停产</Option>
              </Select>
              <Select placeholder="仓库筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="主仓库">主仓库</Option>
                <Option value="分仓库A">分仓库A</Option>
                <Option value="分仓库B">分仓库B</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ImportOutlined />}>导入</Button>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('product')}
              >
                新建商品
              </Button>
            </Space>
          </div>
          <Table 
            columns={productColumns} 
            dataSource={products || []}
            loading={loading}
            pagination={{
              current: productsPagination.current,
              pageSize: productsPagination.pageSize,
              total: productsPagination.total,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
              onChange: (page, pageSize) => {
                setProductsPagination(prev => ({
                  ...prev,
                  current: page,
                  pageSize: pageSize || 10
                }));
                loadProducts(page, productsSearch, pageSize);
              }
            }}
            scroll={{ x: 1400 }}
          />
        </>
      )
    },
    {
      key: 'warehouse',
      label: '仓库管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索仓库..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
                value={warehousesSearch}
                onChange={(e) => setWarehousesSearch(e.target.value)}
                onPressEnter={() => loadWarehouses()}
              />
              <Select placeholder="类型筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="主仓库">主仓库</Option>
                <Option value="分仓库">分仓库</Option>
                <Option value="临时仓库">临时仓库</Option>
              </Select>
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="正常">正常</Option>
                <Option value="维护中">维护中</Option>
                <Option value="停用">停用</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('warehouse')}
              >
                新建仓库
              </Button>
            </Space>
          </div>
          <Table 
            columns={warehouseColumns} 
            dataSource={warehouses || []}
            loading={loading}
            pagination={{
              current: warehousesPagination.current,
              pageSize: warehousesPagination.pageSize,
              total: warehousesPagination.total,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
              onChange: (page, pageSize) => {
                setWarehousesPagination(prev => ({
                  ...prev,
                  current: page,
                  pageSize: pageSize || 10
                }));
                loadWarehouses(page, warehousesSearch, pageSize);
              }
            }}
            scroll={{ x: 1200 }}
          />
        </>
      )
    },
    {
      key: 'movements',
      label: '库存变动',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索变动记录..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
                value={movementsSearch}
                onChange={(e) => setMovementsSearch(e.target.value)}
                onPressEnter={() => loadStockMovements()}
              />
              <Select placeholder="类型筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="入库">入库</Option>
                <Option value="出库">出库</Option>
                <Option value="调拨">调拨</Option>
                <Option value="盘点">盘点</Option>
              </Select>
              <Select placeholder="仓库筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="主仓库">主仓库</Option>
                <Option value="分仓库A">分仓库A</Option>
                <Option value="分仓库B">分仓库B</Option>
              </Select>
              <DatePicker placeholder="变动日期" style={{ width: 150 }} />
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('movement')}
              >
                新建变动
              </Button>
            </Space>
          </div>
          <Table 
            columns={movementColumns} 
            dataSource={stockMovements || []}
            loading={loading}
            pagination={{
              current: movementsPagination.current,
              pageSize: movementsPagination.pageSize,
              total: movementsPagination.total,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
              onChange: (page, pageSize) => {
                setMovementsPagination(prev => ({
                  ...prev,
                  current: page,
                  pageSize: pageSize || 10
                }));
                loadStockMovements(page, movementsSearch, pageSize);
              }
            }}
            scroll={{ x: 1500 }}
          />
        </>
      )
    }
  ];

  return (
    <div style={{ padding: '0 8px' }}>
      <div style={{ marginBottom: 24 }}>
        <Title level={2} style={{ margin: 0, color: '#1f2937' }}>
          📦 库存管理
        </Title>
        <Text type="secondary">管理商品库存、仓库和库存变动</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="商品总数"
              value={totalProducts}
              prefix={<InboxOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="库存预警"
              value={lowStockProducts}
              prefix={<ExclamationCircleOutlined style={{ color: '#ff4d4f' }} />}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="仓库数量"
              value={totalWarehouse}
              prefix={<HomeOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="库存变动"
              value={totalMovements}
              prefix={<SwapOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <Tabs defaultActiveKey="products" items={tabItems} />
      </Card>
    </div>
  );
}

export default withAuth(InventoryPage);