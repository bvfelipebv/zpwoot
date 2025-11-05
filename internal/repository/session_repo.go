package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"zpwoot/internal/model"
)

// SessionRepository gerencia operações de sessão no banco usando database/sql nativo
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository cria um novo repositório
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create cria uma nova sessão
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

// GetByID busca sessão por ID (UUID)
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

// GetByDeviceJID busca sessão por Device JID
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

// List lista todas as sessões
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

// ListConnected lista sessões conectadas
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

// Update atualiza uma sessão
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

// Delete deleta uma sessão permanentemente
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

// UpdateStatus atualiza o status e flag de conexão da sessão
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

// Count retorna o total de sessões
func (r *SessionRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM sessions`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	return count, nil
}
