"use client";
import { useState } from "react";
import { authService } from "@/services/authService";
import { useAuthStore } from "@/store/useAuthStore";
import { useRouter } from "next/navigation";
import { AlertCircle, Mail, Lock, Loader2 } from "lucide-react";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const setAuth = useAuthStore((state) => state.setAuth);
  const router = useRouter();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await authService.login({ email, password });
      if (!response.data?.accessToken || !response.data.user) {
        throw new Error("Respons login tidak lengkap");
      }
      setAuth(response.data.user, response.data.accessToken);
      router.push("/");
      router.refresh();
    } catch (err: unknown) {
      const response = (err as { response?: { data?: { message?: string }; status?: number } }).response;
      setError(
        response?.status === 401
          ? "Email atau password salah!"
          : response?.data?.message || "Login gagal. Silakan coba lagi."
      );
    } finally {
      setLoading(false);
    }
  };

  // Function untuk menampilkan error popup
  const renderErrorPopup = () => {
    if (!error) return null;

    return (
      <div className="fixed top-4 right-4 z-50 animate-slideIn">
        <div className="bg-red-900/90 backdrop-blur-sm border border-red-700 rounded-lg shadow-xl p-4 max-w-sm">
          <div className="flex items-start gap-3">
            <AlertCircle className="text-red-400 mt-0.5 flex-shrink-0" />
            <div className="flex-1">
              <h3 className="font-semibold text-white">Login Gagal</h3>
              <p className="text-red-200 text-sm mt-1">{error}</p>
              <button onClick={() => setError(null)} className="mt-2 text-xs bg-red-800 hover:bg-red-700 px-3 py-1 rounded transition">
                Tutup
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  };

  return (
    <>
      {renderErrorPopup()}

      <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-zinc-950 to-black p-4">
        <div className="w-full max-w-md">
          {/* Logo/Header */}
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-green-500 to-emerald-600 rounded-full mb-4">
              <svg className="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z" />
              </svg>
            </div>
            <h1 className="text-3xl font-bold text-white">Selamat Datang</h1>
            <p className="text-zinc-400 mt-2">Masuk ke akun Melodify Anda</p>
          </div>

          {/* Form */}
          <form onSubmit={handleLogin} className="bg-zinc-900/50 backdrop-blur-sm rounded-2xl border border-zinc-800 p-8 shadow-2xl">
            <div className="space-y-6">
              {/* Email Input */}
              <div>
                <label className="flex items-center gap-2 text-sm font-medium text-zinc-300 mb-2">
                  <Mail size={16} />
                  Email
                </label>
                <input
                  type="email"
                  placeholder="nama@email.com"
                  className="w-full p-3 rounded-lg bg-zinc-800/50 text-white border border-zinc-700 focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-500/30 transition"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  disabled={loading}
                  autoComplete="email"
                />
              </div>

              {/* Password Input */}
              <div>
                <label className="flex items-center gap-2 text-sm font-medium text-zinc-300 mb-2">
                  <Lock size={16} />
                  Password
                </label>
                <input
                  type="password"
                  placeholder="••••••••"
                  className="w-full p-3 rounded-lg bg-zinc-800/50 text-white border border-zinc-700 focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-500/30 transition"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  disabled={loading}
                  autoComplete="current-password"
                />
              </div>

              {/* Error Message (inline) */}
              {error && (
                <div className="bg-red-900/30 border border-red-800 rounded-lg p-3">
                  <div className="flex items-center gap-2 text-red-300">
                    <AlertCircle size={16} />
                    <span className="text-sm">{error}</span>
                  </div>
                </div>
              )}

              {/* Submit Button */}
              <button
                type="submit"
                disabled={loading}
                className="w-full bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 text-white font-semibold py-3 rounded-lg transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                {loading ? (
                  <>
                    <Loader2 className="animate-spin" size={20} />
                    <span>Memproses...</span>
                  </>
                ) : (
                  <span>Masuk</span>
                )}
              </button>

              {/* Divider */}
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-zinc-800"></div>
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-zinc-900/50 px-2 text-zinc-500">atau</span>
                </div>
              </div>

              {/* Register Link */}
              <div className="text-center">
                <p className="text-zinc-500 text-sm">
                  Belum punya akun?{" "}
                  <a href="/register" className="text-green-500 hover:text-green-400 font-medium transition hover:underline">
                    Daftar sekarang
                  </a>
                </p>
              </div>
            </div>
          </form>

          {/* Footer */}
          <div className="mt-8 text-center">
            <p className="text-xs text-zinc-600">© 2024 Melodify. All rights reserved.</p>
          </div>
        </div>
      </div>
    </>
  );
}
