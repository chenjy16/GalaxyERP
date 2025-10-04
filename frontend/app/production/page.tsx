'use client';

import { useState } from 'react';
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
  Progress,
  InputNumber,
  Tree,
  Divider,
  Alert
} from 'antd';
import { 
  PlusOutlined, 
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  MoreOutlined,
  BuildOutlined,
  ToolOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  CloseCircleOutlined,
  ExportOutlined,
  ImportOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  StopOutlined,
  SettingOutlined,
  FileTextOutlined,
  CalendarOutlined,
  TeamOutlined,
  BoxPlotOutlined,
  MenuOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;
const { Option } = Select;

export default function ProductionPage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'workOrder' | 'bom' | 'plan'>('workOrder');
  const [form] = Form.useForm();

  // 模拟工单数据
  const workOrderData = [
    {
      key: '1',
      orderNo: 'WO2025001',
      productName: 'iPhone 15 Pro',
      productCode: 'IP15P001',
      quantity: 1000,
      plannedStartDate: '2025-01-15',
      plannedEndDate: '2025-01-25',
      actualStartDate: '2025-01-15',
      actualEndDate: null,
      status: 'in_progress',
      priority: 'high',
      progress: 65,
      responsible: '张三',
      workshop: '组装车间A',
      completedQuantity: 650,
      defectQuantity: 15,
      cost: 850000
    },
    {
      key: '2',
      orderNo: 'WO2025002',
      productName: 'MacBook Pro 14',
      productCode: 'MBP14001',
      quantity: 500,
      plannedStartDate: '2025-01-20',
      plannedEndDate: '2025-02-05',
      actualStartDate: null,
      actualEndDate: null,
      status: 'planned',
      priority: 'medium',
      progress: 0,
      responsible: '李四',
      workshop: '组装车间B',
      completedQuantity: 0,
      defectQuantity: 0,
      cost: 1200000
    },
    {
      key: '3',
      orderNo: 'WO2025003',
      productName: 'iPad Air',
      productCode: 'IPA001',
      quantity: 800,
      plannedStartDate: '2025-01-10',
      plannedEndDate: '2025-01-20',
      actualStartDate: '2025-01-10',
      actualEndDate: '2025-01-19',
      status: 'completed',
      priority: 'low',
      progress: 100,
      responsible: '王五',
      workshop: '组装车间C',
      completedQuantity: 800,
      defectQuantity: 8,
      cost: 640000
    }
  ];

  // 模拟物料清单数据
  const bomData = [
    {
      key: '1',
      productName: 'iPhone 15 Pro',
      productCode: 'IP15P001',
      version: 'V1.0',
      status: 'active',
      createDate: '2024-12-01',
      materials: [
        { code: 'MAT001', name: 'A17 Pro芯片', quantity: 1, unit: '个', cost: 800 },
        { code: 'MAT002', name: '6.1英寸屏幕', quantity: 1, unit: '个', cost: 300 },
        { code: 'MAT003', name: '钛合金外壳', quantity: 1, unit: '个', cost: 200 },
        { code: 'MAT004', name: '电池', quantity: 1, unit: '个', cost: 50 }
      ],
      totalCost: 1350
    },
    {
      key: '2',
      productName: 'MacBook Pro 14',
      productCode: 'MBP14001',
      version: 'V2.1',
      status: 'active',
      createDate: '2024-11-15',
      materials: [
        { code: 'MAT005', name: 'M3 Pro芯片', quantity: 1, unit: '个', cost: 1200 },
        { code: 'MAT006', name: '14英寸Liquid Retina XDR显示屏', quantity: 1, unit: '个', cost: 600 },
        { code: 'MAT007', name: '铝合金机身', quantity: 1, unit: '个', cost: 300 },
        { code: 'MAT008', name: '键盘', quantity: 1, unit: '个', cost: 100 }
      ],
      totalCost: 2200
    }
  ];

  // 模拟生产计划数据
  const planData = [
    {
      key: '1',
      planNo: 'PP2025001',
      planName: '2025年Q1生产计划',
      startDate: '2025-01-01',
      endDate: '2025-03-31',
      status: 'executing',
      progress: 45,
      totalOrders: 15,
      completedOrders: 7,
      responsible: '生产部',
      description: '第一季度主要产品生产计划'
    },
    {
      key: '2',
      planNo: 'PP2025002',
      planName: '春节特别生产计划',
      startDate: '2025-02-01',
      endDate: '2025-02-28',
      status: 'planned',
      progress: 0,
      totalOrders: 8,
      completedOrders: 0,
      responsible: '生产部',
      description: '春节期间特殊产品生产安排'
    }
  ];

  const getStatusColor = (status: string) => {
    const colors: { [key: string]: string } = {
      planned: 'blue',
      in_progress: 'orange',
      completed: 'green',
      cancelled: 'red',
      paused: 'purple',
      active: 'green',
      inactive: 'gray',
      executing: 'orange'
    };
    return colors[status] || 'default';
  };

  const getStatusText = (status: string, type: string) => {
    if (type === 'workOrder') {
      const texts: { [key: string]: string } = {
        planned: '计划中',
        in_progress: '进行中',
        completed: '已完成',
        cancelled: '已取消',
        paused: '已暂停'
      };
      return texts[status] || status;
    } else if (type === 'bom') {
      const texts: { [key: string]: string } = {
        active: '启用',
        inactive: '停用'
      };
      return texts[status] || status;
    } else if (type === 'plan') {
      const texts: { [key: string]: string } = {
        planned: '计划中',
        executing: '执行中',
        completed: '已完成',
        cancelled: '已取消'
      };
      return texts[status] || status;
    }
    return status;
  };

  const getPriorityColor = (priority: string) => {
    const colors: { [key: string]: string } = {
      high: 'red',
      medium: 'orange',
      low: 'green'
    };
    return colors[priority] || 'default';
  };

  const getPriorityText = (priority: string) => {
    const texts: { [key: string]: string } = {
      high: '高',
      medium: '中',
      low: '低'
    };
    return texts[priority] || priority;
  };

  const workOrderColumns = [
    {
      title: '工单信息',
      key: 'workOrder',
      render: (record: any) => (
        <div>
          <Text strong style={{ display: 'block' }}>{record.orderNo}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.productName} ({record.productCode})
          </Text>
        </div>
      ),
    },
    {
      title: '数量',
      key: 'quantity',
      render: (record: any) => (
        <div>
          <Text strong>{record.quantity ? record.quantity.toLocaleString() : '0'}</Text>
          <br />
          <Text type="secondary" style={{ fontSize: 12 }}>
            已完成: {record.completedQuantity ? record.completedQuantity.toLocaleString() : '0'}
          </Text>
        </div>
      ),
    },
    {
      title: '进度',
      key: 'progress',
      render: (record: any) => (
        <div style={{ width: 120 }}>
          <Progress 
            percent={record.progress} 
            size="small"
            status={record.status === 'completed' ? 'success' : 'active'}
          />
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.progress}%
          </Text>
        </div>
      ),
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      key: 'priority',
      render: (priority: string) => (
        <Tag color={getPriorityColor(priority)}>
          {getPriorityText(priority)}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag 
          color={getStatusColor(status)}
          icon={
            status === 'completed' ? <CheckCircleOutlined /> :
            status === 'in_progress' ? <ClockCircleOutlined /> :
            status === 'cancelled' ? <CloseCircleOutlined /> :
            <ExclamationCircleOutlined />
          }
        >
          {getStatusText(status, 'workOrder')}
        </Tag>
      ),
    },
    {
      title: '负责人',
      dataIndex: 'responsible',
      key: 'responsible',
    },
    {
      title: '车间',
      dataIndex: 'workshop',
      key: 'workshop',
      render: (workshop: string) => <Tag color="blue">{workshop}</Tag>,
    },
    {
      title: '计划时间',
      key: 'plannedTime',
      render: (record: any) => (
        <div>
          <Text style={{ fontSize: 12 }}>开始: {record.plannedStartDate}</Text>
          <br />
          <Text style={{ fontSize: 12 }}>结束: {record.plannedEndDate}</Text>
        </div>
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
              ...(record.status === 'planned' ? [{
                key: 'start',
                label: '开始生产',
                icon: <PlayCircleOutlined />,
              }] : []),
              ...(record.status === 'in_progress' ? [{
                key: 'pause',
                label: '暂停',
                icon: <PauseCircleOutlined />,
              }] : []),
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

  const bomColumns = [
    {
      title: '产品信息',
      key: 'product',
      render: (record: any) => (
        <div>
          <Text strong style={{ display: 'block' }}>{record.productName}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.productCode} | {record.version}
          </Text>
        </div>
      ),
    },
    {
      title: '物料数量',
      key: 'materialCount',
      render: (record: any) => (
        <Text strong>{record.materials.length} 种物料</Text>
      ),
    },
    {
      title: '总成本',
      dataIndex: 'totalCost',
      key: 'totalCost',
      render: (cost: number) => (
        <Text strong style={{ color: '#52c41a' }}>¥{cost ? cost.toLocaleString() : '0.00'}</Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>
          {getStatusText(status, 'bom')}
        </Tag>
      ),
    },
    {
      title: '创建日期',
      dataIndex: 'createDate',
      key: 'createDate',
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Space>
          <Tooltip title="查看详情">
            <Button type="text" icon={<EyeOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="编辑">
            <Button type="text" icon={<EditOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="复制">
            <Button type="text" icon={<FileTextOutlined />} size="small" />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const planColumns = [
    {
      title: '计划信息',
      key: 'plan',
      render: (record: any) => (
        <div>
          <Text strong style={{ display: 'block' }}>{record.planName}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.planNo}
          </Text>
        </div>
      ),
    },
    {
      title: '时间范围',
      key: 'timeRange',
      render: (record: any) => (
        <div>
          <Text style={{ fontSize: 12 }}>开始: {record.startDate}</Text>
          <br />
          <Text style={{ fontSize: 12 }}>结束: {record.endDate}</Text>
        </div>
      ),
    },
    {
      title: '工单进度',
      key: 'orderProgress',
      render: (record: any) => (
        <div>
          <Progress 
            percent={Math.round((record.completedOrders / record.totalOrders) * 100)} 
            size="small"
          />
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.completedOrders}/{record.totalOrders} 个工单
          </Text>
        </div>
      ),
    },
    {
      title: '整体进度',
      dataIndex: 'progress',
      key: 'progress',
      render: (progress: number) => (
        <div style={{ width: 100 }}>
          <Progress percent={progress} size="small" />
        </div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>
          {getStatusText(status, 'plan')}
        </Tag>
      ),
    },
    {
      title: '负责部门',
      dataIndex: 'responsible',
      key: 'responsible',
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Space>
          <Tooltip title="查看详情">
            <Button type="text" icon={<EyeOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="编辑">
            <Button type="text" icon={<EditOutlined />} size="small" />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const handleModalOk = () => {
    form.validateFields().then(values => {
      // 表单提交处理
      message.success(`${modalType === 'workOrder' ? '工单' : modalType === 'bom' ? '物料清单' : '生产计划'}创建成功！`);
      setIsModalVisible(false);
      form.resetFields();
    });
  };

  const showModal = (type: 'workOrder' | 'bom' | 'plan') => {
    setModalType(type);
    setIsModalVisible(true);
  };

  // 计算统计数据
  const totalWorkOrders = workOrderData.length;
  const activeWorkOrders = workOrderData.filter(w => w.status === '进行中').length;
  const totalBOMs = bomData.length;
  const totalPlans = planData.length;

  // 定义Tabs的items
  const tabItems = [
    {
      key: 'workOrders',
      label: '工单管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索工单..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="待开始">待开始</Option>
                <Option value="进行中">进行中</Option>
                <Option value="已完成">已完成</Option>
                <Option value="已取消">已取消</Option>
              </Select>
              <Select placeholder="优先级筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="高">高</Option>
                <Option value="中">中</Option>
                <Option value="低">低</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('workOrder')}
              >
                新建工单
              </Button>
            </Space>
          </div>
          <Table 
            columns={workOrderColumns} 
            dataSource={workOrderData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
            scroll={{ x: 1200 }}
          />
        </>
      )
    },
    {
      key: 'bom',
      label: '物料清单',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索物料清单..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="启用">启用</Option>
                <Option value="禁用">禁用</Option>
              </Select>
              <Select placeholder="版本筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="v1.0">v1.0</Option>
                <Option value="v2.0">v2.0</Option>
                <Option value="v3.0">v3.0</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ImportOutlined />}>导入</Button>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('bom')}
              >
                新建BOM
              </Button>
            </Space>
          </div>
          <Table 
            columns={bomColumns} 
            dataSource={bomData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
            scroll={{ x: 1200 }}
          />
        </>
      )
    },
    {
      key: 'plans',
      label: '生产计划',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索生产计划..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="计划中">计划中</Option>
                <Option value="执行中">执行中</Option>
                <Option value="已完成">已完成</Option>
              </Select>
              <DatePicker placeholder="开始日期" style={{ width: 150 }} />
              <DatePicker placeholder="结束日期" style={{ width: 150 }} />
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('plan')}
              >
                新建计划
              </Button>
            </Space>
          </div>
          <Table 
            columns={planColumns} 
            dataSource={planData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
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
          🏭 生产管理
        </Title>
        <Text type="secondary">管理工单、物料清单和生产计划</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="工单总数"
              value={totalWorkOrders}
              prefix={<FileTextOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="进行中工单"
              value={activeWorkOrders}
              prefix={<PlayCircleOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="物料清单"
              value={totalBOMs}
              prefix={<MenuOutlined style={{ color: '#faad14' }} />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="生产计划"
              value={totalPlans}
              prefix={<CalendarOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <Tabs defaultActiveKey="workOrders" items={tabItems} />
      </Card>

      {/* 创建/编辑模态框 */}
      <Modal
        title={`${modalType === 'workOrder' ? '新建工单' : modalType === 'bom' ? '新建物料清单' : '新建生产计划'}`}
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={() => {
          setIsModalVisible(false);
          form.resetFields();
        }}
        width={800}
      >
        <Form form={form} layout="vertical">
          {modalType === 'workOrder' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="orderNo" label="工单编号" rules={[{ required: true }]}>
                    <Input placeholder="请输入工单编号" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="priority" label="优先级" rules={[{ required: true }]}>
                    <Select placeholder="选择优先级">
                      <Option value="high">高</Option>
                      <Option value="medium">中</Option>
                      <Option value="low">低</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="productCode" label="产品编码" rules={[{ required: true }]}>
                    <Select placeholder="选择产品">
                      <Option value="IP15P001">iPhone 15 Pro</Option>
                      <Option value="MBP14001">MacBook Pro 14</Option>
                      <Option value="IPA001">iPad Air</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="quantity" label="生产数量" rules={[{ required: true }]}>
                    <InputNumber placeholder="请输入数量" style={{ width: '100%' }} min={1} />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="plannedStartDate" label="计划开始日期" rules={[{ required: true }]}>
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="plannedEndDate" label="计划结束日期" rules={[{ required: true }]}>
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="responsible" label="负责人" rules={[{ required: true }]}>
                    <Select placeholder="选择负责人">
                      <Option value="张三">张三</Option>
                      <Option value="李四">李四</Option>
                      <Option value="王五">王五</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="workshop" label="生产车间" rules={[{ required: true }]}>
                    <Select placeholder="选择车间">
                      <Option value="组装车间A">组装车间A</Option>
                      <Option value="组装车间B">组装车间B</Option>
                      <Option value="组装车间C">组装车间C</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
            </>
          )}
          {modalType === 'bom' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="productCode" label="产品编码" rules={[{ required: true }]}>
                    <Input placeholder="请输入产品编码" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="productName" label="产品名称" rules={[{ required: true }]}>
                    <Input placeholder="请输入产品名称" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="version" label="版本号" rules={[{ required: true }]}>
                    <Input placeholder="请输入版本号" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="status" label="状态" rules={[{ required: true }]}>
                    <Select placeholder="选择状态">
                      <Option value="active">启用</Option>
                      <Option value="inactive">停用</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="description" label="描述">
                <Input.TextArea placeholder="请输入描述" rows={3} />
              </Form.Item>
            </>
          )}
          {modalType === 'plan' && (
            <>
              <Form.Item name="planName" label="计划名称" rules={[{ required: true }]}>
                <Input placeholder="请输入计划名称" />
              </Form.Item>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="startDate" label="开始日期" rules={[{ required: true }]}>
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="endDate" label="结束日期" rules={[{ required: true }]}>
                    <DatePicker style={{ width: '100%' }} />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="responsible" label="负责部门" rules={[{ required: true }]}>
                <Select placeholder="选择负责部门">
                  <Option value="生产部">生产部</Option>
                  <Option value="技术部">技术部</Option>
                  <Option value="质量部">质量部</Option>
                </Select>
              </Form.Item>
              <Form.Item name="description" label="计划描述">
                <Input.TextArea placeholder="请输入计划描述" rows={3} />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>
    </div>
  );
}