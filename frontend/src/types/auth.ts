export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface LoginResponse {
  message?: string;
  data?: {
    accessToken: string;
    tokenType: string;
    expiresIn: number;
    user: User;
  };
  error?: string;
}

export interface UserResponse {
  message?: string;
  data?: {
    user: User;
  };
  error?: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  createdAt: string;
}

export interface ApiError {
  error: string;
  message?: string;
  statusCode?: number;
}

// ---- Edit Profile ----
export interface UpdateProfileRequest {
  username: string;
  email: string;
}

export interface UpdateProfileResponse {
  message?: string;
  data?: {
    user: User;
  };
  error?: string;
}

// Bentuk response GET /me (flat, tidak dibungkus "data")
export interface MeResponse {
  id: string;
  username: string;
  email: string;
  created: string;
}
