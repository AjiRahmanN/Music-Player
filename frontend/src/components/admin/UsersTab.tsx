"use client";
import { useEffect, useState } from "react";
import { adminService, type UpdateUserData } from "@/services/adminService";
import type { AdminUser } from "@/types/admin";
import { useAuthStore } from "@/store/useAuthStore";
import { Loader2, AlertCircle, Pencil, Trash2, X, Check, ShieldCheck } from "lucide-react";

export default function UsersTab() {
  const currentUser = useAuthStore((s) => s.user);

  const [users, setUsers] = useState<AdminUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<UpdateUserData>({ username: "", email: "", isAdmin: false });
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const fetchUsers = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await adminService.listUsers();
      setUsers(res.data?.users ?? []);
    } catch (err: any) {
      setError(err.response?.data?.message || "Gagal memuat daftar user.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const openEdit = (user: AdminUser) => {
    setEditingId(user.id);
    setForm({ username: user.username, email: user.email, isAdmin: user.isAdmin });
    setFormError(null);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setFormError(null);
  };

  const handleSave = async (id: string) => {
    setFormError(null);
    if (!form.username.trim() || !form.email.trim()) {
      setFormError("Username dan email wajib diisi.");
      return;
    }

    setSaving(true);
    try {
      await adminService.updateUser(id, form);
      setEditingId(null);
      fetchUsers();
    } catch (err: any) {
      setFormError(err.response?.data?.message || "Gagal menyimpan perubahan.");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (user: AdminUser) => {
    if (!confirm(`Hapus user "${user.username}"? Aksi ini tidak bisa dibatalkan.`)) return;
    setDeletingId(user.id);
    try {
      await adminService.deleteUser(user.id);
      fetchUsers();
    } catch (err: any) {
      alert(err.response?.data?.message || "Gagal menghapus user.");
    } finally {
      setDeletingId(null);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center py-16 gap-3 text-zinc-500">
        <Loader2 className="animate-spin" size={28} />
        <p>Memuat user...</p>
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

  return (
    <div className="flex flex-col gap-1">
      {users.map((user) => {
        const isEditing = editingId === user.id;
        const isSelf = currentUser?.id === user.id;

        return (
          <div key={user.id} className="rounded-lg hover:bg-zinc-900 transition">
            {isEditing ? (
              <div className="p-3 bg-zinc-900 rounded-lg space-y-3">
                <div className="grid grid-cols-2 gap-3">
                  <input
                    type="text"
                    value={form.username}
                    onChange={(e) => setForm({ ...form, username: e.target.value })}
                    placeholder="Username"
                    className="p-2 rounded-lg bg-black/40 text-white border border-zinc-700 focus:outline-none focus:ring-2 focus:ring-green-500/30 focus:border-green-500 text-sm"
                  />
                  <input
                    type="email"
                    value={form.email}
                    onChange={(e) => setForm({ ...form, email: e.target.value })}
                    placeholder="Email"
                    className="p-2 rounded-lg bg-black/40 text-white border border-zinc-700 focus:outline-none focus:ring-2 focus:ring-green-500/30 focus:border-green-500 text-sm"
                  />
                </div>

                <label className="flex items-center gap-2 text-sm text-zinc-300">
                  <input type="checkbox" checked={form.isAdmin} onChange={(e) => setForm({ ...form, isAdmin: e.target.checked })} disabled={isSelf} className="accent-green-500" />
                  Admin
                  {isSelf && <span className="text-xs text-zinc-500">(tidak bisa ubah status sendiri)</span>}
                </label>

                {formError && (
                  <div className="flex items-center gap-2 text-red-300 text-xs">
                    <AlertCircle size={14} />
                    {formError}
                  </div>
                )}

                <div className="flex gap-2">
                  <button onClick={cancelEdit} disabled={saving} className="flex items-center gap-1.5 border border-zinc-700 hover:bg-zinc-800 text-zinc-300 px-3 py-1.5 rounded-full text-xs transition">
                    <X size={14} />
                    Batal
                  </button>
                  <button onClick={() => handleSave(user.id)} disabled={saving} className="flex items-center gap-1.5 bg-green-500 hover:bg-green-400 text-black font-bold px-3 py-1.5 rounded-full text-xs transition disabled:opacity-50">
                    {saving ? <Loader2 className="animate-spin" size={14} /> : <Check size={14} />}
                    Simpan
                  </button>
                </div>
              </div>
            ) : (
              <div className="flex items-center gap-3 px-3 py-2.5">
                <div className="min-w-0 flex-1">
                  <p className="text-sm text-white truncate flex items-center gap-1.5">
                    {user.username}
                    {user.isAdmin && <ShieldCheck size={13} className="text-green-500" />}
                    {isSelf && <span className="text-xs text-zinc-500">(kamu)</span>}
                  </p>
                  <p className="text-xs text-zinc-400 truncate">{user.email}</p>
                </div>
                <button onClick={() => openEdit(user)} className="text-zinc-400 hover:text-white transition p-1.5" aria-label={`Edit ${user.username}`}>
                  <Pencil size={16} />
                </button>
                <button
                  onClick={() => handleDelete(user)}
                  disabled={isSelf || deletingId === user.id}
                  className="text-zinc-400 hover:text-red-400 transition p-1.5 disabled:opacity-30 disabled:cursor-not-allowed"
                  aria-label={`Hapus ${user.username}`}
                  title={isSelf ? "Tidak bisa menghapus akun sendiri" : undefined}
                >
                  {deletingId === user.id ? <Loader2 className="animate-spin" size={16} /> : <Trash2 size={16} />}
                </button>
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
