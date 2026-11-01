"use client";
import { useEffect, useRef } from "react";
import { usePlayerStore } from "@/store/usePlayerStore";

// Komponen ini tidak menampilkan apa pun secara visual.
// Tugasnya cuma satu: jadi jembatan antara state di usePlayerStore
// dengan elemen <audio> HTML asli yang benar-benar memutar suara.
export default function PlayerAudioEngine() {
  const audioRef = useRef<HTMLAudioElement | null>(null);

  const currentTrack = usePlayerStore((s) => s.currentTrack);
  const isPlaying = usePlayerStore((s) => s.isPlaying);
  const volume = usePlayerStore((s) => s.volume);
  const repeatMode = usePlayerStore((s) => s.repeatMode);
  const progress = usePlayerStore((s) => s.progress);
  const seekRequestId = usePlayerStore((s) => s.seekRequestId);
  const setProgress = usePlayerStore((s) => s.setProgress);
  const setIsPlaying = usePlayerStore((s) => s.setIsPlaying);
  const setPlaybackError = usePlayerStore((s) => s.setPlaybackError);
  const playNext = usePlayerStore((s) => s.playNext);

  // Ganti sumber audio setiap kali lagu yang aktif berubah
  useEffect(() => {
    const audio = audioRef.current;
    if (!audio || !currentTrack) return;

    setPlaybackError(null);

    audio.src = currentTrack.audioUrl;
    audio.currentTime = 0;

    if (isPlaying) {
      audio.play().catch(() => {
        // Browser modern kadang menolak autoplay tanpa interaksi user dulu,
        // ini bukan bug — cukup set isPlaying kembali ke false.
        setIsPlaying(false);
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentTrack?.id]);

  // Sinkronkan tombol play/pause dari store ke elemen audio
  useEffect(() => {
    const audio = audioRef.current;
    if (!audio || !currentTrack) return;

    if (isPlaying) {
      audio.play().catch(() => setIsPlaying(false));
    } else {
      audio.pause();
    }
  }, [isPlaying, currentTrack, setIsPlaying]);

  // Sinkronkan volume
  useEffect(() => {
    if (audioRef.current) {
      audioRef.current.volume = volume;
    }
  }, [volume]);

  // Terapkan seek manual (klik/geser progress bar) ke elemen audio.
  // Sengaja dipisah dari onTimeUpdate supaya tidak saling menimpa:
  // efek ini HANYA jalan saat seekRequestId berubah (user benar-benar seek),
  // bukan setiap kali audio berjalan normal.
  useEffect(() => {
    const audio = audioRef.current;
    if (!audio || seekRequestId === 0) return;
    audio.currentTime = progress;
    // Audio yang sudah selesai (ended) tidak jalan lagi hanya karena
    // currentTime diubah; kalau store bilang sedang play, mulai lagi
    // (dipakai untuk restart lagu yang sama).
    if (usePlayerStore.getState().isPlaying && audio.paused) {
      audio.play().catch(() => setIsPlaying(false));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [seekRequestId]);

  return (
    <audio
      ref={audioRef}
      onTimeUpdate={(e) => setProgress(e.currentTarget.currentTime)}
      onEnded={() => {
        const audio = audioRef.current;
        // Repeat satu lagu: ulang dari awal, bukan lanjut ke lagu berikutnya
        if (repeatMode === "one" && audio) {
          audio.currentTime = 0;
          audio.play().catch(() => setIsPlaying(false));
          return;
        }
        playNext();
      }}
      onError={() => {
        // URL rusak / format tidak didukung / file hilang: beri tahu user
        // dan hentikan state "playing" supaya UI tidak menipu.
        setIsPlaying(false);
        setPlaybackError("Lagu tidak dapat diputar. File audio mungkin rusak atau tidak ditemukan.");
      }}
      className="hidden"
    />
  );
}
