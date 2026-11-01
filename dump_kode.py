#!/usr/bin/env python3
"""
Menulis ulang semua kode sumber dalam sebuah folder ke satu file .txt,
kecuali folder library/dependensi, file .env, file lock, dan file biner.

Pemakaian:
    python3 dump_kode.py                      # folder saat ini -> semua_kode.txt
    python3 dump_kode.py /path/proyek         # folder tertentu
    python3 dump_kode.py . -o hasil.txt       # nama output sendiri
    python3 dump_kode.py . --tambah-skip-dir uploads data
"""

import argparse
import fnmatch
import os
import sys

# Folder yang dilewati seluruhnya (library, build, cache, dll)
SKIP_DIRS = {
    "node_modules", "vendor", "venv", ".venv", "env", "__pycache__",
    ".git", ".idea", ".vscode", ".next", ".nuxt", "dist", "build", "out",
    "target", "bin", "obj", ".cache", ".pytest_cache", ".mypy_cache",
    "coverage", ".gradle", ".dart_tool", "site-packages",
}

# Pola nama file yang dilewati (env, lock, rahasia, biner umum)
SKIP_FILE_PATTERNS = [
    ".env", ".env.*", "*.env",
    "package-lock.json", "yarn.lock", "pnpm-lock.yaml", "go.sum",
    "composer.lock", "Cargo.lock", "poetry.lock", "Pipfile.lock",
    "*.pem", "*.key", "*.crt", "*.p12",
    "*.png", "*.jpg", "*.jpeg", "*.gif", "*.webp", "*.ico", "*.svg",
    "*.mp3", "*.mp4", "*.wav", "*.flac", "*.ogg", "*.webm",
    "*.pdf", "*.zip", "*.tar", "*.gz", "*.rar", "*.7z",
    "*.exe", "*.dll", "*.so", "*.dylib", "*.o", "*.a", "*.class", "*.jar",
    "*.woff", "*.woff2", "*.ttf", "*.eot", "*.otf",
    "*.sqlite", "*.db", "*.log", "*.min.js", "*.min.css", "*.map",
    ".DS_Store",
]

MAX_UKURAN_BYTE = 1_000_000  # lewati file teks lebih dari ~1 MB


def harus_skip_file(nama):
    return any(fnmatch.fnmatch(nama, pola) for pola in SKIP_FILE_PATTERNS)


def file_biner(path):
    try:
        with open(path, "rb") as f:
            return b"\0" in f.read(2048)
    except OSError:
        return True


def kumpulkan_file(root, skip_dirs, file_output):
    hasil = []
    for dirpath, dirnames, filenames in os.walk(root):
        # edit di tempat agar os.walk tidak masuk ke folder yang di-skip
        dirnames[:] = sorted(d for d in dirnames if d not in skip_dirs)
        for nama in sorted(filenames):
            path = os.path.join(dirpath, nama)
            if os.path.abspath(path) == file_output:
                continue
            if harus_skip_file(nama):
                continue
            if os.path.getsize(path) > MAX_UKURAN_BYTE:
                print(f"[lewati - terlalu besar] {path}")
                continue
            if file_biner(path):
                continue
            hasil.append(path)
    return hasil


def main():
    parser = argparse.ArgumentParser(
        description="Gabungkan semua kode sumber folder ke satu file txt."
    )
    parser.add_argument("folder", nargs="?", default=".", help="folder proyek (default: .)")
    parser.add_argument("-o", "--output", default="semua_kode.txt", help="file output")
    parser.add_argument(
        "--tambah-skip-dir", nargs="*", default=[],
        help="nama folder tambahan yang dilewati",
    )
    args = parser.parse_args()

    root = os.path.abspath(args.folder)
    if not os.path.isdir(root):
        sys.exit(f"Folder tidak ditemukan: {root}")

    file_output = os.path.abspath(args.output)
    skip_dirs = SKIP_DIRS | set(args.tambah_skip_dir)

    daftar = kumpulkan_file(root, skip_dirs, file_output)

    with open(file_output, "w", encoding="utf-8") as out:
        out.write(f"# Proyek: {os.path.basename(root)}\n")
        out.write(f"# Jumlah file: {len(daftar)}\n\n")

        out.write("# DAFTAR FILE\n")
        for path in daftar:
            out.write(f"- {os.path.relpath(path, root)}\n")
        out.write("\n")

        for path in daftar:
            rel = os.path.relpath(path, root)
            out.write("=" * 80 + "\n")
            out.write(f"FILE: {rel}\n")
            out.write("=" * 80 + "\n")
            try:
                with open(path, "r", encoding="utf-8", errors="replace") as f:
                    isi = f.read()
            except OSError as e:
                isi = f"[gagal dibaca: {e}]"
            out.write(isi)
            if not isi.endswith("\n"):
                out.write("\n")
            out.write("\n")

    print(f"Selesai: {len(daftar)} file ditulis ke {file_output}")


if __name__ == "__main__":
    main()
