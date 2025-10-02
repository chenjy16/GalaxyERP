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
  Switch,
  InputNumber,
  Divider,
  Alert,
  List,
  Tree,
  Transfer,
  Checkbox,
  Radio,
  Upload,
  Descriptions
} from 'antd';
import { 
  PlusOutlined, 
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  MoreOutlined,
  UserOutlined,
  TeamOutlined,
  SettingOutlined,
  SecurityScanOutlined,
  KeyOutlined,
  LockOutlined,
  UnlockOutlined,
  ExportOutlined,
  ImportOutlined,
  ReloadOutlined,
  SaveOutlined,
  CopyOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
  UploadOutlined,
  DownloadOutlined,
  DatabaseOutlined,
  CloudOutlined,
  MailOutlined,
  PhoneOutlined,
  GlobalOutlined,
  BellOutlined,
  FileTextOutlined,
  HistoryOutlined,
  SafetyOutlined,
  CrownOutlined,
  ApiOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;
const { Option } = Select;
const { TreeNode } = Tree;

export default function SystemPage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'user' | 'role' | 'permission' | 'config'>('user');
  const [form] = Form.useForm();

  // 模拟用户数据
  const userData = [
    {
      key: '1',
      username: 'admin',
      realName: '系统管理员',
      email: 'admin@galaxy.com',
      phone: '13800138000',
      department: '信息技术部',
      role: '超级管理员',
      status: 'active',
      lastLogin: '2025-01-20 10:30:00',
      createTime: '2024-01-01 00:00:00',
      avatar: null
    },
    {
      key: '2',
      username: 'zhangsan',
      realName: '张三',
      email: 'zhangsan@galaxy.com',
      phone: '13800138001',
      department: '销售部',
      role: '销售经理',
      status: 'active',
      lastLogin: '2025-01-20 09:15:00',
      createTime: '2024-02-15 10:00:00',
      avatar: null
    },
    {
      key: '3',
      username: 'lisi',
      realName: '李四',
      email: 'lisi@galaxy.com',
      phone: '13800138002',
      department: '采购部',
      role: '采购员',
      status: 'inactive',
      lastLogin: '2025-01-18 16:45:00',
      createTime: '2024-03-10 14:30:00',
      avatar: null
    }
  ];

  // 模拟角色数据
  const roleData = [
    {
      key: '1',
      roleName: '超级管理员',
      roleCode: 'SUPER_ADMIN',
      description: '系统最高权限，可以管理所有功能',
      userCount: 1,
      permissions: ['user:read', 'user:write', 'role:read', 'role:write', 'system:read', 'system:write'],
      status: 'active',
      createTime: '2024-01-01 00:00:00'
    },
    {
      key: '2',
      roleName: '销售经理',
      roleCode: 'SALES_MANAGER',
      description: '销售部门管理权限',
      userCount: 5,
      permissions: ['sales:read', 'sales:write', 'customer:read', 'customer:write'],
      status: 'active',
      createTime: '2024-01-01 00:00:00'
    },
    {
      key: '3',
      roleName: '采购员',
      roleCode: 'PURCHASER',
      description: '采购相关操作权限',
      userCount: 3,
      permissions: ['purchase:read', 'purchase:write', 'supplier:read'],
      status: 'active',
      createTime: '2024-01-01 00:00:00'
    },
    {
      key: '4',
      roleName: '财务专员',
      roleCode: 'ACCOUNTANT',
      description: '财务模块操作权限',
      userCount: 2,
      permissions: ['accounting:read', 'accounting:write', 'report:read'],
      status: 'active',
      createTime: '2024-01-01 00:00:00'
    }
  ];

  // 模拟权限数据
  const permissionData = [
    {
      key: '1',
      permissionName: '用户管理',
      permissionCode: 'user:read',
      module: '系统管理',
      type: 'menu',
      description: '查看用户列表',
      status: 'active'
    },
    {
      key: '2',
      permissionName: '用户编辑',
      permissionCode: 'user:write',
      module: '系统管理',
      type: 'button',
      description: '编辑用户信息',
      status: 'active'
    },
    {
      key: '3',
      permissionName: '销售管理',
      permissionCode: 'sales:read',
      module: '销售管理',
      type: 'menu',
      description: '查看销售数据',
      status: 'active'
    },
    {
      key: '4',
      permissionName: '销售编辑',
      permissionCode: 'sales:write',
      module: '销售管理',
      type: 'button',
      description: '编辑销售数据',
      status: 'active'
    }
  ];

  // 模拟系统配置数据
  const configData = [
    {
      key: '1',
      configKey: 'system.title',
      configValue: 'Galaxy ERP 企业资源规划系统',
      description: '系统标题',
      category: '基础配置',
      type: 'string',
      required: true
    },
    {
      key: '2',
      configKey: 'system.logo',
      configValue: '/logo.png',
      description: '系统Logo',
      category: '基础配置',
      type: 'file',
      required: false
    },
    {
      key: '3',
      configKey: 'email.smtp.host',
      configValue: 'smtp.galaxy.com',
      description: 'SMTP服务器地址',
      category: '邮件配置',
      type: 'string',
      required: true
    },
    {
      key: '4',
      configKey: 'email.smtp.port',
      configValue: '587',
      description: 'SMTP端口',
      category: '邮件配置',
      type: 'number',
      required: true
    }
  ];

  const getStatusColor = (status: string) => {
    const colors: { [key: string]: string } = {
      active: 'green',
      inactive: 'red',
      pending: 'orange',
      locked: 'gray'
    };
    return colors[status] || 'default';
  };

  const getStatusText = (status: string) => {
    const texts: { [key: string]: string } = {
      active: '正常',
      inactive: '禁用',
      pending: '待审核',
      locked: '锁定'
    };
    return texts[status] || status;
  };

  const userColumns = [
    {
      title: '用户信息',
      key: 'user',
      render: (record: any) => (
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <Avatar 
            size="large" 
            icon={<UserOutlined />} 
            src={record.avatar}
            style={{ marginRight: 12 }}
          />
          <div>
            <Text strong style={{ display: 'block' }}>{record.realName}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              @{record.username}
            </Text>
          </div>
        </div>
      ),
    },
    {
      title: '联系方式',
      key: 'contact',
      render: (record: any) => (
        <div>
          <div style={{ marginBottom: 4 }}>
            <MailOutlined style={{ marginRight: 4, color: '#1890ff' }} />
            <Text style={{ fontSize: 12 }}>{record.email}</Text>
          </div>
          <div>
            <PhoneOutlined style={{ marginRight: 4, color: '#52c41a' }} />
            <Text style={{ fontSize: 12 }}>{record.phone}</Text>
          </div>
        </div>
      ),
    },
    {
      title: '部门',
      dataIndex: 'department',
      key: 'department',
    },
    {
      title: '角色',
      key: 'role',
      render: (record: any) => (
        <Tag color="blue" icon={<CrownOutlined />}>
          {record.role}
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
          icon={status === 'active' ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
        >
          {getStatusText(status)}
        </Tag>
      ),
    },
    {
      title: '最后登录',
      dataIndex: 'lastLogin',
      key: 'lastLogin',
      render: (time: string) => (
        <Text style={{ fontSize: 12 }}>{time}</Text>
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
                key: 'reset',
                label: '重置密码',
                icon: <KeyOutlined />,
              },
              {
                key: 'lock',
                label: record.status === 'active' ? '锁定' : '解锁',
                icon: record.status === 'active' ? <LockOutlined /> : <UnlockOutlined />,
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

  const roleColumns = [
    {
      title: '角色信息',
      key: 'role',
      render: (record: any) => (
        <div>
          <Text strong style={{ display: 'block' }}>{record.roleName}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.roleCode}
          </Text>
        </div>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '用户数量',
      key: 'userCount',
      render: (record: any) => (
        <Badge count={record.userCount} style={{ backgroundColor: '#52c41a' }} />
      ),
    },
    {
      title: '权限数量',
      key: 'permissionCount',
      render: (record: any) => (
        <Badge count={record.permissions.length} style={{ backgroundColor: '#1890ff' }} />
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
      title: '创建时间',
      dataIndex: 'createTime',
      key: 'createTime',
      render: (time: string) => (
        <Text style={{ fontSize: 12 }}>{time}</Text>
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
          <Tooltip title="权限配置">
            <Button type="text" icon={<SecurityScanOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="删除">
            <Button type="text" icon={<DeleteOutlined />} size="small" danger />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const permissionColumns = [
    {
      title: '权限名称',
      dataIndex: 'permissionName',
      key: 'permissionName',
    },
    {
      title: '权限代码',
      dataIndex: 'permissionCode',
      key: 'permissionCode',
      render: (code: string) => (
        <Text code>{code}</Text>
      ),
    },
    {
      title: '所属模块',
      dataIndex: 'module',
      key: 'module',
      render: (module: string) => (
        <Tag color="purple">{module}</Tag>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={type === 'menu' ? 'blue' : 'orange'}>
          {type === 'menu' ? '菜单' : '按钮'}
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
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Space>
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

  const configColumns = [
    {
      title: '配置项',
      key: 'config',
      render: (record: any) => (
        <div>
          <Text strong style={{ display: 'block' }}>{record.configKey}</Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            {record.description}
          </Text>
        </div>
      ),
    },
    {
      title: '当前值',
      key: 'value',
      render: (record: any) => {
        if (record.type === 'boolean') {
          return <Switch checked={record.configValue === 'true'} disabled />;
        }
        if (record.type === 'file') {
          return <Text type="secondary">文件路径</Text>;
        }
        return <Text code>{record.configValue}</Text>;
      },
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
      render: (category: string) => (
        <Tag color="geekblue">{category}</Tag>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color="cyan">{type}</Tag>
      ),
    },
    {
      title: '必填',
      dataIndex: 'required',
      key: 'required',
      render: (required: boolean) => (
        <Tag color={required ? 'red' : 'default'}>
          {required ? '必填' : '可选'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Space>
          <Tooltip title="编辑">
            <Button type="text" icon={<EditOutlined />} size="small" />
          </Tooltip>
          <Tooltip title="重置">
            <Button type="text" icon={<ReloadOutlined />} size="small" />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const handleModalOk = () => {
    form.validateFields().then(values => {
      console.log('Form values:', values);
      message.success(`${modalType === 'user' ? '用户' : modalType === 'role' ? '角色' : modalType === 'permission' ? '权限' : '配置'}操作成功！`);
      setIsModalVisible(false);
      form.resetFields();
    });
  };

  const showModal = (type: 'user' | 'role' | 'permission' | 'config') => {
    setModalType(type);
    setIsModalVisible(true);
  };

  // 计算统计数据
  const totalUsers = userData.length;
  const activeUsers = userData.filter(u => u.status === 'active').length;
  const totalRoles = roleData.length;
  const totalPermissions = permissionData.length;

  // 定义Tabs的items
  const tabItems = [
    {
      key: 'users',
      label: '用户管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索用户..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="部门筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="it">信息技术部</Option>
                <Option value="sales">销售部</Option>
                <Option value="purchase">采购部</Option>
                <Option value="finance">财务部</Option>
              </Select>
              <Select placeholder="角色筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="admin">管理员</Option>
                <Option value="manager">经理</Option>
                <Option value="employee">员工</Option>
              </Select>
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="active">正常</Option>
                <Option value="inactive">禁用</Option>
                <Option value="locked">锁定</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ImportOutlined />}>导入</Button>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('user')}
              >
                新建用户
              </Button>
            </Space>
          </div>
          <Table 
            columns={userColumns} 
            dataSource={userData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
            scroll={{ x: 1200 }}
          />
        </>
      )
    },
    {
      key: 'roles',
      label: '角色管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索角色..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="active">正常</Option>
                <Option value="inactive">禁用</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('role')}
              >
                新建角色
              </Button>
            </Space>
          </div>
          <Table 
            columns={roleColumns} 
            dataSource={roleData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
          />
        </>
      )
    },
    {
      key: 'permissions',
      label: '权限管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索权限..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="模块筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="system">系统管理</Option>
                <Option value="sales">销售管理</Option>
                <Option value="purchase">采购管理</Option>
                <Option value="inventory">库存管理</Option>
              </Select>
              <Select placeholder="类型筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="menu">菜单</Option>
                <Option value="button">按钮</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('permission')}
              >
                新建权限
              </Button>
            </Space>
          </div>
          <Table 
            columns={permissionColumns} 
            dataSource={permissionData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
          />
        </>
      )
    },
    {
      key: 'config',
      label: '系统配置',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索配置..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="分类筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="基础配置">基础配置</Option>
                <Option value="邮件配置">邮件配置</Option>
                <Option value="安全配置">安全配置</Option>
              </Select>
              <Select placeholder="类型筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="string">字符串</Option>
                <Option value="number">数字</Option>
                <Option value="boolean">布尔值</Option>
                <Option value="file">文件</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<SaveOutlined />} type="primary">保存配置</Button>
            </Space>
          </div>
          <Table 
            columns={configColumns} 
            dataSource={configData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
          />
        </>
      )
    },
    {
      key: 'backup',
      label: '备份管理',
      children: (
        <div style={{ textAlign: 'center', padding: '60px 20px' }}>
          <DatabaseOutlined style={{ fontSize: '48px', color: '#d9d9d9', marginBottom: '16px' }} />
          <Title level={3} style={{ color: '#999', marginBottom: '8px' }}>
            备份管理
          </Title>
          <Text type="secondary">功能暂未实现</Text>
        </div>
      )
    },
    {
      key: 'monitoring',
      label: '系统监控',
      children: (
        <div style={{ textAlign: 'center', padding: '60px 20px' }}>
          <ApiOutlined style={{ fontSize: '48px', color: '#d9d9d9', marginBottom: '16px' }} />
          <Title level={3} style={{ color: '#999', marginBottom: '8px' }}>
            系统监控
          </Title>
          <Text type="secondary">功能暂未实现</Text>
        </div>
      )
    },
    {
      key: 'audit',
      label: '审计日志',
      children: (
        <div style={{ textAlign: 'center', padding: '60px 20px' }}>
          <HistoryOutlined style={{ fontSize: '48px', color: '#d9d9d9', marginBottom: '16px' }} />
          <Title level={3} style={{ color: '#999', marginBottom: '8px' }}>
            审计日志
          </Title>
          <Text type="secondary">功能暂未实现</Text>
        </div>
      )
    },
    {
      key: 'workflow',
      label: '审批流程',
      children: (
        <div style={{ textAlign: 'center', padding: '60px 20px' }}>
          <SafetyOutlined style={{ fontSize: '48px', color: '#d9d9d9', marginBottom: '16px' }} />
          <Title level={3} style={{ color: '#999', marginBottom: '8px' }}>
            审批流程
          </Title>
          <Text type="secondary">功能暂未实现</Text>
        </div>
      )
    }
  ];

  return (
    <div style={{ padding: '0 8px' }}>
      <div style={{ marginBottom: 24 }}>
        <Title level={2} style={{ margin: 0, color: '#1f2937' }}>
          ⚙️ 系统设置
        </Title>
        <Text type="secondary">管理用户、权限和系统配置</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="用户总数"
              value={totalUsers}
              prefix={<UserOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="活跃用户"
              value={activeUsers}
              prefix={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="角色数量"
              value={totalRoles}
              prefix={<TeamOutlined style={{ color: '#faad14' }} />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="权限数量"
              value={totalPermissions}
              prefix={<SecurityScanOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <Tabs defaultActiveKey="users" items={tabItems} />
      </Card>

      {/* 创建/编辑模态框 */}
      <Modal
        title={`${modalType === 'user' ? '新建用户' : modalType === 'role' ? '新建角色' : modalType === 'permission' ? '新建权限' : '编辑配置'}`}
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={() => {
          setIsModalVisible(false);
          form.resetFields();
        }}
        width={800}
      >
        <Form form={form} layout="vertical">
          {modalType === 'user' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="username" label="用户名" rules={[{ required: true }]}>
                    <Input placeholder="请输入用户名" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="realName" label="真实姓名" rules={[{ required: true }]}>
                    <Input placeholder="请输入真实姓名" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="email" label="邮箱" rules={[{ required: true, type: 'email' }]}>
                    <Input placeholder="请输入邮箱" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="phone" label="手机号" rules={[{ required: true }]}>
                    <Input placeholder="请输入手机号" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="department" label="部门" rules={[{ required: true }]}>
                    <Select placeholder="选择部门">
                      <Option value="信息技术部">信息技术部</Option>
                      <Option value="销售部">销售部</Option>
                      <Option value="采购部">采购部</Option>
                      <Option value="财务部">财务部</Option>
                      <Option value="人力资源部">人力资源部</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="role" label="角色" rules={[{ required: true }]}>
                    <Select placeholder="选择角色">
                      <Option value="超级管理员">超级管理员</Option>
                      <Option value="销售经理">销售经理</Option>
                      <Option value="采购员">采购员</Option>
                      <Option value="财务专员">财务专员</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="password" label="密码" rules={[{ required: true, min: 6 }]}>
                <Input.Password placeholder="请输入密码" />
              </Form.Item>
            </>
          )}
          {modalType === 'role' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="roleName" label="角色名称" rules={[{ required: true }]}>
                    <Input placeholder="请输入角色名称" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="roleCode" label="角色代码" rules={[{ required: true }]}>
                    <Input placeholder="请输入角色代码" />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="description" label="角色描述">
                <Input.TextArea placeholder="请输入角色描述" rows={3} />
              </Form.Item>
              <Form.Item name="permissions" label="权限配置">
                <Checkbox.Group>
                  <Row>
                    <Col span={8}>
                      <Checkbox value="user:read">用户查看</Checkbox>
                    </Col>
                    <Col span={8}>
                      <Checkbox value="user:write">用户编辑</Checkbox>
                    </Col>
                    <Col span={8}>
                      <Checkbox value="role:read">角色查看</Checkbox>
                    </Col>
                    <Col span={8}>
                      <Checkbox value="role:write">角色编辑</Checkbox>
                    </Col>
                    <Col span={8}>
                      <Checkbox value="sales:read">销售查看</Checkbox>
                    </Col>
                    <Col span={8}>
                      <Checkbox value="sales:write">销售编辑</Checkbox>
                    </Col>
                  </Row>
                </Checkbox.Group>
              </Form.Item>
            </>
          )}
          {modalType === 'permission' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="permissionName" label="权限名称" rules={[{ required: true }]}>
                    <Input placeholder="请输入权限名称" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="permissionCode" label="权限代码" rules={[{ required: true }]}>
                    <Input placeholder="请输入权限代码" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="module" label="所属模块" rules={[{ required: true }]}>
                    <Select placeholder="选择模块">
                      <Option value="系统管理">系统管理</Option>
                      <Option value="销售管理">销售管理</Option>
                      <Option value="采购管理">采购管理</Option>
                      <Option value="库存管理">库存管理</Option>
                      <Option value="财务管理">财务管理</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="type" label="权限类型" rules={[{ required: true }]}>
                    <Select placeholder="选择类型">
                      <Option value="menu">菜单</Option>
                      <Option value="button">按钮</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="description" label="权限描述">
                <Input.TextArea placeholder="请输入权限描述" rows={3} />
              </Form.Item>
            </>
          )}
          {modalType === 'config' && (
            <>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="configKey" label="配置键" rules={[{ required: true }]}>
                    <Input placeholder="请输入配置键" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="category" label="配置分类" rules={[{ required: true }]}>
                    <Select placeholder="选择分类">
                      <Option value="基础配置">基础配置</Option>
                      <Option value="邮件配置">邮件配置</Option>
                      <Option value="安全配置">安全配置</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="type" label="数据类型" rules={[{ required: true }]}>
                    <Select placeholder="选择类型">
                      <Option value="string">字符串</Option>
                      <Option value="number">数字</Option>
                      <Option value="boolean">布尔值</Option>
                      <Option value="file">文件</Option>
                      <Option value="time">时间</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="required" label="是否必填" valuePropName="checked">
                    <Switch />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="configValue" label="配置值" rules={[{ required: true }]}>
                <Input placeholder="请输入配置值" />
              </Form.Item>
              <Form.Item name="description" label="配置描述">
                <Input.TextArea placeholder="请输入配置描述" rows={3} />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>
    </div>
  );
}