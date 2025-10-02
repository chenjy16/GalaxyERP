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
  Divider,
  Upload
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
  const loadItems = async (page = 1, search = '', pageSize?: number, category?: string, status?: string) => {
    try {
      setLoading(true);
      const pageSizeValue = pageSize || itemsPagination.pageSize;
      const response = await ItemService.getItems({
        page,
        pageSize: pageSizeValue,
        search: search || undefined,
        category: category && category !== 'all' ? category : undefined,
        status: status && status !== 'all' ? status : undefined
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
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'green' : 'red'}>
          {isActive ? '启用' : '禁用'}
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
            onClick: ({ key }) => {
              if (key === 'view') {
                handleView(record);
              } else if (key === 'edit') {
                handleEdit(record);
              } else if (key === 'inbound' || key === 'outbound') {
                handleStockOperation(record, key as 'inbound' | 'outbound');
              } else if (key === 'delete') {
                handleDelete(record, '物料');
              }
            },
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
      title: '库存ID',
      dataIndex: 'id',
      key: 'id',
      render: (id: string) => <Text strong>{id}</Text>,
    },
    {
      title: '物料信息',
      key: 'item',
      render: (record: Stock) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Avatar 
            size={40}
            style={{ backgroundColor: '#1890ff' }}
            icon={<InboxOutlined />}
          />
          <div>
            <Text strong style={{ display: 'block' }}>{record.item?.name || '未知物料'}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {record.item?.code || 'N/A'}
            </Text>
          </div>
        </div>
      ),
    },
    {
      title: '仓库',
      key: 'warehouse',
      render: (record: Stock) => (
        <div>
          <Text strong>{record.warehouse?.name || '未知仓库'}</Text>
          <br />
          <Text type="secondary" style={{ fontSize: 12 }}>{record.warehouse?.code || 'N/A'}</Text>
        </div>
      ),
    },
    {
      title: '库存数量',
      dataIndex: 'quantity',
      key: 'quantity',
      render: (quantity: number) => (
        <Text strong style={{ fontSize: 16 }}>
          {Number(quantity || 0).toLocaleString()}
        </Text>
      ),
    },
    {
      title: '可用数量',
      dataIndex: 'available_qty',
      key: 'available_qty',
      render: (availableQuantity: number) => (
        <Text style={{ color: '#52c41a' }}>
          {Number(availableQuantity || 0).toLocaleString()}
        </Text>
      ),
    },
    {
      title: '预留数量',
      dataIndex: 'reserved_qty',
      key: 'reserved_qty',
      render: (reservedQuantity: number) => (
        <Text style={{ color: '#ff4d4f' }}>
          {Number(reservedQuantity || 0).toLocaleString()}
        </Text>
      ),
    },
    {
      title: '单位成本',
      key: 'unitCost',
      render: (record: Stock) => (
        <Text>¥{(record.item?.unit_cost || 0).toFixed(2)}</Text>
      ),
    },
    {
      title: '总价值',
      key: 'totalValue',
      render: (record: Stock) => {
        const totalValue = (record.quantity || 0) * (record.item?.unit_cost || 0);
        return (
          <Text strong style={{ color: '#1890ff' }}>¥{totalValue.toFixed(2)}</Text>
        );
      },
    },
    {
      title: '库存状态',
      key: 'stockStatus',
      render: (record: Stock) => {
        const status = getStockStatus(record.quantity, record.item?.min_stock || 0, record.item?.max_stock || 1000);
         return (
           <Tag color={getStockStatusColor(status)}>
             {getStockStatusText(status)}
           </Tag>
         );
       },
     },
     {
       title: '最后更新',
       dataIndex: 'updated_at',
       key: 'updated_at',
       render: (lastUpdated: string) => {
         if (!lastUpdated) return '-';
         try {
           return new Date(lastUpdated).toLocaleString();
         } catch {
           return lastUpdated;
         }
       },
     },
     {
       title: '操作',
       key: 'action',
       render: (record: Stock) => (
         <Dropdown
           menu={{
             items: [
               {
                 key: 'view',
                 label: '查看详情',
                 icon: <EyeOutlined />,
                 onClick: () => handleView(record),
               },
               {
                 key: 'inbound',
                 label: '入库',
                 icon: <InboxOutlined />,
                 onClick: () => {
                   const stockItem: Item = {
                     id: record.item?.id || 0,
                     code: record.item?.code || '',
                     name: record.item?.name || '',
                     description: record.item?.description || '',
                     category: record.item?.category?.name || '',
                     unit: record.item?.unit?.symbol || '',
                     cost: record.item?.unit_cost || 0,
                     price: record.item?.sale_price || 0,
                     reorderLevel: record.item?.min_stock || 0,
                     isActive: record.item?.is_active || true,
                     createdAt: record.item?.created_at || '',
                     updatedAt: record.updated_at || ''
                   };
                   handleStockOperation(stockItem, 'inbound');
                 },
               },
               {
                 key: 'outbound',
                 label: '出库',
                 icon: <SendOutlined />,
                 onClick: () => {
                   const stockItem: Item = {
                     id: record.item?.id || 0,
                     code: record.item?.code || '',
                     name: record.item?.name || '',
                     description: record.item?.description || '',
                     category: record.item?.category?.name || '',
                     unit: record.item?.unit?.symbol || '',
                     cost: record.item?.unit_cost || 0,
                     price: record.item?.sale_price || 0,
                     reorderLevel: record.item?.min_stock || 0,
                     isActive: record.item?.is_active || true,
                     createdAt: record.item?.created_at || '',
                     updatedAt: record.updated_at || ''
                   };
                   handleStockOperation(stockItem, 'outbound');
                 },
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
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'green' : 'red'}>
          {isActive ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
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
        <Dropdown
          menu={{
            items: [
              {
                key: 'view',
                label: '查看详情',
                icon: <EyeOutlined />,
                onClick: () => handleView(record),
              },
              {
                key: 'edit',
                label: '编辑',
                icon: <EditOutlined />,
                onClick: () => showModal('warehouse', record),
              },
              {
                key: 'delete',
                label: '删除',
                icon: <DeleteOutlined />,
                onClick: () => handleDelete(record, 'warehouse'),
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
        <Button 
          type="text" 
          icon={<EyeOutlined />} 
          size="small"
          onClick={() => handleView(record)}
        >
          查看详情
        </Button>
      ),
    },
  ];

  // 编辑状态
  const [editingRecord, setEditingRecord] = useState<any>(null);
  const [viewingRecord, setViewingRecord] = useState<any>(null);
  const [isViewModalVisible, setIsViewModalVisible] = useState(false);
  const [isImportModalVisible, setIsImportModalVisible] = useState(false);
  const [importFile, setImportFile] = useState<File | null>(null);
  const [importLoading, setImportLoading] = useState(false);

  // 筛选状态
  const [itemCategoryFilter, setItemCategoryFilter] = useState('all');
  const [itemStatusFilter, setItemStatusFilter] = useState('all');
  const [stockTypeFilter, setStockTypeFilter] = useState('all');
  const [stockStatusFilter, setStockStatusFilter] = useState('all');
  const [warehouseTypeFilter, setWarehouseTypeFilter] = useState('all');
  const [warehouseStatusFilter, setWarehouseStatusFilter] = useState('all');
  const [movementTypeFilter, setMovementTypeFilter] = useState('all');

  // 处理模态框确认
  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      if (modalType === 'item') {
        if (editingRecord) {
          // 编辑物料
          await ItemService.updateItem(editingRecord.id, values);
          message.success('物料更新成功！');
        } else {
          // 创建物料
          await ItemService.createItem(values);
          message.success('物料创建成功！');
        }
        loadItems();
      } else if (modalType === 'warehouse') {
        if (editingRecord) {
          // 编辑仓库
          await WarehouseService.updateWarehouse(editingRecord.id, values);
          message.success('仓库更新成功！');
        } else {
          // 创建仓库
          await WarehouseService.createWarehouse(values);
          message.success('仓库创建成功！');
        }
        loadWarehouses();
      } else if (modalType === 'inbound' || modalType === 'outbound') {
         // 创建库存记录
         if (modalType === 'inbound') {
           await InventoryService.stockIn({
             itemId: values.itemId,
             warehouseId: values.warehouseId,
             quantity: values.quantity,
             reason: values.reason || '入库操作',
             reference: values.reference
           });
         } else {
           await InventoryService.stockOut({
             itemId: values.itemId,
             warehouseId: values.warehouseId,
             quantity: values.quantity,
             reason: values.reason || '出库操作',
             reference: values.reference
           });
         }
         message.success(`${modalType === 'inbound' ? '入库' : '出库'}单创建成功！`);
         loadStocks();
      } else if (modalType === 'movement') {
        // 创建库存变动
        await InventoryService.createStockMovement(values);
        message.success('库存变动创建成功！');
        loadStockMovements();
      }

      setIsModalVisible(false);
      form.resetFields();
      setEditingRecord(null);
    } catch (error) {
      console.error('操作失败:', error);
      message.error('操作失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const showModal = (type: 'item' | 'inbound' | 'outbound' | 'warehouse' | 'movement', record?: any) => {
    setModalType(type);
    setEditingRecord(record);
    if (record) {
      form.setFieldsValue(record);
    }
    setIsModalVisible(true);
  };

  // 处理查看详情
  const handleView = (record: any) => {
    setViewingRecord(record);
    setIsViewModalVisible(true);
  };

  // 处理删除
  const handleDelete = async (record: any, type: string) => {
    Modal.confirm({
      title: `确认删除${type}？`,
      content: `确定要删除${type}"${record.name || record.code}"吗？此操作不可撤销。`,
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        try {
          setLoading(true);
          if (type === '物料') {
            await ItemService.deleteItem(record.id);
            message.success('物料删除成功！');
            loadItems();
          } else if (type === '仓库') {
            await WarehouseService.deleteWarehouse(record.id);
            message.success('仓库删除成功！');
            loadWarehouses();
          }
        } catch (error) {
          console.error('删除失败:', error);
          message.error('删除失败，请重试');
        } finally {
          setLoading(false);
        }
      }
    });
  };

  // 处理编辑操作
  const handleEdit = (record: Item) => {
    form.setFieldsValue({
      name: record.name,
      code: record.code,
      category: record.category,
      unit: record.unit,
      unitPrice: record.price,
      description: record.description
    });
    setEditingRecord(record);
    showModal('item');
  };

  // 处理入库/出库操作
  const handleStockOperation = (record: Item, type: 'inbound' | 'outbound') => {
    form.setFieldsValue({
      itemId: record.id,
      itemName: record.name,
      type: type
    });
    showModal(type);
  };

  // 导入功能
  const handleImportItems = () => {
    setIsImportModalVisible(true);
  };

  const handleImportConfirm = async () => {
    if (!importFile) {
      message.error('请选择要导入的文件');
      return;
    }

    setImportLoading(true);
    try {
      const result = await ItemService.importItems(importFile);
      if (result.success) {
        message.success(`物料导入成功！共导入 ${result.imported} 条记录`);
        if (result.errors && result.errors.length > 0) {
          message.warning(`有 ${result.errors.length} 条记录导入失败，请检查数据格式`);
        }
      } else {
        message.error('导入失败，请检查文件格式');
      }
      setIsImportModalVisible(false);
      setImportFile(null);
      loadItems(); // 重新加载物料列表
    } catch (error) {
      console.error('导入失败:', error);
      message.error('导入失败，请检查文件格式');
    } finally {
      setImportLoading(false);
    }
  };

  const handleFileChange = (file: File) => {
    // 检查文件类型
    const allowedTypes = [
      'application/vnd.ms-excel',
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      'text/csv'
    ];
    
    if (!allowedTypes.includes(file.type)) {
      message.error('只支持 Excel (.xlsx, .xls) 和 CSV (.csv) 文件');
      return false;
    }

    // 检查文件大小 (10MB)
    if (file.size > 10 * 1024 * 1024) {
      message.error('文件大小不能超过 10MB');
      return false;
    }

    setImportFile(file);
    return false; // 阻止自动上传
  };

  // 导出功能
  const handleExportItems = () => {
    try {
      const csvContent = [
        ['物料编码', '物料名称', '描述', '分类', '单位', '成本价', '售价', '安全库存', '状态', '创建时间'],
        ...items.map(item => [
          item.code,
          item.name,
          item.description || '',
          typeof item.category === 'object' ? (item.category as any)?.name || '' : item.category || '',
          typeof item.unit === 'object' ? (item.unit as any)?.name || '' : item.unit || '',
          item.cost || 0,
          item.price || 0,
          item.reorderLevel || 0,
          item.isActive ? '启用' : '禁用',
          new Date(item.createdAt).toLocaleDateString()
        ])
      ];

      const csvString = csvContent.map(row => row.join(',')).join('\n');
      const blob = new Blob(['\uFEFF' + csvString], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = `物料数据_${new Date().toISOString().split('T')[0]}.csv`;
      link.click();
      message.success('物料数据导出成功！');
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败，请重试');
    }
  };

  const handleExportWarehouses = () => {
    try {
      const csvContent = [
        ['仓库编码', '仓库名称', '地址', '负责人', '状态', '创建时间'],
        ...warehouses.map(warehouse => [
          warehouse.code,
          warehouse.name,
          warehouse.address || '',
          warehouse.manager ? (typeof warehouse.manager === 'object' ? warehouse.manager.firstName + ' ' + warehouse.manager.lastName : warehouse.manager) : '',
          warehouse.is_active ? '启用' : '禁用',
          new Date(warehouse.created_at).toLocaleDateString()
        ])
      ];

      const csvString = csvContent.map(row => row.join(',')).join('\n');
      const blob = new Blob(['\uFEFF' + csvString], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = `仓库数据_${new Date().toISOString().split('T')[0]}.csv`;
      link.click();
      message.success('仓库数据导出成功！');
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败，请重试');
    }
  };

  const handleExportStocks = () => {
    try {
      const csvContent = [
        ['库存ID', '物料编码', '物料名称', '仓库编码', '仓库名称', '库存数量', '可用数量', '预留数量', '单位成本', '总价值', '最后更新时间'],
        ...stocks.map(stock => [
          stock.id,
          stock.item?.code || '',
          stock.item?.name || '',
          stock.warehouse?.code || '',
          stock.warehouse?.name || '',
          stock.quantity || 0,
          stock.available_qty || 0,
          stock.reserved_qty || 0,
          stock.item?.unit_cost || 0,
          (stock.quantity || 0) * (stock.item?.unit_cost || 0),
          stock.updated_at ? new Date(stock.updated_at).toLocaleString() : ''
        ])
      ];

      const csvString = csvContent.map(row => row.join(',')).join('\n');
      const blob = new Blob(['\uFEFF' + csvString], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = `库存记录_${new Date().toISOString().split('T')[0]}.csv`;
      link.click();
      message.success('库存记录导出成功！');
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败，请重试');
    }
  };

  const handleExportMovements = () => {
    try {
      const csvContent = [
        ['变动编码', '物料名称', '变动类型', '数量变动', '仓库', '库位', '变动时间', '备注'],
        ...stockMovements.map(movement => [
          movement.item?.code || '',
          movement.item?.name || '',
          movement.type === 'in' ? '入库' : movement.type === 'out' ? '出库' : movement.type === 'transfer' ? '调拨' : '调整',
          (movement.type === 'in' ? '+' : '-') + Math.abs(movement.quantity || 0),
          movement.warehouse?.name || '',
          movement.location?.name || '',
          movement.created_at ? new Date(movement.created_at).toLocaleString() : '',
          movement.notes || ''
        ])
      ];

      const csvString = csvContent.map(row => row.join(',')).join('\n');
      const blob = new Blob(['\uFEFF' + csvString], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = `库存变动_${new Date().toISOString().split('T')[0]}.csv`;
      link.click();
      message.success('库存变动导出成功！');
    } catch (error) {
      console.error('导出失败:', error);
      message.error('导出失败，请重试');
    }
  };

  // 计算统计数据
  const totalItems = items?.length || 0;
  const lowStockItems = stocks?.filter(s => s.quantity < (s.item?.min_stock || 0))?.length || 0;
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
              <Button icon={<ExportOutlined />} onClick={handleExportStocks}>导出</Button>
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
                onPressEnter={() => loadItems(1, itemsSearch, 10, itemCategoryFilter, itemStatusFilter)}
              />
              <Select 
                placeholder="分类筛选" 
                style={{ width: 120 }}
                value={itemCategoryFilter}
                onChange={(value) => {
                  setItemCategoryFilter(value);
                  loadItems(1, itemsSearch, 10, value, itemStatusFilter);
                }}
              >
                <Option value="all">全部</Option>
                <Option value="原材料">原材料</Option>
                <Option value="半成品">半成品</Option>
                <Option value="成品">成品</Option>
                <Option value="包装材料">包装材料</Option>
                <Option value="辅助材料">辅助材料</Option>
              </Select>
              <Select 
                placeholder="状态筛选" 
                style={{ width: 120 }}
                value={itemStatusFilter}
                onChange={(value) => {
                  setItemStatusFilter(value);
                  loadItems(1, itemsSearch, 10, itemCategoryFilter, value);
                }}
              >
                <Option value="all">全部</Option>
                <Option value="active">启用</Option>
                <Option value="inactive">禁用</Option>
              </Select>
              <Button 
                onClick={() => {
                  setItemsSearch('');
                  setItemCategoryFilter('all');
                  setItemStatusFilter('all');
                  loadItems(1, '', 10, 'all', 'all');
                }}
              >
                重置
              </Button>
            </Space>
            <Space>
              <Button icon={<ImportOutlined />} onClick={handleImportItems}>导入</Button>
              <Button icon={<ExportOutlined />} onClick={handleExportItems}>导出</Button>
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
                loadItems(page, itemsSearch, pageSize, itemCategoryFilter, itemStatusFilter);
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
              <Button icon={<ExportOutlined />} onClick={handleExportWarehouses}>导出</Button>
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
              <Button icon={<ExportOutlined />} onClick={handleExportMovements}>导出</Button>
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

      {/* 创建/编辑模态框 */}
      <Modal
        title={
          modalType === 'item' ? (editingRecord ? '编辑物料' : '新建物料') :
          modalType === 'warehouse' ? (editingRecord ? '编辑仓库' : '新建仓库') :
          modalType === 'movement' ? '新建库存变动' :
          modalType === 'inbound' ? '入库操作' :
          modalType === 'outbound' ? '出库操作' : '操作'
        }
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={() => {
          setIsModalVisible(false);
          setEditingRecord(null);
          form.resetFields();
        }}
        confirmLoading={loading}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={editingRecord}
        >
          {modalType === 'item' && (
            <>
              <Form.Item
                name="code"
                label="物料编码"
                rules={[{ required: true, message: '请输入物料编码' }]}
              >
                <Input placeholder="请输入物料编码" />
              </Form.Item>
              <Form.Item
                name="name"
                label="物料名称"
                rules={[{ required: true, message: '请输入物料名称' }]}
              >
                <Input placeholder="请输入物料名称" />
              </Form.Item>
              <Form.Item
                name="category"
                label="物料分类"
                rules={[{ required: true, message: '请选择物料分类' }]}
              >
                <Select placeholder="请选择物料分类">
                  <Option value="原材料">原材料</Option>
                  <Option value="半成品">半成品</Option>
                  <Option value="成品">成品</Option>
                  <Option value="包装材料">包装材料</Option>
                  <Option value="辅助材料">辅助材料</Option>
                </Select>
              </Form.Item>
              <Form.Item
                name="specification"
                label="规格型号"
              >
                <Input placeholder="请输入规格型号" />
              </Form.Item>
              <Form.Item
                name="unit"
                label="计量单位"
                rules={[{ required: true, message: '请输入计量单位' }]}
              >
                <Input placeholder="请输入计量单位" />
              </Form.Item>
              <Form.Item
                name="unitPrice"
                label="单价"
              >
                <Input type="number" placeholder="请输入单价" />
              </Form.Item>
              <Form.Item
                name="description"
                label="描述"
              >
                <Input.TextArea rows={3} placeholder="请输入物料描述" />
              </Form.Item>
            </>
          )}

          {modalType === 'warehouse' && (
            <>
              <Form.Item
                name="code"
                label="仓库编码"
                rules={[{ required: true, message: '请输入仓库编码' }]}
              >
                <Input placeholder="请输入仓库编码" />
              </Form.Item>
              <Form.Item
                name="name"
                label="仓库名称"
                rules={[{ required: true, message: '请输入仓库名称' }]}
              >
                <Input placeholder="请输入仓库名称" />
              </Form.Item>
              <Form.Item
                name="type"
                label="仓库类型"
                rules={[{ required: true, message: '请选择仓库类型' }]}
              >
                <Select placeholder="请选择仓库类型">
                  <Option value="主仓库">主仓库</Option>
                  <Option value="分仓库">分仓库</Option>
                  <Option value="临时仓库">临时仓库</Option>
                </Select>
              </Form.Item>
              <Form.Item
                name="address"
                label="仓库地址"
              >
                <Input placeholder="请输入仓库地址" />
              </Form.Item>
              <Form.Item
                name="manager"
                label="负责人"
              >
                <Input placeholder="请输入负责人" />
              </Form.Item>
              <Form.Item
                name="phone"
                label="联系电话"
              >
                <Input placeholder="请输入联系电话" />
              </Form.Item>
              <Form.Item
                name="description"
                label="描述"
              >
                <Input.TextArea rows={3} placeholder="请输入仓库描述" />
              </Form.Item>
            </>
          )}

          {(modalType === 'inbound' || modalType === 'outbound') && (
            <>
              <Form.Item
                name="itemId"
                label="物料"
                rules={[{ required: true, message: '请选择物料' }]}
              >
                <Select placeholder="请选择物料" showSearch>
                  {items?.map(item => (
                    <Option key={item.id} value={item.id}>
                      {item.code} - {item.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
              <Form.Item
                name="warehouseId"
                label="仓库"
                rules={[{ required: true, message: '请选择仓库' }]}
              >
                <Select placeholder="请选择仓库">
                  {warehouses?.map(warehouse => (
                    <Option key={warehouse.id} value={warehouse.id}>
                      {warehouse.code} - {warehouse.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
              <Form.Item
                name="quantity"
                label="数量"
                rules={[{ required: true, message: '请输入数量' }]}
              >
                <Input type="number" placeholder="请输入数量" />
              </Form.Item>
              <Form.Item
                name="unitCost"
                label="单位成本"
              >
                <Input type="number" placeholder="请输入单位成本" />
              </Form.Item>
              <Form.Item
                name="reason"
                label="操作原因"
              >
                <Input placeholder="请输入操作原因" />
              </Form.Item>
              <Form.Item
                name="notes"
                label="备注"
              >
                <Input.TextArea rows={3} placeholder="请输入备注" />
              </Form.Item>
            </>
          )}

          {modalType === 'movement' && (
            <>
              <Form.Item
                name="type"
                label="变动类型"
                rules={[{ required: true, message: '请选择变动类型' }]}
              >
                <Select placeholder="请选择变动类型">
                  <Option value="入库">入库</Option>
                  <Option value="出库">出库</Option>
                  <Option value="调拨">调拨</Option>
                  <Option value="盘点">盘点</Option>
                </Select>
              </Form.Item>
              <Form.Item
                name="itemId"
                label="物料"
                rules={[{ required: true, message: '请选择物料' }]}
              >
                <Select placeholder="请选择物料" showSearch>
                  {items?.map(item => (
                    <Option key={item.id} value={item.id}>
                      {item.code} - {item.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
              <Form.Item
                name="warehouseId"
                label="仓库"
                rules={[{ required: true, message: '请选择仓库' }]}
              >
                <Select placeholder="请选择仓库">
                  {warehouses?.map(warehouse => (
                    <Option key={warehouse.id} value={warehouse.id}>
                      {warehouse.code} - {warehouse.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
              <Form.Item
                name="quantity"
                label="变动数量"
                rules={[{ required: true, message: '请输入变动数量' }]}
              >
                <Input type="number" placeholder="请输入变动数量" />
              </Form.Item>
              <Form.Item
                name="reason"
                label="变动原因"
              >
                <Input placeholder="请输入变动原因" />
              </Form.Item>
              <Form.Item
                name="notes"
                label="备注"
              >
                <Input.TextArea rows={3} placeholder="请输入备注" />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>

      {/* 查看详情模态框 */}
      <Modal
        title={
          viewingRecord?.type ? "库存变动详情" : 
          viewingRecord?.address ? "仓库详情" : 
          viewingRecord?.quantity !== undefined ? "库存记录详情" : 
          "物料详情"
        }
        open={isViewModalVisible}
        onCancel={() => {
          setIsViewModalVisible(false);
          setViewingRecord(null);
        }}
        footer={[
          <Button key="close" onClick={() => setIsViewModalVisible(false)}>
            关闭
          </Button>
        ]}
        width={700}
      >
        {viewingRecord && (
          <div style={{ padding: '16px 0' }}>
            {/* 库存变动详情 */}
            {viewingRecord.type ? (
              <>
                <Row gutter={[16, 16]}>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#1890ff' }}>变动信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>变动类型: </Text>
                        <Tag color={
                          viewingRecord.type === 'in' ? 'green' : 
                          viewingRecord.type === 'out' ? 'red' : 
                          viewingRecord.type === 'transfer' ? 'blue' : 'orange'
                        }>
                          {viewingRecord.type === 'in' ? '入库' : 
                           viewingRecord.type === 'out' ? '出库' : 
                           viewingRecord.type === 'transfer' ? '调拨' : '调整'}
                        </Tag>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>变动数量: </Text>
                        <Text style={{ 
                          color: viewingRecord.type === 'in' ? '#52c41a' : '#ff4d4f',
                          fontSize: '16px',
                          fontWeight: 'bold'
                        }}>
                          {viewingRecord.type === 'in' ? '+' : '-'}{Math.abs(Number(viewingRecord.quantity) || 0)}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>变动时间: </Text>
                        <Text>
                          {viewingRecord.created_at ? 
                            new Date(viewingRecord.created_at).toLocaleString() : '-'}
                        </Text>
                      </div>
                    </div>
                  </Col>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#52c41a' }}>物料信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>物料编码: </Text>
                        <Text>{viewingRecord.item?.code || '-'}</Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>物料名称: </Text>
                        <Text>{viewingRecord.item?.name || '-'}</Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>物料分类: </Text>
                        <Text>
                          {typeof viewingRecord.item?.category === 'object' && viewingRecord.item?.category !== null ? 
                            (viewingRecord.item.category.name || viewingRecord.item.category.code || viewingRecord.item.category.id || '-') : 
                            (viewingRecord.item?.category || '-')}
                        </Text>
                      </div>
                    </div>
                  </Col>
                </Row>
                
                <div style={{ marginTop: 16 }}>
                  <Text strong style={{ color: '#722ed1' }}>仓储信息</Text>
                  <Divider style={{ margin: '8px 0' }} />
                  <Row gutter={[16, 8]}>
                    <Col span={12}>
                      <Text strong>仓库: </Text>
                      <Text>
                        {typeof viewingRecord.warehouse === 'object' && viewingRecord.warehouse !== null ? 
                          (viewingRecord.warehouse.name || viewingRecord.warehouse.code || viewingRecord.warehouse.id || '-') : 
                          (viewingRecord.warehouse || '-')}
                      </Text>
                    </Col>
                    <Col span={12}>
                      <Text strong>库位: </Text>
                      <Text>
                        {typeof viewingRecord.location === 'object' && viewingRecord.location !== null ? 
                          (viewingRecord.location.name || viewingRecord.location.code || viewingRecord.location.id || '-') : 
                          (viewingRecord.location || '-')}
                      </Text>
                    </Col>
                  </Row>
                </div>

                <div style={{ marginTop: 16 }}>
                  <Text strong style={{ color: '#13c2c2' }}>备注信息</Text>
                  <Divider style={{ margin: '8px 0' }} />
                  <div style={{ 
                    background: '#f5f5f5', 
                    padding: '12px', 
                    borderRadius: '6px',
                    minHeight: '60px'
                  }}>
                    <Text>{viewingRecord.notes || '无备注'}</Text>
                  </div>
                </div>
              </>
            ) : viewingRecord.address ? (
              /* 仓库详情 */
              <>
                <Row gutter={[16, 16]}>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#1890ff' }}>基本信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>仓库编码: </Text>
                        <Text style={{ color: '#1890ff', fontSize: '16px', fontWeight: 'bold' }}>
                          {viewingRecord.code}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>仓库名称: </Text>
                        <Text style={{ fontSize: '16px', fontWeight: 'bold' }}>
                          {viewingRecord.name}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>状态: </Text>
                        <Tag color={viewingRecord.is_active ? 'green' : 'red'}>
                          {viewingRecord.is_active ? '启用' : '禁用'}
                        </Tag>
                      </div>
                    </div>
                  </Col>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#52c41a' }}>位置信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>地址: </Text>
                        <Text>{viewingRecord.address}</Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>联系人: </Text>
                        <Text>{viewingRecord.contact_person || '未设置'}</Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>联系电话: </Text>
                        <Text>{viewingRecord.contact_phone || '未设置'}</Text>
                      </div>
                    </div>
                  </Col>
                </Row>
                
                <div style={{ marginTop: 16 }}>
                  <Text strong style={{ color: '#722ed1' }}>描述信息</Text>
                  <Divider style={{ margin: '8px 0' }} />
                  <div style={{ 
                    background: '#f5f5f5', 
                    padding: '12px', 
                    borderRadius: '6px',
                    minHeight: '60px'
                  }}>
                    <Text>{viewingRecord.description || '暂无描述'}</Text>
                  </div>
                </div>

                <div style={{ marginTop: 16 }}>
                  <Text strong style={{ color: '#13c2c2' }}>时间信息</Text>
                  <Divider style={{ margin: '8px 0' }} />
                  <Row gutter={[16, 8]}>
                    <Col span={12}>
                      <Text strong>创建时间: </Text>
                      <Text>{viewingRecord.created_at ? new Date(viewingRecord.created_at).toLocaleString() : '-'}</Text>
                    </Col>
                    <Col span={12}>
                      <Text strong>更新时间: </Text>
                      <Text>{viewingRecord.updated_at ? new Date(viewingRecord.updated_at).toLocaleString() : '-'}</Text>
                    </Col>
                  </Row>
                </div>
              </>
            ) : viewingRecord.quantity !== undefined ? (
              /* 库存记录详情 */
              <>
                <Row gutter={[16, 16]}>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#1890ff' }}>物料信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>物料编码: </Text>
                        <Text style={{ color: '#1890ff', fontSize: '16px', fontWeight: 'bold' }}>
                          {viewingRecord.item?.code || '-'}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>物料名称: </Text>
                        <Text style={{ fontSize: '16px', fontWeight: 'bold' }}>
                          {viewingRecord.item?.name || '-'}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>分类: </Text>
                        <Text>
                          {typeof viewingRecord.item?.category === 'object' && viewingRecord.item?.category !== null ? 
                            (viewingRecord.item.category.name || viewingRecord.item.category.code || viewingRecord.item.category.id || '-') : 
                            (viewingRecord.item?.category || '-')}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>单位: </Text>
                        <Text>
                          {typeof viewingRecord.item?.unit === 'object' && viewingRecord.item?.unit !== null ? 
                            (viewingRecord.item.unit.name || viewingRecord.item.unit.symbol || viewingRecord.item.unit.id || '-') : 
                            (viewingRecord.item?.unit || '-')}
                        </Text>
                      </div>
                    </div>
                  </Col>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#52c41a' }}>库存信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>当前库存: </Text>
                        <Text style={{ color: '#52c41a', fontSize: '18px', fontWeight: 'bold' }}>
                          {viewingRecord.quantity}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>可用库存: </Text>
                        <Text style={{ color: '#1890ff', fontSize: '16px', fontWeight: 'bold' }}>
                          {viewingRecord.available_qty || viewingRecord.quantity}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>预留库存: </Text>
                        <Text style={{ color: '#fa8c16', fontSize: '16px', fontWeight: 'bold' }}>
                          {viewingRecord.reserved_qty || 0}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>最小库存: </Text>
                        <Text>{viewingRecord.min_stock || viewingRecord.item?.min_stock || '未设置'}</Text>
                      </div>
                    </div>
                  </Col>
                </Row>

                <Row gutter={[16, 16]}>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#722ed1' }}>仓库信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>仓库: </Text>
                        <Text>
                          {typeof viewingRecord.warehouse === 'object' && viewingRecord.warehouse !== null ? 
                            (viewingRecord.warehouse.name || viewingRecord.warehouse.code || viewingRecord.warehouse.id || '-') : 
                            (viewingRecord.warehouse || '-')}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>库位: </Text>
                        <Text>
                          {typeof viewingRecord.location === 'object' && viewingRecord.location !== null ? 
                            (viewingRecord.location.name || viewingRecord.location.code || viewingRecord.location.id || '-') : 
                            (viewingRecord.location || '-')}
                        </Text>
                      </div>
                    </div>
                  </Col>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#13c2c2' }}>时间信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>最后更新: </Text>
                        <Text>{viewingRecord.updated_at ? new Date(viewingRecord.updated_at).toLocaleString() : '-'}</Text>
                      </div>
                    </div>
                  </Col>
                </Row>
              </>
            ) : (
              /* 物料详情 */
              <>
                <Row gutter={[16, 16]}>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#1890ff' }}>基本信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>物料编码: </Text>
                        <Text>{viewingRecord.code}</Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>物料名称: </Text>
                        <Text>{viewingRecord.name}</Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>物料分类: </Text>
                        <Tag color="blue">
                          {typeof viewingRecord.category === 'object' && viewingRecord.category !== null ? 
                            (viewingRecord.category.name || viewingRecord.category.code || viewingRecord.category.id || 'Unknown') : 
                            (viewingRecord.category || '-')}
                        </Tag>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>计量单位: </Text>
                        <Text>
                          {typeof viewingRecord.unit === 'object' && viewingRecord.unit !== null ? 
                            (viewingRecord.unit.name || viewingRecord.unit.symbol || viewingRecord.unit.id || '-') : 
                            (viewingRecord.unit || '-')}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>状态: </Text>
                        <Tag color={viewingRecord.isActive ? 'green' : 'red'}>
                          {viewingRecord.isActive ? '启用' : '禁用'}
                        </Tag>
                      </div>
                    </div>
                  </Col>
                  <Col span={12}>
                    <div style={{ marginBottom: 16 }}>
                      <Text strong style={{ color: '#52c41a' }}>价格信息</Text>
                      <Divider style={{ margin: '8px 0' }} />
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>单价: </Text>
                        <Text style={{ color: '#f5222d', fontSize: '16px', fontWeight: 'bold' }}>
                          ¥{viewingRecord.price?.toFixed(2) || '0.00'}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>成本: </Text>
                        <Text style={{ color: '#fa8c16', fontSize: '16px', fontWeight: 'bold' }}>
                          ¥{viewingRecord.cost?.toFixed(2) || '0.00'}
                        </Text>
                      </div>
                      <div style={{ marginBottom: 8 }}>
                        <Text strong>重订点: </Text>
                        <Text>{viewingRecord.reorderLevel || '未设置'}</Text>
                      </div>
                    </div>
                  </Col>
                </Row>
                
                <div style={{ marginTop: 16 }}>
                  <Text strong style={{ color: '#722ed1' }}>描述信息</Text>
                  <Divider style={{ margin: '8px 0' }} />
                  <div style={{ 
                    background: '#f5f5f5', 
                    padding: '12px', 
                    borderRadius: '6px',
                    minHeight: '60px'
                  }}>
                    <Text>{viewingRecord.description || '暂无描述'}</Text>
                  </div>
                </div>

                <div style={{ marginTop: 16 }}>
                  <Text strong style={{ color: '#13c2c2' }}>时间信息</Text>
                  <Divider style={{ margin: '8px 0' }} />
                  <Row gutter={[16, 8]}>
                    <Col span={12}>
                      <Text strong>创建时间: </Text>
                      <Text>{new Date(viewingRecord.createdAt).toLocaleString()}</Text>
                    </Col>
                    <Col span={12}>
                      <Text strong>更新时间: </Text>
                      <Text>{new Date(viewingRecord.updatedAt).toLocaleString()}</Text>
                    </Col>
                  </Row>
                </div>
              </>
            )}
          </div>
        )}
      </Modal>

      {/* 导入模态框 */}
      <Modal
        title="导入物料"
        open={isImportModalVisible}
        onOk={handleImportConfirm}
        onCancel={() => {
          setIsImportModalVisible(false);
          setImportFile(null);
        }}
        confirmLoading={importLoading}
        okText="开始导入"
        cancelText="取消"
        width={600}
      >
        <div style={{ marginBottom: 16 }}>
          <Alert
            message="导入说明"
            description={
              <div>
                <p>1. 支持 Excel (.xlsx, .xls) 和 CSV (.csv) 格式文件</p>
                <p>2. 文件大小不能超过 10MB</p>
                <p>3. 请确保文件包含以下列：物料编码、物料名称、分类、单位、单价、描述</p>
                <p>4. 第一行应为表头，从第二行开始为数据</p>
              </div>
            }
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
        </div>
        
        <Upload.Dragger
          name="file"
          multiple={false}
          beforeUpload={handleFileChange}
          fileList={importFile ? [{ uid: '1', name: importFile.name, status: 'done' }] : []}
          onRemove={() => setImportFile(null)}
          accept=".xlsx,.xls,.csv"
        >
          <p className="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
          <p className="ant-upload-hint">
            支持单个文件上传，严格禁止上传公司数据或其他敏感文件
          </p>
        </Upload.Dragger>

        {importFile && (
          <div style={{ marginTop: 16 }}>
            <Text strong>已选择文件：</Text>
            <Text>{importFile.name}</Text>
            <Text type="secondary"> ({(importFile.size / 1024 / 1024).toFixed(2)} MB)</Text>
          </div>
        )}
      </Modal>
    </div>
  );
}

export default withAuth(InventoryPage);