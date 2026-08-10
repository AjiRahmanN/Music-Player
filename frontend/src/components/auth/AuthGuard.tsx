"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { useRouter, usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import { authService } from "@/services/authService";

export default function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, accessToken, hasHydrated, setAuth, logout } = useAuthStore();
  const router = useRouter();
  const pathname = usePathname();
  const [isChecking, setIsChecking] = useState(true);
  const isPublicRoute = pathname === "/login" || pathname === "/register";

  useEffect(() => {
    if (isPublicRoute) {
      setIsChecking(false);
      return;
    }
    if (!hasHydrated) return;
    if (!accessToken) {
      router.replace("/login");
      return;
    }
    if (user) {
      setIsChecking(false);
      return;
    }

    const restoreSession = async () => {
      try {
        const response = await authService.getMe();
        const restoredUser = response.data?.user;
        if (!restoredUser) throw new Error("Missing user data");
        setAuth(restoredUser, accessToken);
        setIsChecking(false);
      } catch {
        logout();
        router.replace("/login");
      }
    };
    void restoreSession();
  }, [accessToken, hasHydrated, isPublicRoute, logout, router, setAuth, user]);

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
