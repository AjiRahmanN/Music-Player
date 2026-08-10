export interface Song {
  id: string;
  title: string;
  artist: string;
  album?: string | null;
  coverUrl: string;
  audioUrl: string;
  duration: number; // dalam detik
  createdAt: string;
}

export interface SongListResponse {
  data?: {
    songs: Song[];
  };
  error?: string;
  message?: string;
}

export interface SongDetailResponse {
  data?: {
    song: Song;
  };
  error?: string;
  message?: string;
}