import { beforeEach, describe, expect, it } from "vitest";
import { usePlayerStore } from "@/store/usePlayerStore";
import type { Track } from "@/types/player";

const t = (id: string): Track => ({
  id,
  title: `Lagu ${id}`,
  artist: "Artis",
  coverUrl: "",
  audioUrl: `/${id}.mp3`,
  duration: 200,
});
const A = t("a"), B = t("b"), C = t("c");

const s = () => usePlayerStore.getState();

beforeEach(() => {
  usePlayerStore.setState(usePlayerStore.getInitialState(), true);
});

describe("usePlayerStore", () => {
  it("playTrack mengatur queue, index, dan mulai play", () => {
    s().playTrack(B, [A, B, C]);
    expect(s().currentTrack?.id).toBe("b");
    expect(s().currentIndex).toBe(1);
    expect(s().isPlaying).toBe(true);
  });

  it("playNext maju, dan berhenti di akhir queue saat repeat off", () => {
    s().playTrack(B, [A, B, C]);
    s().playNext();
    expect(s().currentTrack?.id).toBe("c");
    s().playNext();
    expect(s().currentTrack?.id).toBe("c");
    expect(s().isPlaying).toBe(false);
  });

  it("repeat all kembali ke lagu pertama di akhir queue", () => {
    s().playTrack(C, [A, B, C]);
    s().cycleRepeatMode(); // off -> all
    expect(s().repeatMode).toBe("all");
    s().playNext();
    expect(s().currentTrack?.id).toBe("a");
  });

  it("queue 1 lagu + repeat all meminta restart (seekRequestId naik, progress 0)", () => {
    s().playTrack(A, [A]);
    s().cycleRepeatMode(); // all
    s().setProgress(150);
    const before = s().seekRequestId;
    s().playNext();
    expect(s().currentTrack?.id).toBe("a");
    expect(s().seekRequestId).toBe(before + 1);
    expect(s().progress).toBe(0);
    expect(s().isPlaying).toBe(true);
  });

  it("playPrevious >3 detik me-restart lagu lewat seek", () => {
    s().playTrack(B, [A, B, C]);
    s().setProgress(10);
    const before = s().seekRequestId;
    s().playPrevious();
    expect(s().currentTrack?.id).toBe("b");
    expect(s().progress).toBe(0);
    expect(s().seekRequestId).toBe(before + 1);
  });

  it("playPrevious di index 0 (<=3 detik) me-restart lagu yang sama", () => {
    s().playTrack(A, [A, B, C]);
    s().setProgress(1);
    const before = s().seekRequestId;
    s().playPrevious();
    expect(s().currentTrack?.id).toBe("a");
    expect(s().seekRequestId).toBe(before + 1);
  });

  it("playPrevious <=3 detik pindah ke lagu sebelumnya", () => {
    s().playTrack(B, [A, B, C]);
    s().setProgress(2);
    s().playPrevious();
    expect(s().currentTrack?.id).toBe("a");
  });

  it("klik lagu yang sedang diputar me-restart-nya", () => {
    s().playTrack(A, [A, B]);
    s().setProgress(90);
    const before = s().seekRequestId;
    s().playTrack(A, [A, B]);
    expect(s().progress).toBe(0);
    expect(s().seekRequestId).toBe(before + 1);
  });

  it("pindah lagu menghapus playbackError", () => {
    s().playTrack(A, [A, B]);
    s().setPlaybackError("rusak");
    s().playNext();
    expect(s().playbackError).toBeNull();
  });

  it("shuffle menjaga lagu aktif di posisi pertama dan semua lagu tetap ada", () => {
    s().playTrack(B, [A, B, C]);
    s().toggleShuffle();
    expect(s().isShuffled).toBe(true);
    expect(s().shuffledQueue[0].id).toBe("b");
    expect(s().shuffledQueue.map((x) => x.id).sort()).toEqual(["a", "b", "c"]);
    expect(s().currentTrack?.id).toBe("b");
    s().toggleShuffle();
    expect(s().isShuffled).toBe(false);
    expect(s().queue[s().currentIndex].id).toBe("b");
  });

  it("repeat mode berputar off -> all -> one -> off", () => {
    const modes = [s().repeatMode];
    for (let i = 0; i < 3; i++) {
      s().cycleRepeatMode();
      modes.push(s().repeatMode);
    }
    expect(modes).toEqual(["off", "all", "one", "off"]);
  });

  it("volume dibatasi 0..1 dan mute/unmute mengembalikan volume", () => {
    s().setVolume(5);
    expect(s().volume).toBe(1);
    s().setVolume(-1);
    expect(s().volume).toBe(0);
    s().setVolume(0.6);
    s().toggleMute();
    expect(s().volume).toBe(0);
    s().toggleMute();
    expect(s().volume).toBeGreaterThan(0);
  });
});
