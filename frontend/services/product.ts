import { apiClient } from '@/lib/api';
import { Product, CreateProductRequest, PaginatedResponse } from '@/types/api';

export class ProductService {
  // 获取产品列表
  static async getProducts(params?: {
    page?: number;
    limit?: number;
    search?: string;
    category?: string;
    status?: string;
  }): Promise<PaginatedResponse<Product>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.category) searchParams.append('category', params.category);
    if (params?.status) searchParams.append('status', params.status);
    
    const query = searchParams.toString();
    const endpoint = query ? `/products/?${query}` : '/products/';
    
    return apiClient.get<PaginatedResponse<Product>>(endpoint);
  }

  // 获取单个产品
  static async getProduct(id: number): Promise<Product> {
    return apiClient.get<Product>(`/products/${id}`);
  }

  // 创建产品
  static async createProduct(productData: CreateProductRequest): Promise<Product> {
    return apiClient.post<Product>('/products/', productData);
  }

  // 更新产品
  static async updateProduct(id: number, productData: Partial<CreateProductRequest>): Promise<Product> {
    return apiClient.put<Product>(`/products/${id}`, productData);
  }

  // 删除产品
  static async deleteProduct(id: number): Promise<void> {
    return apiClient.delete<void>(`/products/${id}`);
  }

  // 获取产品统计信息
  static async getProductStats(): Promise<{
    total: number;
    lowStock: number;
    outOfStock: number;
    categories: number;
  }> {
    return apiClient.get('/items/stats');
  }

  // 获取产品分类
  static async getCategories(): Promise<string[]> {
    return apiClient.get('/items/categories');
  }

  // 批量更新产品
  static async bulkUpdateProducts(updates: Array<{
    id: number;
    data: Partial<CreateProductRequest>;
  }>): Promise<Product[]> {
    return apiClient.put('/items/bulk', { updates });
  }

  // 导入产品
  static async importProducts(file: File): Promise<{
    success: number;
    failed: number;
    errors: string[];
  }> {
    const formData = new FormData();
    formData.append('file', file);
    
    return apiClient.post('/items/import', formData);
  }

  // 导出产品
  static async exportProducts(params?: {
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
    
    return apiClient.get(endpoint);
  }
}

export default ProductService;