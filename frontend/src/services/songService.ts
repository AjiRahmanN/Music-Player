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
  async list(query?: string): Promise<SongListResponse> {
    const res = await api.get("/songs", {
      params: query ? { q: query } : undefined,
    });
    return res.data;
  },

  async getById(id: string): Promise<SongDetailResponse> {
    const res = await api.get(`/songs/${id}`);
    return res.data;
  },
};
