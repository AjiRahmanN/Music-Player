import api from "@/lib/axios";
import type { SongListResponse, SongDetailResponse, Song } from "@/types/song";
import type { Track } from "@/types/player";

// Song (bentuk dari backend) dan Track (bentuk yang dipakai usePlayerStore)
// sengaja dipisah tipenya — Song sedikit lebih detail (ada createdAt, album
// bisa null). Fungsi ini yang menjembatani keduanya.
export function mapSongToTrack(song: Song): Track {
  return {
    id: song.id,
    title: song.title,
    artist: song.artist,
    album: song.album ?? undefined,
    coverUrl: song.coverUrl,
    audioUrl: song.audioUrl,
    duration: song.duration,
  };
}

export const songService = {
  // limit/offset opsional. Tanpa keduanya backend memakai default (100 lagu,
  // maks 200 per request).
  async list(query?: string, page?: { limit?: number; offset?: number }): Promise<SongListResponse> {
    const params: Record<string, string | number> = {};
    if (query) params.q = query;
    if (page?.limit !== undefined) params.limit = page.limit;
    if (page?.offset !== undefined) params.offset = page.offset;
    const res = await api.get("/songs", {
      params: Object.keys(params).length ? params : undefined,
    });
    return res.data;
  },

  async getById(id: string): Promise<SongDetailResponse> {
    const res = await api.get(`/songs/${id}`);
    return res.data;
  },
};
