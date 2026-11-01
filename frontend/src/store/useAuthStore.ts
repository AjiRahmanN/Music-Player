import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import type { User } from "@/types/auth";

interface AuthState {
  user: User | null;
  // Access token HANYA disimpan di memori (tidak di-persist) supaya tidak bisa
  // dicuri lewat XSS dari localStorage. Setelah reload halaman, token dipulihkan
  // lewat refresh token (cookie httpOnly) -- lihat AuthGuard.
  accessToken: string | null;
  // True setelah state persisted (user) selesai dimuat dari storage.
  hasHydrated: boolean;
  // True setelah user logout di sesi ini, supaya AuthGuard tidak mencoba
  // memulihkan sesi lagi lewat refresh token.
  isLoggedOut: boolean;

  setAuth: (user: User, token: string) => void;
  setToken: (token: string) => void;
  setUser: (user: User) => void;
  logout: () => void;
  setHasHydrated: (state: boolean) => void;

  // Computed property
  isAuthenticated: () => boolean;
}

type PersistedAuth = Pick<AuthState, "user">;

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      hasHydrated: false,
      isLoggedOut: false,

      setAuth: (user, token) => {
        set({ user, accessToken: token, isLoggedOut: false });
      },

      setToken: (token) => {
        set({ accessToken: token });
      },

      setUser: (user) => set({ user }),

      // hasHydrated sengaja TIDAK di-reset: flag itu hanya soal storage sudah
      // dibaca atau belum, bukan soal status login.
      logout: () => {
        set({ user: null, accessToken: null, isLoggedOut: true });
      },

      setHasHydrated: (state) => set({ hasHydrated: state }),

      // Getter untuk isAuthenticated
      isAuthenticated: () => {
        return !!get().accessToken;
      },
    }),
    {
      name: "auth-storage",
      storage: createJSONStorage(() => {
        if (typeof window === "undefined") {
          return {
            getItem: () => null,
            setItem: () => {},
            removeItem: () => {},
          };
        }
        return localStorage;
      }),
      // Hanya data profil non-sensitif yang disimpan.
      partialize: (state): PersistedAuth => ({ user: state.user }),
      // Abaikan field lain dari storage -- termasuk accessToken yang mungkin
      // masih tersimpan oleh versi aplikasi sebelumnya.
      merge: (persisted, current) => ({
        ...current,
        user: (persisted as Partial<PersistedAuth> | undefined)?.user ?? null,
      }),
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.setHasHydrated(true);
        }
      },
      // Skip hydration jika error
      skipHydration: false,
    }
  )
);
