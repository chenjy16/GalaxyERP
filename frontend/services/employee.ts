import { apiClient } from '@/lib/api';
import { Employee, CreateEmployeeRequest, PaginatedResponse } from '@/types/api';

export class EmployeeService {
  // 获取员工列表
  static async getEmployees(params?: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    departmentId?: number;
  }): Promise<PaginatedResponse<Employee>> {
    const searchParams = new URLSearchParams();
    
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('page_size', params.limit.toString());
    if (params?.search) searchParams.append('search', params.search);
    if (params?.status) searchParams.append('status', params.status);
    if (params?.departmentId) searchParams.append('departmentId', params.departmentId.toString());
    
    const query = searchParams.toString();
    const endpoint = query ? `/employees/?${query}` : '/employees/';
    
    return apiClient.get<PaginatedResponse<Employee>>(endpoint);
  }

  // 获取单个员工
  static async getEmployee(id: number): Promise<Employee> {
    return apiClient.get<Employee>(`/employees/${id}`);
  }

  // 创建员工
  static async createEmployee(employeeData: CreateEmployeeRequest): Promise<Employee> {
    return apiClient.post<Employee>('/employees/', employeeData);
  }

  // 更新员工
  static async updateEmployee(id: number, employeeData: Partial<CreateEmployeeRequest>): Promise<Employee> {
    return apiClient.put<Employee>(`/employees/${id}`, employeeData);
  }

  // 删除员工
  static async deleteEmployee(id: number): Promise<void> {
    return apiClient.delete<void>(`/employees/${id}`);
  }

  // 获取员工统计信息
  static async getEmployeeStats(): Promise<{
    total: number;
    active: number;
    inactive: number;
    byDepartment: Record<string, number>;
  }> {
    return apiClient.get('/employees/stats');
  }
}

export default EmployeeService;