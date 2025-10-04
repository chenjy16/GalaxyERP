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
  Calendar,
  TimePicker
} from 'antd';
import { 
  PlusOutlined, 
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  UserOutlined,
  EyeOutlined,
  MoreOutlined,
  TeamOutlined,
  ClockCircleOutlined,
  DollarOutlined,
  CalendarOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
  ExportOutlined,
  ImportOutlined,
  PhoneOutlined,
  MailOutlined,
  HomeOutlined,
  IdcardOutlined,
  CalculatorOutlined
} from '@ant-design/icons';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Option } = Select;
const { RangePicker } = DatePicker;

export default function HRPage() {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'employee' | 'attendance' | 'salary'>('employee');
  const [form] = Form.useForm();

  // 模拟员工数据
  const employeeData = [
    {
      key: '1',
      id: 'EMP001',
      name: '张三',
      department: '技术部',
      position: '高级工程师',
      phone: '13800138001',
      email: 'zhangsan@company.com',
      hireDate: '2023-01-15',
      salary: 15000,
      status: 'active',
      avatar: null,
      address: '深圳市南山区科技园',
      emergencyContact: '李四 - 13900139001'
    },
    {
      key: '2',
      id: 'EMP002',
      name: '李四',
      department: '销售部',
      position: '销售经理',
      phone: '13800138002',
      email: 'lisi@company.com',
      hireDate: '2022-08-20',
      salary: 12000,
      status: 'active',
      avatar: null,
      address: '深圳市福田区中心区',
      emergencyContact: '王五 - 13900139002'
    },
    {
      key: '3',
      id: 'EMP003',
      name: '王五',
      department: '财务部',
      position: '财务专员',
      phone: '13800138003',
      email: 'wangwu@company.com',
      hireDate: '2023-03-10',
      salary: 8000,
      status: 'active',
      avatar: null,
      address: '深圳市宝安区新安街道',
      emergencyContact: '赵六 - 13900139003'
    },
    {
      key: '4',
      id: 'EMP004',
      name: '赵六',
      department: '人事部',
      position: '人事主管',
      phone: '13800138004',
      email: 'zhaoliu@company.com',
      hireDate: '2021-12-01',
      salary: 10000,
      status: 'leave',
      avatar: null,
      address: '深圳市龙岗区布吉街道',
      emergencyContact: '孙七 - 13900139004'
    }
  ];

  // 模拟考勤数据
  const attendanceData = [
    {
      key: '1',
      employeeId: 'EMP001',
      employeeName: '张三',
      date: '2025-01-10',
      checkIn: '09:00',
      checkOut: '18:30',
      workHours: 8.5,
      overtime: 0.5,
      status: 'normal',
      location: '公司总部'
    },
    {
      key: '2',
      employeeId: 'EMP002',
      employeeName: '李四',
      date: '2025-01-10',
      checkIn: '08:45',
      checkOut: '17:45',
      workHours: 8,
      overtime: 0,
      status: 'normal',
      location: '公司总部'
    },
    {
      key: '3',
      employeeId: 'EMP003',
      employeeName: '王五',
      date: '2025-01-10',
      checkIn: '09:15',
      checkOut: '18:00',
      workHours: 7.75,
      overtime: 0,
      status: 'late',
      location: '公司总部'
    },
    {
      key: '4',
      employeeId: 'EMP001',
      employeeName: '张三',
      date: '2025-01-09',
      checkIn: null,
      checkOut: null,
      workHours: 0,
      overtime: 0,
      status: 'absent',
      location: null
    }
  ];

  // 模拟薪资数据
  const salaryData = [
    {
      key: '1',
      employeeId: 'EMP001',
      employeeName: '张三',
      month: '2025-01',
      baseSalary: 15000,
      overtime: 500,
      bonus: 2000,
      deduction: 200,
      insurance: 1500,
      tax: 1800,
      netSalary: 14000,
      status: 'paid'
    },
    {
      key: '2',
      employeeId: 'EMP002',
      employeeName: '李四',
      month: '2025-01',
      baseSalary: 12000,
      overtime: 300,
      bonus: 1500,
      deduction: 100,
      insurance: 1200,
      tax: 1200,
      netSalary: 11300,
      status: 'pending'
    },
    {
      key: '3',
      employeeId: 'EMP003',
      employeeName: '王五',
      month: '2025-01',
      baseSalary: 8000,
      overtime: 0,
      bonus: 800,
      deduction: 0,
      insurance: 800,
      tax: 600,
      netSalary: 7400,
      status: 'pending'
    }
  ];

  const getStatusColor = (status: string) => {
    const colors: { [key: string]: string } = {
      active: 'green',
      leave: 'orange',
      resigned: 'red',
      normal: 'green',
      late: 'orange',
      absent: 'red',
      paid: 'green',
      pending: 'orange',
      cancelled: 'red'
    };
    return colors[status] || 'default';
  };

  const getStatusText = (status: string, type: string) => {
    if (type === 'employee') {
      const texts: { [key: string]: string } = {
        active: '在职',
        leave: '请假',
        resigned: '离职'
      };
      return texts[status] || status;
    } else if (type === 'attendance') {
      const texts: { [key: string]: string } = {
        normal: '正常',
        late: '迟到',
        absent: '缺勤'
      };
      return texts[status] || status;
    } else if (type === 'salary') {
      const texts: { [key: string]: string } = {
        paid: '已发放',
        pending: '待发放',
        cancelled: '已取消'
      };
      return texts[status] || status;
    }
    return status;
  };

  const employeeColumns = [
    {
      title: '员工信息',
      key: 'employee',
      render: (record: any) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <Avatar 
            size={40}
            style={{ backgroundColor: '#1890ff' }}
            icon={<UserOutlined />}
          />
          <div>
            <Text strong style={{ display: 'block' }}>{record.name}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {record.id} | {record.position}
            </Text>
          </div>
        </div>
      ),
    },
    {
      title: '部门',
      dataIndex: 'department',
      key: 'department',
      render: (department: string) => <Tag color="blue">{department}</Tag>,
    },
    {
      title: '联系方式',
      key: 'contact',
      render: (record: any) => (
        <div>
          <div style={{ marginBottom: 4 }}>
            <PhoneOutlined style={{ marginRight: 4, color: '#666' }} />
            <Text>{record.phone}</Text>
          </div>
          <div>
            <MailOutlined style={{ marginRight: 4, color: '#666' }} />
            <Text type="secondary" style={{ fontSize: 12 }}>{record.email}</Text>
          </div>
        </div>
      ),
    },
    {
      title: '入职日期',
      dataIndex: 'hireDate',
      key: 'hireDate',
    },
    {
      title: '薪资',
      dataIndex: 'salary',
      key: 'salary',
      render: (salary: number) => (
        <Text strong style={{ color: '#52c41a' }}>¥{salary ? salary.toLocaleString() : '0.00'}</Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>
          {getStatusText(status, 'employee')}
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
                key: 'attendance',
                label: '考勤记录',
                icon: <ClockCircleOutlined />,
              },
              {
                key: 'salary',
                label: '薪资记录',
                icon: <DollarOutlined />,
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

  const attendanceColumns = [
    {
      title: '员工',
      key: 'employee',
      render: (record: any) => (
        <div>
          <Text strong>{record.employeeName}</Text>
          <br />
          <Text type="secondary" style={{ fontSize: 12 }}>{record.employeeId}</Text>
        </div>
      ),
    },
    {
      title: '日期',
      dataIndex: 'date',
      key: 'date',
    },
    {
      title: '上班时间',
      dataIndex: 'checkIn',
      key: 'checkIn',
      render: (time: string) => time || '-',
    },
    {
      title: '下班时间',
      dataIndex: 'checkOut',
      key: 'checkOut',
      render: (time: string) => time || '-',
    },
    {
      title: '工作时长',
      dataIndex: 'workHours',
      key: 'workHours',
      render: (hours: number) => `${hours}小时`,
    },
    {
      title: '加班时长',
      dataIndex: 'overtime',
      key: 'overtime',
      render: (hours: number) => hours > 0 ? `${hours}小时` : '-',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag 
          color={getStatusColor(status)}
          icon={
            status === 'normal' ? <CheckCircleOutlined /> :
            status === 'late' ? <ExclamationCircleOutlined /> :
            <CloseCircleOutlined />
          }
        >
          {getStatusText(status, 'attendance')}
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
          <Tooltip title="编辑">
            <Button type="text" icon={<EditOutlined />} size="small" />
          </Tooltip>
        </Space>
      ),
    },
  ];

  const salaryColumns = [
    {
      title: '员工',
      key: 'employee',
      render: (record: any) => (
        <div>
          <Text strong>{record.employeeName}</Text>
          <br />
          <Text type="secondary" style={{ fontSize: 12 }}>{record.employeeId}</Text>
        </div>
      ),
    },
    {
      title: '月份',
      dataIndex: 'month',
      key: 'month',
    },
    {
      title: '基本工资',
      dataIndex: 'baseSalary',
      key: 'baseSalary',
      render: (amount: number) => `¥${amount ? amount.toLocaleString() : '0.00'}`,
    },
    {
      title: '加班费',
      dataIndex: 'overtime',
      key: 'overtime',
      render: (amount: number) => amount > 0 ? `¥${amount.toLocaleString()}` : '-',
    },
    {
      title: '奖金',
      dataIndex: 'bonus',
      key: 'bonus',
      render: (amount: number) => amount > 0 ? `¥${amount.toLocaleString()}` : '-',
    },
    {
      title: '扣款',
      dataIndex: 'deduction',
      key: 'deduction',
      render: (amount: number) => amount > 0 ? `-¥${amount.toLocaleString()}` : '-',
    },
    {
      title: '实发工资',
      dataIndex: 'netSalary',
      key: 'netSalary',
      render: (amount: number) => (
        <Text strong style={{ color: '#52c41a' }}>¥{amount ? amount.toLocaleString() : '0.00'}</Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>
          {getStatusText(status, 'salary')}
        </Tag>
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
          {record.status === 'pending' && (
            <Tooltip title="发放">
              <Button type="text" icon={<CheckCircleOutlined />} size="small" style={{ color: '#52c41a' }} />
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  const handleModalOk = () => {
    form.validateFields().then(values => {
      // 表单提交处理
      message.success(`${modalType === 'employee' ? '员工' : modalType === 'attendance' ? '考勤记录' : '薪资记录'}创建成功！`);
      setIsModalVisible(false);
      form.resetFields();
    });
  };

  const showModal = (type: 'employee' | 'attendance' | 'salary') => {
    setModalType(type);
    setIsModalVisible(true);
  };

  // 计算统计数据
  const totalEmployees = employeeData.length;
  const activeEmployees = employeeData.filter(e => e.status === '在职').length;
  const totalAttendance = attendanceData.length;
  const totalSalaries = salaryData.length;

  // 定义Tabs的items
  const tabItems = [
    {
      key: 'employees',
      label: '员工管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索员工..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="部门筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="技术部">技术部</Option>
                <Option value="销售部">销售部</Option>
                <Option value="财务部">财务部</Option>
                <Option value="人事部">人事部</Option>
              </Select>
              <Select placeholder="职位筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="经理">经理</Option>
                <Option value="主管">主管</Option>
                <Option value="专员">专员</Option>
                <Option value="助理">助理</Option>
              </Select>
              <Select placeholder="状态筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="在职">在职</Option>
                <Option value="离职">离职</Option>
                <Option value="试用期">试用期</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ImportOutlined />}>导入</Button>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('employee')}
              >
                新建员工
              </Button>
            </Space>
          </div>
          <Table 
            columns={employeeColumns} 
            dataSource={employeeData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
            scroll={{ x: 1200 }}
          />
        </>
      )
    },
    {
      key: 'attendance',
      label: '考勤管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索员工..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="部门筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="技术部">技术部</Option>
                <Option value="销售部">销售部</Option>
                <Option value="财务部">财务部</Option>
                <Option value="人事部">人事部</Option>
              </Select>
              <DatePicker placeholder="考勤日期" style={{ width: 150 }} />
              <Select placeholder="考勤状态" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="正常">正常</Option>
                <Option value="迟到">迟到</Option>
                <Option value="早退">早退</Option>
                <Option value="缺勤">缺勤</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('attendance')}
              >
                录入考勤
              </Button>
            </Space>
          </div>
          <Table 
            columns={attendanceColumns} 
            dataSource={attendanceData}
            pagination={{ pageSize: 10, showSizeChanger: true }}
            scroll={{ x: 1200 }}
          />
        </>
      )
    },
    {
      key: 'salary',
      label: '薪资管理',
      children: (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Space>
              <Input
                placeholder="搜索员工..."
                prefix={<SearchOutlined />}
                style={{ width: 300 }}
              />
              <Select placeholder="部门筛选" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="技术部">技术部</Option>
                <Option value="销售部">销售部</Option>
                <Option value="财务部">财务部</Option>
                <Option value="人事部">人事部</Option>
              </Select>
              <DatePicker placeholder="薪资月份" picker="month" style={{ width: 150 }} />
              <Select placeholder="发放状态" style={{ width: 120 }}>
                <Option value="all">全部</Option>
                <Option value="已发放">已发放</Option>
                <Option value="待发放">待发放</Option>
                <Option value="已暂停">已暂停</Option>
              </Select>
            </Space>
            <Space>
              <Button icon={<CalculatorOutlined />}>批量计算</Button>
              <Button icon={<ExportOutlined />}>导出</Button>
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => showModal('salary')}
              >
                新建薪资
              </Button>
            </Space>
          </div>
          <Table 
            columns={salaryColumns} 
            dataSource={salaryData}
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
          👥 人力资源
        </Title>
        <Text type="secondary">管理员工、考勤和薪资</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="员工总数"
              value={totalEmployees}
              prefix={<UserOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="在职员工"
              value={activeEmployees}
              prefix={<TeamOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="考勤记录"
              value={totalAttendance}
              prefix={<ClockCircleOutlined style={{ color: '#faad14' }} />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={6}>
          <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
            <Statistic
              title="薪资记录"
              value={totalSalaries}
              prefix={<DollarOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Card style={{ borderRadius: 12, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <Tabs defaultActiveKey="employees" items={tabItems} />
      </Card>
    </div>
  );
}