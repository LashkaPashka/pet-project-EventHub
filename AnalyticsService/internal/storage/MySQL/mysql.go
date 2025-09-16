package mysql

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/LashkaPashka/AnalyticsService/internal/configs"
	"github.com/LashkaPashka/AnalyticsService/internal/model"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
	cfg *configs.Configs
	logger *slog.Logger
}

func New(cfg *configs.Configs, logger *slog.Logger) *Storage {
	const op = "AnalyticsService.storage.mysql.New"
	
	db, err := sql.Open("mysql", cfg.StoragePath)
	if err != nil {
		logger.Error("Invalid open session to MySQL", slog.String("Error: ", op))
		return nil
	}

	if err := db.Ping(); err != nil {
		logger.Error("Error ping Db", slog.String("Error: ", op))
		return nil
	}

	return &Storage{
		db: db,
		cfg: cfg,
		logger: logger,
	}
}

func (s *Storage) AddPost(payload model.UserPostCreated) {
	const op = "AnalyticsSerivce.storage.mysql.AddPost"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.Begin()
	if err != nil {
		s.logger.Error("Error begin transaction", slog.String("Error: ", op))
		return
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO user_posts_stats(email, total_posts, last_post_at) VALUES (?, ?, ?)`)
	if err != nil {
		s.logger.Error("Error prepare transaction", slog.String("Error: ", op))
		return
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		s.logger.Error("Error comeplete query SQL", slog.String("Error: ", op))
		return
	}

	//id, _ := res.LastInsertId()

	if err = tx.Commit(); err != nil {
		s.logger.Error("Invalid transaction commit", slog.String("Error: ", op))
		return
	}

}