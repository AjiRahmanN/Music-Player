"use client";
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { playlistService } from "@/services/playlistService";
import { songService, mapSongToTrack } from "@/services/songService";
import TrackList from "@/components/player/TrackList";
import type { PlaylistDetail } from "@/types/playlist";
import type { Song } from "@/types/song";
import { Loader2, AlertCircle, ListMusic, Trash2, Plus, Check } from "lucide-react";

export default function PlaylistDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();

  const [playlist, setPlaylist] = useState<PlaylistDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [allSongs, setAllSongs] = useState<Song[]>([]);
  const [showAddSongs, setShowAddSongs] = useState(false);
  const [addingSongId, setAddingSongId] = useState<string | null>(null);
  const [deleting, setDeleting] = useState(false);

  const fetchPlaylist = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await playlistService.getDetail(id);
      setPlaylist(res.data?.playlist ?? null);
    } catch (err: any) {
      setError(err.response?.data?.message || err.response?.data?.error || "Gagal memuat playlist.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPlaylist();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  const openAddSongs = async () => {
    setShowAddSongs(true);
    if (allSongs.length === 0) {
      try {
        const res = await songService.list();
        setAllSongs(res.data?.songs ?? []);
      } catch {
        // biarkan silent, seksi tambah lagu akan tampil kosong
      }
    }
  };

  const handleAddSong = async (songId: string) => {
    setAddingSongId(songId);
    try {
      await playlistService.addSong(id, songId);
      await fetchPlaylist();
    } catch (err: any) {
      alert(err.response?.data?.message || "Gagal menambahkan lagu.");
    } finally {
      setAddingSongId(null);
    }
  };

  const handleRemoveSong = async (songId: string) => {
    try {
      await playlistService.removeSong(id, songId);
      fetchPlaylist();
    } catch (err: any) {
      alert(err.response?.data?.message || "Gagal menghapus lagu dari playlist.");
    }
  };

  const handleDeletePlaylist = async () => {
    if (!confirm(`Hapus playlist "${playlist?.name}"? Aksi ini tidak bisa dibatalkan.`)) return;
    setDeleting(true);
    try {
      await playlistService.remove(id);
      router.push("/playlists");
    } catch (err: any) {
      alert(err.response?.data?.message || "Gagal menghapus playlist.");
      setDeleting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-20 gap-3 text-zinc-500">
        <Loader2 className="animate-spin" size={28} />
        <p>Memuat playlist...</p>
      </div>
    );
  }

  if (error || !playlist) {
    return (
      <div className="max-w-2xl mx-auto bg-red-900/30 border border-red-800 rounded-lg p-4 flex items-center gap-3">
        <AlertCircle className="text-red-400 flex-shrink-0" size={20} />
        <p className="text-sm text-red-300">{error ?? "Playlist tidak ditemukan."}</p>
      </div>
    );
  }

  const songIdsInPlaylist = new Set(playlist.songs.map((s) => s.id));
  const availableSongs = allSongs.filter((s) => !songIdsInPlaylist.has(s.id));

  return (
    <div className="max-w-4xl mx-auto">
      {/* Header */}
      <div className="flex items-end gap-6 mb-8">
        <div className="w-40 h-40 bg-zinc-800 rounded-xl flex items-center justify-center flex-shrink-0 overflow-hidden">
          {playlist.coverUrl ? <img src={playlist.coverUrl} alt={playlist.name} className="w-full h-full object-cover" /> : <ListMusic className="text-zinc-600" size={48} />}
        </div>
        <div className="min-w-0">
          <p className="text-xs text-zinc-400 uppercase font-bold tracking-wider mb-1">Playlist</p>
          <h1 className="text-4xl font-bold truncate">{playlist.name}</h1>
          {playlist.description && <p className="text-zinc-400 mt-2">{playlist.description}</p>}
          <p className="text-sm text-zinc-500 mt-2">{playlist.songs.length} lagu</p>
        </div>
      </div>

      <div className="flex items-center gap-3 mb-6">
        <button onClick={openAddSongs} className="flex items-center gap-2 bg-green-500 hover:bg-green-400 text-black font-bold px-4 py-2 rounded-full transition text-sm">
          <Plus size={16} />
          Tambah Lagu
        </button>
        <button
          onClick={handleDeletePlaylist}
          disabled={deleting}
          className="flex items-center gap-2 border border-zinc-700 hover:bg-red-900/30 hover:border-red-800 hover:text-red-300 text-zinc-400 px-4 py-2 rounded-full transition text-sm disabled:opacity-50"
        >
          <Trash2 size={16} />
          {deleting ? "Menghapus..." : "Hapus Playlist"}
        </button>
      </div>

      {/* Panel tambah lagu */}
      {showAddSongs && (
        <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-4 mb-8">
          <p className="text-sm font-medium text-zinc-300 mb-3">Tambah dari koleksi lagu</p>
          {availableSongs.length === 0 ? (
            <p className="text-sm text-zinc-500">Semua lagu sudah ada di playlist ini.</p>
          ) : (
            <div className="flex flex-col gap-1 max-h-64 overflow-y-auto">
              {availableSongs.map((song) => (
                <div key={song.id} className="flex items-center justify-between gap-3 px-2 py-2 rounded-md hover:bg-zinc-800/60">
                  <div className="flex items-center gap-3 min-w-0">
                    <img src={song.coverUrl} alt={song.title} className="w-9 h-9 rounded object-cover flex-shrink-0" />
                    <div className="min-w-0">
                      <p className="text-sm text-white truncate">{song.title}</p>
                      <p className="text-xs text-zinc-400 truncate">{song.artist}</p>
                    </div>
                  </div>
                  <button onClick={() => handleAddSong(song.id)} disabled={addingSongId === song.id} className="flex-shrink-0 text-zinc-400 hover:text-green-500 transition disabled:opacity-50" aria-label={`Tambah ${song.title}`}>
                    {addingSongId === song.id ? <Loader2 className="animate-spin" size={18} /> : <Plus size={18} />}
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Daftar lagu di playlist */}
      {playlist.songs.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-16 gap-3 text-zinc-500">
          <ListMusic size={32} />
          <p>Belum ada lagu di playlist ini.</p>
        </div>
      ) : (
        <TrackList tracks={playlist.songs.map(mapSongToTrack)} />
      )}
    </div>
  );
}
