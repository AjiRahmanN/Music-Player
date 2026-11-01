package handlers

import (
	"bytes"
	"testing"
)

func pad(b []byte) []byte { return append(b, make([]byte, 64)...) }

func TestSniffAllowed(t *testing.T) {
	png := pad([]byte("\x89PNG\r\n\x1a\n"))
	jpeg := pad([]byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F'})
	mp3 := pad([]byte("ID3\x03\x00\x00\x00\x00\x00\x00"))
	ogg := pad([]byte("OggS\x00\x02"))
	wav := pad(append(append([]byte("RIFF\x24\x00\x00\x00WAVEfmt "), make([]byte, 4)...), 0))
	m4a := pad([]byte("\x00\x00\x00\x20ftypM4A \x00\x00\x00\x00M4A mp42isom"))
	html := pad([]byte("<html><script>alert(1)</script></html>"))
	svg := pad([]byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`))

	cases := []struct {
		name     string
		fileType string
		ext      string
		data     []byte
		want     bool
	}{
		{"png ok", "cover", ".png", png, true},
		{"jpeg as .jpg ok", "cover", ".jpg", jpeg, true},
		{"jpeg as .jpeg ok", "cover", ".jpeg", jpeg, true},
		{"png disguised as .jpg rejected", "cover", ".jpg", png, false},
		{"html as .png rejected", "cover", ".png", html, false},
		{"svg as .png rejected", "cover", ".png", svg, false},
		{"mp3 ok", "audio", ".mp3", mp3, true},
		{"ogg ok", "audio", ".ogg", ogg, true},
		{"wav ok", "audio", ".wav", wav, true},
		{"m4a ok", "audio", ".m4a", m4a, true},
		{"html as .mp3 rejected", "audio", ".mp3", html, false},
		{"image as audio rejected", "audio", ".mp3", png, false},
		{"audio as cover rejected", "cover", ".png", mp3, false},
		{"empty rejected", "audio", ".mp3", nil, false},
		{"unknown ext rejected", "audio", ".exe", mp3, false},
		{"unknown type rejected", "video", ".mp4", m4a, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := sniffAllowed(tc.fileType, tc.ext, tc.data); got != tc.want {
				t.Fatalf("sniffAllowed(%s,%s) = %v, want %v (head=%q)", tc.fileType, tc.ext, got, tc.want, bytes.TrimRight(tc.data[:min(12, len(tc.data))], "\x00"))
			}
		})
	}
}
