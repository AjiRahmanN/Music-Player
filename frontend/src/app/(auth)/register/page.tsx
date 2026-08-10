"use client";
import { useState } from "react";
import { authService } from "@/services/authService";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { AlertCircle, Mail, Lock, User, Loader2, Eye, EyeOff, Check, X, CheckCircle2 } from "lucide-react";

// ------- Aturan password, disamakan persis dengan backend (internal/utils/password.go) -------
const passwordRules = [
  { label: "Minimal 8 karakter", test: (p: string) => p.length >= 8 },
  { label: "Ada huruf besar (A-Z)", test: (p: string) => /[A-Z]/.test(p) },
  { label: "Ada huruf kecil (a-z)", test: (p: string) => /[a-z]/.test(p) },
  { label: "Ada angka (0-9)", test: (p: string) => /[0-9]/.test(p) },
];

const usernameValid = (u: string) => /^[a-zA-Z0-9_]{3,50}$/.test(u);
const emailValid = (e: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(e);

export default function RegisterPage() {
  const [formData, setFormData] = useState({ username: "", email: "", password: "" });
  const [showPassword, setShowPassword] = useState(false);
  const [touched, setTouched] = useState({ username: false, email: false, password: false });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const router = useRouter();

  const passwordChecks = passwordRules.map((rule) => ({
    ...rule,
    passed: rule.test(formData.password),
  }));
  const passwordOk = passwordChecks.every((c) => c.passed);
  const usernameOk = usernameValid(formData.username);
  const emailOk = emailValid(formData.email);
  const formValid = usernameOk && emailOk && passwordOk;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setTouched({ username: true, email: true, password: true });
    setError(null);

    if (!formValid) {
      setError("Periksa lagi data yang kamu masukkan.");
      return;
    }

    setLoading(true);
    try {
      await authService.register(formData);
      setSuccess(true);
      setTimeout(() => router.push("/login"), 1200);
    } catch (err: any) {
      setError(err.response?.data?.message || err.response?.data?.error || "Registrasi gagal. Silakan coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-br from-zinc-950 to-black p-4">
      <div className="w-full max-w-md">
        {/* Logo/Header — konsisten dengan halaman login */}
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-green-500 to-emerald-600 rounded-full mb-4">
            <svg className="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z" />
            </svg>
          </div>
          <h1 className="text-3xl font-bold text-white">Buat Akun</h1>
          <p className="text-zinc-400 mt-2">Gabung dan mulai dengarkan musik favoritmu</p>
        </div>

        <form onSubmit={handleSubmit} className="bg-zinc-900/50 backdrop-blur-sm rounded-2xl border border-zinc-800 p-8 shadow-2xl">
          <div className="space-y-5">
            {/* Success state — menggantikan alert() */}
            {success ? (
              <div className="flex flex-col items-center text-center py-6 gap-3">
                <CheckCircle2 className="text-green-500" size={48} />
                <p className="text-white font-semibold">Registrasi berhasil!</p>
                <p className="text-zinc-400 text-sm">Mengarahkan ke halaman login...</p>
              </div>
            ) : (
              <>
                {/* Username */}
                <div>
                  <label className="flex items-center gap-2 text-sm font-medium text-zinc-300 mb-2">
                    <User size={16} />
                    Username
                  </label>
                  <input
                    type="text"
                    placeholder="username_kamu"
                    value={formData.username}
                    onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                    onBlur={() => setTouched((t) => ({ ...t, username: true }))}
                    disabled={loading}
                    autoComplete="username"
                    className={`w-full p-3 rounded-lg bg-zinc-800/50 text-white border focus:outline-none focus:ring-2 transition ${
                      touched.username && !usernameOk ? "border-red-700 focus:ring-red-500/30 focus:border-red-500" : "border-zinc-700 focus:ring-green-500/30 focus:border-green-500"
                    }`}
                  />
                  {touched.username && !usernameOk && <p className="text-red-400 text-xs mt-1.5">3-50 karakter, hanya huruf, angka, dan underscore.</p>}
                </div>

                {/* Email */}
                <div>
                  <label className="flex items-center gap-2 text-sm font-medium text-zinc-300 mb-2">
                    <Mail size={16} />
                    Email
                  </label>
                  <input
                    type="email"
                    placeholder="nama@email.com"
                    value={formData.email}
                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                    onBlur={() => setTouched((t) => ({ ...t, email: true }))}
                    disabled={loading}
                    autoComplete="email"
                    className={`w-full p-3 rounded-lg bg-zinc-800/50 text-white border focus:outline-none focus:ring-2 transition ${
                      touched.email && !emailOk ? "border-red-700 focus:ring-red-500/30 focus:border-red-500" : "border-zinc-700 focus:ring-green-500/30 focus:border-green-500"
                    }`}
                  />
                  {touched.email && !emailOk && <p className="text-red-400 text-xs mt-1.5">Format email tidak valid.</p>}
                </div>

                {/* Password */}
                <div>
                  <label className="flex items-center gap-2 text-sm font-medium text-zinc-300 mb-2">
                    <Lock size={16} />
                    Password
                  </label>
                  <div className="relative">
                    <input
                      type={showPassword ? "text" : "password"}
                      placeholder="••••••••"
                      value={formData.password}
                      onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                      onFocus={() => setTouched((t) => ({ ...t, password: true }))}
                      disabled={loading}
                      autoComplete="new-password"
                      className={`w-full p-3 pr-11 rounded-lg bg-zinc-800/50 text-white border focus:outline-none focus:ring-2 transition ${
                        touched.password && !passwordOk && formData.password.length > 0 ? "border-red-700 focus:ring-red-500/30 focus:border-red-500" : "border-zinc-700 focus:ring-green-500/30 focus:border-green-500"
                      }`}
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword((v) => !v)}
                      tabIndex={-1}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-300 transition"
                      aria-label={showPassword ? "Sembunyikan password" : "Tampilkan password"}
                    >
                      {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                    </button>
                  </div>

                  {/* Checklist syarat password — live feedback, ini yang tadinya hilang */}
                  {touched.password && (
                    <ul className="mt-2.5 space-y-1">
                      {passwordChecks.map((c) => (
                        <li key={c.label} className={`flex items-center gap-1.5 text-xs transition ${c.passed ? "text-green-400" : "text-zinc-500"}`}>
                          {c.passed ? <Check size={13} /> : <X size={13} />}
                          {c.label}
                        </li>
                      ))}
                    </ul>
                  )}
                </div>

                {/* Error inline — menggantikan alert() */}
                {error && (
                  <div className="bg-red-900/30 border border-red-800 rounded-lg p-3">
                    <div className="flex items-center gap-2 text-red-300">
                      <AlertCircle size={16} className="flex-shrink-0" />
                      <span className="text-sm">{error}</span>
                    </div>
                  </div>
                )}

                {/* Submit */}
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
                    <span>Daftar Sekarang</span>
                  )}
                </button>

                <div className="text-center">
                  <p className="text-zinc-500 text-sm">
                    Sudah punya akun?{" "}
                    <Link href="/login" className="text-green-500 hover:text-green-400 font-medium transition hover:underline">
                      Masuk
                    </Link>
                  </p>
                </div>
              </>
            )}
          </div>
        </form>

        <div className="mt-8 text-center">
          <p className="text-xs text-zinc-600">© 2024 Melodify. All rights reserved.</p>
        </div>
      </div>
    </div>
  );
}
