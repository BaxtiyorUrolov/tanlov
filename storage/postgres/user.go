package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"it-tanlov/pkg/logger"
	"it-tanlov/storage"

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

func (u *userRepo) Create(ctx context.Context, id int)  error {

	_, err := u.db.Exec(ctx, `
		INSERT INTO users (id) 
		VALUES ($1)
		`,
		id,
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

func (u *userRepo) IUserTelegramIDExist(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := u.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)
	`, id).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			// Bu holatda hech qanday qator topilmadi
			fmt.Println("No rows found for the provided id")
			return false, nil
		}
		fmt.Println("Error while checking id existence:", err)
		return false, err
	}

	return exists, nil
}
