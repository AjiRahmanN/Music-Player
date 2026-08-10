"use client";
import { useEffect, useState } from "react";
import { songService, mapSongToTrack } from "@/services/songService";
import TrackList from "@/components/player/TrackList";
import type { Track } from "@/types/player";
import { Loader2, AlertCircle, Music, Search, X } from "lucide-react";

export default function HomePage() {
  const [tracks, setTracks] = useState<Track[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchInput, setSearchInput] = useState("");
  // Nilai yang sudah "didiamkan" 400ms setelah user berhenti mengetik,
  // supaya tidak fetch API di setiap ketukan huruf.
  const [debouncedQuery, setDebouncedQuery] = useState("");

  useEffect(() => {
    const timeout = setTimeout(() => setDebouncedQuery(searchInput.trim()), 400);
    return () => clearTimeout(timeout);
  }, [searchInput]);

  useEffect(() => {
    const fetchSongs = async () => {
      setLoading(true);
      setError(null);
      try {
        const res = await songService.list(debouncedQuery || undefined);
        const songs = res.data?.songs ?? [];
        setTracks(songs.map(mapSongToTrack));
      } catch (err: any) {
        setError(err.response?.data?.message || err.response?.data?.error || "Gagal memuat daftar lagu.");
      } finally {
        setLoading(false);
      }
    };

    fetchSongs();
  }, [debouncedQuery]);

  return (
    <div className="max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-6 gap-4">
        <h1 className="text-3xl font-bold">{debouncedQuery ? "Hasil Pencarian" : "Semua Lagu"}</h1>

        <div className="relative w-72">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-500" size={16} />
          <input
            type="text"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            placeholder="Cari judul atau artis..."
            className="w-full bg-zinc-800/60 border border-zinc-700 rounded-full py-2 pl-9 pr-9 text-sm text-white placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-green-500/30 focus:border-green-500 transition"
          />
          {searchInput && (
            <button onClick={() => setSearchInput("")} className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-white transition" aria-label="Bersihkan pencarian">
              <X size={16} />
            </button>
          )}
        </div>
      </div>

      {loading && (
        <div className="flex flex-col items-center justify-center py-20 gap-3 text-zinc-500">
          <Loader2 className="animate-spin" size={28} />
          <p>Memuat lagu...</p>
        </div>
      )}

      {!loading && error && (
        <div className="bg-red-900/30 border border-red-800 rounded-lg p-4 flex items-center gap-3">
          <AlertCircle className="text-red-400 flex-shrink-0" size={20} />
          <p className="text-sm text-red-300">{error}</p>
        </div>
      )}

      {!loading && !error && tracks.length === 0 && (
        <div className="flex flex-col items-center justify-center py-20 gap-3 text-zinc-500">
          <Music size={32} />
          <p>{debouncedQuery ? `Tidak ada lagu untuk "${debouncedQuery}"` : "Belum ada lagu."}</p>
        </div>
      )}

      {!loading && !error && tracks.length > 0 && <TrackList tracks={tracks} />}
    </div>
  );
}
