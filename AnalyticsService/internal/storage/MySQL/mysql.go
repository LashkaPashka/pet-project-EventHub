package mysql

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/LashkaPashka/AnalyticsService/internal/configs"
	modelforstats "github.com/LashkaPashka/AnalyticsService/internal/model/model-for-stats"
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

func (s *Storage) AddPayload(payload any) {
	const op = "AnalyticsSerivce.storage.mysql.AddPost"
	
	query, args := PrepareDataForAdd(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.Begin()
	if err != nil {
		s.logger.Error("Error begin transaction", slog.String("Error: ", op))
		return
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		s.logger.Error("Error prepare transaction", slog.String("Error: ", op))
		return
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
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

func (s *Storage) GetQuantityPosts(email string) int {
	const op = "AnlyticsService.storage.mysql.GetStats"

	rows, err := s.db.Query(`SELECT total_posts FROM user_posts_stats WHERE email = ?`, email)
	if err != nil {
		s.logger.Error("Invalid executes query", slog.String("Error: ", op))
		return -1
	}

	defer rows.Close()

	var totalCost int
	err = rows.Scan(&totalCost)
	if err != nil {
		s.logger.Error("Invalid scan the columns", slog.String("Error: ", op))
		return -1
	}

	return totalCost
}

func (s *Storage) UpdateStatPost(updatedValue int, email string) {
	const op = "AnalyticsSerivce.storage.mysql.AddPost"
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.Begin()
	if err != nil {
		s.logger.Error("Error begin transaction", slog.String("Error: ", op))
		return
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE user_posts_stats SET total_cost = ? where email = ?`)
	if err != nil {
		s.logger.Error("Error prepare transaction", slog.String("Error: ", op))
		return
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, updatedValue, email)
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


func PrepareDataForAdd(payload any) (string, []any) {
	var query string
	var args []any

	switch v := payload.(type) {
		case modelforstats.UserPostsStats:
			query = `INSERT INTO user_posts_stats(email, total_posts, last_post_at) VALUES (?, ?, ?)`
			args = append(args, v.Email, v.TotatPosts, v.LastPostAt)
		case modelforstats.PostLikesAt:
			query = `INSERT INTO posts_likes_stats(post_id, email, total_likes, last_liked_at) VALUES (?, ?, ?, ?)`
			args = append(args, v.PostID, v.Email, v.TotalLikes, v.LastLikedAt)
		case modelforstats.UserOrdersStats:
			query = `INSERT INTO user_orders_stats(email, total_orders, total_amount, last_order_at) VALUES (?, ?, ?, ?)`
			args = append(args, v.Email, v.TotalOrders, v.TotalAmmount, v.LastOrderAt)
	}

	return  query, args
}