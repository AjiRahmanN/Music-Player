export interface AdminUser {
  id: string;
  username: string;
  email: string;
  isAdmin: boolean;
  createdAt: string;
}

export interface AdminUserListResponse {
  data?: { users: AdminUser[] };
  error?: string;
  message?: string;
}

export interface AdminUserUpdateResponse {
  data?: { user: AdminUser };
  message?: string;
  error?: string;
}

export interface AdminPlaylist {
  id: string;
  userId: string;
  name: string;
  description?: string | null;
  coverUrl?: string | null;
  createdAt: string;
  songCount: number;
  ownerUsername: string;
}

export interface AdminPlaylistListResponse {
  data?: { playlists: AdminPlaylist[] };
  error?: string;
  message?: string;
}
