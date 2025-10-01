// 通用 API 响应类型
export interface ApiResponse<T = any> {
  data: T;
  message?: string;
  success?: boolean;
}

// 分页响应类型
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// 认证相关类型
export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  phone?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface User {
  id: number;
  email: string;
  firstName: string;
  lastName: string;
  phone?: string;
  createdAt: string;
  updatedAt: string;
}

// 客户相关类型
export interface Customer {
  id: number;
  code: string;
  name: string;
  email: string;
  phone: string;
  address: string;
  contactPerson: string;
  creditLimit: number;
  paymentTerms: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateCustomerRequest {
  code: string;
  name: string;
  email: string;
  phone: string;
  address: string;
  contactPerson: string;
  creditLimit: number;
  paymentTerms: string;
  status: string;
}

// 供应商相关类型
export interface Supplier {
  id: number;
  code: string;
  name: string;
  contactName: string;
  email: string;
  phone: string;
  address: string;
  taxNumber: string;
  paymentTerms: string;
  creditLimit: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateSupplierRequest {
  code: string;
  name: string;
  contactName?: string;
  email?: string;
  phone?: string;
  address?: string;
  taxNumber?: string;
  paymentTerms?: string;
  creditLimit?: number;
}

// 采购请求相关类型
export interface PurchaseRequest {
  id: number;
  number: string;
  title: string;
  description?: string;
  priority: string;
  status: string;
  requiredDate: string;
  totalAmount: number;
  items: PurchaseRequestItem[];
  createdBy: User;
  approvedBy?: User;
  createdAt: string;
  updatedAt: string;
}

export interface PurchaseRequestItem {
  id: number;
  quantity: number;
  unitPrice: number;
  amount: number;
  notes?: string;
  item: Product;
}

export interface CreatePurchaseRequestRequest {
  title: string;
  description?: string;
  priority: string;
  requiredDate: string;
  items: CreatePurchaseRequestItemRequest[];
}

export interface CreatePurchaseRequestItemRequest {
  itemId: number;
  quantity: number;
  unitPrice?: number;
  notes?: string;
}

// 采购订单相关类型
export interface PurchaseOrder {
  id: number;
  orderNumber: string;
  supplierId: number;
  orderDate: string;
  expectedDate: string;
  deliveryDate?: string;
  status: string;
  currency: string;
  exchangeRate: number;
  paymentTerms?: string;
  deliveryAddress?: string;
  billingAddress?: string;
  terms?: string;
  notes?: string;
  purchaseRequestId?: number;
  subTotal: number;
  totalDiscount: number;
  totalTax: number;
  totalAmount: number;
  supplier?: Supplier;
  request?: PurchaseRequest;
  items: PurchaseOrderItem[];
  createdBy?: User;
  approvedBy?: User;
  createdAt: string;
  updatedAt: string;
}

export interface PurchaseOrderItem {
  id: number;
  quantity: number;
  unitPrice: number;
  taxRate: number;
  taxAmount: number;
  amount: number;
  receivedQty: number;
  notes?: string;
  item: Product;
}

export interface CreatePurchaseOrderRequest {
  supplierId: number;
  requestId?: number;
  orderDate: string;
  expectedDate: string;
  deliveryDate?: string;
  status?: string;
  paymentTerms?: string;
  terms?: string;
  notes?: string;
  items: CreatePurchaseOrderItemRequest[];
}

export interface CreatePurchaseOrderItemRequest {
  itemId: number;
  quantity: number;
  unitPrice: number;
  taxRate?: number;
  notes?: string;
}

// 产品相关类型
export interface Product {
  id: number;
  code: string;
  name: string;
  description: string;
  category: string;
  unit: string;
  price: number;
  cost: number;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateProductRequest {
  code: string;
  name: string;
  description: string;
  category: string;
  unit: string;
  price: number;
  cost: number;
  status: string;
}

// 仓库相关类型
export interface Warehouse {
  id: number;
  code: string;
  name: string;
  address: string;
  manager: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateWarehouseRequest {
  code: string;
  name: string;
  address: string;
  manager: string;
  status: string;
}

// 库存相关类型
export interface Stock {
  id: number;
  itemId: number;
  itemCode: string;
  itemName: string;
  warehouseId: number;
  warehouseCode: string;
  warehouseName: string;
  quantity: number;
  reservedQuantity: number;
  availableQuantity: number;
  unitCost: number;
  totalValue: number;
  reorderLevel: number;
  maxLevel: number;
  lastUpdated: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateStockRequest {
  itemId: number;
  warehouseId: number;
  quantity: number;
  unitCost: number;
  reorderLevel?: number;
  maxLevel?: number;
}

export interface StockMovement {
  id: number;
  stockId: number;
  itemId: number;
  itemCode: string;
  itemName: string;
  warehouseId: number;
  warehouseCode: string;
  warehouseName: string;
  movementType: string; // 'IN' | 'OUT' | 'TRANSFER' | 'ADJUSTMENT'
  quantity: number;
  unitCost: number;
  totalValue: number;
  referenceType: string; // 'PURCHASE' | 'SALES' | 'TRANSFER' | 'ADJUSTMENT'
  referenceId?: number;
  notes: string;
  createdBy: number;
  createdAt: string;
}

export interface CreateStockMovementRequest {
  stockId: number;
  movementType: string;
  quantity: number;
  unitCost: number;
  referenceType: string;
  referenceId?: number;
  notes?: string;
}

export interface InventoryReport {
  itemId: number;
  itemCode: string;
  itemName: string;
  category: string;
  totalQuantity: number;
  totalValue: number;
  warehouses: {
    warehouseId: number;
    warehouseName: string;
    quantity: number;
    value: number;
  }[];
}

export interface InventoryAnalytics {
  totalItems: number;
  totalValue: number;
  lowStockItems: number;
  outOfStockItems: number;
  topMovingItems: {
    itemId: number;
    itemName: string;
    movementCount: number;
  }[];
  categoryDistribution: {
    category: string;
    itemCount: number;
    totalValue: number;
  }[];
}

// 员工相关类型
export interface Employee {
  id: number;
  code: string;
  firstName: string;
  lastName: string;
  fullName: string;
  email: string;
  phone: string;
  dateOfBirth: string;
  gender: string;
  hireDate: string;
  departmentId?: number;
  positionId?: number;
  managerId?: number;
  status: string;
  emergencyContact: string;
  idNumber: string;
  address: string;
  bankAccount: string;
  createdAt: string;
  updatedAt: string;
}

// 报价相关类型
export interface Quotation {
  id: number;
  quotationNumber: string;
  customerId: number;
  customer: Customer;
  date: string;
  validTill: string;
  status: string;
  subject: string;
  items: QuotationItem[];
  totalAmount: number;
  discountAmount: number;
  taxAmount: number;
  grandTotal: number;
  terms: string;
  notes: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface QuotationItem {
  id: number;
  quotationId: number;
  itemId: number;
  itemName: string;
  description: string;
  quantity: number;
  unitPrice: number;
  discountPercent: number;
  discountAmount: number;
  taxPercent: number;
  taxAmount: number;
  totalAmount: number;
}

export interface CreateQuotationRequest {
  customerId: number;
  date: string;
  validTill: string;
  subject: string;
  items: CreateQuotationItemRequest[];
  discountAmount?: number;
  taxAmount?: number;
  terms?: string;
  notes?: string;
}

export interface CreateQuotationItemRequest {
  itemId: number;
  itemName: string;
  description?: string;
  quantity: number;
  unitPrice: number;
  discountPercent?: number;
  taxPercent?: number;
}

// 销售订单相关类型
export interface SalesOrder {
  id: number;
  orderNumber: string;
  customerId: number;
  customer: Customer;
  quotationId?: number;
  quotation?: Quotation;
  orderDate: string;
  deliveryDate: string;
  status: string;
  priority: string;
  items: SalesOrderItem[];
  totalAmount: number;
  discountAmount: number;
  taxAmount: number;
  grandTotal: number;
  terms: string;
  notes: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface SalesOrderItem {
  id: number;
  salesOrderId: number;
  itemId: number;
  itemName: string;
  description: string;
  quantity: number;
  unitPrice: number;
  discountPercent: number;
  discountAmount: number;
  taxPercent: number;
  taxAmount: number;
  totalAmount: number;
  deliveredQuantity: number;
}

export interface CreateSalesOrderRequest {
  customerId: number;
  quotationId?: number;
  orderDate: string;
  deliveryDate: string;
  priority?: string;
  items: CreateSalesOrderItemRequest[];
  discountAmount?: number;
  taxAmount?: number;
  terms?: string;
  notes?: string;
}

export interface CreateSalesOrderItemRequest {
  itemId: number;
  itemName: string;
  description?: string;
  quantity: number;
  unitPrice: number;
  discountPercent?: number;
  taxPercent?: number;
}

export interface CreateEmployeeRequest {
  code: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  dateOfBirth: string;
  gender: string;
  hireDate: string;
  departmentId?: number;
  positionId?: number;
  managerId?: number;
  status: string;
  emergencyContact: string;
  idNumber: string;
  address: string;
  bankAccount: string;
}

// 账户相关类型
export interface Account {
  id: number;
  code: string;
  name: string;
  type: string;
  parentId?: number;
  balance: number;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAccountRequest {
  code: string;
  name: string;
  type: string;
  parentId?: number;
  status: string;
}

// 错误类型
export interface ApiError {
  message: string;
  code?: string;
  details?: any;
}