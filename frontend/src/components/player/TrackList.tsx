"use client";
import { usePlayerStore } from "@/store/usePlayerStore";
import type { Track } from "@/types/player";
import { Play, Pause } from "lucide-react";

function formatDuration(seconds: number) {
  if (!isFinite(seconds) || seconds < 0) return "-:--";
  const m = Math.floor(seconds / 60);
  const s = Math.floor(seconds % 60);
  return `${m}:${s.toString().padStart(2, "0")}`;
}

interface TrackListProps {
  tracks: Track[];
}

// Daftar lagu ala Spotify: nomor urut, cover, judul+artis, durasi.
// Klik baris = putar lagu itu dengan seluruh `tracks` sebagai antrian.
// Klik baris yang SEDANG diputar = toggle play/pause, bukan mengulang dari awal.
export default function TrackList({ tracks }: TrackListProps) {
  const currentTrack = usePlayerStore((s) => s.currentTrack);
  const isPlaying = usePlayerStore((s) => s.isPlaying);
  const playTrack = usePlayerStore((s) => s.playTrack);
  const togglePlay = usePlayerStore((s) => s.togglePlay);

  const handleRowClick = (track: Track) => {
    const isThisTrackActive = currentTrack?.id === track.id;

    if (isThisTrackActive) {
      togglePlay();
    } else {
      playTrack(track, tracks);
    }
  };

  return (
    <div className="w-full">
      {/* Header kolom */}
      <div className="grid grid-cols-[40px_1fr_80px] gap-4 px-4 py-2 text-xs text-zinc-400 uppercase border-b border-zinc-800">
        <span className="text-center">#</span>
        <span>Judul</span>
        <span className="text-right">Durasi</span>
      </div>

      {tracks.map((track, index) => {
        const isActive = currentTrack?.id === track.id;
        const isActiveAndPlaying = isActive && isPlaying;

        return (
          <div key={track.id} onClick={() => handleRowClick(track)} className="group grid grid-cols-[40px_1fr_80px] gap-4 px-4 py-2 rounded-md items-center hover:bg-zinc-800/60 cursor-pointer transition">
            {/* Nomor urut / tombol play-pause saat hover */}
            <div className="flex items-center justify-center text-sm text-zinc-400">
              <span className={`${isActive ? "hidden" : "group-hover:hidden"} ${isActive ? "text-green-500" : ""}`}>{index + 1}</span>
              <span className={`${isActive ? "flex" : "hidden group-hover:flex"} text-white`}>{isActiveAndPlaying ? <Pause size={14} fill="currentColor" /> : <Play size={14} fill="currentColor" />}</span>
            </div>

            {/* Cover + judul + artis */}
            <div className="flex items-center gap-3 min-w-0">
              <img src={track.coverUrl} alt={track.title} className="w-10 h-10 rounded object-cover flex-shrink-0" />
              <div className="min-w-0">
                <p className={`text-sm font-medium truncate ${isActive ? "text-green-500" : "text-white"}`}>{track.title}</p>
                <p className="text-xs text-zinc-400 truncate">{track.artist}</p>
              </div>
            </div>

            {/* Durasi */}
            <span className="text-sm text-zinc-400 text-right tabular-nums">{formatDuration(track.duration)}</span>
          </div>
        );
      })}
    </div>
  );
}
