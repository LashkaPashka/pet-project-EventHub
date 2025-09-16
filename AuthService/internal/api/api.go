package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	requser "github.com/LashkaPashka/EventHub/AuthService/internal/api/req_user"
	"github.com/LashkaPashka/EventHub/AuthService/internal/configs"
	serve "github.com/LashkaPashka/EventHub/AuthService/internal/service"
	"github.com/LashkaPashka/EventHub/AuthService/internal/storage"
	"github.com/LashkaPashka/EventHub/AuthService/pkg/jwt"
	"github.com/LashkaPashka/EventHub/AuthService/pkg/req"
	"github.com/LashkaPashka/EventHub/AuthService/pkg/res"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)


type API struct {
	router *mux.Router
	logger *slog.Logger
	pool   *pgxpool.Pool
	cfg    *configs.Configs
	serve  *http.Server
}

func NewAPI(pool *pgxpool.Pool, cfg *configs.Configs, logger *slog.Logger) *API{
	api := API{
		router: mux.NewRouter(),
		logger: logger,
		cfg: cfg,
		pool: pool,
	}
	api.serve = &http.Server{Handler: api.router}

	api.endPoints()

	return &api
}

func (api *API) Run(addr string) error {
	api.serve.Addr = addr

	return api.serve.ListenAndServe()
}

func (api *API) Stop(ctx context.Context) error {
	return api.serve.Shutdown(ctx)
}


func (api *API) endPoints() {

	api.router.Handle("/register", register(api.pool, api.logger)).Methods(http.MethodPost)
	api.router.Handle("/login", login(api.pool, api.cfg, api.logger)).Methods(http.MethodPost)
	api.router.Handle("/logout", logout()).Methods(http.MethodDelete)
}

func register(pool *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := req.HandleBody[requser.RegisterRequest](w, r, logger)
		if body == nil {
			return 
		}

		storage := storage.NewStorage(pool, logger)

		service := serve.NewAuthSerice(storage)
		
		email := service.Register(body.Email, body.Password, body.Username, logger)
		
		payload := &requser.RegisterResponse{
			Detail: fmt.Sprintf("Email: %s - created", email),
		}

		res.Json(w, payload, http.StatusOK)
	}
}

func login(pool *pgxpool.Pool, cfg *configs.Configs, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := req.HandleBody[requser.LoginRequest](w, r, logger)
		if body == nil { return }

		storage := storage.NewStorage(pool, logger)

		service := serve.NewAuthSerice(storage)

		email := service.Login(body.Email, body.Password, logger)

		if email == "" {
			return
		}

		token, err := jwt.NewJWT(cfg.JWT_Secret, logger).CreateJWT(email)
		if err != nil {
			logger.Error("Error's creating JWT")
			return 
		}

		res.Json(w, &requser.LoginResponse{
			Detail: "You're successfuly authenticated!",
			Token: token,
		}, http.StatusOK)
	}
}

func logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}