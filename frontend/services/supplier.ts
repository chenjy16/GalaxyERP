import { apiClient } from '@/lib/api';
import { Supplier, CreateSupplierRequest, PaginatedResponse, BackendPaginatedResponse } from '@/types/api';

export class SupplierService {
  // 获取供应商列表
  static async getSuppliers(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
  }): Promise<PaginatedResponse<Supplier>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.status) searchParams.append('status', params.status);
    
    const query = searchParams.toString();
    const endpoint = query ? `/suppliers/?${query}` : '/suppliers/';
    
    const response = await apiClient.getPaginated<BackendPaginatedResponse<Supplier>>(endpoint);
    
    // 转换后端响应格式为前端期望的格式
    return {
      data: response.data,
      total: response.pagination.total,
      page: response.pagination.page,
      limit: response.pagination.page_size,
      totalPages: response.pagination.total_pages,
    };
  }

  // 获取单个供应商
  static async getSupplier(id: number): Promise<Supplier> {
    return apiClient.get<Supplier>(`/suppliers/${id}`);
  }

  // 创建供应商
  static async createSupplier(supplierData: CreateSupplierRequest): Promise<Supplier> {
    return apiClient.post<Supplier>('/suppliers/', supplierData);
  }

  // 更新供应商
  static async updateSupplier(id: number, supplierData: Partial<CreateSupplierRequest>): Promise<Supplier> {
    return apiClient.put<Supplier>(`/suppliers/${id}`, supplierData);
  }

  // 删除供应商
  static async deleteSupplier(id: number): Promise<void> {
    return apiClient.delete<void>(`/suppliers/${id}`);
  }

  // 获取供应商统计信息
  static async getSupplierStats(): Promise<{
    total: number;
    active: number;
    inactive: number;
  }> {
    return apiClient.get('/suppliers/stats');
  }
}

export default SupplierService;