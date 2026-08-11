import api from "@/lib/axios";
import type { Song } from "@/types/song";
import type { AdminUserListResponse, AdminUserUpdateResponse, AdminPlaylistListResponse } from "@/types/admin";

export interface SongFormData {
  title: string;
  artist: string;
  album: string;
  coverUrl: string;
  audioUrl: string;
  duration: number;
}

export interface UpdateUserData {
  username: string;
  email: string;
  isAdmin: boolean;
}

export const adminService = {
  // Melempar error 403 kalau bukan admin — dipakai admin page untuk cek akses
  async checkAccess(): Promise<{ isAdmin: boolean }> {
    const res = await api.get("/admin/check");
    return res.data;
  },

  async uploadFile(file: File, type: "cover" | "audio"): Promise<{ data?: { url: string } }> {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("type", type);

    const res = await api.post("/admin/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    return res.data;
  },

  async createSong(data: SongFormData): Promise<{ data?: { song: Song }; message?: string }> {
    const res = await api.post("/admin/songs", data);
    return res.data;
  },

  async updateSong(id: string, data: SongFormData): Promise<{ data?: { song: Song }; message?: string }> {
    const res = await api.put(`/admin/songs/${id}`, data);
    return res.data;
  },

  async deleteSong(id: string): Promise<{ message?: string }> {
    const res = await api.delete(`/admin/songs/${id}`);
    return res.data;
  },

  // ---- Users ----
  async listUsers(): Promise<AdminUserListResponse> {
    const res = await api.get("/admin/users");
    return res.data;
  },

  async updateUser(id: string, data: UpdateUserData): Promise<AdminUserUpdateResponse> {
    const res = await api.put(`/admin/users/${id}`, data);
    return res.data;
  },

  async deleteUser(id: string): Promise<{ message?: string }> {
    const res = await api.delete(`/admin/users/${id}`);
    return res.data;
  },

  // ---- Playlists ----
  async listAllPlaylists(): Promise<AdminPlaylistListResponse> {
    const res = await api.get("/admin/playlists");
    return res.data;
  },

  async deletePlaylist(id: string): Promise<{ message?: string }> {
    const res = await api.delete(`/admin/playlists/${id}`);
    return res.data;
  },
};
