import { apiClient } from '@/lib/api';

// 报价单版本相关类型定义
export interface QuotationVersion {
  id: number;
  quotation_id: number;
  version_number: number;
  version_name?: string;
  change_reason?: string;
  is_active: boolean;
  version_data: string;
  created_by: number;
  created_at: string;
  updated_at: string;
}

export interface QuotationVersionCreateRequest {
  quotation_id: number;
  version_name?: string;
  change_reason?: string;
}

export interface QuotationVersionCompareRequest {
  quotation_id: number;
  from_version_id: number;
  to_version_id: number;
}

export interface QuotationVersionComparisonResponse {
  field_name: string;
  old_value: any;
  new_value: any;
  change_type: string; // added, modified, deleted
  description?: string;
}

export interface QuotationVersionRollbackRequest {
  quotation_id: number;
  version_id: number;
  reason?: string;
}

export interface QuotationVersionHistoryResponse {
  id: number;
  version_number: number;
  version_name?: string;
  change_reason?: string;
  is_active: boolean;
  created_by: number;
  created_at: string;
  creator_name?: string;
}

export class QuotationVersionService {
  // 创建报价单版本
  static async createVersion(data: QuotationVersionCreateRequest): Promise<QuotationVersion> {
    return apiClient.post<QuotationVersion>('/quotation-versions', data);
  }

  // 获取报价单版本详情
  static async getVersion(versionId: number): Promise<QuotationVersion> {
    return apiClient.get<QuotationVersion>(`/quotation-versions/${versionId}`);
  }

  // 获取报价单的所有版本
  static async getVersionsByQuotation(quotationId: number): Promise<QuotationVersion[]> {
    return apiClient.get<QuotationVersion[]>(`/quotations/${quotationId}/versions`);
  }

  // 设置活跃版本
  static async setActiveVersion(quotationId: number, versionNumber: number): Promise<void> {
    return apiClient.put<void>(`/quotations/${quotationId}/versions/${versionNumber}/set-active`);
  }

  // 比较版本
  static async compareVersions(data: QuotationVersionCompareRequest): Promise<QuotationVersionComparisonResponse[]> {
    return apiClient.post<QuotationVersionComparisonResponse[]>('/quotation-versions/compare', data);
  }

  // 获取版本历史
  static async getVersionHistory(quotationId: number): Promise<QuotationVersionHistoryResponse[]> {
    return apiClient.get<QuotationVersionHistoryResponse[]>(`/quotations/${quotationId}/version-history`);
  }

  // 回滚到指定版本
  static async rollbackToVersion(data: QuotationVersionRollbackRequest): Promise<QuotationVersion> {
    return apiClient.post<QuotationVersion>('/quotation-versions/rollback', data);
  }

  // 删除版本
  static async deleteVersion(versionId: number): Promise<void> {
    return apiClient.delete<void>(`/quotation-versions/${versionId}`);
  }
}

export default QuotationVersionService;