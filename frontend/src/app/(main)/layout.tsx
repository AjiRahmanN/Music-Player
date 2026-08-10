import AuthGuard from "@/components/auth/AuthGuard";
import Navbar from "@/components/layout/Navbar";
import PlayerAudioEngine from "@/components/player/PlayerAudioEngine";
import PlayerBar from "@/components/player/PlayerBar";

export default function MainLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthGuard>
      <div className="flex flex-col min-h-screen bg-zinc-950">
        <Navbar />
        <PlayerAudioEngine />
        <main className="flex-1 p-6 pb-24">{children}</main>
        <PlayerBar />
      </div>
    </AuthGuard>
  );
}
