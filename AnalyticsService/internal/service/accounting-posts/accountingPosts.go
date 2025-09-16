package accountingposts

import (
	"log/slog"
	"time"

	"github.com/LashkaPashka/AnalyticsService/internal/model"
	modelforstats "github.com/LashkaPashka/AnalyticsService/internal/model/model-for-stats"
)

type Storager interface {
	AddStat(modelforstats.UserPostsStats)
	GetQuantityPosts(string) *modelforstats.UserPostsStats
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
	var postsStats *modelforstats.UserPostsStats

	if postsStats = a.storage.GetQuantityPosts(payload.DataM.Email); postsStats == nil {
		a.storage.AddStat(modelforstats.UserPostsStats{
			Email: payload.DataM.Email,
			TotatPosts: 1,
			LastPostAt: time.Now(),
		})

		a.logger.Info("Model UserPostsStats was created successfully in MySQL")
		return
	}

	newTotalPosts := postsStats.TotatPosts + 1
	a.storage.UpdateStat(newTotalPosts)
	a.logger.Info("Variable total_posts was updated successully!")
}