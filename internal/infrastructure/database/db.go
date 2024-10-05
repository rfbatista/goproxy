package database

import (
	"context"
	"database/sql"
	"goproxy/internal/shared"
	"path"

	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(lc fx.Lifecycle, log *zap.Logger) *sql.DB {
	root := shared.FindProjectRoot()
	filePath := path.Join(root, "./db/blocked_ips.db")
	log.Info("trying to connect to database...")
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	log.Info("connected to database")
	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				log.Info("closing database connection")
				err := db.Close()
				if err != nil {
					log.Warn("failed to close database connection", zap.Error(err))
					return err
				}
				return nil
			},
		},
	)
	return db
}
