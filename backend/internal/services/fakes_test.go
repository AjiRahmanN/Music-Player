package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"music-player-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

// ---- fake UserRepository (hanya method yang dipakai AuthService) ----

type fakeUserRepo struct {
	byEmail map[string]*models.User
	byID    map[string]*models.User
}

func newFakeUserRepo(users ...*models.User) *fakeUserRepo {
	r := &fakeUserRepo{byEmail: map[string]*models.User{}, byID: map[string]*models.User{}}
	for _, u := range users {
		r.byEmail[u.Email] = u
		r.byID[u.ID] = u
	}
	return r
}

func (r *fakeUserRepo) Create(ctx context.Context, u *models.User) (*models.User, error) {
	r.byEmail[u.Email], r.byID[u.ID] = u, u
	return u, nil
}
func (r *fakeUserRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}
func (r *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}
func (r *fakeUserRepo) Update(ctx context.Context, u *models.User) error { return nil }
func (r *fakeUserRepo) Delete(ctx context.Context, id string) error      { return nil }
func (r *fakeUserRepo) ExistsByEmail(ctx context.Context, e string) (bool, error) {
	_, ok := r.byEmail[e]
	return ok, nil
}
func (r *fakeUserRepo) ExistsByUsername(ctx context.Context, n string) (bool, error) {
	for _, u := range r.byID {
		if u.Username == n {
			return true, nil
		}
	}
	return false, nil
}
func (r *fakeUserRepo) Count(ctx context.Context) (int64, error) { return int64(len(r.byID)), nil }
func (r *fakeUserRepo) FindAll(ctx context.Context, l, o int) ([]*models.User, error) {
	return nil, nil
}
func (r *fakeUserRepo) FindByIDs(ctx context.Context, ids []string) ([]*models.User, error) {
	return nil, nil
}

// ---- fake RefreshTokenRepository ----

type fakeRTRepo struct {
	mu     sync.Mutex
	tokens map[string]rtRow
}
type rtRow struct {
	userID  string
	expires time.Time
}

func newFakeRTRepo() *fakeRTRepo { return &fakeRTRepo{tokens: map[string]rtRow{}} }

func (r *fakeRTRepo) Create(ctx context.Context, userID, hash string, exp time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[hash] = rtRow{userID, exp}
	return nil
}
func (r *fakeRTRepo) FindByHash(ctx context.Context, hash string) (string, time.Time, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	row, ok := r.tokens[hash]
	if !ok {
		return "", time.Time{}, pgx.ErrNoRows
	}
	return row.userID, row.expires, nil
}
func (r *fakeRTRepo) DeleteByHash(ctx context.Context, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tokens, hash)
	return nil
}
func (r *fakeRTRepo) DeleteByUser(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for h, row := range r.tokens {
		if row.userID == userID {
			delete(r.tokens, h)
		}
	}
	return nil
}
func (r *fakeRTRepo) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.tokens)
}

// ---- fake PlaylistRepository ----

type fakePlaylistRepo struct {
	playlists map[string]*models.Playlist
	songs     map[string][]string // playlistID -> songIDs
}

func newFakePlaylistRepo(ps ...*models.Playlist) *fakePlaylistRepo {
	r := &fakePlaylistRepo{playlists: map[string]*models.Playlist{}, songs: map[string][]string{}}
	for _, p := range ps {
		r.playlists[p.ID] = p
	}
	return r
}
func (r *fakePlaylistRepo) Create(ctx context.Context, p *models.Playlist) error {
	p.ID = "new-id"
	r.playlists[p.ID] = p
	return nil
}
func (r *fakePlaylistRepo) ListByUser(ctx context.Context, userID string) ([]models.Playlist, error) {
	var out []models.Playlist
	for _, p := range r.playlists {
		if p.UserID == userID {
			out = append(out, *p)
		}
	}
	return out, nil
}
func (r *fakePlaylistRepo) ListAll(ctx context.Context) ([]models.Playlist, error) {
	var out []models.Playlist
	for _, p := range r.playlists {
		out = append(out, *p)
	}
	return out, nil
}
func (r *fakePlaylistRepo) FindByID(ctx context.Context, id string) (*models.Playlist, error) {
	return r.playlists[id], nil // nil jika tidak ada
}
func (r *fakePlaylistRepo) Delete(ctx context.Context, id string) error {
	delete(r.playlists, id)
	return nil
}
func (r *fakePlaylistRepo) AddSong(ctx context.Context, pid, sid string) error {
	r.songs[pid] = append(r.songs[pid], sid)
	return nil
}
func (r *fakePlaylistRepo) RemoveSong(ctx context.Context, pid, sid string) error {
	var keep []string
	for _, s := range r.songs[pid] {
		if s != sid {
			keep = append(keep, s)
		}
	}
	r.songs[pid] = keep
	return nil
}
func (r *fakePlaylistRepo) GetSongs(ctx context.Context, pid string) ([]models.Song, error) {
	return []models.Song{}, nil
}

// ---- fake SongRepository ----

type fakeSongRepo struct {
	songs      map[string]*models.Song
	lastLimit  int
	lastOffset int
	created    *models.Song
}

func newFakeSongRepo(ss ...*models.Song) *fakeSongRepo {
	r := &fakeSongRepo{songs: map[string]*models.Song{}}
	for _, s := range ss {
		r.songs[s.ID] = s
	}
	return r
}
func (r *fakeSongRepo) List(ctx context.Context, limit, offset int) ([]models.Song, error) {
	r.lastLimit, r.lastOffset = limit, offset
	return []models.Song{}, nil
}
func (r *fakeSongRepo) Search(ctx context.Context, q string, limit, offset int) ([]models.Song, error) {
	r.lastLimit, r.lastOffset = limit, offset
	return []models.Song{}, nil
}
func (r *fakeSongRepo) FindByID(ctx context.Context, id string) (*models.Song, error) {
	return r.songs[id], nil
}
func (r *fakeSongRepo) Create(ctx context.Context, s *models.Song) error {
	s.ID = "song-new"
	r.created = s
	r.songs[s.ID] = s
	return nil
}
func (r *fakeSongRepo) Update(ctx context.Context, s *models.Song) error {
	r.songs[s.ID] = s
	return nil
}
func (r *fakeSongRepo) Delete(ctx context.Context, id string) error { delete(r.songs, id); return nil }

// ---- fake AdminUserRepository ----

type fakeAdminRepo struct {
	users     map[string]*models.AdminUser
	updateErr error
}

func (r *fakeAdminRepo) List(ctx context.Context) ([]models.AdminUser, error) { return nil, nil }
func (r *fakeAdminRepo) FindByID(ctx context.Context, id string) (*models.AdminUser, error) {
	if u, ok := r.users[id]; ok {
		c := *u
		return &c, nil
	}
	return nil, nil
}
func (r *fakeAdminRepo) Update(ctx context.Context, u *models.AdminUser) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	c := *u
	r.users[u.ID] = &c
	return nil
}
func (r *fakeAdminRepo) Delete(ctx context.Context, id string) error {
	delete(r.users, id)
	return nil
}
func (r *fakeAdminRepo) CountAdmins(ctx context.Context) (int, error) {
	n := 0
	for _, u := range r.users {
		if u.IsAdmin {
			n++
		}
	}
	return n, nil
}
