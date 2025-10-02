import { apiClient } from '@/lib/api';
import { Item, CreateItemRequest, PaginatedResponse } from '@/types/api';

export class ItemService {
  // 获取物料列表
  static async getItems(params?: {
    page?: number;
    pageSize?: number;
    search?: string;
    category?: string;
    status?: string;
  }): Promise<PaginatedResponse<Item>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.pageSize) searchParams.append('page_size', params.pageSize.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.category) searchParams.append('category', params.category);
    if (params?.status) searchParams.append('status', params.status);
    
    const query = searchParams.toString();
    const endpoint = query ? `/items/?${query}` : '/items/';
    
    return apiClient.get<PaginatedResponse<Item>>(endpoint);
  }

  // 获取单个物料
  static async getItem(id: number): Promise<Item> {
    return apiClient.get<Item>(`/items/${id}`);
  }

  // 创建物料
  static async createItem(itemData: CreateItemRequest): Promise<Item> {
    return apiClient.post<Item>('/items/', itemData);
  }

  // 更新物料
  static async updateItem(id: number, itemData: Partial<CreateItemRequest>): Promise<Item> {
    return apiClient.put<Item>(`/items/${id}`, itemData);
  }

  // 删除物料
  static async deleteItem(id: number): Promise<void> {
    return apiClient.delete<void>(`/items/${id}`);
  }

  // 搜索物料
  static async searchItems(params: {
    keyword: string;
    page?: number;
    pageSize?: number;
  }): Promise<PaginatedResponse<Item>> {
    const searchParams = new URLSearchParams();
    searchParams.append('keyword', params.keyword);
    if (params.page) searchParams.append('page', params.page.toString());
    if (params.pageSize) searchParams.append('page_size', params.pageSize.toString());
    
    return apiClient.post<PaginatedResponse<Item>>('/items/search', searchParams);
  }

  // 获取物料统计
  static async getItemStats(): Promise<{
    totalItems: number;
    activeItems: number;
    inactiveItems: number;
    lowStockItems: number;
    categories: Array<{
      category: string;
      count: number;
    }>;
  }> {
    return apiClient.get('/items/stats');
  }

  // 批量更新物料
  static async bulkUpdateItems(updates: Array<{
    id: number;
    data: Partial<CreateItemRequest>;
  }>): Promise<{ success: boolean; updated: number }> {
    return apiClient.put('/items/bulk-update', { updates });
  }

  // 导入物料
  static async importItems(file: File): Promise<{
    success: boolean;
    imported: number;
    errors: Array<{ row: number; message: string }>;
  }> {
    const formData = new FormData();
    formData.append('file', file);
    
    // 对于文件上传，我们需要直接使用fetch
    const token = localStorage.getItem('auth_token');
    const response = await fetch('/api/v1/items/import', {
      method: 'POST',
      headers: {
        'Authorization': token ? `Bearer ${token}` : '',
      },
      body: formData,
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return response.json();
  }

  // 导出物料
  static async exportItems(params?: {
    category?: string;
    status?: string;
    format?: 'csv' | 'xlsx';
  }): Promise<Blob> {
    const searchParams = new URLSearchParams();
    if (params?.category) searchParams.append('category', params.category);
    if (params?.status) searchParams.append('status', params.status);
    if (params?.format) searchParams.append('format', params.format);
    
    const query = searchParams.toString();
    const endpoint = query ? `/items/export?${query}` : '/items/export';
    
    // 对于blob响应，我们需要直接使用fetch
    const token = localStorage.getItem('auth_token');
    const response = await fetch(`/api/v1${endpoint}`, {
      method: 'GET',
      headers: {
        'Authorization': token ? `Bearer ${token}` : '',
      },
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return response.blob();
  }
}

export default ItemService;