import { apiClient } from '../lib/api';

export interface AuditLog {
  id: number;
  user_id: number;
  username: string;
  action: string;
  resource_type: string;
  resource_id: string;
  description: string;
  old_data?: any;
  new_data?: any;
  ip_address: string;
  user_agent: string;
  created_at: string;
}

export interface AuditLogListResponse {
  audit_logs: AuditLog[];
  total: number;
  page: number;
  page_size: number;
}

export interface AuditLogParams {
  page?: number;
  page_size?: number;
  user_id?: number;
  action?: string;
  resource_type?: string;
  start_date?: string;
  end_date?: string;
}

export const auditLogService = {
  // 获取审计日志列表
  getAuditLogs: async (params?: AuditLogParams): Promise<AuditLogListResponse> => {
    const queryString = params ? '?' + new URLSearchParams(params as any).toString() : '';
    return await apiClient.get(`/audit-logs${queryString}`);
  },

  // 根据ID获取审计日志
  getAuditLogById: async (id: number): Promise<AuditLog> => {
    return await apiClient.get(`/audit-logs/${id}`);
  },

  // 获取用户审计日志
  getUserAuditLogs: async (userId: number, params?: AuditLogParams): Promise<AuditLogListResponse> => {
    const queryString = params ? '?' + new URLSearchParams(params as any).toString() : '';
    return await apiClient.get(`/audit-logs/user/${userId}${queryString}`);
  },

  // 获取资源审计日志
  getResourceAuditLogs: async (resourceType: string, resourceId: string, params?: AuditLogParams): Promise<AuditLogListResponse> => {
    const queryString = params ? '?' + new URLSearchParams(params as any).toString() : '';
    return await apiClient.get(`/audit-logs/resource/${resourceType}/${resourceId}${queryString}`);
  },

  // 清理旧日志
  cleanupOldLogs: async (days: number): Promise<{ message: string; deleted_count: number }> => {
    return await apiClient.delete(`/audit-logs/cleanup?days=${days}`);
  },
};