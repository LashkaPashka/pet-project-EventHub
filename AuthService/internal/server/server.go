package server

import (
	"log/slog"
	"time"

	"github.com/LashkaPashka/EventHub/AuthService/internal/api"
	"github.com/LashkaPashka/EventHub/AuthService/internal/configs"
	"github.com/LashkaPashka/EventHub/AuthService/internal/db"
	"golang.org/x/net/context"
)


type Server struct {
	api *api.API
	db *db.Db
	
}

func NewServer(cfg *configs.Configs, logger *slog.Logger) *Server {
	server := Server{}

	server.db = db.NewDb(cfg.Storage_path, logger)
	server.api = api.NewAPI(server.db.Pool, cfg, logger)

	return  &server
}

func (s *Server) Run() {
	s.api.Run(":8080")
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	if err := s.api.Stop(ctx); err != nil {
		panic(err)
	}
}
