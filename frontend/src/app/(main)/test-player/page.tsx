"use client";

// Sesuaikan path import ini dengan lokasi file usePlayerStore kamu
import { usePlayerStore } from "@/store/usePlayerStore";
import type { Track } from "@/types/player";

// Dummy queue berisi 3 lagu, supaya efek shuffle & next/previous kelihatan.
// audioUrl pakai file contoh publik dari SoundHelix, cuma buat testing.
const dummyQueue: Track[] = [
  {
    id: "1",
    title: "Lagu A",
    artist: "Artis Percobaan 1",
    audioUrl: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3",
    coverUrl: "https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=150&h=150&fit=crop",
    duration: 372,
  },
  {
    id: "2",
    title: "Lagu B",
    artist: "Artis Percobaan 2",
    audioUrl: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-2.mp3",
    coverUrl: "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=150&h=150&fit=crop",
    duration: 320,
  },
  {
    id: "3",
    title: "Lagu C",
    artist: "Artis Percobaan 3",
    audioUrl: "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-3.mp3",
    coverUrl: "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=150&h=150&fit=crop",
    duration: 240,
  },
];

export default function TestPlayerPage() {
  const playTrack = usePlayerStore((state) => state.playTrack);

  const handleTestPlay = (track: Track) => {
    console.log("1. Tombol diklik, playTrack dipanggil:", track.title);
    playTrack(track, dummyQueue);
  };

  return (
    <div className="p-10 text-white flex flex-col items-start gap-4">
      <h1 className="text-3xl font-bold">Halaman Test Player</h1>
      <p className="text-zinc-400">Klik salah satu lagu di bawah untuk mulai memutar. Ketiganya masuk ke satu antrian yang sama, jadi tombol next/previous/shuffle di PlayerBar bisa dites beneran.</p>

      <div className="flex flex-col gap-2 w-full max-w-sm">
        {dummyQueue.map((track) => (
          <button key={track.id} onClick={() => handleTestPlay(track)} className="flex items-center gap-3 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 rounded-lg p-3 text-left transition">
            <img src={track.coverUrl} alt={track.title} className="w-12 h-12 rounded object-cover" />
            <div>
              <p className="font-medium">{track.title}</p>
              <p className="text-sm text-zinc-400">{track.artist}</p>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
}
