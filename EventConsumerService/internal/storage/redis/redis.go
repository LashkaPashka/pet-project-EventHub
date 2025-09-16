package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LashkaPashka/EventConsumerService/internal/configs"
	"github.com/LashkaPashka/EventConsumerService/internal/lib/decode"
	"github.com/LashkaPashka/EventConsumerService/internal/lib/encode"
	"github.com/LashkaPashka/EventConsumerService/internal/service/model"
	"github.com/go-redis/redis/v8"
)

const (
	db = 0
)

type RedisSt struct {
	RDb *redis.Client
	cfg *configs.Configs
	logger *slog.Logger
}

func New(cfg *configs.Configs, logger *slog.Logger) *RedisSt {
	var rdClient = RedisSt{
		cfg: cfg,
		logger: logger,
	}
	
	rdClient.RDb = checkConnection(rdClient.cfg, rdClient.logger)
	if rdClient.RDb == nil {
		return nil
	}

	return &rdClient
}

func checkConnection(cfg *configs.Configs, logger *slog.Logger) *redis.Client {
	const op = "EventConsumerService.storage.redis.checkConnection"
	
	rdClient := redis.NewClient(&redis.Options{
		Addr: cfg.Raddr,
		Password: cfg.Rpassword,
		DB: db,
	})

	if res := rdClient.Ping(context.Background()); res.Err() != nil {
		logger.Error("Invalid ping redis Db", slog.String("Error: ", op))
		return nil
	}

	return rdClient
}

func (r *RedisSt) GetPostsInMemory(key string) []model.UserPostCreated {
 	var posts []model.UserPostCreated

	iter := r.RDb.Scan(context.Background(), 0, key+":*", 0).Iterator()

	for iter.Next(context.Background()) {
		key := iter.Val()
		post, err := r.RDb.Get(context.Background(), key).Result()
		if err == nil {
			fmt.Printf("%s => %s\n", key, post)
		}

		postDecode := decode.DecodeJSON[model.UserPostCreated]([]byte(post), r.logger)
		posts = append(posts, postDecode)
	}

	r.logger.Info("Data was recieved successfuly", slog.Any("Data", posts))

	return posts
}


func (r *RedisSt) SaveInMemory(data any, key string, id string) {
	const op = "EventConsumerService.storage.redis.SaveInMemory"

	bytes := encode.EncodeBytes(data, r.logger)
	if bytes == nil {
		return 
	}
	
	cKey := fmt.Sprintf("%s:%s", key, id)

	if err := r.RDb.Set(context.Background(), cKey, string(bytes), 10*time.Minute).Err(); err != nil {
		r.logger.Error("Invalid set cache in Redis", slog.String("Error: ", op))
		return
	}

	r.logger.Info("Data was added sucessfully in Redis")
}