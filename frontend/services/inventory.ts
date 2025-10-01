import { apiClient } from '@/lib/api';
import { Stock, CreateStockRequest, StockMovement, CreateStockMovementRequest, PaginatedResponse } from '@/types/api';

export class InventoryService {
  // 获取库存列表
  static async getStocks(params?: {
    page?: number;
    limit?: number;
    search?: string;
    warehouseId?: number;
    itemId?: number;
    lowStock?: boolean;
  }): Promise<PaginatedResponse<Stock>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.warehouseId) searchParams.append('warehouse_id', params.warehouseId.toString());
    if (params?.itemId) searchParams.append('item_id', params.itemId.toString());
    if (params?.lowStock) searchParams.append('low_stock', params.lowStock.toString());
    
    const query = searchParams.toString();
    const endpoint = query ? `/stocks/?${query}` : '/stocks/';
    
    return apiClient.get<PaginatedResponse<Stock>>(endpoint);
  }

  // 获取单个库存记录
  static async getStock(id: number): Promise<Stock> {
    return apiClient.get<Stock>(`/stocks/${id}`);
  }

  // 创建或更新库存
  static async createOrUpdateStock(stockData: CreateStockRequest): Promise<Stock> {
    return apiClient.post<Stock>('/stocks/', stockData);
  }

  // 更新库存
  static async updateStock(id: number, stockData: Partial<CreateStockRequest>): Promise<Stock> {
    return apiClient.put<Stock>(`/stocks/${id}`, stockData);
  }

  // 删除库存记录
  static async deleteStock(id: number): Promise<void> {
    return apiClient.delete<void>(`/stocks/${id}`);
  }

  // 获取库存移动记录
  static async getStockMovements(params?: {
    page?: number;
    limit?: number;
    search?: string;
    warehouseId?: number;
    itemId?: number;
    movementType?: string;
    startDate?: string;
    endDate?: string;
  }): Promise<PaginatedResponse<StockMovement>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.warehouseId) searchParams.append('warehouse_id', params.warehouseId.toString());
    if (params?.itemId) searchParams.append('item_id', params.itemId.toString());
    if (params?.movementType) searchParams.append('movement_type', params.movementType);
    if (params?.startDate) searchParams.append('start_date', params.startDate);
    if (params?.endDate) searchParams.append('end_date', params.endDate);
    
    const query = searchParams.toString();
    const endpoint = query ? `/stock-movements/?${query}` : '/stock-movements/';
    
    return apiClient.get<PaginatedResponse<StockMovement>>(endpoint);
  }

  // 创建库存移动记录
  static async createStockMovement(movementData: CreateStockMovementRequest): Promise<StockMovement> {
    return apiClient.post<StockMovement>('/stock-movements/', movementData);
  }

  // 入库操作
  static async stockIn(data: {
    itemId: number;
    warehouseId: number;
    quantity: number;
    reason: string;
    reference?: string;
  }): Promise<StockMovement> {
    return apiClient.post<StockMovement>('/stock-movements/in', data);
  }

  // 出库操作
  static async stockOut(data: {
    itemId: number;
    warehouseId: number;
    quantity: number;
    reason: string;
    reference?: string;
  }): Promise<StockMovement> {
    return apiClient.post<StockMovement>('/stock-movements/out', data);
  }

  // 库存调整
  static async stockAdjustment(data: {
    itemId: number;
    warehouseId: number;
    quantity: number;
    reason: string;
    reference?: string;
  }): Promise<StockMovement> {
    return apiClient.post<StockMovement>('/stock-movements/adjustment', data);
  }

  // 库存转移
  static async stockTransfer(data: {
    itemId: number;
    fromWarehouseId: number;
    toWarehouseId: number;
    quantity: number;
    reason: string;
    reference?: string;
  }): Promise<StockMovement[]> {
    return apiClient.post<StockMovement[]>('/stock-movements/transfer', data);
  }

  // 获取库存统计信息
  static async getInventoryStats(): Promise<{
    totalItems: number;
    totalValue: number;
    lowStockItems: number;
    outOfStockItems: number;
    totalMovements: number;
    recentMovements: StockMovement[];
  }> {
    return apiClient.get('/inventory/stats');
  }

  // 获取库存报告
  static async getInventoryReport(params?: {
    warehouseId?: number;
    category?: string;
    startDate?: string;
    endDate?: string;
  }): Promise<{
    items: Array<{
      itemId: number;
      itemCode: string;
      itemName: string;
      category: string;
      unit: string;
      quantity: number;
      cost: number;
      totalValue: number;
      reorderLevel: number;
      status: string;
    }>;
    summary: {
      totalItems: number;
      totalValue: number;
      lowStockCount: number;
      outOfStockCount: number;
    };
  }> {
    const searchParams = new URLSearchParams();
    
    if (params?.warehouseId) searchParams.append('warehouse_id', params.warehouseId.toString());
    if (params?.category) searchParams.append('category', params.category);
    if (params?.startDate) searchParams.append('start_date', params.startDate);
    if (params?.endDate) searchParams.append('end_date', params.endDate);
    
    const query = searchParams.toString();
    const endpoint = query ? `/inventory/report?${query}` : '/inventory/report';
    
    return apiClient.get(endpoint);
  }

  // 获取ABC分析
  static async getABCAnalysis(params?: {
    warehouseId?: number;
    period?: string;
  }): Promise<{
    categoryA: Array<{ itemId: number; itemName: string; value: number; percentage: number; }>;
    categoryB: Array<{ itemId: number; itemName: string; value: number; percentage: number; }>;
    categoryC: Array<{ itemId: number; itemName: string; value: number; percentage: number; }>;
  }> {
    const searchParams = new URLSearchParams();
    
    if (params?.warehouseId) searchParams.append('warehouse_id', params.warehouseId.toString());
    if (params?.period) searchParams.append('period', params.period);
    
    const query = searchParams.toString();
    const endpoint = query ? `/inventory/abc-analysis?${query}` : '/inventory/abc-analysis';
    
    return apiClient.get(endpoint);
  }

  // 导出库存报告
  static async exportInventoryReport(params?: {
    warehouseId?: number;
    category?: string;
    format?: 'csv' | 'xlsx';
  }): Promise<Blob> {
    const searchParams = new URLSearchParams();
    
    if (params?.warehouseId) searchParams.append('warehouse_id', params.warehouseId.toString());
    if (params?.category) searchParams.append('category', params.category);
    if (params?.format) searchParams.append('format', params.format);
    
    const query = searchParams.toString();
    const endpoint = query ? `/inventory/export?${query}` : '/inventory/export';
    
    return apiClient.get(endpoint);
  }
}

export default InventoryService;