package service

import (
	"log/slog"

	"github.com/LashkaPashka/AnalyticsService/internal/configs"
)

const (
	userPostCreated = "user.post.created"
	userPostLiked = "user.post.liked"
	orderPaid = "order.paid"
)

type AccountingPoster interface {
	AddStatsAboutPosts()
}

type AccountingLiker interface {
	AddStatsAboutLikes()
}

type AccountingOrder interface {
	AddStatsAboutOrders()
}

type Accountings interface {
	AccountingPoster
	AccountingLiker
	AccountingOrder
}

type AnalyticsService struct {
	cfg *configs.Configs
	logger *slog.Logger
	Accountings
}

func New(accountings Accountings, cfg *configs.Configs, logger *slog.Logger) *AnalyticsService {
	return &AnalyticsService{
		cfg: cfg,
		logger: logger,
		Accountings: accountings,
	}
}

func (a *AnalyticsService) StatsAccounting(topic string, msg []byte) {

	switch topic {
		case userPostCreated:
			a.Accountings.AddStatsAboutPosts()
		case userPostLiked:
			a.Accountings.AddStatsAboutLikes()
		case orderPaid:
			a.Accountings.AddStatsAboutOrders()
		default:
			a.logger.Warn("Unknown topic")
	}


}
