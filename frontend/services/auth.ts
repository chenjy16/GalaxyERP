import { apiClient } from '@/lib/api';
import { LoginRequest, RegisterRequest, AuthResponse, User } from '@/types/api';

export class AuthService {
  // 用户登录
  static async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/login', credentials);
    
    // 保存 token
    if (response.token) {
      apiClient.setToken(response.token);
    }
    
    return response;
  }

  // 用户注册
  static async register(userData: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/register', userData);
    
    // 保存 token
    if (response.token) {
      apiClient.setToken(response.token);
    }
    
    return response;
  }

  // 获取当前用户信息
  static async getCurrentUser(): Promise<User> {
    return apiClient.get<User>('/auth/me');
  }

  // 用户登出
  static async logout(): Promise<void> {
    try {
      await apiClient.post('/auth/logout');
    } finally {
      // 无论请求是否成功，都清除本地 token
      apiClient.clearToken();
    }
  }

  // 刷新 token
  static async refreshToken(): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/refresh');
    
    if (response.token) {
      apiClient.setToken(response.token);
    }
    
    return response;
  }

  // 检查是否已登录
  static isAuthenticated(): boolean {
    if (typeof window === 'undefined') return false;
    return !!localStorage.getItem('auth_token');
  }
}

export default AuthService;