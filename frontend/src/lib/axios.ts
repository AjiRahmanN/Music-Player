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

// ---------------------------------------------------------------------------
// Single-flight refresh
//
// Backend merotasi refresh token: setiap /auth/refresh MENGHAPUS token lama.
// Kalau beberapa request kena 401 bersamaan lalu masing-masing memanggil
// /auth/refresh, hanya yang pertama berhasil dan sisanya gagal -> user ter-logout
// paksa. Karena itu semua pemanggil berbagi SATU promise refresh yang sedang
// berjalan.
// ---------------------------------------------------------------------------
let refreshPromise: Promise<string> | null = null;

export function refreshAccessToken(): Promise<string> {
  if (!refreshPromise) {
    refreshPromise = api
      .post("/auth/refresh", undefined, { skipAuthRefresh: true } as AuthRequestConfig)
      .then((res) => {
        const token: string | undefined = res.data?.data?.accessToken ?? res.data?.accessToken;
        if (!token) throw new Error("Refresh response has no access token");
        useAuthStore.getState().setToken(token);
        return token;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }
  return refreshPromise;
}

function redirectToLogin() {
  useAuthStore.getState().logout();
  if (typeof window !== "undefined" && window.location.pathname !== "/login") {
    window.location.assign("/login");
  }
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config as AuthRequestConfig | undefined;

    const shouldTryRefresh =
      error.response?.status === 401 &&
      originalRequest &&
      !originalRequest._retry &&
      !originalRequest.skipAuthRefresh &&
      !originalRequest.url?.startsWith("/auth/");

    if (!shouldTryRefresh || !originalRequest) {
      return Promise.reject(error);
    }

    originalRequest._retry = true;

    // Kalau token di store sudah lebih baru daripada yang dipakai request ini
    // (request lain sudah selesai me-refresh), cukup ulangi dengan token baru.
    const currentToken = useAuthStore.getState().accessToken;
    const sentAuth = originalRequest.headers?.Authorization;
    let token: string;

    if (currentToken && sentAuth !== `Bearer ${currentToken}`) {
      token = currentToken;
    } else {
      try {
        token = await refreshAccessToken();
      } catch {
        // Hanya kegagalan REFRESH yang berarti sesi berakhir.
        redirectToLogin();
        return Promise.reject(error);
      }
    }

    // Kegagalan retry (mis. 500/404) dikembalikan apa adanya; bukan alasan logout.
    originalRequest.headers.Authorization = `Bearer ${token}`;
    return api(originalRequest);
  }
);

export default api;
