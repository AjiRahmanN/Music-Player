"use client";
import { useEffect, useState } from "react";
import { adminService } from "@/services/adminService";
import type { AdminPlaylist } from "@/types/admin";
import { Loader2, AlertCircle, Trash2, ListMusic } from "lucide-react";

export default function PlaylistsTab() {
  const [playlists, setPlaylists] = useState<AdminPlaylist[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const fetchPlaylists = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await adminService.listAllPlaylists();
      setPlaylists(res.data?.playlists ?? []);
    } catch (err: any) {
      setError(err.response?.data?.message || "Gagal memuat daftar playlist.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPlaylists();
  }, []);

  const handleDelete = async (playlist: AdminPlaylist) => {
    if (!confirm(`Hapus playlist "${playlist.name}" milik ${playlist.ownerUsername}?`)) return;
    setDeletingId(playlist.id);
    try {
      await adminService.deletePlaylist(playlist.id);
      fetchPlaylists();
    } catch (err: any) {
      alert(err.response?.data?.message || "Gagal menghapus playlist.");
    } finally {
      setDeletingId(null);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-16 gap-3 text-zinc-500">
        <Loader2 className="animate-spin" size={28} />
        <p>Memuat playlist...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-900/30 border border-red-800 rounded-lg p-4 flex items-center gap-3">
        <AlertCircle className="text-red-400 flex-shrink-0" size={20} />
        <p className="text-sm text-red-300">{error}</p>
      </div>
    );
  }

  if (playlists.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 gap-3 text-zinc-500">
        <ListMusic size={32} />
        <p>Belum ada playlist dari user mana pun.</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-1">
      {playlists.map((playlist) => (
        <div key={playlist.id} className="flex items-center gap-3 px-3 py-2.5 rounded-lg hover:bg-zinc-900 transition">
          <div className="w-10 h-10 bg-zinc-800 rounded flex items-center justify-center flex-shrink-0 overflow-hidden">
            {playlist.coverUrl ? <img src={playlist.coverUrl} alt={playlist.name} className="w-full h-full object-cover" /> : <ListMusic className="text-zinc-600" size={18} />}
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-sm text-white truncate">{playlist.name}</p>
            <p className="text-xs text-zinc-400 truncate">
              oleh {playlist.ownerUsername} · {playlist.songCount} lagu
            </p>
          </div>
          <button onClick={() => handleDelete(playlist)} disabled={deletingId === playlist.id} className="text-zinc-400 hover:text-red-400 transition p-1.5 disabled:opacity-50" aria-label={`Hapus playlist ${playlist.name}`}>
            {deletingId === playlist.id ? <Loader2 className="animate-spin" size={16} /> : <Trash2 size={16} />}
          </button>
        </div>
      ))}
    </div>
  );
}
