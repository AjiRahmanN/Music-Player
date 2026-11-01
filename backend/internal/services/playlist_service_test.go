package services

import (
	"context"
	stdErrors "errors"
	"testing"

	apperrors "music-player-backend/internal/errors"
	"music-player-backend/internal/models"
)

func isNotFound(err error) bool {
	var app *apperrors.AppError
	return stdErrors.As(err, &app) && app.Code == "NOT_FOUND"
}

func newPlaylistFixture() (*PlaylistService, *fakePlaylistRepo) {
	repo := newFakePlaylistRepo(&models.Playlist{ID: "p1", UserID: "owner", Name: "Punya Owner"})
	songs := newFakeSongRepo(&models.Song{ID: "s1", Title: "Lagu"})
	return NewPlaylistService(repo, songs), repo
}

func TestGetOwnedPlaylistHidesOthersPlaylists(t *testing.T) {
	svc, _ := newPlaylistFixture()
	ctx := context.Background()

	if _, err := svc.GetDetail(ctx, "owner", "p1"); err != nil {
		t.Fatalf("pemilik harus bisa melihat: %v", err)
	}
	// Milik orang lain dan tidak ada harus menghasilkan error yang SAMA (404),
	// supaya keberadaan playlist tidak bocor.
	_, errOther := svc.GetDetail(ctx, "intruder", "p1")
	_, errMissing := svc.GetDetail(ctx, "intruder", "tidak-ada")
	if !isNotFound(errOther) || !isNotFound(errMissing) {
		t.Fatalf("errOther=%v errMissing=%v, keduanya harus NOT_FOUND", errOther, errMissing)
	}
}

func TestNonOwnerCannotMutatePlaylist(t *testing.T) {
	svc, repo := newPlaylistFixture()
	ctx := context.Background()

	if err := svc.Delete(ctx, "intruder", "p1"); !isNotFound(err) {
		t.Fatalf("Delete oleh non-owner: err = %v", err)
	}
	if _, ok := repo.playlists["p1"]; !ok {
		t.Fatal("playlist tidak boleh terhapus oleh non-owner")
	}

	req := &models.AddSongToPlaylistRequest{SongID: "s1"}
	if err := svc.AddSong(ctx, "intruder", "p1", req); !isNotFound(err) {
		t.Fatalf("AddSong oleh non-owner: err = %v", err)
	}
	if len(repo.songs["p1"]) != 0 {
		t.Fatal("lagu tidak boleh ditambahkan oleh non-owner")
	}
	if err := svc.RemoveSong(ctx, "intruder", "p1", "s1"); !isNotFound(err) {
		t.Fatalf("RemoveSong oleh non-owner: err = %v", err)
	}
}

func TestOwnerCanAddRemoveAndDelete(t *testing.T) {
	svc, repo := newPlaylistFixture()
	ctx := context.Background()

	if err := svc.AddSong(ctx, "owner", "p1", &models.AddSongToPlaylistRequest{SongID: "s1"}); err != nil {
		t.Fatal(err)
	}
	if len(repo.songs["p1"]) != 1 {
		t.Fatal("lagu seharusnya tertambah")
	}
	if err := svc.AddSong(ctx, "owner", "p1", &models.AddSongToPlaylistRequest{SongID: "tidak-ada"}); !isNotFound(err) {
		t.Fatalf("lagu tak ada: err = %v, want NOT_FOUND", err)
	}
	if err := svc.RemoveSong(ctx, "owner", "p1", "s1"); err != nil || len(repo.songs["p1"]) != 0 {
		t.Fatalf("RemoveSong gagal: %v", err)
	}
	if err := svc.Delete(ctx, "owner", "p1"); err != nil {
		t.Fatal(err)
	}
	if _, ok := repo.playlists["p1"]; ok {
		t.Fatal("playlist seharusnya terhapus")
	}
}

func TestAdminDeleteIgnoresOwnership(t *testing.T) {
	svc, repo := newPlaylistFixture()
	ctx := context.Background()
	if err := svc.AdminDelete(ctx, "p1"); err != nil {
		t.Fatal(err)
	}
	if _, ok := repo.playlists["p1"]; ok {
		t.Fatal("AdminDelete seharusnya menghapus playlist siapa pun")
	}
	if err := svc.AdminDelete(ctx, "p1"); !isNotFound(err) {
		t.Fatalf("AdminDelete playlist tak ada: err = %v", err)
	}
}
