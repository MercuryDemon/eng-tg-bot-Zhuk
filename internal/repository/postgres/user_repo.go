package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
)

type UserRepo struct {
	db     *sql.DB
	logger *zerolog.Logger
}

func NewUserRepo(db *sql.DB, parentLogger *zerolog.Logger) *UserRepo {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "user_repo").Logger()

	return &UserRepo{
		db:     db,
		logger: &logger,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, id int64, username string) error {
	const op = "CreateUser"

	const query = `
		INSERT INTO users (tg_id, username, created_at)
		VALUES ($1, NULLIF($2, ''), now())
		ON CONFLICT (tg_id) DO UPDATE
		SET username = COALESCE(EXCLUDED.username, users.username),
		created_at = COALESCE(users.created_at, EXCLUDED.created_at);
	`

	if _, err := r.db.ExecContext(ctx, query, id, username); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id int64) error {
	const op = "DeleteUser"

	const query = `
		DELETE FROM users
		WHERE tg_id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rows == 0 {
		r.logger.Warn().
			Int64("tg_id", id).
			Msgf("%s: user not found", op)
	}

	return nil
}

func (r *UserRepo) SetActiveDictionaryID(ctx context.Context, userID int64, dictionaryID string) error {
	const op = "SetActiveDictionaryID"

	const query = `
		UPDATE users
		SET active_dictionary_id = $2
		WHERE tg_id = $1;
	`

	_, err := r.db.ExecContext(ctx, query, userID, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *UserRepo) GetActiveDictionaryID(ctx context.Context, userID int64) (string, error) {
	const op = "GetActiveDictionaryID"

	const query = `
		SELECT COALESCE(active_dictionary_id::text, '')
		FROM users
		WHERE tg_id = $1;
	`

	var dictionaryID string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&dictionaryID)
	if err != nil {
		// TODO: return error
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return dictionaryID, nil
}

func (r *UserRepo) ClearActiveDictionaryID(ctx context.Context, userID int64) error {
	const op = "ClearActiveDictionaryID"

	const query = `
		UPDATE users
		SET active_dictionary_id = NULL
		WHERE tg_id = $1;
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
