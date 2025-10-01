import { apiClient } from '@/lib/api';
import { Quotation, CreateQuotationRequest, PaginatedResponse } from '@/types/api';

export class QuotationService {
  // 获取报价列表
  static async getQuotations(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    customerId?: number;
  }): Promise<PaginatedResponse<Quotation>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.status) searchParams.append('status', params.status);
    if (params?.customerId) searchParams.append('customer_id', params.customerId.toString());
    
    const query = searchParams.toString();
    const endpoint = query ? `/quotations/?${query}` : '/quotations/';
    
    return apiClient.get<PaginatedResponse<Quotation>>(endpoint);
  }

  // 获取单个报价
  static async getQuotation(id: number): Promise<Quotation> {
    return apiClient.get<Quotation>(`/quotations/${id}`);
  }

  // 创建报价
  static async createQuotation(quotationData: CreateQuotationRequest): Promise<Quotation> {
    return apiClient.post<Quotation>('/quotations/', quotationData);
  }

  // 更新报价
  static async updateQuotation(id: number, quotationData: Partial<CreateQuotationRequest>): Promise<Quotation> {
    return apiClient.put<Quotation>(`/quotations/${id}`, quotationData);
  }

  // 删除报价
  static async deleteQuotation(id: number): Promise<void> {
    return apiClient.delete<void>(`/quotations/${id}`);
  }

  // 提交报价
  static async submitQuotation(id: number): Promise<Quotation> {
    return apiClient.post<Quotation>(`/quotations/${id}/submit`);
  }

  // 接受报价
  static async acceptQuotation(id: number): Promise<Quotation> {
    return apiClient.post<Quotation>(`/quotations/${id}/accept`);
  }

  // 拒绝报价
  static async rejectQuotation(id: number): Promise<Quotation> {
    return apiClient.post<Quotation>(`/quotations/${id}/reject`);
  }

  // 获取报价统计信息
  static async getQuotationStats(): Promise<{
    total: number;
    draft: number;
    submitted: number;
    accepted: number;
    rejected: number;
    expired: number;
    totalAmount: number;
  }> {
    return apiClient.get('/quotations/stats');
  }
}

export default QuotationService;