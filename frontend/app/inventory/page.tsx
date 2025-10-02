'use client';

import { useState, useEffect } from 'react';
import { ItemService } from '@/services/item';
import { InventoryService } from '@/services/inventory';
import { WarehouseService } from '@/services/warehouse';
import { withAuth } from '@/contexts/AuthContext';
import { 
  Item, 
  Stock, 
  StockMovement, 
  Warehouse, 
  CreateItemRequest,
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
const { Option } = Select;

function InventoryPage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'item' | 'inbound' | 'outbound' | 'warehouse' | 'movement'>('item');
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  // API数据状态
  const [items, setItems] = useState<Item[]>([]);
  const [stocks, setStocks] = useState<Stock[]>([]);
  const [stockMovements, setStockMovements] = useState<StockMovement[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);

  // 分页状态
  const [itemsPagination, setItemsPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
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
  const [itemsSearch, setItemsSearch] = useState('');
  const [stocksSearch, setStocksSearch] = useState('');
  const [movementsSearch, setMovementsSearch] = useState('');
  const [warehousesSearch, setWarehousesSearch] = useState('');

  // 数据加载函数
  const loadItems = async (page = 1, search = '', pageSize?: number) => {
    try {
      setLoading(true);
      const pageSizeValue = pageSize || itemsPagination.pageSize;
      const response = await ItemService.getItems({
        page,
        pageSize: pageSizeValue,
        search: search || undefined
      });
      setItems(response.data || []);
      setItemsPagination({
        current: response.page,
        pageSize: response.limit,
        total: response.total
      });
    } catch (error) {
      console.error('加载产品数据失败:', error);
      message.error('加载产品数据失败');
      setItems([]);
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
    loadItems();
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

  const itemColumns = [
    {
      title: '物料信息',
      key: 'item',
      render: (record: Item) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Avatar 
            size={40}
            style={{ backgroundColor: '#1890ff' }}
            icon={<InboxOutlined />}
          />
          <div>
            <Text strong style={{ display: 'block' }}>{record.name}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {record.code} | {typeof record.category === 'object' && record.category !== null ? 
                ((record.category as any).name || (record.category as any).code || (record.category as any).id || 'Unknown') : 
                (record.category || '-')}
            </Text>
          </div>
        </div>
      ),
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      render: (category: any) => {
        // Fix: Properly handle category object rendering
        if (typeof category === 'object' && category !== null) {
          return <Tag color="blue">{category.name || category.code || category.id || 'Unknown'}</Tag>;
        }
        return <Tag color="blue">{category || '-'}</Tag>;
      },
    },
    {
      title: '单位',
      dataIndex: 'unit',
      key: 'unit',
      render: (unit: any) => {
        // Fix: Properly handle unit object rendering
        if (typeof unit === 'object' && unit !== null) {
          return <Text>{unit.name || unit.symbol || unit.id || '-'}</Text>;
        }
        return <Text>{unit || '-'}</Text>;
      },
    },
    {
      title: '价格信息',
      key: 'price',
      render: (record: Item) => (
        <div>
          <div style={{ marginBottom: 4 }}>
            <Text type="secondary">售价: </Text>
            <Text strong style={{ color: '#52c41a' }}>
              ¥{record.price ? record.price.toLocaleString() : '0.00'}
            </Text>
          </div>
          <div>
            <Text type="secondary">成本: </Text>
            <Text>¥{record.cost ? record.cost.toLocaleString() : '0.00'}</Text>
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
      render: (record: Item) => (
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
      key: 'product',
      render: (record: any) => {
        // 支持后端返回 item 或 product 对象，统一展示名称
        // Fix: Properly handle item/product object rendering
        if (record?.item && typeof record.item === 'object' && record.item !== null) {
          return record.item.name || record.item.code || record.item.id || '-';
        }
        if (record?.product && typeof record.product === 'object' && record.product !== null) {
          return record.product.name || record.product.code || record.product.id || '-';
        }
        return record?.item?.name || record?.product?.name || '-';
      },
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
      key: 'warehouse',
      render: (record: any) => {
        // 后端返回 warehouse 为对象时，避免直接渲染对象
        // Fix: Properly handle warehouse object rendering
        if (record?.warehouse && typeof record.warehouse === 'object' && record.warehouse !== null) {
          return record.warehouse.name || record.warehouse.code || record.warehouse.id || '-';
        }
        return record?.warehouse?.name || '-';
      },
    },
    {
      title: '操作人',
      dataIndex: 'operator',
      key: 'operator',
      render: (operator: any) => {
        // Fix: Properly handle operator object rendering
        if (typeof operator === 'object' && operator !== null) {
          if (operator.name) return operator.name;
          if (operator.firstName || operator.lastName) return `${operator.firstName || ''} ${operator.lastName || ''}`.trim();
          return operator.id || '-';
        }
        return operator || '-';
      },
    },
    {
      title: '操作时间',
      dataIndex: 'time',
      key: 'time',
      render: (time: any) => {
        // Fix: Properly handle time object rendering
        if (typeof time === 'object' && time !== null) {
          // 如果是对象，尝试获取时间字符串
          const timeStr = time?.created_at || time?.updated_at;
          if (timeStr) {
            try {
              return new Date(timeStr).toLocaleString();
            } catch {
              return timeStr;
            }
          }
          return '-';
        }
        if (time) {
          try {
            return new Date(time).toLocaleString();
          } catch {
            return time;
          }
        }
        return '-';
      },
    },
    {
      title: '原因',
      dataIndex: 'reason',
      key: 'reason',
      render: (reason: any) => {
        // Fix: Properly handle reason object rendering
        if (typeof reason === 'object' && reason !== null) {
          return reason.name || reason.id || '-';
        }
        return reason || '-';
      },
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
      render: (address: any) => {
        // Fix: Properly handle address object rendering
        if (typeof address === 'object' && address !== null) {
          // Check if it's a nested object with name property
          if (address.name) return <Text>{address.name}</Text>;
          // Check if it's a nested object with other properties
          if (address.id) return <Text>{address.id}</Text>;
          // Fallback to dash if no valid property found
          return <Text>-</Text>;
        }
        return <Text>{address || '-'}</Text>;
      },
    },
    {
      title: '负责人',
      dataIndex: 'manager',
      key: 'manager',
      render: (manager: any) => {
        // Fix: Properly handle manager object rendering
        if (typeof manager === 'object' && manager !== null) {
          // Check if it's a nested object with name property
          if (manager.name) return <Text>{manager.name}</Text>;
          // Check if it's a nested object with other properties
          if (manager.id) return <Text>{manager.id}</Text>;
          // Fallback to dash if no valid property found
          return <Text>-</Text>;
        }
        return <Text>{manager || '-'}</Text>;
      },
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
      render: (date: string) => {
        if (!date) return '-';
        try {
          return new Date(date).toLocaleDateString();
        } catch (error) {
          return date;
        }
      },
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
      render: (record: any) => {
        // 正确处理后端返回的嵌套对象结构
        const itemCode = record.item?.code || '-';
        const itemName = record.item?.name || '-';
        
        return (
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <Avatar 
              size={40}
              style={{ backgroundColor: record.type === 'in' ? '#52c41a' : record.type === 'out' ? '#ff4d4f' : '#1890ff' }}
              icon={<SwapOutlined />}
            />
            <div>
              <Text strong style={{ display: 'block' }}>{itemCode}</Text>
              <Text type="secondary" style={{ fontSize: 12 }}>
                {itemName}
              </Text>
            </div>
          </div>
        );
      },
    },
    {
      title: '变动类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => {
        const typeMap: { [key: string]: { text: string; color: string } } = {
          'in': { text: '入库', color: 'green' },
          'out': { text: '出库', color: 'red' },
          'transfer': { text: '调拨', color: 'blue' },
          'adjustment': { text: '调整', color: 'orange' }
        };
        const config = typeMap[type] || { text: type || '-', color: 'default' };
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
      render: (quantity: number, record: any) => {
        const numQuantity = Number(quantity) || 0;
        const movementType = record.type;
        
        return (
          <div>
            <Text strong style={{ color: movementType === 'in' ? '#52c41a' : '#ff4d4f' }}>
              {movementType === 'in' ? '+' : '-'}{Math.abs(numQuantity)}
            </Text>
          </div>
        );
      },
    },
    {
      title: '仓库',
      key: 'warehouse',
      render: (record: any) => {
        // Fix: Properly handle warehouse object rendering
        if (record.warehouse && typeof record.warehouse === 'object' && record.warehouse !== null) {
          return record.warehouse.name || record.warehouse.code || record.warehouse.id || '-';
        }
        return record.warehouse || '-';
      },
    },
    {
      title: '库位',
      key: 'location',
      render: (record: any) => {
        // Fix: Properly handle location object rendering
        if (record.location && typeof record.location === 'object' && record.location !== null) {
          return record.location.name || record.location.code || record.location.id || '-';
        }
        return record.location || '-';
      },
    },
    {
      title: '变动时间',
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
      title: '备注',
      dataIndex: 'notes',
      key: 'notes',
      render: (notes: string) => {
        return (
          <Tooltip title={notes || '无备注'}>
            <Text ellipsis style={{ maxWidth: 100 }}>
              {notes || '-'}
            </Text>
          </Tooltip>
        );
      },
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
      message.success(`${modalType === 'item' ? '物料' : modalType === 'inbound' ? '入库单' : modalType === 'outbound' ? '出库单' : '仓库'}创建成功！`);
      setIsModalVisible(false);
      form.resetFields();
    });
  };

  const showModal = (type: 'item' | 'inbound' | 'outbound' | 'warehouse' | 'movement') => {
    setModalType(type);
    setIsModalVisible(true);
  };

  // 计算统计数据
  const totalItems = items?.length || 0;
  const lowStockItems = stocks?.filter(s => s.quantity < s.reorderLevel)?.length || 0;
  const totalWarehouse = warehouses?.length || 0;
  const totalMovements = stockMovements?.length || 0;

  // 定义Tabs的items
  const tabItems = [
    {
      key: 'stocks',
      label: '库存记录',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索库存记录..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
                value={stocksSearch}
                onChange={(e) => setStocksSearch(e.target.value)}
                onPressEnter={() => loadStocks()}
              />
              <Select placeholder="类型筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="inbound">入库</Option>
                <Option value="outbound">出库</Option>
              </Select>
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="completed">已完成</Option>
                <Option value="pending">待处理</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('inbound')}
              >
                新建入库
              </Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('outbound')}
              >
                新建出库
              </Button>
            </Space>
          </div>
          <Table 
            columns={stockColumns} 
            dataSource={stocks || []}
            rowKey="id"
            loading={loading}
            pagination={{
              current: stocksPagination.current,
              pageSize: stocksPagination.pageSize,
              total: stocksPagination.total,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
              onChange: (page, pageSize) => {
                setStocksPagination(prev => ({
                  ...prev,
                  current: page,
                  pageSize: pageSize || 10
                }));
                loadStocks(page, stocksSearch, pageSize);
              }
            }}
            scroll={{ x: 1400 }}
          />
        </>
      )
    },
    {
      key: 'items',
      label: '物料管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索物料..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
                value={itemsSearch}
                onChange={(e) => setItemsSearch(e.target.value)}
                onPressEnter={() => loadItems()}
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
                onClick={() => showModal('item')}
              >
                  新建物料
                </Button>
            </Space>
          </div>
          <Table 
             columns={itemColumns} 
            dataSource={items || []}
            rowKey="id"
            loading={loading}
            pagination={{
              current: itemsPagination.current,
              pageSize: itemsPagination.pageSize,
              total: itemsPagination.total,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
              onChange: (page, pageSize) => {
                setItemsPagination(prev => ({
                  ...prev,
                  current: page,
                  pageSize: pageSize || 10
                }));
                loadItems(page, itemsSearch, pageSize);
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
            rowKey="id"
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
            rowKey="id"
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
               title="物料总数"
              value={totalItems}
              prefix={<InboxOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="库存预警"
              value={lowStockItems}
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
        <Tabs defaultActiveKey="items" items={tabItems} />
      </Card>
    </div>
  );
}

export default withAuth(InventoryPage);