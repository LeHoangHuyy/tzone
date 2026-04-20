import client from './client';
import type {
  ApiResponse,
  FavoriteListResponse,
  AddFavoriteRequest,
  SyncFavoritesRequest,
} from '../types';

export const favoritesApi = {
  getAll: () =>
    client.get<ApiResponse<FavoriteListResponse>>('/api/v1/favorites'),

  add: (data: AddFavoriteRequest) =>
    client.post<ApiResponse<FavoriteListResponse>>('/api/v1/favorites', data),

  remove: (deviceId: string) =>
    client.delete<ApiResponse<FavoriteListResponse>>(`/api/v1/favorites/${deviceId}`),

  sync: (data: SyncFavoritesRequest) =>
    client.post<ApiResponse<FavoriteListResponse>>('/api/v1/favorites/sync', data),
};

