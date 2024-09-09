package postgres

import (
	"context"
	"fmt"
	"it-tanlov/pkg/logger"
	"it-tanlov/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewUserRepo(db *pgxpool.Pool, log logger.ILogger) storage.IUserStorage {
	return &userRepo{
		db:  db,
		log: log,
	}
}

func (u *userRepo) Create(ctx context.Context, phone string)  error {
	uid := uuid.New()

	_, err := u.db.Exec(ctx, `
		INSERT INTO users (id, phone) 
		VALUES ($1, $2)
		`,
		uid,
		phone,
	)
	if err != nil {
		 logger.Error(err)
	}

	return  nil
}

func (u *userRepo) AddScore(ctx context.Context, id string)  error {
	
	query := `UPDATE partners SET score = score + 1, updated_at = now() WHERE id = $1`

	if _, err := u.db.Exec(ctx, query, &id); err != nil {
		u.log.Error("error is while updating", logger.Error(err))
		return  err
	}
	return  nil
}

func (u *userRepo) PhoneExist(ctx context.Context, phone string) (bool, error) {
	var exists bool
	err := u.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM users WHERE phone = $1)
	`, phone).Scan(&exists)
	if err != nil {
		fmt.Println("error while checking phone existence:", err)
		return false, err
	}

	return exists, nil
}