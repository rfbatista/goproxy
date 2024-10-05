package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"regexp"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var re = regexp.MustCompile("[^0-9]")

type BlockedIpsRepository struct {
	db *sql.DB
}

func NewBlocketIpsRepository(
	lc fx.Lifecycle,
	db *sql.DB,
	log *zap.Logger,
) (*BlockedIpsRepository, error) {
	if db == nil {
		log.Error("no db connection provided")
		return nil, errors.New("no db connection provided")
	}
	repo := &BlockedIpsRepository{db: db}
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				err := repo.Migrate()
				if err != nil {
					log.Error("failed to create migration", zap.Error(err))
					return err
				}
				return nil
			},
		},
	)
	return repo, nil
}

func (b *BlockedIpsRepository) ListBlockedIPs() (map[string]bool, error) {
	ips := make(map[string]bool)
	rows, err := b.db.Query("SELECT ip FROM blocked_ips")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return ips, err
		}
		ips[re.ReplaceAllString(ip, "")] = true
	}
	return ips, nil
}

func (b *BlockedIpsRepository) InsertBlockedIP(ip string) error {
	stmt, err := b.db.Prepare("INSERT INTO blocked_ips (ip) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ip)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlockedIpsRepository) RemoveBlockedIP(ip string) error {
	stmt, err := b.db.Prepare("DELETE FROM blocked_ips WHERE ip = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ip)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlockedIpsRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS blocked_ips (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL UNIQUE
	);`
	_, err := b.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
