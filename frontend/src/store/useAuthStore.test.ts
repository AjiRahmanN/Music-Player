import { afterEach, describe, expect, it, vi } from "vitest";

function memoryStorage(seed: Record<string, string> = {}) {
  const data = { ...seed };
  return {
    data,
    getItem: (k: string) => data[k] ?? null,
    setItem: (k: string, v: string) => void (data[k] = v),
    removeItem: (k: string) => void delete data[k],
  };
}

async function load(seed?: Record<string, string>) {
  vi.resetModules();
  const storage = memoryStorage(seed);
  vi.stubGlobal("window", {});
  vi.stubGlobal("localStorage", storage);
  const { useAuthStore } = await import("@/store/useAuthStore");
  return { useAuthStore, storage };
}

const user = { id: "1", username: "aji", email: "aji@x.com", createdAt: "" };

afterEach(() => vi.unstubAllGlobals());

describe("useAuthStore", () => {
  it("access token TIDAK ditulis ke localStorage, hanya user", async () => {
    const { useAuthStore, storage } = await load();
    useAuthStore.getState().setAuth(user, "rahasia-token");
    const saved = JSON.parse(storage.data["auth-storage"]);
    expect(saved.state.user.username).toBe("aji");
    expect(JSON.stringify(saved)).not.toContain("rahasia-token");
  });

  it("accessToken lama dari storage versi sebelumnya diabaikan saat hydrate", async () => {
    const old = JSON.stringify({ state: { user, accessToken: "token-lama" }, version: 0 });
    const { useAuthStore } = await load({ "auth-storage": old });
    expect(useAuthStore.getState().accessToken).toBeNull();
    expect(useAuthStore.getState().user?.username).toBe("aji");
    expect(useAuthStore.getState().hasHydrated).toBe(true);
  });

  it("logout membersihkan sesi tanpa mereset hasHydrated", async () => {
    const { useAuthStore } = await load();
    useAuthStore.getState().setAuth(user, "t");
    useAuthStore.getState().logout();
    const s = useAuthStore.getState();
    expect(s.user).toBeNull();
    expect(s.accessToken).toBeNull();
    expect(s.isLoggedOut).toBe(true);
    expect(s.hasHydrated).toBe(true);
  });

  it("setAuth setelah logout menghapus penanda isLoggedOut", async () => {
    const { useAuthStore } = await load();
    useAuthStore.getState().logout();
    useAuthStore.getState().setAuth(user, "t");
    expect(useAuthStore.getState().isLoggedOut).toBe(false);
    expect(useAuthStore.getState().isAuthenticated()).toBe(true);
  });
});
