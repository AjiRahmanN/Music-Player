"use client";
import { useEffect, useRef, useState } from "react";
import { usePlayerStore } from "@/store/usePlayerStore";
import { Play, Pause, SkipBack, SkipForward, Shuffle, Repeat, Repeat1, Volume2, Volume1, VolumeX } from "lucide-react";

function formatTime(seconds: number) {
  if (!isFinite(seconds) || seconds < 0) return "0:00";
  const m = Math.floor(seconds / 60);
  const s = Math.floor(seconds % 60);
  return `${m}:${s.toString().padStart(2, "0")}`;
}

export default function PlayerBar() {
  const currentTrack = usePlayerStore((s) => s.currentTrack);
  const isPlaying = usePlayerStore((s) => s.isPlaying);
  const progress = usePlayerStore((s) => s.progress);
  const isShuffled = usePlayerStore((s) => s.isShuffled);
  const repeatMode = usePlayerStore((s) => s.repeatMode);
  const playbackError = usePlayerStore((s) => s.playbackError);

  const togglePlay = usePlayerStore((s) => s.togglePlay);
  const playNext = usePlayerStore((s) => s.playNext);
  const playPrevious = usePlayerStore((s) => s.playPrevious);
  const toggleShuffle = usePlayerStore((s) => s.toggleShuffle);
  const cycleRepeatMode = usePlayerStore((s) => s.cycleRepeatMode);
  const seekTo = usePlayerStore((s) => s.seekTo);
  const volume = usePlayerStore((s) => s.volume);
  const setVolume = usePlayerStore((s) => s.setVolume);
  const toggleMute = usePlayerStore((s) => s.toggleMute);

  const hasTrack = !!currentTrack;
  const duration = currentTrack?.duration ?? 0;

  // ---- Seek bar interaktif ----
  const barRef = useRef<HTMLDivElement>(null);
  const [isDragging, setIsDragging] = useState(false);
  // Persentase sementara selama drag, supaya tampilan responsif tanpa
  // menunggu audio benar-benar pindah posisi.
  const [dragPercent, setDragPercent] = useState(0);

  const percentFromClientX = (clientX: number) => {
    const bar = barRef.current;
    if (!bar) return 0;
    const rect = bar.getBoundingClientRect();
    const ratio = (clientX - rect.left) / rect.width;
    return Math.min(100, Math.max(0, ratio * 100));
  };

  const handlePointerDown = (e: React.PointerEvent) => {
    if (!hasTrack) return;
    setIsDragging(true);
    setDragPercent(percentFromClientX(e.clientX));
  };

  // Dengarkan gerakan mouse/jari di seluruh window selama drag berlangsung,
  // supaya seek tetap jalan meski kursor keluar dari area bar.
  useEffect(() => {
    if (!isDragging) return;

    const handleMove = (e: PointerEvent) => setDragPercent(percentFromClientX(e.clientX));
    const handleUp = (e: PointerEvent) => {
      const finalPercent = percentFromClientX(e.clientX);
      seekTo((finalPercent / 100) * duration);
      setIsDragging(false);
    };

    window.addEventListener("pointermove", handleMove);
    window.addEventListener("pointerup", handleUp);
    return () => {
      window.removeEventListener("pointermove", handleMove);
      window.removeEventListener("pointerup", handleUp);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isDragging, duration]);

  const storedPercent = hasTrack && duration > 0 ? Math.min(100, (progress / duration) * 100) : 0;
  const progressPercent = isDragging ? dragPercent : storedPercent;
  const displayedSeconds = isDragging ? (dragPercent / 100) * duration : progress;

  return (
    <div className="h-24 bg-black border-t border-zinc-800 px-4 flex items-center justify-between fixed bottom-0 left-0 right-0 z-40">
      {/* Info lagu */}
      <div className="flex items-center gap-3 w-1/4 min-w-[180px]">
        {hasTrack ? (
          <>
            <img src={currentTrack.coverUrl} alt={currentTrack.title} className="w-14 h-14 rounded object-cover" />
            <div className="min-w-0">
              <p className="text-sm font-medium text-white truncate">{currentTrack.title}</p>
              {playbackError ? (
                <p role="alert" className="text-xs text-red-400 line-clamp-2">
                  {playbackError}
                </p>
              ) : (
                <p className="text-xs text-zinc-400 truncate">{currentTrack.artist}</p>
              )}
            </div>
          </>
        ) : (
          <p className="text-sm text-zinc-500">Tidak ada lagu yang diputar</p>
        )}
      </div>

      {/* Kontrol tengah */}
      <div className="flex flex-col items-center gap-2 flex-1 max-w-xl">
        <div className="flex items-center gap-5">
          <button onClick={toggleShuffle} disabled={!hasTrack} className={`transition disabled:opacity-30 disabled:cursor-not-allowed ${isShuffled ? "text-green-500" : "text-zinc-400 hover:text-white"}`} aria-label="Shuffle">
            <Shuffle size={18} />
          </button>

          <button onClick={playPrevious} disabled={!hasTrack} className="text-zinc-300 hover:text-white transition disabled:opacity-30 disabled:cursor-not-allowed" aria-label="Previous">
            <SkipBack size={20} fill="currentColor" />
          </button>

          <button
            onClick={togglePlay}
            disabled={!hasTrack}
            className="w-9 h-9 rounded-full bg-white text-black flex items-center justify-center hover:scale-105 transition disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100"
            aria-label={isPlaying ? "Pause" : "Play"}
          >
            {isPlaying ? <Pause size={18} fill="currentColor" /> : <Play size={18} fill="currentColor" className="ml-0.5" />}
          </button>

          <button onClick={playNext} disabled={!hasTrack} className="text-zinc-300 hover:text-white transition disabled:opacity-30 disabled:cursor-not-allowed" aria-label="Next">
            <SkipForward size={20} fill="currentColor" />
          </button>

          <button onClick={cycleRepeatMode} disabled={!hasTrack} className={`transition disabled:opacity-30 disabled:cursor-not-allowed ${repeatMode !== "off" ? "text-green-500" : "text-zinc-400 hover:text-white"}`} aria-label="Repeat">
            {repeatMode === "one" ? <Repeat1 size={18} /> : <Repeat size={18} />}
          </button>
        </div>

        {/* Progress bar — bisa diklik & digeser */}
        <div className="flex items-center gap-2 w-full">
          <span className="text-[11px] text-zinc-400 w-9 text-right tabular-nums">{formatTime(displayedSeconds)}</span>

          <div ref={barRef} onPointerDown={handlePointerDown} className={`group relative flex-1 h-3 flex items-center ${hasTrack ? "cursor-pointer" : "cursor-not-allowed"}`}>
            {/* Track (jalur) */}
            <div className="w-full h-1 bg-zinc-700 rounded-full overflow-hidden group-hover:h-1.5 transition-all">
              <div className={`h-full rounded-full ${isDragging ? "bg-green-500" : "bg-white group-hover:bg-green-500"}`} style={{ width: `${progressPercent}%` }} />
            </div>

            {/* Thumb (bulatan penanda posisi), muncul saat hover atau sedang di-drag */}
            {hasTrack && <div className={`absolute w-3 h-3 bg-white rounded-full shadow -translate-x-1/2 transition-opacity ${isDragging ? "opacity-100" : "opacity-0 group-hover:opacity-100"}`} style={{ left: `${progressPercent}%` }} />}
          </div>

          <span className="text-[11px] text-zinc-400 w-9 tabular-nums">{formatTime(duration)}</span>
        </div>
      </div>

      {/* Volume control */}
      <div className="flex items-center gap-2 w-1/4 min-w-[180px] justify-end">
        <button onClick={toggleMute} className="text-zinc-400 hover:text-white transition" aria-label={volume === 0 ? "Unmute" : "Mute"}>
          {volume === 0 ? <VolumeX size={18} /> : volume < 0.5 ? <Volume1 size={18} /> : <Volume2 size={18} />}
        </button>

        <input type="range" min={0} max={100} value={Math.round(volume * 100)} onChange={(e) => setVolume(Number(e.target.value) / 100)} className="w-24 accent-white hover:accent-green-500 cursor-pointer" aria-label="Volume" />
      </div>
    </div>
  );
}
