package accountinglikes

import (
	"log/slog"
	"time"

	"github.com/LashkaPashka/AnalyticsService/internal/model"
	modelforstats "github.com/LashkaPashka/AnalyticsService/internal/model/model-for-stats"
)

type Storager interface {
	AddStats(modelforstats.PostLikesAt)
	GetQuantityLikes(string) *modelforstats.PostLikesAt
	UpdateValue(int)
}

type AccLikes struct {
	logger *slog.Logger
	storage Storager
}

func New(storage Storager, logger *slog.Logger) *AccLikes {
	return &AccLikes{
		logger: logger,
		storage: storage,
	}
}

func (a *AccLikes) AddStatsAboutLikes(payload model.UserPostLiked) {
	var postLikedAt *modelforstats.PostLikesAt

	if postLikedAt = a.storage.GetQuantityLikes(payload.DataM.PostID); postLikedAt == nil {
		a.storage.AddStats(modelforstats.PostLikesAt{
			Email: payload.DataM.EmailAuthor,
			PostID: payload.DataM.PostID,
			TotalLikes: 1,
			LastLikedAt: time.Now(),
		})

		a.logger.Info("Model PostLikesAt was created successffuly in MySQL")
		return
	}

	new_total_likes := postLikedAt.TotalLikes + 1
	a.storage.UpdateValue(new_total_likes)
	a.logger.Info("Field total_likes was updated successffully!")
}