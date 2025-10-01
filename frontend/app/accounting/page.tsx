'use client';

import React, { useState, useEffect } from 'react';
import {
  Card,
  Tabs,
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  DatePicker,
  InputNumber,
  Space,
  Tag,
  Statistic,
  Row,
  Col,
  message,
  Divider,
  Typography,
  Progress,
  Alert,
  Tooltip,
  Badge
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  FileTextOutlined,
  BarChartOutlined,
  DollarOutlined,
  CalculatorOutlined,
  GlobalOutlined,
  AuditOutlined,
  TrophyOutlined,
  RiseOutlined,
  FallOutlined
} from '@ant-design/icons';
import dayjs from 'dayjs';
import type { ColumnsType } from 'antd/es/table';
const { Option } = Select;
const { Title, Text } = Typography;
const { RangePicker } = DatePicker;

// 类型定义
interface Account {
  id: number;
  code: string;
  name: string;
  type: 'asset' | 'liability' | 'equity' | 'revenue' | 'expense';
  balance: number;
  parentId: number | null;
}

interface JournalEntry {
  id: number;
  date: string;
  reference: string;
  description: string;
  totalDebit: number;
  totalCredit: number;
  status: 'draft' | 'posted' | 'cancelled';
  entries: {
    accountCode: string;
    accountName: string;
    debit: number;
    credit: number;
  }[];
}

interface ExchangeRate {
  id: number;
  fromCurrency: string;
  toCurrency: string;
  rate: number;
  date: string;
  type: 'manual' | 'auto';
}

// 模拟数据
const mockAccounts: Account[] = [
  { id: 1, code: '1001', name: '库存现金', type: 'asset', balance: 50000, parentId: null },
  { id: 2, code: '1002', name: '银行存款', type: 'asset', balance: 1200000, parentId: null },
  { id: 3, code: '2001', name: '应付账款', type: 'liability', balance: 300000, parentId: null },
  { id: 4, code: '3001', name: '实收资本', type: 'equity', balance: 1000000, parentId: null },
  { id: 5, code: '6001', name: '主营业务收入', type: 'revenue', balance: 2000000, parentId: null },
];

const mockJournalEntries: JournalEntry[] = [
  {
    id: 1,
    date: '2024-01-15',
    reference: 'JE001',
    description: '销售商品收入',
    totalDebit: 100000,
    totalCredit: 100000,
    status: 'posted',
    entries: [
      { accountCode: '1002', accountName: '银行存款', debit: 100000, credit: 0 },
      { accountCode: '6001', accountName: '主营业务收入', debit: 0, credit: 100000 }
    ]
  },
  {
    id: 2,
    date: '2024-01-16',
    reference: 'JE002',
    description: '采购原材料',
    totalDebit: 50000,
    totalCredit: 50000,
    status: 'draft',
    entries: [
      { accountCode: '1401', accountName: '原材料', debit: 50000, credit: 0 },
      { accountCode: '2001', accountName: '应付账款', debit: 0, credit: 50000 }
    ]
  }
];

const mockExchangeRates: ExchangeRate[] = [
  { id: 1, fromCurrency: 'USD', toCurrency: 'CNY', rate: 7.2345, date: '2024-01-15', type: 'manual' },
  { id: 2, fromCurrency: 'EUR', toCurrency: 'CNY', rate: 7.8901, date: '2024-01-15', type: 'auto' },
  { id: 3, fromCurrency: 'JPY', toCurrency: 'CNY', rate: 0.0489, date: '2024-01-15', type: 'auto' },
];

