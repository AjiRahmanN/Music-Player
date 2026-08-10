"use client";
import { useAuthStore } from "@/store/useAuthStore";
import { authService } from "@/services/authService";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { User, LogOut, Music } from "lucide-react";

export default function Navbar() {
  const { user, logout } = useAuthStore();
  const router = useRouter();

  const handleLogout = async () => {
    try {
      await authService.logout();
    } finally {
      logout();
    }
    router.push("/login");
  };

  return (
    <nav className="flex justify-between items-center p-4 bg-zinc-950 border-b border-zinc-800 sticky top-0 z-50">
      <Link href="/" className="flex items-center gap-2 font-bold text-xl text-green-500">
        <Music size={28} />
        <span>Melodify</span>
      </Link>

      <div className="flex items-center gap-6">
        {user ? (
          <>
            <Link href="/profile" className="flex items-center gap-2 text-zinc-300 hover:text-white transition">
              <User size={20} />
              <span>{user.username}</span>
            </Link>
            <button onClick={handleLogout} className="flex items-center gap-2 text-zinc-400 hover:text-red-400 transition">
              <LogOut size={20} />
              <span className="hidden sm:inline">Logout</span>
            </button>
          </>
        ) : (
          <Link href="/login" className="bg-white text-black px-5 py-2 rounded-full font-bold hover:scale-105 transition">
            Log In
          </Link>
        )}
      </div>
    </nav>
  );
}
