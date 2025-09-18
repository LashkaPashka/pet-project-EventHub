package filter

import (
	"log/slog"
	"strings"

	"github.com/LashkaPashka/EventConsumerService/internal/model"
)

func FilterPostsByTag(q string, posts []model.UserPostCreated, logger *slog.Logger) []model.UserPostCreated {
	var filtered []model.UserPostCreated
	
	for _, post := range posts {
		for _, tag := range post.DataM.Tags {
			if strings.EqualFold(q, tag) {
				filtered = append(filtered, post)
			}
		}
	}

	return filtered
}

func FilterPostsByTitle(q string, posts []model.UserPostCreated, logger *slog.Logger) []model.UserPostCreated {
	var filtered []model.UserPostCreated
	
	for _, post := range posts {
		title := post.DataM.Title
		if strings.Contains(strings.ToLower(title), strings.ToLower(q)) {
			filtered = append(filtered, post)
		}
	}

	return filtered
}