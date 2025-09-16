package modelforstats

import "time"

type UserPostsStats struct {
	Email string `json:"email"`
	TotatPosts int `json:"total_posts"`
	LastPostAt time.Time `json:"last_post_at"`
}

type PostLikesAt struct {
	PostID string `json:"post_id"`
	Email string `json:"email"`
	TotalLikes int `json:"total_likes"`
	LastLikedAt time.Time `json:"last_liked_at"`
}

type UserOrdersStats struct {
	Email string `json:"email"`
	TotalOrders int `json:"total_orders"`
	TotalAmmount float64 `json:"total_ammount"`
	LastOrderAt time.Time `json:"last_order_at"`
}