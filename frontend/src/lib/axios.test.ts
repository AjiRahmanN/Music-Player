import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import type { AxiosAdapter, InternalAxiosRequestConfig } from "axios";

// Adapter palsu: menggantikan jaringan. Setiap panggilan dicatat.
type Handler = (config: InternalAxiosRequestConfig) => { status: number; data?: unknown };

async function setup(handler: Handler) {
  vi.resetModules();
  const assign = vi.fn();
  vi.stubGlobal("window", { location: { pathname: "/", assign } });
  const mem: Record<string, string> = {};
  vi.stubGlobal("localStorage", {
    getItem: (k: string) => mem[k] ?? null,
    setItem: (k: string, v: string) => void (mem[k] = v),
    removeItem: (k: string) => void delete mem[k],
  });
  const calls: string[] = [];
  const { default: api } = await import("@/lib/axios");
  const { useAuthStore } = await import("@/store/useAuthStore");
  const adapter: AxiosAdapter = async (config) => {
    calls.push(`${config.method} ${config.url}`);
    await new Promise((r) => setTimeout(r, 5)); // beri kesempatan request lain berjalan paralel
    const { status, data } = handler(config);
    const response = { data, status, statusText: "", headers: {}, config, request: {} };
    if (status >= 200 && status < 300) return response;
    const err = Object.assign(new Error(`status ${status}`), { isAxiosError: true, config, response });
    throw err;
  };
  api.defaults.adapter = adapter;
  useAuthStore.getState().setAuth({ id: "1", username: "u", email: "e@x.com", createdAt: "" }, "old");
  return { api, useAuthStore, calls, assign };
}

afterEach(() => vi.unstubAllGlobals());

describe("axios refresh interceptor", () => {
  it("request paralel yang kena 401 hanya memicu SATU refresh", async () => {
    let refreshCount = 0;
    const { api, calls, useAuthStore } = await setup((config) => {
      if (config.url === "/auth/refresh") {
        refreshCount++;
        // Meniru backend: rotasi token, panggilan refresh kedua akan gagal.
        return refreshCount === 1
          ? { status: 200, data: { data: { accessToken: "new" } } }
          : { status: 401 };
      }
      return config.headers.Authorization === "Bearer new" ? { status: 200, data: { ok: true } } : { status: 401 };
    });

    const results = await Promise.all(["/songs", "/playlists", "/me", "/songs/1", "/profile"].map((u) => api.get(u)));

    expect(results.every((r) => r.status === 200)).toBe(true);
    expect(refreshCount).toBe(1);
    expect(calls.filter((c) => c.includes("/auth/refresh"))).toHaveLength(1);
    expect(useAuthStore.getState().accessToken).toBe("new");
  });

  it("request yang 401 setelah refresh selesai cukup diulang dengan token terbaru (tanpa refresh lagi)", async () => {
    let refreshCount = 0;
    const { api, useAuthStore } = await setup((config) => {
      if (config.url === "/auth/refresh") {
        refreshCount++;
        return { status: 200, data: { data: { accessToken: "new" } } };
      }
      return config.headers.Authorization === "Bearer new" ? { status: 200 } : { status: 401 };
    });
    useAuthStore.getState().setToken("new"); // store sudah lebih baru daripada token yang dikirim
    const req = api.get("/songs", { headers: { Authorization: "Bearer old" } });
    // request-interceptor menimpa header dengan token store, jadi langsung sukses
    expect((await req).status).toBe(200);
    expect(refreshCount).toBe(0);
  });

  it("kegagalan refresh -> logout dan redirect ke /login", async () => {
    const { api, useAuthStore, assign } = await setup((config) =>
      config.url === "/auth/refresh" ? { status: 401 } : { status: 401 }
    );
    await expect(api.get("/songs")).rejects.toBeTruthy();
    expect(useAuthStore.getState().accessToken).toBeNull();
    expect(useAuthStore.getState().isLoggedOut).toBe(true);
    expect(assign).toHaveBeenCalledWith("/login");
  });

  it("error non-401 pada retry TIDAK menyebabkan logout", async () => {
    const { api, useAuthStore, assign } = await setup((config) => {
      if (config.url === "/auth/refresh") return { status: 200, data: { data: { accessToken: "new" } } };
      return config.headers.Authorization === "Bearer new" ? { status: 500 } : { status: 401 };
    });
    await expect(api.get("/songs")).rejects.toMatchObject({ response: { status: 500 } });
    expect(useAuthStore.getState().accessToken).toBe("new");
    expect(assign).not.toHaveBeenCalled();
  });

  it("401 dari endpoint /auth/* (mis. login salah) tidak memicu refresh", async () => {
    let refreshCount = 0;
    const { api } = await setup((config) => {
      if (config.url === "/auth/refresh") refreshCount++;
      return { status: 401 };
    });
    await expect(api.post("/auth/login", {})).rejects.toBeTruthy();
    expect(refreshCount).toBe(0);
  });
});

beforeEach(() => vi.useRealTimers());
