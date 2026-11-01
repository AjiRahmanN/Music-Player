import { create } from "zustand";
import type { Track, RepeatMode } from "@/types/player";

// Fisher-Yates shuffle — tidak mengubah array asli
function shuffleArray<T>(items: T[]): T[] {
  const result = [...items];
  for (let i = result.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [result[i], result[j]] = [result[j], result[i]];
  }
  return result;
}

// Bikin urutan acak dengan currentTrack tetap di posisi pertama
function buildShuffledQueue(queue: Track[], currentTrack: Track | null): Track[] {
  if (!currentTrack) return shuffleArray(queue);
  const rest = queue.filter((t) => t.id !== currentTrack.id);
  return [currentTrack, ...shuffleArray(rest)];
}

// Patch state untuk "putar ulang dari awal lagu yang sama". Dipakai saat
// lagu berikutnya/sebelumnya ternyata lagu yang sama dengan yang sedang aktif
// (queue 1 lagu + repeat all, previous di index 0, klik lagu yang sedang
// diputar): id track tidak berubah sehingga PlayerAudioEngine tidak akan
// memuat ulang sumber audio, jadi kita minta restart lewat seek ke 0.
function restartPatch(seekRequestId: number) {
  return { progress: 0, seekRequestId: seekRequestId + 1, isPlaying: true, playbackError: null };
}

interface PlayerState {
  // Lagu yang sedang aktif
  currentTrack: Track | null;
  // Antrian lagu (playlist yang sedang diputar)
  queue: Track[];
  // Versi acak dari queue, dipakai saat isShuffled true. Lagu yang sedang
  // diputar selalu jadi elemen pertama, sisanya diacak — supaya toggle
  // shuffle tidak membuat lagu yang lagi jalan tiba-tiba berpindah/berhenti.
  shuffledQueue: Track[];
  // Index currentTrack di dalam queue AKTIF (queue biasa atau shuffledQueue,
  // tergantung isShuffled)
  currentIndex: number;

  isPlaying: boolean;
  progress: number; // detik saat ini
  // Bertambah setiap kali user melakukan seek manual (klik/geser progress bar).
  // Dipakai PlayerAudioEngine untuk membedakan "user minta pindah posisi"
  // dari "audio berjalan otomatis" (setProgress), supaya keduanya tidak bentrok.
  seekRequestId: number;
  volume: number; // 0 - 1
  // Menyimpan volume terakhir sebelum di-mute, supaya tombol mute bisa
  // mengembalikan ke level semula (bukan selalu balik ke 100%).
  previousVolume: number;
  isShuffled: boolean;
  repeatMode: RepeatMode;
  // Pesan error pemutaran (mis. URL audio rusak); null jika tidak ada.
  playbackError: string | null;

  // Actions
  setPlaybackError: (message: string | null) => void;
  playTrack: (track: Track, queue?: Track[]) => void;
  togglePlay: () => void;
  setIsPlaying: (playing: boolean) => void;
  setProgress: (seconds: number) => void;
  seekTo: (seconds: number) => void;
  setVolume: (volume: number) => void;
  toggleMute: () => void;
  playNext: () => void;
  playPrevious: () => void;
  toggleShuffle: () => void;
  cycleRepeatMode: () => void;
}

