package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zpwoot/internal/model"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (
			name, device_jid, status, connected,
			qr_code, proxy_config, webhook_config,
			apikey, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, NOW(), NOW()
		) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		session.Name, session.DeviceJID, session.Status, session.Connected,
		session.QRCode, session.ProxyConfig, session.WebhookConfig,
		session.APIKey,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *SessionRepository) GetByID(ctx context.Context, id string) (*model.Session, error) {
	query := `
		SELECT
			id, name, device_jid, status, connected,
			qr_code, proxy_config, webhook_config,
			apikey, created_at, updated_at
		FROM sessions
		WHERE id = $1
	`

	session := &model.Session{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.Name, &session.DeviceJID, &session.Status, &session.Connected,
		&session.QRCode, &session.ProxyConfig, &session.WebhookConfig,
		&session.APIKey, &session.CreatedAt, &session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) GetByDeviceJID(ctx context.Context, deviceJID string) (*model.Session, error) {
	query := `
		SELECT
			id, name, device_jid, status, connected,
			qr_code, proxy_config, webhook_config,
			apikey, created_at, updated_at
		FROM sessions
		WHERE device_jid = $1
	`

	session := &model.Session{}

	err := r.db.QueryRowContext(ctx, query, deviceJID).Scan(
		&session.ID, &session.Name, &session.DeviceJID, &session.Status, &session.Connected,
		&session.QRCode, &session.ProxyConfig, &session.WebhookConfig,
		&session.APIKey, &session.CreatedAt, &session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) List(ctx context.Context) ([]*model.Session, error) {
	query := `
		SELECT
			id, name, device_jid, status, connected,
			qr_code, proxy_config, webhook_config,
			apikey, created_at, updated_at
		FROM sessions
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	sessions := []*model.Session{}
	for rows.Next() {
		session := &model.Session{}

		err := rows.Scan(
			&session.ID, &session.Name, &session.DeviceJID, &session.Status, &session.Connected,
			&session.QRCode, &session.ProxyConfig, &session.WebhookConfig,
			&session.APIKey, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

func (r *SessionRepository) ListConnected(ctx context.Context) ([]*model.Session, error) {
	query := `
		SELECT
			id, name, device_jid, status, connected,
			qr_code, proxy_config, webhook_config,
			apikey, created_at, updated_at
		FROM sessions
		WHERE connected = true AND status = 'connected'
		ORDER BY updated_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list connected sessions: %w", err)
	}
	defer rows.Close()

	sessions := []*model.Session{}
	for rows.Next() {
		session := &model.Session{}

		err := rows.Scan(
			&session.ID, &session.Name, &session.DeviceJID, &session.Status, &session.Connected,
			&session.QRCode, &session.ProxyConfig, &session.WebhookConfig,
			&session.APIKey, &session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

func (r *SessionRepository) Update(ctx context.Context, session *model.Session) error {
	query := `
		UPDATE sessions SET
			name = $1,
			device_jid = $2,
			status = $3,
			connected = $4,
			qr_code = $5,
			proxy_config = $6,
			webhook_config = $7,
			apikey = $8,
			updated_at = NOW()
		WHERE id = $9
	`

	result, err := r.db.ExecContext(ctx, query,
		session.Name, session.DeviceJID, session.Status, session.Connected,
		session.QRCode, session.ProxyConfig, session.WebhookConfig,
		session.APIKey, session.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	session.UpdatedAt = time.Now()
	return nil
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, id string, status string, connected bool) error {
	query := `
		UPDATE sessions
		SET status = $1, connected = $2, updated_at = NOW()
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, status, connected, id)
	if err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *SessionRepository) UpdateQRCode(ctx context.Context, id string, qrCode string) error {
	query := `
		UPDATE sessions
		SET qr_code = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, qrCode, id)
	if err != nil {
		return fmt.Errorf("failed to update QR code: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *SessionRepository) UpdateDeviceJID(ctx context.Context, id string, deviceJID string) error {
	query := `
		UPDATE sessions
		SET device_jid = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, deviceJID, id)
	if err != nil {
		return fmt.Errorf("failed to update device JID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *SessionRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM sessions`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	return count, nil
}
