package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"pet-platform/backend/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	CreateSession(ctx context.Context, session *domain.UserSession) error
	GetSessionByTokenHash(ctx context.Context, hash string) (*domain.UserSession, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RegisterDevice(ctx context.Context, device *domain.UserDevice) error
	GetNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error)
	MarkNotificationRead(ctx context.Context, notifID, userID uuid.UUID) error
	CreateNotification(ctx context.Context, n *domain.Notification) error
}

type sqlUserRepo struct {
	db *DB
}

func NewUserRepository(db *DB) UserRepository {
	return &sqlUserRepo{db: db}
}

func (r *sqlUserRepo) CreateUser(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, phone_number, role, timezone, allow_direct_calls, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	u.IsActive = true
	if u.Role == "" {
		u.Role = domain.RolePetOwner
	}
	if u.Timezone == "" {
		u.Timezone = "UTC"
	}

	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName,
		u.PhoneNumber, u.Role, u.Timezone, u.AllowDirectCalls, u.IsActive, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *sqlUserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone_number, phone_verified, email_verified,
		       role, avatar_url, timezone, allow_direct_calls, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	var u domain.User
	var phone, avatar sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &phone,
		&u.PhoneVerified, &u.EmailVerified, &u.Role, &avatar, &u.Timezone,
		&u.AllowDirectCalls, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if phone.Valid {
		u.PhoneNumber = phone.String
	}
	if avatar.Valid {
		u.AvatarURL = avatar.String
	}
	return &u, nil
}

func (r *sqlUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone_number, phone_verified, email_verified,
		       role, avatar_url, timezone, allow_direct_calls, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	var u domain.User
	var phone, avatar sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &phone,
		&u.PhoneVerified, &u.EmailVerified, &u.Role, &avatar, &u.Timezone,
		&u.AllowDirectCalls, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if phone.Valid {
		u.PhoneNumber = phone.String
	}
	if avatar.Valid {
		u.AvatarURL = avatar.String
	}
	return &u, nil
}

func (r *sqlUserRepo) CreateSession(ctx context.Context, s *domain.UserSession) error {
	query := `
		INSERT INTO user_sessions (id, user_id, refresh_token_hash, user_agent, client_ip, expires_at, is_revoked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.CreatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query, s.ID, s.UserID, s.RefreshTokenHash, s.UserAgent, s.ClientIP, s.ExpiresAt, s.IsRevoked, s.CreatedAt)
	return err
}

func (r *sqlUserRepo) GetSessionByTokenHash(ctx context.Context, hash string) (*domain.UserSession, error) {
	query := `
		SELECT id, user_id, refresh_token_hash, user_agent, client_ip, expires_at, is_revoked, created_at
		FROM user_sessions
		WHERE refresh_token_hash = $1 AND is_revoked = FALSE
	`
	var s domain.UserSession
	err := r.db.QueryRowContext(ctx, query, hash).Scan(
		&s.ID, &s.UserID, &s.RefreshTokenHash, &s.UserAgent, &s.ClientIP, &s.ExpiresAt, &s.IsRevoked, &s.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &s, err
}

func (r *sqlUserRepo) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `UPDATE user_sessions SET is_revoked = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, sessionID)
	return err
}

func (r *sqlUserRepo) RegisterDevice(ctx context.Context, d *domain.UserDevice) error {
	query := `
		INSERT INTO user_devices (id, user_id, device_token, platform, device_name, app_version, is_active, last_seen_at)
		VALUES ($1, $2, $3, $4, $5, $6, TRUE, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, device_token) 
		DO UPDATE SET platform = EXCLUDED.platform, device_name = EXCLUDED.device_name, 
		              app_version = EXCLUDED.app_version, is_active = TRUE, last_seen_at = CURRENT_TIMESTAMP
	`
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	_, err := r.db.ExecContext(ctx, query, d.ID, d.UserID, d.DeviceToken, d.Platform, d.DeviceName, d.AppVersion)
	return err
}

func (r *sqlUserRepo) GetNotifications(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error) {
	query := `
		SELECT id, user_id, pet_id, notification_type, title, body, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []domain.Notification
	for rows.Next() {
		var n domain.Notification
		var petID sql.NullString
		if err := rows.Scan(&n.ID, &n.UserID, &petID, &n.NotificationType, &n.Title, &n.Body, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		if petID.Valid {
			pid, _ := uuid.Parse(petID.String)
			n.PetID = &pid
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (r *sqlUserRepo) MarkNotificationRead(ctx context.Context, notifID, userID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = TRUE, read_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, notifID, userID)
	return err
}

func (r *sqlUserRepo) CreateNotification(ctx context.Context, n *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, pet_id, notification_type, title, body, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, FALSE, CURRENT_TIMESTAMP)
	`
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	_, err := r.db.ExecContext(ctx, query, n.ID, n.UserID, n.PetID, n.NotificationType, n.Title, n.Body)
	return err
}
