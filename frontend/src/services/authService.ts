import api from "@/lib/axios";
import { LoginRequest, LoginResponse, RegisterRequest, MeResponse, UpdateProfileRequest, UpdateProfileResponse } from "@/types/auth";

export const authService = {
  async login(data: LoginRequest): Promise<LoginResponse> {
    const res = await api.post("/auth/login", data);
    return res.data;
  },

  async register(data: RegisterRequest) {
    const res = await api.post("/auth/register", data);
    return res.data;
  },

  // GET /me di backend mengembalikan objek flat ({id, username, email, created}),
  // bukan dibungkus { data: { user } }. Tipe MeResponse mencerminkan bentuk itu.
  async getMe(): Promise<MeResponse> {
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
