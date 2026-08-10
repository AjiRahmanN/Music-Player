import api from "@/lib/axios";
import { LoginRequest, LoginResponse, RegisterRequest, UserResponse, UpdateProfileRequest, UpdateProfileResponse } from "@/types/auth";

export const authService = {
  async login(data: LoginRequest): Promise<LoginResponse> {
    const res = await api.post("/auth/login", data);
    return res.data;
  },

  async register(data: RegisterRequest) {
    const res = await api.post("/auth/register", data);
    return res.data;
  },

  // Catatan: endpoint GET /me di backend sebenarnya mengembalikan objek flat
  // ({id, username, email, created}), bukan dibungkus { data: { user } } seperti
  // tipe UserResponse. Tipe ini dipertahankan apa adanya (perilaku lama) supaya
  // AuthGuard.tsx yang sudah memakainya tidak berubah. Di halaman profile,
  // hasilnya dibaca lewat tipe MeResponse (lihat src/types/auth.ts).
  async getMe(): Promise<UserResponse> {
    const res = await api.get("/me");
    return res.data;
  },

  async updateProfile(data: UpdateProfileRequest): Promise<UpdateProfileResponse> {
    const res = await api.put("/profile", data);
    return res.data;
  },

  async logout() {
    return api.post("/auth/logout");
  },

  async refreshToken() {
    return api.post("/auth/refresh");
  },
};
