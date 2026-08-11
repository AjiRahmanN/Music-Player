"use client";
import { useEffect, useState } from "react";
import { songService } from "@/services/songService";
import { adminService, type SongFormData } from "@/services/adminService";
import type { Song } from "@/types/song";
import { Loader2, AlertCircle, Plus, Pencil, Trash2, X, Check, Music, ImageIcon } from "lucide-react";

interface FormState {
  title: string;
  artist: string;
  album: string;
  duration: number;
}

const emptyForm: FormState = { title: "", artist: "", album: "", duration: 0 };

function formatDuration(seconds: number) {
  if (!seconds) return "-";
  const m = Math.floor(seconds / 60);
  const s = Math.floor(seconds % 60);
  return `${m}:${s.toString().padStart(2, "0")}`;
}

export default function SongsTab() {
  const [songs, setSongs] = useState<Song[]>([]);
  const [loadingSongs, setLoadingSongs] = useState(true);
  const [listError, setListError] = useState<string | null>(null);

  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm);

  const [coverFile, setCoverFile] = useState<File | null>(null);
  const [audioFile, setAudioFile] = useState<File | null>(null);
  const [existingCoverUrl, setExistingCoverUrl] = useState("");
  const [existingAudioUrl, setExistingAudioUrl] = useState("");
  const [coverPreview, setCoverPreview] = useState<string>("");
  const [detectingDuration, setDetectingDuration] = useState(false);

  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const fetchSongs = async () => {
    setLoadingSongs(true);
    setListError(null);
    try {
      const res = await songService.list();
      setSongs(res.data?.songs ?? []);
    } catch (err: any) {
      setListError(err.response?.data?.message || "Gagal memuat daftar lagu.");
    } finally {
      setLoadingSongs(false);
    }
  };

  useEffect(() => {
    fetchSongs();
  }, []);

  const resetForm = () => {
    setForm(emptyForm);
    setCoverFile(null);
    setAudioFile(null);
    setExistingCoverUrl("");
    setExistingAudioUrl("");
    setCoverPreview("");
    setFormError(null);
  };

  const openCreateForm = () => {
    setEditingId(null);
    resetForm();
    setShowForm(true);
  };

  const openEditForm = (song: Song) => {
    setEditingId(song.id);
    setForm({ title: song.title, artist: song.artist, album: song.album ?? "", duration: song.duration });
    setCoverFile(null);
    setAudioFile(null);
    setExistingCoverUrl(song.coverUrl);
    setExistingAudioUrl(song.audioUrl);
    setCoverPreview(song.coverUrl);
    setFormError(null);
    setShowForm(true);
  };

  const closeForm = () => {
    setShowForm(false);
    setEditingId(null);
    resetForm();
  };

  const handleCoverChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setCoverFile(file);
    setCoverPreview(URL.createObjectURL(file));
  };

  const handleAudioChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setAudioFile(file);
    setDetectingDuration(true);

    const objectUrl = URL.createObjectURL(file);
    const audio = new Audio();
    audio.src = objectUrl;
    audio.addEventListener("loadedmetadata", () => {
      setForm((f) => ({ ...f, duration: Math.round(audio.duration) }));
      setDetectingDuration(false);
      URL.revokeObjectURL(objectUrl);
    });
    audio.addEventListener("error", () => {
      setDetectingDuration(false);
      setFormError("Gagal membaca file audio. Pastikan formatnya benar (mp3/wav/ogg/m4a).");
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);

    if (!form.title.trim() || !form.artist.trim()) {
      setFormError("Judul dan artis wajib diisi.");
      return;
    }
    if (!editingId && !audioFile) {
      setFormError("File audio wajib dipilih untuk lagu baru.");
      return;
    }
    if (audioFile && form.duration <= 0) {
      setFormError("Durasi belum terbaca dari file audio, tunggu sebentar lalu coba lagi.");
      return;
    }

    setSaving(true);
    try {
      let coverUrl = existingCoverUrl;
      let audioUrl = existingAudioUrl;

      if (coverFile) {
        const res = await adminService.uploadFile(coverFile, "cover");
        coverUrl = res.data?.url ?? coverUrl;
      }
      if (audioFile) {
        const res = await adminService.uploadFile(audioFile, "audio");
        audioUrl = res.data?.url ?? audioUrl;
      }

      const payload: SongFormData = {
        title: form.title,
        artist: form.artist,
        album: form.album,
        coverUrl,
        audioUrl,
        duration: form.duration,
      };

      if (editingId) {
        await adminService.updateSong(editingId, payload);
      } else {
        await adminService.createSong(payload);
      }

      closeForm();
      fetchSongs();
    } catch (err: any) {
      setFormError(err.response?.data?.message || "Gagal menyimpan lagu.");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (song: Song) => {
    if (!confirm(`Hapus lagu "${song.title}"?`)) return;
    setDeletingId(song.id);
    try {
      await adminService.deleteSong(song.id);
      fetchSongs();
    } catch (err: any) {
      alert(err.response?.data?.message || "Gagal menghapus lagu.");
    } finally {
      setDeletingId(null);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-end mb-4">
        <button onClick={openCreateForm} className="flex items-center gap-2 bg-green-500 hover:bg-green-400 text-black font-bold px-4 py-2 rounded-full transition">
          <Plus size={18} />
          Tambah Lagu
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6 mb-8 space-y-5">
          <div className="flex items-center justify-between">
            <p className="font-medium text-white">{editingId ? "Edit Lagu" : "Lagu Baru"}</p>
            <button type="button" onClick={closeForm} className="text-zinc-500 hover:text-white">
              <X size={18} />
            </button>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm text-zinc-300 mb-1.5 block">
                Cover <span className="text-zinc-500">(opsional)</span>
              </label>
              <label className="flex flex-col items-center justify-center gap-2 border border-dashed border-zinc-700 rounded-lg h-32 cursor-pointer hover:border-green-500 transition overflow-hidden">
                {coverPreview ? (
                  <img src={coverPreview} alt="Preview cover" className="w-full h-full object-cover" />
                ) : (
                  <>
                    <ImageIcon className="text-zinc-500" size={24} />
                    <span className="text-xs text-zinc-500">Pilih gambar</span>
                  </>
                )}
                <input type="file" accept="image/png,image/jpeg,image/webp" onChange={handleCoverChange} className="hidden" />
              </label>
            </div>

            <div>
              <label className="text-sm text-zinc-300 mb-1.5 block">File Audio {!editingId && <span className="text-red-400">*</span>}</label>
              <label className="flex flex-col items-center justify-center gap-2 border border-dashed border-zinc-700 rounded-lg h-32 cursor-pointer hover:border-green-500 transition px-2 text-center">
                <Music className="text-zinc-500" size={24} />
                <span className="text-xs text-zinc-500 truncate max-w-full">{audioFile ? audioFile.name : editingId ? "Ganti file audio" : "Pilih file audio"}</span>
                <input type="file" accept="audio/mpeg,audio/wav,audio/ogg,audio/mp4,audio/x-m4a" onChange={handleAudioChange} className="hidden" />
              </label>
              <p className="text-xs text-zinc-500 mt-1.5 flex items-center gap-1.5">
                {detectingDuration ? (
                  <>
                    <Loader2 className="animate-spin" size={12} /> Membaca durasi...
                  </>
                ) : (
                  <>Durasi: {formatDuration(form.duration)}</>
                )}
              </p>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="col-span-2 sm:col-span-1">
              <label className="text-sm text-zinc-300 mb-1.5 block">Judul</label>
              <input
                type="text"
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
                className="w-full p-2.5 rounded-lg bg-black/40 text-white border border-zinc-700 focus:outline-none focus:ring-2 focus:ring-green-500/30 focus:border-green-500"
              />
            </div>
            <div className="col-span-2 sm:col-span-1">
              <label className="text-sm text-zinc-300 mb-1.5 block">Artis</label>
              <input
                type="text"
                value={form.artist}
                onChange={(e) => setForm({ ...form, artist: e.target.value })}
                className="w-full p-2.5 rounded-lg bg-black/40 text-white border border-zinc-700 focus:outline-none focus:ring-2 focus:ring-green-500/30 focus:border-green-500"
              />
            </div>
            <div className="col-span-2">
              <label className="text-sm text-zinc-300 mb-1.5 block">
                Album <span className="text-zinc-500">(opsional)</span>
              </label>
              <input
                type="text"
                value={form.album}
                onChange={(e) => setForm({ ...form, album: e.target.value })}
                className="w-full p-2.5 rounded-lg bg-black/40 text-white border border-zinc-700 focus:outline-none focus:ring-2 focus:ring-green-500/30 focus:border-green-500"
              />
            </div>
          </div>

          {formError && (
            <div className="flex items-center gap-2 text-red-300 text-sm">
              <AlertCircle size={16} />
              {formError}
            </div>
          )}

          <button type="submit" disabled={saving || detectingDuration} className="flex items-center gap-2 bg-green-500 hover:bg-green-400 text-black font-bold px-5 py-2.5 rounded-full transition disabled:opacity-50">
            {saving ? <Loader2 className="animate-spin" size={16} /> : <Check size={16} />}
            {saving ? "Mengunggah & menyimpan..." : editingId ? "Simpan Perubahan" : "Tambah Lagu"}
          </button>
        </form>
      )}

      {loadingSongs && (
        <div className="flex flex-col items-center justify-center py-16 gap-3 text-zinc-500">
          <Loader2 className="animate-spin" size={28} />
          <p>Memuat lagu...</p>
        </div>
      )}

      {!loadingSongs && listError && (
        <div className="bg-red-900/30 border border-red-800 rounded-lg p-4 flex items-center gap-3">
          <AlertCircle className="text-red-400 flex-shrink-0" size={20} />
          <p className="text-sm text-red-300">{listError}</p>
        </div>
      )}

      {!loadingSongs && !listError && (
        <div className="flex flex-col gap-1">
          {songs.map((song) => (
            <div key={song.id} className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-zinc-900 transition">
              <img src={song.coverUrl} alt={song.title} className="w-10 h-10 rounded object-cover flex-shrink-0 bg-zinc-800" />
              <div className="min-w-0 flex-1">
                <p className="text-sm text-white truncate">{song.title}</p>
                <p className="text-xs text-zinc-400 truncate">{song.artist}</p>
              </div>
              <span className="text-xs text-zinc-500 tabular-nums">{formatDuration(song.duration)}</span>
              <button onClick={() => openEditForm(song)} className="text-zinc-400 hover:text-white transition p-1.5" aria-label={`Edit ${song.title}`}>
                <Pencil size={16} />
              </button>
              <button onClick={() => handleDelete(song)} disabled={deletingId === song.id} className="text-zinc-400 hover:text-red-400 transition p-1.5 disabled:opacity-50" aria-label={`Hapus ${song.title}`}>
                {deletingId === song.id ? <Loader2 className="animate-spin" size={16} /> : <Trash2 size={16} />}
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
