package postgres

import (
	"context"
	"fmt"
	"it-tanlov/config"
	"it-tanlov/pkg/logger"
	"it-tanlov/storage"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database"          //database is needed for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //postgres is used for database
	_ "github.com/golang-migrate/migrate/v4/source/file"       //file is needed for migration url
	_ "github.com/lib/pq"
)

type Store struct {
	pool *pgxpool.Pool
	log  logger.ILogger
	cfg  config.Config
}

func New(ctx context.Context, cfg config.Config, log logger.ILogger) (storage.IStorage, error) {
	url := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("error while parsing config", logger.Error(err))
		return nil, err
	}

	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("error while connecting to db", logger.Error(err))
		return nil, err
	}

	//migration
	m, err := migrate.New("file://migrations/postgres/", url)
	if err != nil {
		log.Error("error while migrating", logger.Error(err))
		return nil, err
	}

	log.Info("???? came")

	if err = m.Up(); err != nil {
		log.Warning("migration up", logger.Error(err))
		if !strings.Contains(err.Error(), "no change") {
			fmt.Println("entered")
			version, dirty, err := m.Version()
			log.Info("version and dirty", logger.Any("version", version), logger.Any("dirty", dirty))
			if err != nil {
				log.Error("err in checking version and dirty", logger.Error(err))
				return nil, err
			}

			if dirty {
				version--
				if err = m.Force(int(version)); err != nil {
					log.Error("ERR in making force", logger.Error(err))
					return nil, err
				}
			}
			log.Warning("WARNING in migrating", logger.Error(err))
			return nil, err
		}
	}

	log.Info("!!!!! came here")

	return Store{
		pool: pool,
		log:  log,
		cfg:  cfg,
	}, nil
}

func (s Store) Close() {
	s.pool.Close()
}

func (s Store) Partner() storage.IPartnerStorage {
	return NewPartnerRepo(s.pool, s.log)
}

func (s Store) User() storage.IUserStorage {
	return NewUserRepo(s.pool, s.log)
}
