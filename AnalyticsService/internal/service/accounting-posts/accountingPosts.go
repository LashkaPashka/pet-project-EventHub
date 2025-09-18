package accountingposts

import (
	"log/slog"
	"time"

	"github.com/LashkaPashka/AnalyticsService/internal/model"
	modelforstats "github.com/LashkaPashka/AnalyticsService/internal/model/model-for-stats"
)

type Storager interface {
	AddPayload(payload any)
	GetQuantityPosts(email string) int
	UpdateStat(int)
}

type AccPosts struct {
	logger *slog.Logger
	storage Storager
}

func New(storage Storager, logger *slog.Logger) *AccPosts{
	return &AccPosts{
		logger: logger,
		storage: storage,
	}
}

func (a *AccPosts) AddStatsAboutPosts(payload model.UserPostCreated) {
	var totalCost int

	if totalCost = a.storage.GetQuantityPosts(payload.DataM.Email); totalCost == 0 {
		a.storage.AddPayload(modelforstats.UserPostsStats{
			Email: payload.DataM.Email,
			TotatPosts: 1,
			LastPostAt: time.Now(),
		})

		a.logger.Info("Model UserPostsStats was created successfully in MySQL")
		return
	}

	newTotalPosts := totalCost + 1
	a.storage.UpdateStat(newTotalPosts)
	a.logger.Info("Variable total_posts was updated successully!")
}