import { apiClient } from '@/lib/api';
import { Customer, CreateCustomerRequest, PaginatedResponse } from '@/types/api';

export class CustomerService {
  // 获取客户列表
  static async getCustomers(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
  }): Promise<PaginatedResponse<Customer>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.status) searchParams.append('status', params.status);
    
    const query = searchParams.toString();
    const endpoint = query ? `/customers/?${query}` : '/customers/';
    
    return apiClient.get<PaginatedResponse<Customer>>(endpoint);
  }

  // 获取单个客户
  static async getCustomer(id: number): Promise<Customer> {
    return apiClient.get<Customer>(`/customers/${id}`);
  }

  // 创建客户
  static async createCustomer(customerData: CreateCustomerRequest): Promise<Customer> {
    return apiClient.post<Customer>('/customers/', customerData);
  }

  // 更新客户
  static async updateCustomer(id: number, customerData: Partial<CreateCustomerRequest>): Promise<Customer> {
    return apiClient.put<Customer>(`/customers/${id}`, customerData);
  }

  // 删除客户
  static async deleteCustomer(id: number): Promise<void> {
    return apiClient.delete<void>(`/customers/${id}`);
  }

  // 获取客户统计信息
  static async getCustomerStats(): Promise<{
    total: number;
    active: number;
    inactive: number;
    totalCreditLimit: number;
  }> {
    return apiClient.get('/customers/stats');
  }
}

export default CustomerService;