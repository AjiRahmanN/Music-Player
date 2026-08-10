import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import type { User } from "@/types/auth";

interface AuthState {
  user: User | null;
  accessToken: string | null;
  hasHydrated: boolean;

  setAuth: (user: User, token: string) => void;
  setToken: (token: string) => void;
  setUser: (user: User) => void;
  logout: () => void;
  setHasHydrated: (state: boolean) => void;

  // Computed property
  isAuthenticated: () => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      hasHydrated: false,

      setAuth: (user, token) => {
        set({ user, accessToken: token });
      },

      setToken: (token) => {
        set({ accessToken: token });
      },

      setUser: (user) => set({ user }),

      logout: () => {
        set({ user: null, accessToken: null, hasHydrated: false });
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
