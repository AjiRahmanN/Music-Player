export interface Track {
  id: string;
  title: string;
  artist: string;
  album?: string;
  coverUrl: string;
  audioUrl: string;
  duration: number; // dalam detik
}

export type RepeatMode = "off" | "all" | "one";
