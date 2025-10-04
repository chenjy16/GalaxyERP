// API 基础配置
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// HTTP 客户端配置
class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
    // 从 localStorage 获取 token
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('auth_token');
    }
  }

  // 设置认证 token
  setToken(token: string) {
    this.token = token;
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token);
    }
  }

  // 清除认证 token
  clearToken() {
    this.token = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  }

  // 通用请求方法
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    };

    // 添加认证头
    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const config: RequestInit = {
      ...options,
      headers,
    };

    try {
      const response = await fetch(url, config);
      
      // 处理认证失败
      if (response.status === 401) {
        this.clearToken();
        throw new Error('认证失败，请重新登录');
      }

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      const jsonResponse = await response.json();
      
      // 处理新的API响应格式 {success: true, data: [...]}
      if (jsonResponse && typeof jsonResponse === 'object' && 'success' in jsonResponse) {
        if (!jsonResponse.success) {
          throw new Error(jsonResponse.message || 'API request failed');
        }
        // 返回data字段，如果没有data字段则返回整个响应
        return jsonResponse.data !== undefined ? jsonResponse.data : jsonResponse;
      }
      
      // 兼容旧的响应格式
      return jsonResponse;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // GET 请求
  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  // GET 请求 - 用于分页响应，返回完整的响应结构
  async getPaginated<T>(endpoint: string): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    // 添加认证头
    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const config: RequestInit = {
      method: 'GET',
      headers,
    };

    try {
      const response = await fetch(url, config);
      
      // 处理认证失败
      if (response.status === 401) {
        this.clearToken();
        throw new Error('认证失败，请重新登录');
      }

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      const jsonResponse = await response.json();
      
      // 对于分页响应，直接返回完整的响应结构
      return jsonResponse;
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // POST 请求
  async post<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  // PUT 请求
  async put<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  // DELETE 请求
  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

// 创建全局 API 客户端实例
export const apiClient = new ApiClient(API_BASE_URL);

// 导出 API 基础 URL
export { API_BASE_URL };