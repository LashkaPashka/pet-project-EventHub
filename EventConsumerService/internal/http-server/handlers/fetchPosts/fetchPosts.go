package fetchposts

import (
	"log/slog"
	"net/http"

	"github.com/LashkaPashka/EventConsumerService/internal/lib/filter"
	keysforredis "github.com/LashkaPashka/EventConsumerService/internal/lib/keys-for-redis"
	"github.com/LashkaPashka/EventConsumerService/internal/lib/res"
	"github.com/LashkaPashka/EventConsumerService/internal/service/model"
)

type RedisStorage interface {
	GetPostsInMemory(key string) []model.UserPostCreated
}

func New(rStorage RedisStorage, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts := rStorage.GetPostsInMemory(keysforredis.FeedLatestPost)
		
		var filteredPosts []model.UserPostCreated

		query := r.URL.Query()

		switch {
		case query.Has("tag"):
			tag := query.Get("tag")
			filteredPosts = filter.FilterPostsByTag(tag, posts, logger)

		case query.Has("title"):
			title := query.Get("title")
			filteredPosts = filter.FilterPostsByTitle(title, posts, logger)
		
		default:
			filteredPosts = append([]model.UserPostCreated(nil), posts...)
		}

		res.Encode(w, filteredPosts)
	}
}