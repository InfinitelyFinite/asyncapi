package apiserver

import (
	"asyncapi/config"
	"asyncapi/store"
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

type ApiServer struct {
	config     *config.Config
	logger     *slog.Logger
	store      *store.Store
	jwtManager *JwtManager
}

func New(config *config.Config, logger *slog.Logger, store *store.Store, jwtManager *JwtManager) *ApiServer {
	return &ApiServer{config: config, logger: logger, store: store, jwtManager: jwtManager}
}

func (s *ApiServer) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (s *ApiServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", s.ping)
	mux.HandleFunc("POST /auth/signup", s.signupHandler())
	mux.HandleFunc("POST /auth/signin", s.signinHandler())

	middleware := NewLoggerMiddleware(s.logger)
	middleware = NewAuthMiddleware(s.jwtManager, s.store.Users)

	server := &http.Server{
		Addr:    net.JoinHostPort(s.config.ApiServerHost, s.config.ApiServerPort),
		Handler: middleware(mux),
	}

	go func() {
		s.logger.Info("apiserver running on", "port ", s.config.ApiServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("apiserver failed to listen and serve")
		}

	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("apiserver failed to shutdown")
		}
	}()

	wg.Wait()
	return nil
}
