// ==================== Domain Models ====================

export interface Network {
  technology?: string;
  bands_2g?: string;
  bands_3g?: string;
  bands_4g?: string;
  bands_5g?: string;
  speed?: string;
}

export interface Launch {
  announced?: string;
  status?: string;
}

export interface Body {
  dimensions?: string;
  weight?: string;
  build?: string;
  sim?: string;
  ip_rating?: string;
}

export interface Display {
  type?: string;
  size?: string;
  resolution?: string;
}

export interface Platform {
  os?: string;
  chipset?: string;
  cpu?: string;
  gpu?: string;
}

export interface Memory {
  card_lot?: string;
  internal?: string;
}

export interface MainCamera {
  triple?: string;
  single?: string;
  features?: string;
  video?: string;
}

export interface SelfieCamera {
  single?: string;
  video?: string;
}

export interface Sound {
  loudspeaker?: string;
  'jack_3.5mm'?: string;
}

export interface Comms {
  wlan?: string;
  bluetooth?: string;
  positioning?: string;
  nfc?: string;
  radio?: string;
  usb?: string;
}

export interface Features {
  sensors?: string;
}

export interface Battery {
  type?: string;
  charging?: string;
}

export interface Misc {
  colors?: string;
  models?: string;
  price?: string;
}

export interface Specifications {
  network?: Network;
  launch?: Launch;
  body?: Body;
  display?: Display;
  platform?: Platform;
  memory?: Memory;
  mainCamera?: MainCamera;
  selfieCamera?: SelfieCamera;
  sound?: Sound;
  comms?: Comms;
  features?: Features;
  battery?: Battery;
  misc?: Misc;
}

export interface Device {
  id?: string;
  brand_id?: string;
  model_name?: string;
  imageUrl?: string;
  specifications?: Specifications;
}

export interface Brand {
  id?: string;
  brand_name?: string;
  devices?: Device[];
}

// ==================== API Response ====================

export interface ApiResponse<T = unknown> {
  success: boolean;
  code: number;
  message?: string;
  data?: T;
  errors?: ErrorResponse[];
}

export interface ErrorResponse {
  field?: string;
  error: string;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface BrandListResponse {
  brands: Brand[];
  total: number;
  pagination: PaginationMeta;
}

export interface DeviceListResponse {
  devices: Device[];
  total: number;
  pagination: PaginationMeta;
}

export interface FavoriteListResponse {
  device_ids: string[];
}

export interface DeviceFinderParams {
  brand_id?: string;
  os?: string;
  chipset?: string;
  cpu?: string;
  gpu?: string;
  memory?: string;
  display_size?: string;
  battery?: string;
  nfc?: string;
  page?: number;
  limit?: number;
}

// ==================== Auth ====================

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  otp: string;
}

export interface SendOtpRequest {
  email: string;
}

export interface ResetPasswordRequest {
  email: string;
  otp: string;
  new_password: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
  otp: string;
}

export interface AuthResponse {
  access_token: string;
  user: {
    id: string;
    email: string;
  };
}

export interface User {
  id: string;
  email: string;
  role?: string;
}

// ==================== Request DTOs ====================

export interface CreateBrandRequest {
  brand_name: string;
}

export interface UpdateBrandRequest {
  brand_name: string;
}

export interface CreateDeviceRequest {
  brand_id: string;
  model_name: string;
  image: File;
  specifications: Specifications;
}

export interface UpdateDeviceRequest {
  brand_id: string;
  model_name: string;
  image?: File;
  specifications: Specifications;
}

export interface AddFavoriteRequest {
  device_id: string;
}

export interface SyncFavoritesRequest {
  device_ids: string[];
}

