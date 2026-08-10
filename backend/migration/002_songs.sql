CREATE TABLE IF NOT EXISTS songs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    cover_url TEXT,
    audio_url TEXT NOT NULL,
    duration INTEGER NOT NULL, -- durasi dalam detik
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_songs_created_at ON songs (created_at DESC);

-- Data contoh supaya bisa langsung dites tanpa insert manual.
-- Silakan hapus/ganti audio_url dengan file lagu asli nanti.
INSERT INTO songs (title, artist, album, cover_url, audio_url, duration)
VALUES
    ('Lagu A', 'Artis Percobaan 1', NULL,
     'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=300&h=300&fit=crop',
     'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3', 372),
    ('Lagu B', 'Artis Percobaan 2', NULL,
     'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=300&h=300&fit=crop',
     'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-2.mp3', 320),
    ('Lagu C', 'Artis Percobaan 3', NULL,
     'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=300&h=300&fit=crop',
     'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-3.mp3', 240)
ON CONFLICT DO NOTHING;