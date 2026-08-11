import api from "@/lib/axios";
import type { PlaylistListResponse, PlaylistDetailResponse, PlaylistCreateResponse } from "@/types/playlist";

export const playlistService = {
  async create(data: { name: string; description?: string }): Promise<PlaylistCreateResponse> {
    const res = await api.post("/playlists", data);
    return res.data;
  },

  async list(): Promise<PlaylistListResponse> {
    const res = await api.get("/playlists");
    return res.data;
  },

  async getDetail(id: string): Promise<PlaylistDetailResponse> {
    const res = await api.get(`/playlists/${id}`);
    return res.data;
  },

  async remove(id: string) {
    const res = await api.delete(`/playlists/${id}`);
    return res.data;
  },

  async addSong(playlistId: string, songId: string) {
    const res = await api.post(`/playlists/${playlistId}/songs`, { songId });
    return res.data;
  },

  async removeSong(playlistId: string, songId: string) {
    const res = await api.delete(`/playlists/${playlistId}/songs/${songId}`);
    return res.data;
  },
};
