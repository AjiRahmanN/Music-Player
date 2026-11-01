"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { useRouter, usePathname } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { authService } from "@/services/authService";
import { refreshAccessToken } from "@/lib/axios";

export default function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, accessToken, hasHydrated, isLoggedOut, setAuth, logout } = useAuthStore();
  const router = useRouter();
  const pathname = usePathname();
  const [isChecking, setIsChecking] = useState(true);
  const restoringRef = useRef(false);
  const isPublicRoute = pathname === "/login" || pathname === "/register";

  useEffect(() => {
    if (isPublicRoute) {
      setIsChecking(false);
      return;
    }
    if (!hasHydrated) return;

    // Sudah login penuh di sesi ini.
    if (accessToken && user) {
      setIsChecking(false);
      return;
    }

    // User baru saja logout: jangan coba pulihkan sesi.
    if (isLoggedOut) {
      router.replace("/login");
      return;
    }

    // Belum ada access token (mis. halaman baru di-reload; token hanya di
    // memori) atau user belum termuat: pulihkan sesi lewat refresh token
    // (cookie httpOnly), lalu ambil profil.
    if (restoringRef.current) return;
    restoringRef.current = true;

    const restoreSession = async () => {
      try {
        const token = accessToken ?? (await refreshAccessToken());
        const me = await authService.getMe(); // GET /me mengembalikan objek flat
        setAuth({ id: me.id, username: me.username, email: me.email, createdAt: me.created }, token);
        setIsChecking(false);
      } catch {
        logout();
        router.replace("/login");
      } finally {
        restoringRef.current = false;
      }
    };
    void restoreSession();
  }, [accessToken, hasHydrated, isLoggedOut, isPublicRoute, logout, router, setAuth, user]);

  if (isChecking) {
    return (
      <div className="h-screen flex flex-col items-center justify-center bg-black text-white">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-green-500 mb-4" />
        <p className="text-lg">Checking authentication...</p>
      </div>
    );
  }

  if (isPublicRoute) return <>{children}</>;
  return accessToken ? <>{children}</> : null;
}
