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
  Timeline,
  Divider,
  Alert,
  List
} from 'antd';
import { 
  PlusOutlined, 
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  MoreOutlined,
  ProjectOutlined,
  TeamOutlined,
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
  UserOutlined,
  FlagOutlined,
  BugOutlined,
  RocketOutlined,
  StarOutlined,
  BarChartOutlined,
  CheckSquareOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;
const { Option } = Select;

export default function ProjectPage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'project' | 'task' | 'milestone'>('project');
  const [form] = Form.useForm();

  // 模拟项目数据
  const projectData = [
    {
      key: '1',
      projectName: 'ERP系统升级',
      projectCode: 'PRJ2025001',
      manager: '张三',
      team: ['张三', '李四', '王五', '赵六'],
      startDate: '2025-01-01',
      endDate: '2025-06-30',
      status: 'in_progress',
      priority: 'high',
      progress: 65,
      budget: 500000,
      spent: 325000,
      description: '企业资源规划系统全面升级改造',
      tasks: 45,
      completedTasks: 29,
      milestones: 8,
      completedMilestones: 5
    },
    {
      key: '2',
      projectName: '移动端APP开发',
      projectCode: 'PRJ2025002',
      manager: '李四',
      team: ['李四', '王五', '钱七'],
      startDate: '2025-02-01',
      endDate: '2025-08-31',
      status: 'planned',
      priority: 'medium',
      progress: 0,
      budget: 300000,
      spent: 0,
      description: '企业移动端应用程序开发',
      tasks: 32,
      completedTasks: 0,
      milestones: 6,
      completedMilestones: 0
    },
    {
      key: '3',
      projectName: '数据中心迁移',
      projectCode: 'PRJ2024003',
      manager: '王五',
      team: ['王五', '赵六', '孙八'],
      startDate: '2024-10-01',
      endDate: '2024-12-31',
      status: 'completed',
      priority: 'high',
      progress: 100,
      budget: 800000,
      spent: 750000,
      description: '企业数据中心整体迁移项目',
      tasks: 28,
      completedTasks: 28,
      milestones: 5,
      completedMilestones: 5
    }
  ];

  // 模拟任务数据
  const taskData = [
    {
      key: '1',
      taskName: '需求分析',
      taskCode: 'TASK001',
      projectName: 'ERP系统升级',
      assignee: '张三',
      priority: 'high',
      status: 'completed',
      startDate: '2025-01-01',
      endDate: '2025-01-15',
      progress: 100,
      estimatedHours: 80,
      actualHours: 75,
      description: '系统需求调研和分析'
    },
    {
      key: '2',
      taskName: '数据库设计',
      taskCode: 'TASK002',
      projectName: 'ERP系统升级',
      assignee: '李四',
      priority: 'high',
      status: 'in_progress',
      startDate: '2025-01-16',
      endDate: '2025-02-15',
      progress: 70,
      estimatedHours: 120,
      actualHours: 84,
      description: '数据库结构设计和优化'
    },
    {
      key: '3',
      taskName: 'UI界面设计',
      taskCode: 'TASK003',
      projectName: '移动端APP开发',
      assignee: '王五',
      priority: 'medium',
      status: 'planned',
      startDate: '2025-02-01',
      endDate: '2025-03-01',
      progress: 0,
      estimatedHours: 100,
      actualHours: 0,
      description: '移动端用户界面设计'
    },
    {
      key: '4',
      taskName: '系统测试',
      taskCode: 'TASK004',
      projectName: 'ERP系统升级',
      assignee: '赵六',
      priority: 'medium',
      status: 'pending',
      startDate: '2025-05-01',
      endDate: '2025-06-15',
      progress: 0,
      estimatedHours: 150,
      actualHours: 0,
      description: '系统功能和性能测试'
    }
  ];

  // 模拟里程碑数据
  const milestoneData = [
    {
      key: '1',
      milestoneName: '需求确认',
      projectName: 'ERP系统升级',
      targetDate: '2025-01-15',
      status: 'completed',
      description: '项目需求分析完成并确认',
      completedDate: '2025-01-14'
    },
    {
      key: '2',
      milestoneName: '原型设计',
      projectName: 'ERP系统升级',
      targetDate: '2025-03-01',
      status: 'in_progress',
      description: '系统原型设计和评审',
      completedDate: null
    },
    {
      key: '3',
      milestoneName: 'Alpha版本',
      projectName: '移动端APP开发',
      targetDate: '2025-05-01',
      status: 'planned',
      description: '移动端APP Alpha版本发布',
      completedDate: null
    }
  ];

  const getStatusColor = (status: string) => {
    const colors: { [key: string]: string } = {
      planned: 'blue',
      in_progress: 'orange',
      completed: 'green',
      cancelled: 'red',
      paused: 'purple',
      pending: 'gray',
      overdue: 'red'
    };
    return colors[status] || 'default';
  };

  const getStatusText = (status: string) => {
    const texts: { [key: string]: string } = {
      planned: '计划中',
      in_progress: '进行中',
      completed: '已完成',
      cancelled: '已取消',
      paused: '已暂停',
      pending: '待开始',
      overdue: '已逾期'
    };
    return texts[status] || status;
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

  const projectColumns = [
    {
      title: '项目信息',
      key: 'project',
      render: (record: any) => (
        <div>
          <Text strong style={{ display: 'block' }}>{record.projectName}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.projectCode}
          </Text>
        </div>
      ),
    },
    {
      title: '项目经理',
      key: 'manager',
      render: (record: any) => (
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <Avatar size="small" icon={<UserOutlined />} style={{ marginRight: 8 }} />
          <Text>{record.manager}</Text>
        </div>
      ),
    },
    {
      title: '团队',
      key: 'team',
      render: (record: any) => (
        <Avatar.Group maxCount={3} size="small">
          {record.team.map((member: string, index: number) => (
            <Tooltip key={index} title={member}>
              <Avatar size="small">{member.charAt(0)}</Avatar>
            </Tooltip>
          ))}
        </Avatar.Group>
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
            {record.completedTasks}/{record.tasks} 任务
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
          {getStatusText(status)}
        </Tag>
      ),
    },
    {
      title: '预算',
      key: 'budget',
      render: (record: any) => (
        <div>
          <Text strong style={{ color: '#52c41a' }}>¥{record.budget ? record.budget.toLocaleString() : '0.00'}</Text>
          <br />
          <Text type="secondary" style={{ fontSize: 12 }}>
            已用: ¥{record.spent ? record.spent.toLocaleString() : '0.00'}
          </Text>
        </div>
      ),
    },
    {
      title: '时间',
      key: 'time',
      render: (record: any) => (
        <div>
          <Text style={{ fontSize: 12 }}>开始: {record.startDate}</Text>
          <br />
          <Text style={{ fontSize: 12 }}>结束: {record.endDate}</Text>
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
              {
                key: 'gantt',
                label: '甘特图',
                icon: <BarChartOutlined />,
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

  const taskColumns = [
    {
      title: '任务信息',
      key: 'task',
      render: (record: any) => (
        <div>
          <Text strong style={{ display: 'block' }}>{record.taskName}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.taskCode} | {record.projectName}
          </Text>
        </div>
      ),
    },
    {
      title: '负责人',
      key: 'assignee',
      render: (record: any) => (
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <Avatar size="small" icon={<UserOutlined />} style={{ marginRight: 8 }} />
          <Text>{record.assignee}</Text>
        </div>
      ),
    },
    {
      title: '进度',
      key: 'progress',
      render: (record: any) => (
        <div style={{ width: 100 }}>
          <Progress 
            percent={record.progress} 
            size="small"
            status={record.status === 'completed' ? 'success' : 'active'}
          />
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
        <Tag color={getStatusColor(status)}>
          {getStatusText(status)}
        </Tag>
      ),
    },
    {
      title: '工时',
      key: 'hours',
      render: (record: any) => (
        <div>
          <Text>预估: {record.estimatedHours}h</Text>
          <br />
          <Text type="secondary" style={{ fontSize: 12 }}>
            实际: {record.actualHours}h
          </Text>
        </div>
      ),
    },
    {
      title: '时间',
      key: 'time',
      render: (record: any) => (
        <div>
          <Text style={{ fontSize: 12 }}>开始: {record.startDate}</Text>
          <br />
          <Text style={{ fontSize: 12 }}>结束: {record.endDate}</Text>
        </div>
      ),
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
          <Tooltip title="开始任务">
            <Button type="text" icon={<PlayCircleOutlined />} size="small" />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const milestoneColumns = [
    {
      title: '里程碑',
      key: 'milestone',
      render: (record: any) => (
        <div>
          <Text strong style={{ display: 'block' }}>{record.milestoneName}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.projectName}
          </Text>
        </div>
      ),
    },
    {
      title: '目标日期',
      dataIndex: 'targetDate',
      key: 'targetDate',
    },
    {
      title: '完成日期',
      dataIndex: 'completedDate',
      key: 'completedDate',
      render: (date: string) => date || '-',
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
            <FlagOutlined />
          }
        >
          {getStatusText(status)}
        </Tag>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
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
      message.success(`${modalType === 'project' ? '项目' : modalType === 'task' ? '任务' : '里程碑'}创建成功！`);
      setIsModalVisible(false);
      form.resetFields();
    });
  };

  const showModal = (type: 'project' | 'task' | 'milestone') => {
    setModalType(type);
    setIsModalVisible(true);
  };

  // 计算统计数据
  const totalProjects = projectData.length;
  const activeProjects = projectData.filter(p => p.status === '进行中').length;
  const totalTasks = taskData.length;
  const totalMilestones = milestoneData.length;

  // 定义Tabs的items
  const tabItems = [
    {
      key: 'projects',
      label: '项目管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索项目..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="计划中">计划中</Option>
                <Option value="进行中">进行中</Option>
                <Option value="已完成">已完成</Option>
                <Option value="已暂停">已暂停</Option>
              </Select>
              <Select placeholder="优先级筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="高">高</Option>
                <Option value="中">中</Option>
                <Option value="低">低</Option>
              </Select>
              <DatePicker placeholder="开始日期" style={{ width: 150 }} />
              <DatePicker placeholder="结束日期" style={{ width: 150 }} />
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('project')}
              >
                新建项目
              </Button>
            </Space>
          </div>
          <Table 
            columns={projectColumns} 
            dataSource={projectData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
            scroll={{ x: 1200 }}
          />
        </>
      )
    },
    {
      key: 'tasks',
      label: '任务管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索任务..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="项目筛选" style={{ width: 150 }}>
                <Option value="all">全部项目</Option>
                <Option value="project1">ERP系统开发</Option>
                <Option value="project2">移动应用开发</Option>
                <Option value="project3">数据分析平台</Option>
              </Select>
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="待开始">待开始</Option>
                <Option value="进行中">进行中</Option>
                <Option value="已完成">已完成</Option>
                <Option value="已延期">已延期</Option>
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
                onClick={() => showModal('task')}
              >
                新建任务
              </Button>
            </Space>
          </div>
          <Table 
            columns={taskColumns} 
            dataSource={taskData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
            scroll={{ x: 1200 }}
          />
        </>
      )
    },
    {
      key: 'milestones',
      label: '里程碑',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索里程碑..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="项目筛选" style={{ width: 150 }}>
                <Option value="all">全部项目</Option>
                <Option value="project1">ERP系统开发</Option>
                <Option value="project2">移动应用开发</Option>
                <Option value="project3">数据分析平台</Option>
              </Select>
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="未开始">未开始</Option>
                <Option value="进行中">进行中</Option>
                <Option value="已完成">已完成</Option>
                <Option value="已延期">已延期</Option>
              </Select>
              <DatePicker placeholder="目标日期" style={{ width: 150 }} />
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('milestone')}
              >
                新建里程碑
              </Button>
            </Space>
          </div>
          <Table 
            columns={milestoneColumns} 
            dataSource={milestoneData}
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
          📋 项目管理
        </Title>
        <Text type="secondary">管理项目、任务和里程碑</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="项目总数"
              value={totalProjects}
              prefix={<ProjectOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="进行中项目"
              value={activeProjects}
              prefix={<PlayCircleOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="任务总数"
              value={totalTasks}
              prefix={<CheckSquareOutlined style={{ color: '#faad14' }} />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="里程碑"
              value={totalMilestones}
              prefix={<FlagOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <Tabs defaultActiveKey="projects" items={tabItems} />
      </Card>

      {/* 创建/编辑模态框 */}
      <Modal
        title={`${modalType === 'project' ? '新建项目' : modalType === 'task' ? '新建任务' : '新建里程碑'}`}
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={() => {
          setIsModalVisible(false);
          form.resetFields();
        }}
        width={800}
      >
        <Form form={form} layout="vertical">
          {modalType === 'project' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="projectName" label="项目名称" rules={[{ required: true }]}>
                    <Input placeholder="请输入项目名称" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="projectCode" label="项目编码" rules={[{ required: true }]}>
                    <Input placeholder="请输入项目编码" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="manager" label="项目经理" rules={[{ required: true }]}>
                    <Select placeholder="选择项目经理">
                      <Option value="张三">张三</Option>
                      <Option value="李四">李四</Option>
                      <Option value="王五">王五</Option>
                    </Select>
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
              <Form.Item name="budget" label="项目预算" rules={[{ required: true }]}>
                <InputNumber 
                  placeholder="请输入预算金额" 
                  style={{ width: '100%' }} 
                  min={0}
                  formatter={value => `¥ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
                />
              </Form.Item>
              <Form.Item name="description" label="项目描述">
                <Input.TextArea placeholder="请输入项目描述" rows={3} />
              </Form.Item>
            </>
          )}
          {modalType === 'task' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="taskName" label="任务名称" rules={[{ required: true }]}>
                    <Input placeholder="请输入任务名称" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="taskCode" label="任务编码" rules={[{ required: true }]}>
                    <Input placeholder="请输入任务编码" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="projectCode" label="所属项目" rules={[{ required: true }]}>
                    <Select placeholder="选择项目">
                      <Option value="PRJ2025001">ERP系统升级</Option>
                      <Option value="PRJ2025002">移动端APP开发</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="assignee" label="负责人" rules={[{ required: true }]}>
                    <Select placeholder="选择负责人">
                      <Option value="张三">张三</Option>
                      <Option value="李四">李四</Option>
                      <Option value="王五">王五</Option>
                      <Option value="赵六">赵六</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="priority" label="优先级" rules={[{ required: true }]}>
                    <Select placeholder="选择优先级">
                      <Option value="high">高</Option>
                      <Option value="medium">中</Option>
                      <Option value="low">低</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="estimatedHours" label="预估工时" rules={[{ required: true }]}>
                    <InputNumber placeholder="请输入工时" style={{ width: '100%' }} min={1} />
                  </Form.Item>
                </Col>
              </Row>
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
              <Form.Item name="description" label="任务描述">
                <Input.TextArea placeholder="请输入任务描述" rows={3} />
              </Form.Item>
            </>
          )}
          {modalType === 'milestone' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="milestoneName" label="里程碑名称" rules={[{ required: true }]}>
                    <Input placeholder="请输入里程碑名称" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="projectCode" label="所属项目" rules={[{ required: true }]}>
                    <Select placeholder="选择项目">
                      <Option value="PRJ2025001">ERP系统升级</Option>
                      <Option value="PRJ2025002">移动端APP开发</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="targetDate" label="目标日期" rules={[{ required: true }]}>
                <DatePicker style={{ width: '100%' }} />
              </Form.Item>
              <Form.Item name="description" label="里程碑描述">
                <Input.TextArea placeholder="请输入里程碑描述" rows={3} />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>
    </div>
  );
}