import { apiClient } from '@/lib/api';
import { Warehouse, CreateWarehouseRequest, PaginatedResponse } from '@/types/api';

export class WarehouseService {
  // 获取仓库列表
  static async getWarehouses(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
  }): Promise<PaginatedResponse<Warehouse>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.status) searchParams.append('status', params.status);
    
    const query = searchParams.toString();
    const endpoint = query ? `/warehouses/?${query}` : '/warehouses/';
    
    return apiClient.get<PaginatedResponse<Warehouse>>(endpoint);
  }

  // 获取单个仓库
  static async getWarehouse(id: number): Promise<Warehouse> {
    return apiClient.get<Warehouse>(`/warehouses/${id}`);
  }

  // 创建仓库
  static async createWarehouse(warehouseData: CreateWarehouseRequest): Promise<Warehouse> {
    return apiClient.post<Warehouse>('/warehouses/', warehouseData);
  }

  // 更新仓库
  static async updateWarehouse(id: number, warehouseData: Partial<CreateWarehouseRequest>): Promise<Warehouse> {
    return apiClient.put<Warehouse>(`/warehouses/${id}`, warehouseData);
  }

  // 删除仓库
  static async deleteWarehouse(id: number): Promise<void> {
    return apiClient.delete<void>(`/warehouses/${id}`);
  }

  // 获取仓库统计信息
  static async getWarehouseStats(): Promise<{
    total: number;
    active: number;
    inactive: number;
  }> {
    return apiClient.get('/warehouses/stats');
  }
}

export default WarehouseService;