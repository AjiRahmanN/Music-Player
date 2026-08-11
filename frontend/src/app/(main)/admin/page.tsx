"use client";
import { useEffect, useState } from "react";
import { adminService } from "@/services/adminService";
import SongsTab from "@/components/admin/SongsTab";
import UsersTab from "@/components/admin/UsersTab";
import PlaylistsTab from "@/components/admin/PlaylistsTab";
import { Loader2, ShieldAlert, Music, Users, ListMusic } from "lucide-react";

type Tab = "songs" | "users" | "playlists";

const tabs: { key: Tab; label: string; icon: typeof Music }[] = [
  { key: "songs", label: "Lagu", icon: Music },
  { key: "users", label: "User", icon: Users },
  { key: "playlists", label: "Playlist", icon: ListMusic },
];

export default function AdminPage() {
  const [isAdmin, setIsAdmin] = useState<boolean | null>(null);
  const [activeTab, setActiveTab] = useState<Tab>("songs");

  useEffect(() => {
    adminService
      .checkAccess()
      .then(() => setIsAdmin(true))
      .catch(() => setIsAdmin(false));
  }, []);

  if (isAdmin === null) {
    return (
      <div className="flex flex-col items-center justify-center py-20 gap-3 text-zinc-500">
        <Loader2 className="animate-spin" size={28} />
        <p>Memeriksa akses...</p>
      </div>
    );
  }

  if (isAdmin === false) {
    return (
      <div className="max-w-md mx-auto text-center py-20">
        <ShieldAlert className="mx-auto text-red-400 mb-4" size={40} />
        <h1 className="text-xl font-bold text-white mb-2">Akses Ditolak</h1>
        <p className="text-zinc-400 text-sm">Halaman ini hanya bisa diakses oleh admin.</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">Panel Admin</h1>

      <div className="flex gap-1 mb-6 border-b border-zinc-800">
        {tabs.map(({ key, label, icon: Icon }) => (
          <button
            key={key}
            onClick={() => setActiveTab(key)}
            className={`flex items-center gap-2 px-4 py-2.5 text-sm font-medium border-b-2 transition ${activeTab === key ? "border-green-500 text-white" : "border-transparent text-zinc-400 hover:text-white"}`}
          >
            <Icon size={16} />
            {label}
          </button>
        ))}
      </div>

      {activeTab === "songs" && <SongsTab />}
      {activeTab === "users" && <UsersTab />}
      {activeTab === "playlists" && <PlaylistsTab />}
    </div>
  );
}
