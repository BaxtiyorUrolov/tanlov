package postgres

import (
	"context"
	"fmt"
	"it-tanlov/api/models"
	"it-tanlov/pkg/logger"
	"it-tanlov/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type partnerRepo struct {
	db  *pgxpool.Pool
	log logger.ILogger
}

func NewPartnerRepo(db *pgxpool.Pool, log logger.ILogger) storage.IPartnerStorage {
	return &partnerRepo{
		db:  db,
		log: log,
	}
}

func (p *partnerRepo) Create(ctx context.Context, createPartner models.CreatePartner) (string, error) {
	uid := uuid.New()

	_, err := p.db.Exec(ctx, `
		INSERT INTO partners (id, full_name, phone, email, video_link) 
		VALUES ($1, $2, $3, $4, $5)
		`,
		uid,
		createPartner.FullName,
		createPartner.Phone,
		createPartner.Email,
		createPartner.VideoLink,
	)
	if err != nil {
		p.log.Error("error is while inserting data", logger.Error(err))
	}

	return uid.String(), err
}

func (p *partnerRepo) GetByID(ctx context.Context, pKey models.PrimaryKey) (models.Partner, error) {
	partner := models.Partner{}

	query := `
		SELECT id, full_name, phone, email, video_link, score 
		FROM partners WHERE id = $1
	`
	if err := p.db.QueryRow(ctx, query, pKey.ID).Scan(
		&partner.ID,
		&partner.FullName,
		&partner.Phone,
		&partner.Email,
		&partner.VideoLink,
		&partner.Score,
	); err != nil {
		p.log.Error("error while selecting partner by id", logger.Error(err))
		return partner, err
	}

	return partner, nil
}

func (p *partnerRepo) GetList(ctx context.Context, request models.GetListRequest) (models.PartnerResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		query, countQuery string
		count             = 0
		partners          = []models.Partner{}
		search            = request.Search
	)

	countQuery = `SELECT count(1) FROM partners WHERE video_verify = true`
	if search != "" {
		countQuery += fmt.Sprintf(` and full_name ilike '%%%s%%'`, search)
	}
	if err := p.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		p.log.Error("error is while scanning count", logger.Error(err))
		return models.PartnerResponse{}, err
	}

	query = `SELECT id, full_name, video_link, score, created_at 
	         FROM partners 
	         WHERE video_verify = true`
	if search != "" {
		query += fmt.Sprintf(` and full_name ilike '%%%s%%'`, search)
	}

	// Sort by score in descending order
	query += ` ORDER BY score DESC LIMIT $1 OFFSET $2`

	rows, err := p.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		p.log.Error("error is while selecting all", logger.Error(err))
		return models.PartnerResponse{}, err
	}

	for rows.Next() {
		partner := models.Partner{}
		if err = rows.Scan(
			&partner.ID,
			&partner.FullName,
			&partner.VideoLink,
			&partner.Score,
			&partner.CreatedAt,
		); err != nil {
			p.log.Error("error is while scanning partner", logger.Error(err))
			return models.PartnerResponse{}, err
		}

		partners = append(partners, partner)
	}

	return models.PartnerResponse{
		Partners: partners,
		Count:    count,
	}, nil
}

func (p *partnerRepo) Update(ctx context.Context, id string) error {
	query := `UPDATE partners SET video_verify = true, updated_at = now() WHERE id = $1`
	if _, err := p.db.Exec(ctx, query, id); err != nil {
		p.log.Error("error is while updating partner with image", logger.Error(err))
		return err
	}
	return nil
}

func (p *partnerRepo) Delete(ctx context.Context, id string) error {

	fmt.Println()
	fmt.Println(id)
	query := `DELETE FROM partners WHERE id = $1`
	if _, err := p.db.Exec(ctx, query, id); err != nil {
		p.log.Error("error is while deleting", logger.Error(err))
		return err
	}
	return nil
}

func (p *partnerRepo) PhoneExist(ctx context.Context, phone string) (bool, error) {
	var exists bool
	err := p.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM partners WHERE phone = $1)
	`, phone).Scan(&exists)
	if err != nil {
		fmt.Println("error while checking phone existence:", err)
		return false, err
	}

	return exists, nil
}

func (p *partnerRepo) EmailExist(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := p.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM partners WHERE email = $1)
	`, email).Scan(&exists)
	if err != nil {
		fmt.Println("error while checking email existence:", err)
		return false, err
	}

	return exists, nil
}

func (p *partnerRepo) VideoLinkExist(ctx context.Context, videoLink string) (bool, error) {
	var exists bool
	err := p.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM partners WHERE video_link = $1)
	`, videoLink).Scan(&exists)
	if err != nil {
		fmt.Println("error while checking video_link existence:", err)
		return false, err
	}

	return exists, nil
}
