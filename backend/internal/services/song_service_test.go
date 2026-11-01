package services

import (
	"context"
	"testing"

	"music-player-backend/internal/models"
)

func TestNormalizePage(t *testing.T) {
	cases := []struct{ inL, inO, wantL, wantO int }{
		{0, 0, DefaultSongLimit, 0},
		{-5, -5, DefaultSongLimit, 0},
		{10, 20, 10, 20},
		{99999, 3, MaxSongLimit, 3},
	}
	for _, c := range cases {
		l, o := NormalizePage(c.inL, c.inO)
		if l != c.wantL || o != c.wantO {
			t.Errorf("NormalizePage(%d,%d) = (%d,%d), want (%d,%d)", c.inL, c.inO, l, o, c.wantL, c.wantO)
		}
	}
}

func TestSongListAppliesPaginationBounds(t *testing.T) {
	repo := newFakeSongRepo()
	svc := NewSongService(repo)
	_, _ = svc.List(context.Background(), 0, -1)
	if repo.lastLimit != DefaultSongLimit || repo.lastOffset != 0 {
		t.Fatalf("List limit/offset = %d/%d", repo.lastLimit, repo.lastOffset)
	}
	_, _ = svc.Search(context.Background(), "x", 1000, 5)
	if repo.lastLimit != MaxSongLimit || repo.lastOffset != 5 {
		t.Fatalf("Search limit/offset = %d/%d", repo.lastLimit, repo.lastOffset)
	}
}

func TestSongCreateStoresRelativeUploadPaths(t *testing.T) {
	models.SetAssetBaseURL("https://api.contoh.com")
	t.Cleanup(func() { models.SetAssetBaseURL("") })

	repo := newFakeSongRepo()
	svc := NewSongService(repo)
	_, err := svc.Create(context.Background(), &models.CreateSongRequest{
		Title: "T", Artist: "A", Duration: 10,
		AudioURL: "https://api.contoh.com/uploads/audios/a.mp3",
		CoverURL: "https://cdn.lain.com/cover.png",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repo.created.AudioURL != "/uploads/audios/a.mp3" {
		t.Errorf("AudioURL tersimpan = %q, want path relatif", repo.created.AudioURL)
	}
	if repo.created.CoverURL != "https://cdn.lain.com/cover.png" {
		t.Errorf("CoverURL eksternal tidak boleh diubah: %q", repo.created.CoverURL)
	}
}

func TestSongCreateValidation(t *testing.T) {
	svc := NewSongService(newFakeSongRepo())
	if _, err := svc.Create(context.Background(), &models.CreateSongRequest{}); appCode(err) != "VALIDATION_ERROR" {
		t.Fatalf("err = %v, want VALIDATION_ERROR", err)
	}
}
