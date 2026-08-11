"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { playlistService } from "@/services/playlistService";
import type { Playlist } from "@/types/playlist";
import { Loader2, AlertCircle, Plus, ListMusic, X } from "lucide-react";

export default function PlaylistsPage() {
  const router = useRouter();
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [creating, setCreating] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const fetchPlaylists = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await playlistService.list();
      setPlaylists(res.data?.playlists ?? []);
    } catch (err: any) {
      setError(err.response?.data?.message || err.response?.data?.error || "Gagal memuat playlist.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPlaylists();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);

    if (!name.trim()) {
      setFormError("Nama playlist tidak boleh kosong.");
      return;
    }

    setCreating(true);
    try {
      const res = await playlistService.create({ name, description });
      const created = res.data?.playlist;
      setShowForm(false);
      setName("");
      setDescription("");
      if (created) {
        router.push(`/playlists/${created.id}`);
      } else {
        fetchPlaylists();
      }
    } catch (err: any) {
      setFormError(err.response?.data?.message || err.response?.data?.error || "Gagal membuat playlist.");
    } finally {
      setCreating(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Playlist Kamu</h1>
        <button onClick={() => setShowForm((v) => !v)} className="flex items-center gap-2 bg-green-500 hover:bg-green-400 text-black font-bold px-4 py-2 rounded-full transition">
          {showForm ? <X size={18} /> : <Plus size={18} />}
          {showForm ? "Batal" : "Buat Playlist"}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6 mb-8 space-y-4">
          <div>
            <label className="text-sm font-medium text-zinc-300 mb-1.5 block">Nama Playlist</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Contoh: Lagu Galau"
              disabled={creating}
              className="w-full p-3 rounded-lg bg-black/40 text-white border border-zinc-700 focus:outline-none focus:ring-2 focus:ring-green-500/30 focus:border-green-500 transition"
            />
          </div>
          <div>
            <label className="text-sm font-medium text-zinc-300 mb-1.5 block">
              Deskripsi <span className="text-zinc-500">(opsional)</span>
            </label>
            <input
              type="text"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Playlist buat nangis di kosan"
              disabled={creating}
              className="w-full p-3 rounded-lg bg-black/40 text-white border border-zinc-700 focus:outline-none focus:ring-2 focus:ring-green-500/30 focus:border-green-500 transition"
            />
          </div>

          {formError && (
            <div className="flex items-center gap-2 text-red-300 text-sm">
              <AlertCircle size={16} />
              {formError}
            </div>
          )}

          <button type="submit" disabled={creating} className="bg-green-500 hover:bg-green-400 text-black font-bold px-5 py-2.5 rounded-full transition disabled:opacity-50 flex items-center gap-2">
            {creating && <Loader2 className="animate-spin" size={16} />}
            {creating ? "Membuat..." : "Buat Playlist"}
          </button>
        </form>
      )}

      {loading && (
        <div className="flex flex-col items-center justify-center py-20 gap-3 text-zinc-500">
          <Loader2 className="animate-spin" size={28} />
          <p>Memuat playlist...</p>
        </div>
      )}

      {!loading && error && (
        <div className="bg-red-900/30 border border-red-800 rounded-lg p-4 flex items-center gap-3">
          <AlertCircle className="text-red-400 flex-shrink-0" size={20} />
          <p className="text-sm text-red-300">{error}</p>
        </div>
      )}

      {!loading && !error && playlists.length === 0 && (
        <div className="flex flex-col items-center justify-center py-20 gap-3 text-zinc-500">
          <ListMusic size={32} />
          <p>Belum ada playlist. Yuk buat yang pertama.</p>
        </div>
      )}

      {!loading && !error && playlists.length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
          {playlists.map((playlist) => (
            <button key={playlist.id} onClick={() => router.push(`/playlists/${playlist.id}`)} className="bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 rounded-xl p-4 text-left transition flex flex-col gap-3">
              <div className="w-full aspect-square bg-zinc-800 rounded-lg flex items-center justify-center overflow-hidden">
                {playlist.coverUrl ? <img src={playlist.coverUrl} alt={playlist.name} className="w-full h-full object-cover" /> : <ListMusic className="text-zinc-600" size={32} />}
              </div>
              <div className="min-w-0">
                <p className="font-medium text-white truncate">{playlist.name}</p>
                <p className="text-xs text-zinc-400">{playlist.songCount} lagu</p>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
