"use client";
import { useEffect, useState } from "react";
import { authService } from "@/services/authService";
import { useAuthStore } from "@/store/useAuthStore";
import type { MeResponse } from "@/types/auth";
import { User as UserIcon, Mail, Calendar, ShieldCheck, Pencil, X, Check, Loader2, AlertCircle, CheckCircle2 } from "lucide-react";

const usernameValid = (u: string) => /^[a-zA-Z0-9_]{3,50}$/.test(u);
const emailValid = (e: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(e);

export default function ProfilePage() {
  const [profile, setProfile] = useState<MeResponse | null>(null);
  const [loading, setLoading] = useState(true);

  const [isEditing, setIsEditing] = useState(false);
  const [form, setForm] = useState({ username: "", email: "" });
  const [touched, setTouched] = useState({ username: false, email: false });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const setStoreUser = useAuthStore((state) => state.setUser);
  const storeUser = useAuthStore((state) => state.user);

  const fetchProfile = async () => {
    try {
      // Backend GET /me mengembalikan objek flat, lihat catatan di authService.getMe()
      const data = (await authService.getMe()) as unknown as MeResponse;
      setProfile(data);
    } catch (err) {
      console.error("Gagal mengambil data profil");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProfile();
  }, []);

  const usernameOk = usernameValid(form.username);
  const emailOk = emailValid(form.email);
  const formValid = usernameOk && emailOk;
  const hasChanges = profile ? form.username !== profile.username || form.email !== profile.email : false;

  const startEdit = () => {
    if (!profile) return;
    setForm({ username: profile.username, email: profile.email });
    setTouched({ username: false, email: false });
    setError(null);
    setIsEditing(true);
  };

  const cancelEdit = () => {
    setIsEditing(false);
    setError(null);
  };

  const handleSave = async () => {
    setTouched({ username: true, email: true });
    setError(null);

    if (!formValid) {
      setError("Periksa lagi data yang kamu masukkan.");
      return;
    }
    if (!hasChanges) {
      setIsEditing(false);
      return;
    }

    setSaving(true);
    try {
      await authService.updateProfile(form);

      // Refresh data profil dari server supaya konsisten
      await fetchProfile();

      // Sinkronkan username/email yang tampil di Navbar (Zustand store)
      if (storeUser) {
        setStoreUser({ ...storeUser, username: form.username, email: form.email });
      }

      setIsEditing(false);
      setSuccess(true);
      setTimeout(() => setSuccess(false), 3000);
    } catch (err: any) {
      setError(err.response?.data?.message || err.response?.data?.error || "Gagal menyimpan perubahan. Silakan coba lagi.");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-20 gap-3 text-zinc-500">
        <Loader2 className="animate-spin" size={28} />
        <p>Memuat profil...</p>
      </div>
    );
  }

  if (!profile) {
    return <div className="text-center py-20 text-zinc-500">Gagal memuat profil.</div>;
  }

  return (
    <div className="max-w-2xl mx-auto mt-10">
      {/* Toast sukses */}
      {success && (
        <div className="fixed top-4 right-4 z-50 animate-slideIn">
          <div className="bg-green-900/90 backdrop-blur-sm border border-green-700 rounded-lg shadow-xl p-4 flex items-center gap-3">
            <CheckCircle2 className="text-green-400 flex-shrink-0" size={20} />
            <p className="text-sm text-white">Profil berhasil diperbarui</p>
          </div>
        </div>
      )}

      <div className="bg-zinc-900 rounded-3xl p-8 border border-zinc-800">
        <div className="flex flex-col items-center mb-10 text-center">
          <div className="w-24 h-24 bg-green-500 rounded-full flex items-center justify-center mb-4 text-black shadow-xl shadow-green-500/20">
            <UserIcon size={48} />
          </div>
          {!isEditing && (
            <>
              <h1 className="text-3xl font-bold">{profile.username}</h1>
              <p className="text-zinc-400">Member Melodify Free</p>
            </>
          )}
        </div>

        {isEditing ? (
          <div className="space-y-5">
            {/* Username */}
            <div>
              <label className="flex items-center gap-2 text-sm font-medium text-zinc-300 mb-2">
                <UserIcon size={16} />
                Username
              </label>
              <input
                type="text"
                value={form.username}
                onChange={(e) => setForm({ ...form, username: e.target.value })}
                onBlur={() => setTouched((t) => ({ ...t, username: true }))}
                disabled={saving}
                className={`w-full p-3 rounded-lg bg-black/40 text-white border focus:outline-none focus:ring-2 transition ${
                  touched.username && !usernameOk ? "border-red-700 focus:ring-red-500/30 focus:border-red-500" : "border-zinc-700 focus:ring-green-500/30 focus:border-green-500"
                }`}
              />
              {touched.username && !usernameOk && <p className="text-red-400 text-xs mt-1.5">3-50 karakter, hanya huruf, angka, dan underscore.</p>}
            </div>

            {/* Email */}
            <div>
              <label className="flex items-center gap-2 text-sm font-medium text-zinc-300 mb-2">
                <Mail size={16} />
                Email
              </label>
              <input
                type="email"
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
                onBlur={() => setTouched((t) => ({ ...t, email: true }))}
                disabled={saving}
                className={`w-full p-3 rounded-lg bg-black/40 text-white border focus:outline-none focus:ring-2 transition ${
                  touched.email && !emailOk ? "border-red-700 focus:ring-red-500/30 focus:border-red-500" : "border-zinc-700 focus:ring-green-500/30 focus:border-green-500"
                }`}
              />
              {touched.email && !emailOk && <p className="text-red-400 text-xs mt-1.5">Format email tidak valid.</p>}
            </div>

            {error && (
              <div className="bg-red-900/30 border border-red-800 rounded-lg p-3">
                <div className="flex items-center gap-2 text-red-300">
                  <AlertCircle size={16} className="flex-shrink-0" />
                  <span className="text-sm">{error}</span>
                </div>
              </div>
            )}

            <div className="flex gap-3 pt-2">
              <button onClick={cancelEdit} disabled={saving} className="flex-1 py-3 rounded-full border border-zinc-700 font-bold hover:bg-zinc-800 transition flex items-center justify-center gap-2 disabled:opacity-50">
                <X size={18} />
                Batal
              </button>
              <button
                onClick={handleSave}
                disabled={saving}
                className="flex-1 py-3 rounded-full bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 font-bold transition flex items-center justify-center gap-2 disabled:opacity-50"
              >
                {saving ? (
                  <>
                    <Loader2 className="animate-spin" size={18} />
                    Menyimpan...
                  </>
                ) : (
                  <>
                    <Check size={18} />
                    Simpan
                  </>
                )}
              </button>
            </div>
          </div>
        ) : (
          <>
            <div className="space-y-6">
              <div className="flex items-center gap-4 p-4 bg-black/40 rounded-2xl">
                <Mail className="text-green-500" />
                <div>
                  <p className="text-xs text-zinc-500 uppercase font-bold tracking-wider">Email</p>
                  <p className="text-lg">{profile.email}</p>
                </div>
              </div>

              <div className="flex items-center gap-4 p-4 bg-black/40 rounded-2xl">
                <ShieldCheck className="text-green-500" />
                <div>
                  <p className="text-xs text-zinc-500 uppercase font-bold tracking-wider">User ID</p>
                  <p className="text-sm font-mono text-zinc-300">{profile.id}</p>
                </div>
              </div>

              <div className="flex items-center gap-4 p-4 bg-black/40 rounded-2xl text-zinc-500">
                <Calendar />
                <p>
                  Bergabung sejak{" "}
                  {profile.created
                    ? new Date(profile.created).toLocaleDateString("id-ID", {
                        day: "numeric",
                        month: "long",
                        year: "numeric",
                      })
                    : "-"}
                </p>
              </div>
            </div>

            <button onClick={startEdit} className="w-full mt-10 py-3 rounded-full border border-zinc-700 font-bold hover:bg-zinc-800 transition flex items-center justify-center gap-2">
              <Pencil size={16} />
              Edit Profil
            </button>
          </>
        )}
      </div>
    </div>
  );
}
