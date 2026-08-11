"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Home, ListMusic, Shield } from "lucide-react";

// Sisipkan komponen ini ke dalam Navbar.tsx kamu yang sudah ada,
// biasanya di antara logo (kiri) dan avatar/profil (kanan).
export default function NavLinks() {
  const pathname = usePathname();

  const links = [
    { href: "/", label: "Beranda", icon: Home },
    { href: "/playlists", label: "Playlist", icon: ListMusic },
    // { href: "/admin", label: "Admin", icon: Shield },
  ];

  return (
    <nav className="flex items-center gap-1">
      {links.map(({ href, label, icon: Icon }) => {
        // Exact match untuk "/", startsWith untuk sisanya (biar /playlists/xxx tetap dianggap aktif)
        const isActive = href === "/" ? pathname === "/" : pathname.startsWith(href);

        return (
          <Link key={href} href={href} className={`flex items-center gap-2 px-3 py-2 rounded-full text-sm font-medium transition ${isActive ? "bg-zinc-800 text-white" : "text-zinc-400 hover:text-white hover:bg-zinc-800/50"}`}>
            <Icon size={16} />
            {label}
          </Link>
        );
      })}
    </nav>
  );
}
