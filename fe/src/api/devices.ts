import client from './client';
import type {
  ApiResponse,
  Device,
  DeviceListResponse,
  CreateDeviceRequest,
  UpdateDeviceRequest,
} from '../types';

export const devicesApi = {
  getAll: (page = 1, limit = 10) =>
    client.get<ApiResponse<DeviceListResponse>>('/api/v1/devices', {
      params: { page, limit },
    }),

  getById: (id: string) =>
    client.get<ApiResponse<Device>>(`/api/v1/devices/${id}`),

  getByBrandId: (brandId: string, page = 1, limit = 10) =>
    client.get<ApiResponse<DeviceListResponse>>(`/api/v1/devices/brand/${brandId}`, {
      params: { page, limit },
    }),

  create: (data: CreateDeviceRequest) =>
    client.post<ApiResponse<Device>>('/api/v1/devices', data),

  update: (id: string, data: UpdateDeviceRequest) =>
    client.put<ApiResponse<Device>>(`/api/v1/devices/${id}`, data),

  delete: (id: string) =>
    client.delete<ApiResponse>(`/api/v1/devices/${id}`),
};
