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

  // Actions
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

  // Mulai putar sebuah track. Kalau dikasih `queue` (misal daftar lagu di
  // satu playlist), track lain di queue itu jadi antrian next/previous.
  playTrack: (track, queue) => {
    const newQueue = queue ?? [track];
    const { isShuffled } = get();
    const shuffledQueue = buildShuffledQueue(newQueue, track);
    const activeQueue = isShuffled ? shuffledQueue : newQueue;
    const index = activeQueue.findIndex((t) => t.id === track.id);

    set({
      currentTrack: track,
      queue: newQueue,
      shuffledQueue,
      currentIndex: index === -1 ? 0 : index,
      progress: 0,
      isPlaying: true,
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
    const { queue, shuffledQueue, isShuffled, currentIndex, repeatMode } = get();
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

    set({
      currentTrack: activeQueue[nextIndex],
      currentIndex: nextIndex,
      progress: 0,
      isPlaying: true,
    });
  },

  playPrevious: () => {
    const { queue, shuffledQueue, isShuffled, currentIndex, progress } = get();
    const activeQueue = isShuffled ? shuffledQueue : queue;
    if (activeQueue.length === 0) return;

    // Kalau lagu sudah berjalan >3 detik, tombol previous mengulang dari
    // awal dulu (perilaku sama seperti Spotify), bukan langsung ke lagu sebelumnya.
    if (progress > 3) {
      set({ progress: 0 });
      return;
    }

    const prevIndex = currentIndex - 1 < 0 ? 0 : currentIndex - 1;
    set({
      currentTrack: activeQueue[prevIndex],
      currentIndex: prevIndex,
      progress: 0,
      isPlaying: true,
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