export default function AccountingPage() {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [accounts, setAccounts] = useState<Account[]>(mockAccounts);
  const [journalEntries, setJournalEntries] = useState<JournalEntry[]>(mockJournalEntries);
  const [exchangeRates, setExchangeRates] = useState<ExchangeRate[]>(mockExchangeRates);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState('');
  const [editingRecord, setEditingRecord] = useState<any>(null);
  const [form] = Form.useForm();

  // 会计科目表列定义
  const accountColumns: ColumnsType<Account> = [
    {
      title: '科目编码',
      dataIndex: 'code',
      key: 'code',
      width: 120,
    },
    {
      title: '科目名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '科目类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: Account['type']) => {
        const typeMap: Record<Account['type'], { color: string; text: string }> = {
          asset: { color: 'blue', text: '资产' },
          liability: { color: 'red', text: '负债' },
          equity: { color: 'green', text: '所有者权益' },
          revenue: { color: 'orange', text: '收入' },
          expense: { color: 'purple', text: '费用' }
        };
        const config = typeMap[type];
        return <Tag color={config.color}>{config.text}</Tag>;
      }
    },
    {
      title: '余额',
      dataIndex: 'balance',
      key: 'balance',
      align: 'right',
      render: (balance: number) => (
        <Text strong style={{ color: balance >= 0 ? '#52c41a' : '#ff4d4f' }}>
          ¥{balance.toLocaleString()}
        </Text>
      )
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_, record: Account) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit('account', record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete('account', record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  // 凭证列定义
  const journalColumns: ColumnsType<JournalEntry> = [
    {
      title: '凭证号',
      dataIndex: 'reference',
      key: 'reference',
      width: 100,
    },
    {
      title: '日期',
      dataIndex: 'date',
      key: 'date',
      width: 120,
    },
    {
      title: '摘要',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '借方金额',
      dataIndex: 'totalDebit',
      key: 'totalDebit',
      align: 'right',
      render: (amount: number) => `¥${amount.toLocaleString()}`
    },
    {
      title: '贷方金额',
      dataIndex: 'totalCredit',
      key: 'totalCredit',
      align: 'right',
      render: (amount: number) => `¥${amount.toLocaleString()}`
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: JournalEntry['status']) => {
        const statusMap: Record<JournalEntry['status'], { color: string; text: string }> = {
          draft: { color: 'orange', text: '草稿' },
          posted: { color: 'green', text: '已过账' },
          cancelled: { color: 'red', text: '已取消' }
        };
        const config = statusMap[status];
        return <Badge status={config.color === 'green' ? 'success' : config.color === 'red' ? 'error' : 'processing'} text={config.text} />;
      }
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record: JournalEntry) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit('journal', record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            icon={<FileTextOutlined />}
            onClick={() => handleViewJournal(record)}
          >
            查看
          </Button>
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete('journal', record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  // 汇率列定义
  const exchangeRateColumns: ColumnsType<ExchangeRate> = [
    {
      title: '货币对',
      key: 'pair',
      render: (_, record: ExchangeRate) => `${record.fromCurrency}/${record.toCurrency}`
    },
    {
      title: '汇率',
      dataIndex: 'rate',
      key: 'rate',
      align: 'right',
      render: (rate: number) => rate.toFixed(4)
    },
    {
      title: '日期',
      dataIndex: 'date',
      key: 'date',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: ExchangeRate['type']) => (
        <Tag color={type === 'auto' ? 'blue' : 'green'}>
          {type === 'auto' ? '自动获取' : '手动录入'}
        </Tag>
      )
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_, record: ExchangeRate) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit('exchangeRate', record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete('exchangeRate', record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  const handleAdd = (type: string) => {
    setModalType(type);
    setEditingRecord(null);
    setIsModalVisible(true);
    form.resetFields();
  };

  const handleEdit = (type: string, record: any) => {
    setModalType(type);
    setEditingRecord(record);
    setIsModalVisible(true);
    form.setFieldsValue(record);
  };

  const handleDelete = (type: string, id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这条记录吗？',
      onOk: () => {
        if (type === 'account') {
          setAccounts(accounts.filter(item => item.id !== id));
        } else if (type === 'journal') {
          setJournalEntries(journalEntries.filter(item => item.id !== id));
        } else if (type === 'exchangeRate') {
          setExchangeRates(exchangeRates.filter(item => item.id !== id));
        }
        message.success('删除成功');
      }
    });
  };

  const handleViewJournal = (record: JournalEntry) => {
    Modal.info({
      title: `凭证详情 - ${record.reference}`,
      width: 800,
      content: (
        <div>
          <p><strong>日期：</strong>{record.date}</p>
          <p><strong>摘要：</strong>{record.description}</p>
          <Table
            dataSource={record.entries}
            pagination={false}
            size="small"
            columns={[
              { title: '科目编码', dataIndex: 'accountCode', key: 'accountCode' },
              { title: '科目名称', dataIndex: 'accountName', key: 'accountName' },
              { title: '借方', dataIndex: 'debit', key: 'debit', align: 'right' as const, render: (val: number) => val ? `¥${val.toLocaleString()}` : '-' },
              { title: '贷方', dataIndex: 'credit', key: 'credit', align: 'right' as const, render: (val: number) => val ? `¥${val.toLocaleString()}` : '-' }
            ]}
          />
        </div>
      )
    });
  };

  const handleModalOk = () => {
    form.validateFields().then(values => {
      if (modalType === 'account') {
        if (editingRecord) {
          setAccounts(accounts.map(item => 
            item.id === editingRecord.id ? { ...item, ...values } : item
          ));
        } else {
          const newAccount: Account = {
            id: Date.now(),
            ...values,
            balance: 0
          };
          setAccounts([...accounts, newAccount]);
        }
      } else if (modalType === 'exchangeRate') {
        if (editingRecord) {
          setExchangeRates(exchangeRates.map(item => 
            item.id === editingRecord.id ? { ...item, ...values } : item
          ));
        } else {
          const newRate: ExchangeRate = {
            id: Date.now(),
            ...values,
            date: values.date.format('YYYY-MM-DD'),
            type: 'manual' as const
          };
          setExchangeRates([...exchangeRates, newRate]);
        }
      }
      setIsModalVisible(false);
      message.success(editingRecord ? '更新成功' : '添加成功');
    });
  };

  // 仪表板内容
  const renderDashboard = () => (
    <div>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总资产"
              value={1250000}
              precision={0}
              valueStyle={{ color: '#3f8600' }}
              prefix={<RiseOutlined />}
              suffix="¥"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="总负债"
              value={300000}
              precision={0}
              valueStyle={{ color: '#cf1322' }}
              prefix={<FallOutlined />}
              suffix="¥"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="所有者权益"
              value={950000}
              precision={0}
              valueStyle={{ color: '#1890ff' }}
              prefix={<TrophyOutlined />}
              suffix="¥"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="本月收入"
              value={200000}
              precision={0}
              valueStyle={{ color: '#722ed1' }}
              prefix={<DollarOutlined />}
              suffix="¥"
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Card title="资产负债率" extra={<Tooltip title="负债总额/资产总额"><Button type="link" icon={<CalculatorOutlined />} /></Tooltip>}>
            <Progress
              percent={24}
              status="active"
              strokeColor={{
                '0%': '#108ee9',
                '100%': '#87d068',
              }}
            />
            <Text type="secondary">当前资产负债率为 24%，财务状况良好</Text>
          </Card>
        </Col>
        <Col span={12}>
          <Card title="近期财务动态">
            <Alert
              message="财务提醒"
              description="本月应收账款回收率达到 95%，现金流状况良好。建议关注下月的应付账款到期情况。"
              type="info"
              showIcon
              style={{ marginBottom: 16 }}
            />
            <Alert
              message="税务提醒"
              description="下月15日为增值税申报截止日期，请及时准备相关资料。"
              type="warning"
              showIcon
            />
          </Card>
        </Col>
      </Row>
    </div>
  );

  // 财务报表内容
  const renderReports = () => (
    <div>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col span={24}>
          <Card>
            <Space>
              <Text strong>报表期间：</Text>
              <RangePicker defaultValue={[dayjs().startOf('month'), dayjs()]} />
              <Button type="primary" icon={<BarChartOutlined />}>生成报表</Button>
            </Space>
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={8}>
          <Card title="资产负债表" extra={<Button type="link">查看详情</Button>}>
            <div style={{ textAlign: 'center' }}>
              <Title level={4}>¥1,250,000</Title>
              <Text type="secondary">总资产</Text>
            </div>
            <Divider />
            <Row>
              <Col span={12}>
                <Text>流动资产</Text>
                <br />
                <Text strong>¥950,000</Text>
              </Col>
              <Col span={12}>
                <Text>非流动资产</Text>
                <br />
                <Text strong>¥300,000</Text>
              </Col>
            </Row>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="利润表" extra={<Button type="link">查看详情</Button>}>
            <div style={{ textAlign: 'center' }}>
              <Title level={4} style={{ color: '#52c41a' }}>¥150,000</Title>
              <Text type="secondary">净利润</Text>
            </div>
            <Divider />
            <Row>
              <Col span={12}>
                <Text>营业收入</Text>
                <br />
                <Text strong>¥2,000,000</Text>
              </Col>
              <Col span={12}>
                <Text>营业成本</Text>
                <br />
                <Text strong>¥1,500,000</Text>
              </Col>
            </Row>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="现金流量表" extra={<Button type="link">查看详情</Button>}>
            <div style={{ textAlign: 'center' }}>
              <Title level={4} style={{ color: '#1890ff' }}>¥80,000</Title>
              <Text type="secondary">经营活动现金流</Text>
            </div>
            <Divider />
            <Row>
              <Col span={12}>
                <Text>现金流入</Text>
                <br />
                <Text strong>¥1,800,000</Text>
              </Col>
              <Col span={12}>
                <Text>现金流出</Text>
                <br />
                <Text strong>¥1,720,000</Text>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  );

  // 税务合规内容
  const renderTaxCompliance = () => (
    <div>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col span={24}>
          <Alert
            message="税务合规状态"
            description="当前税务申报状态正常，所有税种均已按时申报。下次申报日期：2024年2月15日"
            type="success"
            showIcon
            action={
              <Button size="small" type="primary">
                查看详情
              </Button>
            }
          />
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={8}>
          <Card title="增值税" extra={<Tag color="green">已申报</Tag>}>
            <Statistic
              title="本期应纳税额"
              value={25000}
              precision={2}
              suffix="¥"
            />
            <div style={{ marginTop: 16 }}>
              <Text type="secondary">申报期限：每月15日前</Text>
              <br />
              <Text type="secondary">下次申报：2024-02-15</Text>
            </div>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="企业所得税" extra={<Tag color="orange">待申报</Tag>}>
            <Statistic
              title="预计应纳税额"
              value={37500}
              precision={2}
              suffix="¥"
            />
            <div style={{ marginTop: 16 }}>
              <Text type="secondary">申报期限：季度后15日内</Text>
              <br />
              <Text type="secondary">下次申报：2024-04-15</Text>
            </div>
          </Card>
        </Col>
        <Col span={8}>
          <Card title="个人所得税" extra={<Tag color="green">已申报</Tag>}>
            <Statistic
              title="本期代扣代缴"
              value={12000}
              precision={2}
              suffix="¥"
            />
            <div style={{ marginTop: 16 }}>
              <Text type="secondary">申报期限：每月15日前</Text>
              <br />
              <Text type="secondary">下次申报：2024-02-15</Text>
            </div>
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col span={24}>
          <Card title="税务计算器" extra={<Button type="primary" icon={<CalculatorOutlined />}>开始计算</Button>}>
            <Text type="secondary">
              提供增值税、企业所得税、个人所得税等税种的快速计算功能，帮助您准确计算应纳税额。
            </Text>
          </Card>
        </Col>
      </Row>
    </div>
  );

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={2} style={{ margin: 0 }}>
          💼 财务管理
        </Title>
      </div>

      <Tabs 
        activeKey={activeTab} 
        onChange={setActiveTab}
        items={[
          {
            key: 'dashboard',
            label: <span><BarChartOutlined />财务仪表板</span>,
            children: renderDashboard()
          },
          {
            key: 'accounts',
            label: <span><FileTextOutlined />会计科目</span>,
            children: (
              <Card
                title="会计科目表"
                extra={
                  <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => handleAdd('account')}
                  >
                    新增科目
                  </Button>
                }
              >
                <Table
                  dataSource={accounts}
                  columns={accountColumns}
                  rowKey="id"
                  pagination={{ pageSize: 10 }}
                />
              </Card>
            )
          },
          {
            key: 'journals',
            label: <span><EditOutlined />凭证管理</span>,
            children: (
              <Card
                title="记账凭证"
                extra={
                  <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => handleAdd('journal')}
                  >
                    新增凭证
                  </Button>
                }
              >
                <Table
                  dataSource={journalEntries}
                  columns={journalColumns}
                  rowKey="id"
                  pagination={{ pageSize: 10 }}
                />
              </Card>
            )
          },
          {
            key: 'reports',
            label: <span><BarChartOutlined />财务报表</span>,
            children: renderReports()
          },
          {
            key: 'exchange',
            label: <span><GlobalOutlined />汇率管理</span>,
            children: (
              <Card
                title="汇率管理"
                extra={
                  <Space>
                    <Button icon={<RiseOutlined />}>批量更新</Button>
                    <Button
                      type="primary"
                      icon={<PlusOutlined />}
                      onClick={() => handleAdd('exchangeRate')}
                    >
                      新增汇率
                    </Button>
                  </Space>
                }
              >
                <Table
                  dataSource={exchangeRates}
                  columns={exchangeRateColumns}
                  rowKey="id"
                  pagination={{ pageSize: 10 }}
                />
              </Card>
            )
          },
          {
            key: 'tax',
            label: <span><AuditOutlined />税务合规</span>,
            children: renderTaxCompliance()
          }
        ]}
      />

      {/* 通用模态框 */}
      <Modal
        title={
          modalType === 'account' ? (editingRecord ? '编辑科目' : '新增科目') :
          modalType === 'exchangeRate' ? (editingRecord ? '编辑汇率' : '新增汇率') :
          '编辑'
        }
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={() => setIsModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          {modalType === 'account' && (
            <>
              <Form.Item
                name="code"
                label="科目编码"
                rules={[{ required: true, message: '请输入科目编码' }]}
              >
                <Input placeholder="如：1001" />
              </Form.Item>
              <Form.Item
                name="name"
                label="科目名称"
                rules={[{ required: true, message: '请输入科目名称' }]}
              >
                <Input placeholder="如：库存现金" />
              </Form.Item>
              <Form.Item
                name="type"
                label="科目类型"
                rules={[{ required: true, message: '请选择科目类型' }]}
              >
                <Select placeholder="请选择科目类型">
                  <Option value="asset">资产</Option>
                  <Option value="liability">负债</Option>
                  <Option value="equity">所有者权益</Option>
                  <Option value="revenue">收入</Option>
                  <Option value="expense">费用</Option>
                </Select>
              </Form.Item>
            </>
          )}
          
          {modalType === 'exchangeRate' && (
            <>
              <Form.Item
                name="fromCurrency"
                label="基础货币"
                rules={[{ required: true, message: '请选择基础货币' }]}
              >
                <Select placeholder="请选择基础货币">
                  <Option value="USD">美元 (USD)</Option>
                  <Option value="EUR">欧元 (EUR)</Option>
                  <Option value="JPY">日元 (JPY)</Option>
                  <Option value="GBP">英镑 (GBP)</Option>
                </Select>
              </Form.Item>
              <Form.Item
                name="toCurrency"
                label="目标货币"
                rules={[{ required: true, message: '请选择目标货币' }]}
              >
                <Select placeholder="请选择目标货币">
                  <Option value="CNY">人民币 (CNY)</Option>
                  <Option value="USD">美元 (USD)</Option>
                  <Option value="EUR">欧元 (EUR)</Option>
                </Select>
              </Form.Item>
              <Form.Item
                name="rate"
                label="汇率"
                rules={[{ required: true, message: '请输入汇率' }]}
              >
                <InputNumber
                  placeholder="请输入汇率"
                  min={0}
                  step={0.0001}
                  precision={4}
                  style={{ width: '100%' }}
                />
              </Form.Item>
              <Form.Item
                name="date"
                label="日期"
                rules={[{ required: true, message: '请选择日期' }]}
              >
                <DatePicker style={{ width: '100%' }} />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>
    </div>
  );
}