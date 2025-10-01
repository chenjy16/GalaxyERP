import { apiClient } from '@/lib/api';
import { SalesOrder, CreateSalesOrderRequest, PaginatedResponse } from '@/types/api';

export class SalesOrderService {
  // 获取销售订单列表
  static async getSalesOrders(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    customerId?: number;
    priority?: string;
  }): Promise<PaginatedResponse<SalesOrder>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.status) searchParams.append('status', params.status);
    if (params?.customerId) searchParams.append('customer_id', params.customerId.toString());
    if (params?.priority) searchParams.append('priority', params.priority);
    
    const query = searchParams.toString();
    const endpoint = query ? `/sales-orders/?${query}` : '/sales-orders/';
    
    return apiClient.get<PaginatedResponse<SalesOrder>>(endpoint);
  }

  // 获取单个销售订单
  static async getSalesOrder(id: number): Promise<SalesOrder> {
    return apiClient.get<SalesOrder>(`/sales-orders/${id}`);
  }

  // 创建销售订单
  static async createSalesOrder(orderData: CreateSalesOrderRequest): Promise<SalesOrder> {
    return apiClient.post<SalesOrder>('/sales-orders/', orderData);
  }

  // 更新销售订单
  static async updateSalesOrder(id: number, orderData: Partial<CreateSalesOrderRequest>): Promise<SalesOrder> {
    return apiClient.put<SalesOrder>(`/sales-orders/${id}`, orderData);
  }

  // 删除销售订单
  static async deleteSalesOrder(id: number): Promise<void> {
    return apiClient.delete<void>(`/sales-orders/${id}`);
  }

  // 确认销售订单
  static async confirmSalesOrder(id: number): Promise<SalesOrder> {
    return apiClient.post<SalesOrder>(`/sales-orders/${id}/confirm`);
  }

  // 取消销售订单
  static async cancelSalesOrder(id: number): Promise<SalesOrder> {
    return apiClient.post<SalesOrder>(`/sales-orders/${id}/cancel`);
  }

  // 从报价创建销售订单
  static async createFromQuotation(quotationId: number, orderData?: Partial<CreateSalesOrderRequest>): Promise<SalesOrder> {
    return apiClient.post<SalesOrder>(`/quotations/${quotationId}/create-order`, orderData);
  }

  // 获取销售订单统计信息
  static async getSalesOrderStats(): Promise<{
    total: number;
    draft: number;
    confirmed: number;
    shipped: number;
    delivered: number;
    cancelled: number;
    totalAmount: number;
  }> {
    return apiClient.get('/sales-orders/stats');
  }
}

export default SalesOrderService;