export const usePlayerStore = create<PlayerState>()((set, get) => ({
  currentTrack: null,
  queue: [],
  shuffledQueue: [],
  currentIndex: -1,

  isPlaying: false,
  progress: 0,
  seekRequestId: 0,
  volume: 1,
  previousVolume: 1,
  isShuffled: false,
  repeatMode: "off",
  playbackError: null,

  setPlaybackError: (message) => set({ playbackError: message }),

  // Mulai putar sebuah track. Kalau dikasih `queue` (misal daftar lagu di
  // satu playlist), track lain di queue itu jadi antrian next/previous.
  playTrack: (track, queue) => {
    const newQueue = queue ?? [track];
    const { isShuffled } = get();
    const shuffledQueue = buildShuffledQueue(newQueue, track);
    const activeQueue = isShuffled ? shuffledQueue : newQueue;
    const index = activeQueue.findIndex((t) => t.id === track.id);
    const state = get();
    const sameTrack = state.currentTrack?.id === track.id;

    set({
      currentTrack: track,
      queue: newQueue,
      shuffledQueue,
      currentIndex: index === -1 ? 0 : index,
      ...(sameTrack ? restartPatch(state.seekRequestId) : { progress: 0, isPlaying: true, playbackError: null }),
    });
  },

  togglePlay: () => {
    // Kalau belum ada track sama sekali, tombol play tidak melakukan apa-apa
    if (!get().currentTrack) return;
    set((state) => ({ isPlaying: !state.isPlaying }));
  },

  setIsPlaying: (playing) => set({ isPlaying: playing }),

  setProgress: (seconds) => set({ progress: seconds }),

  seekTo: (seconds) =>
    set((state) => ({
      progress: seconds,
      seekRequestId: state.seekRequestId + 1,
    })),

  setVolume: (volume) =>
    set((state) => {
      const clamped = Math.min(1, Math.max(0, volume));
      return {
        volume: clamped,
        // Ikut catat sebagai previousVolume asal bukan nol, supaya tombol
        // mute nanti tahu harus kembali ke level berapa.
        previousVolume: clamped > 0 ? clamped : state.previousVolume,
      };
    }),

  toggleMute: () => set((state) => (state.volume > 0 ? { volume: 0 } : { volume: state.previousVolume > 0 ? state.previousVolume : 1 })),

  playNext: () => {
    const { queue, shuffledQueue, isShuffled, currentIndex, repeatMode, currentTrack, seekRequestId } = get();
    const activeQueue = isShuffled ? shuffledQueue : queue;
    if (activeQueue.length === 0) return;

    let nextIndex = currentIndex + 1;

    if (nextIndex >= activeQueue.length) {
      if (repeatMode === "all") {
        nextIndex = 0;
      } else {
        // Sudah di lagu terakhir dan tidak repeat, berhenti
        set({ isPlaying: false });
        return;
      }
    }

    const next = activeQueue[nextIndex];
    set({
      currentTrack: next,
      currentIndex: nextIndex,
      ...(next.id === currentTrack?.id ? restartPatch(seekRequestId) : { progress: 0, isPlaying: true, playbackError: null }),
    });
  },

  playPrevious: () => {
    const { queue, shuffledQueue, isShuffled, currentIndex, progress, currentTrack, seekRequestId } = get();
    const activeQueue = isShuffled ? shuffledQueue : queue;
    if (activeQueue.length === 0) return;

    // Kalau lagu sudah berjalan >3 detik, tombol previous mengulang dari
    // awal dulu (perilaku sama seperti Spotify), bukan langsung ke lagu sebelumnya.
    // Sebelumnya hanya me-reset progress di store, sehingga audio asli tidak
    // ikut kembali ke awal; sekarang lewat seek supaya benar-benar mengulang.
    if (progress > 3) {
      set({ progress: 0, seekRequestId: seekRequestId + 1 });
      return;
    }

    const prevIndex = currentIndex - 1 < 0 ? 0 : currentIndex - 1;
    const prev = activeQueue[prevIndex];
    set({
      currentTrack: prev,
      currentIndex: prevIndex,
      ...(prev.id === currentTrack?.id ? restartPatch(seekRequestId) : { progress: 0, isPlaying: true, playbackError: null }),
    });
  },

  toggleShuffle: () =>
    set((state) => {
      const nextIsShuffled = !state.isShuffled;

      if (!state.currentTrack) {
        return { isShuffled: nextIsShuffled };
      }

      if (nextIsShuffled) {
        // Baru dinyalakan: acak ulang, currentTrack tetap di posisi 0
        const shuffledQueue = buildShuffledQueue(state.queue, state.currentTrack);
        return { isShuffled: true, shuffledQueue, currentIndex: 0 };
      }

      // Dimatikan: balik cari posisi currentTrack di urutan asli
      const idx = state.queue.findIndex((t) => t.id === state.currentTrack!.id);
      return { isShuffled: false, currentIndex: idx === -1 ? 0 : idx };
    }),

  cycleRepeatMode: () =>
    set((state) => {
      const order: RepeatMode[] = ["off", "all", "one"];
      const next = order[(order.indexOf(state.repeatMode) + 1) % order.length];
      return { repeatMode: next };
    }),
}));
