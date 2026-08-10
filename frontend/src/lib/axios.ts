import axios, { type InternalAxiosRequestConfig } from "axios";
import { useAuthStore } from "@/store/useAuthStore";

type AuthRequestConfig = InternalAxiosRequestConfig & {
  _retry?: boolean;
  skipAuthRefresh?: boolean;
};

const api = axios.create({
  // Backend routes are mounted below /api. Configure the environment variable
  // with that same API root, for example https://api.example.com/api.
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api",
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

api.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().accessToken;
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config as AuthRequestConfig | undefined;

    if (
      error.response?.status === 401 &&
      originalRequest &&
      !originalRequest._retry &&
      !originalRequest.skipAuthRefresh &&
      !originalRequest.url?.startsWith("/auth/")
    ) {
      originalRequest._retry = true;

      try {
        const refreshResponse = await api.post("/auth/refresh", undefined, {
          skipAuthRefresh: true,
        } as AuthRequestConfig);
        const newAccessToken = refreshResponse.data.data?.accessToken || refreshResponse.data.accessToken;

        if (newAccessToken) {
          useAuthStore.getState().setToken(newAccessToken);
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          return api(originalRequest);
        }
        useAuthStore.getState().logout();
        if (typeof window !== "undefined") {
          window.location.assign("/login");
        }
      } catch {
        useAuthStore.getState().logout();
        if (typeof window !== "undefined") {
          window.location.assign("/login");
        }
      }
    }

    return Promise.reject(error);
  }
);

export default api;
