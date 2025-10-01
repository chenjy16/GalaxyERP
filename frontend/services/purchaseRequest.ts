import { apiClient } from '@/lib/api';
import { PurchaseRequest, CreatePurchaseRequestRequest, PaginatedResponse } from '@/types/api';

export class PurchaseRequestService {
  // 获取采购请求列表
  static async getPurchaseRequests(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    department?: string;
  }): Promise<PaginatedResponse<PurchaseRequest>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.status) searchParams.append('status', params.status);
    if (params?.department) searchParams.append('department', params.department);
    
    const query = searchParams.toString();
    const endpoint = query ? `/purchase-requests/?${query}` : '/purchase-requests/';
    
    return apiClient.get<PaginatedResponse<PurchaseRequest>>(endpoint);
  }

  // 获取单个采购请求
  static async getPurchaseRequest(id: number): Promise<PurchaseRequest> {
    return apiClient.get<PurchaseRequest>(`/purchase-requests/${id}`);
  }

  // 创建采购请求
  static async createPurchaseRequest(requestData: CreatePurchaseRequestRequest): Promise<PurchaseRequest> {
    return apiClient.post<PurchaseRequest>('/purchase-requests/', requestData);
  }

  // 更新采购请求
  static async updatePurchaseRequest(id: number, requestData: Partial<CreatePurchaseRequestRequest>): Promise<PurchaseRequest> {
    return apiClient.put<PurchaseRequest>(`/purchase-requests/${id}`, requestData);
  }

  // 删除采购请求
  static async deletePurchaseRequest(id: number): Promise<void> {
    return apiClient.delete<void>(`/purchase-requests/${id}`);
  }

  // 提交采购请求
  static async submitPurchaseRequest(id: number): Promise<PurchaseRequest> {
    return apiClient.post<PurchaseRequest>(`/purchase-requests/${id}/submit`);
  }

  // 批准采购请求
  static async approvePurchaseRequest(id: number): Promise<PurchaseRequest> {
    return apiClient.post<PurchaseRequest>(`/purchase-requests/${id}/approve`);
  }

  // 拒绝采购请求
  static async rejectPurchaseRequest(id: number, reason?: string): Promise<PurchaseRequest> {
    return apiClient.post<PurchaseRequest>(`/purchase-requests/${id}/reject`, { reason });
  }

  // 获取采购请求统计信息
  static async getPurchaseRequestStats(): Promise<{
    total: number;
    pending: number;
    approved: number;
    rejected: number;
  }> {
    return apiClient.get('/purchase-requests/stats');
  }
}

export default PurchaseRequestService;