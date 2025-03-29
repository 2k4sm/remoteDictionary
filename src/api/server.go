package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/2k4sm/remoteDictionary/src/cache"
	"github.com/2k4sm/remoteDictionary/src/config"
	"github.com/2k4sm/remoteDictionary/src/handlers"
)

type Server struct {
	server  *http.Server
	handler *handlers.Handler
	cache   *cache.Cache
	config  *config.Config
}

func NewServer(cfg *config.Config) *Server {
	cache := cache.NewCache(cfg.MaxCacheSize, cfg.MaxKeySize, cfg.MaxValueSize)
	handler := handlers.NewHandler(cache)

	return &Server{
		cache:   cache,
		handler: handler,
		config:  cfg,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/put", s.handler.PutHandler)
	mux.HandleFunc("/get", s.handler.GetHandler)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", s.config.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Starting server on port %s", s.config.Port)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
