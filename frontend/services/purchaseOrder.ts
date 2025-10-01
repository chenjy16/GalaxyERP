import { apiClient } from '@/lib/api';
import { PurchaseOrder, CreatePurchaseOrderRequest, PaginatedResponse } from '@/types/api';

export class PurchaseOrderService {
  // 获取采购订单列表
  static async getPurchaseOrders(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    supplierId?: number;
  }): Promise<PaginatedResponse<PurchaseOrder>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.status) searchParams.append('status', params.status);
    if (params?.supplierId) searchParams.append('supplier_id', params.supplierId.toString());
    
    const query = searchParams.toString();
    const endpoint = query ? `/purchase-orders/?${query}` : '/purchase-orders/';
    
    return apiClient.get<PaginatedResponse<PurchaseOrder>>(endpoint);
  }

  // 获取单个采购订单
  static async getPurchaseOrder(id: number): Promise<PurchaseOrder> {
    return apiClient.get<PurchaseOrder>(`/purchase-orders/${id}`);
  }

  // 创建采购订单
  static async createPurchaseOrder(orderData: CreatePurchaseOrderRequest): Promise<PurchaseOrder> {
    return apiClient.post<PurchaseOrder>('/purchase-orders/', orderData);
  }

  // 更新采购订单
  static async updatePurchaseOrder(id: number, orderData: Partial<CreatePurchaseOrderRequest>): Promise<PurchaseOrder> {
    return apiClient.put<PurchaseOrder>(`/purchase-orders/${id}`, orderData);
  }

  // 删除采购订单
  static async deletePurchaseOrder(id: number): Promise<void> {
    return apiClient.delete<void>(`/purchase-orders/${id}`);
  }

  // 确认采购订单
  static async confirmPurchaseOrder(id: number): Promise<PurchaseOrder> {
    return apiClient.post<PurchaseOrder>(`/purchase-orders/${id}/confirm`);
  }

  // 取消采购订单
  static async cancelPurchaseOrder(id: number, reason?: string): Promise<PurchaseOrder> {
    return apiClient.post<PurchaseOrder>(`/purchase-orders/${id}/cancel`, { reason });
  }

  // 标记为已发货
  static async markAsShipped(id: number, trackingNumber?: string): Promise<PurchaseOrder> {
    return apiClient.post<PurchaseOrder>(`/purchase-orders/${id}/ship`, { trackingNumber });
  }

  // 标记为已收货
  static async markAsReceived(id: number): Promise<PurchaseOrder> {
    return apiClient.post<PurchaseOrder>(`/purchase-orders/${id}/receive`);
  }

  // 从采购请求创建采购订单
  static async createFromRequest(requestId: number, orderData: Omit<CreatePurchaseOrderRequest, 'items'>): Promise<PurchaseOrder> {
    return apiClient.post<PurchaseOrder>(`/purchase-orders/from-request/${requestId}`, orderData);
  }

  // 获取采购订单统计信息
  static async getPurchaseOrderStats(): Promise<{
    total: number;
    pending: number;
    confirmed: number;
    shipped: number;
    received: number;
    cancelled: number;
    totalAmount: number;
  }> {
    return apiClient.get('/purchase-orders/stats');
  }
}

export default PurchaseOrderService;