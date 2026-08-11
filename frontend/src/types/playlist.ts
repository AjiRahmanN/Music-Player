import type { Song } from "./song";

export interface Playlist {
  id: string;
  userId: string;
  name: string;
  description?: string | null;
  coverUrl?: string | null;
  createdAt: string;
  songCount: number;
}

export interface PlaylistDetail extends Playlist {
  songs: Song[];
}

export interface PlaylistListResponse {
  data?: { playlists: Playlist[] };
  error?: string;
  message?: string;
}

export interface PlaylistDetailResponse {
  data?: { playlist: PlaylistDetail };
  error?: string;
  message?: string;
}

export interface PlaylistCreateResponse {
  data?: { playlist: Playlist };
  message?: string;
  error?: string;
}