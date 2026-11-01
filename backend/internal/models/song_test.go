package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAssetURLResolveAndNormalize(t *testing.T) {
	SetAssetBaseURL("https://api.contoh.com/")
	t.Cleanup(func() { SetAssetBaseURL("") })

	if got := ResolveAssetURL("/uploads/audios/a.mp3"); got != "https://api.contoh.com/uploads/audios/a.mp3" {
		t.Fatalf("ResolveAssetURL = %q", got)
	}
	// URL absolut/eksternal tidak diubah
	for _, u := range []string{"https://cdn.lain.com/x.mp3", "http://localhost:8080/uploads/x.mp3", ""} {
		if got := ResolveAssetURL(u); got != u {
			t.Fatalf("ResolveAssetURL(%q) = %q, want tidak berubah", u, got)
		}
	}

	if got := NormalizeAssetURL("https://api.contoh.com/uploads/covers/c.png"); got != "/uploads/covers/c.png" {
		t.Fatalf("NormalizeAssetURL = %q", got)
	}
	if got := NormalizeAssetURL("https://cdn.lain.com/c.png"); got != "https://cdn.lain.com/c.png" {
		t.Fatalf("URL eksternal tidak boleh diubah: %q", got)
	}
	// round trip
	rel := "/uploads/audios/z.mp3"
	if NormalizeAssetURL(ResolveAssetURL(rel)) != rel {
		t.Fatal("round trip Resolve->Normalize harus mengembalikan path relatif")
	}
}

func TestSongMarshalJSONResolvesUploads(t *testing.T) {
	SetAssetBaseURL("https://api.contoh.com")
	t.Cleanup(func() { SetAssetBaseURL("") })

	b, err := json.Marshal(Song{ID: "1", Title: "T", Artist: "A", CoverURL: "/uploads/covers/c.png", AudioURL: "https://x.com/a.mp3", Duration: 10})
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `"coverUrl":"https://api.contoh.com/uploads/covers/c.png"`) {
		t.Fatalf("coverUrl tidak di-resolve: %s", s)
	}
	if !strings.Contains(s, `"audioUrl":"https://x.com/a.mp3"`) {
		t.Fatalf("audioUrl eksternal berubah: %s", s)
	}
}

func TestCreateSongRequestValidate(t *testing.T) {
	ok := CreateSongRequest{Title: "t", Artist: "a", AudioURL: "u", Duration: 1}
	if err := ok.Validate(); err != nil {
		t.Fatal(err)
	}
	for name, mut := range map[string]func(*CreateSongRequest){
		"title":    func(r *CreateSongRequest) { r.Title = "" },
		"artist":   func(r *CreateSongRequest) { r.Artist = "" },
		"audio":    func(r *CreateSongRequest) { r.AudioURL = "" },
		"duration": func(r *CreateSongRequest) { r.Duration = 0 },
	} {
		r := ok
		mut(&r)
		if err := r.Validate(); err == nil {
			t.Fatalf("%s kosong seharusnya error", name)
		}
	}
}
