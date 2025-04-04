package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/2k4sm/remoteDictionary/src/cache"
	"github.com/2k4sm/remoteDictionary/src/config"
	"github.com/2k4sm/remoteDictionary/src/handlers"
	"github.com/rs/cors"
)

type Server struct {
	server  *http.Server
	handler *handlers.Handler
	cache   *cache.Cache
	config  *config.Config
}

func NewServer(cfg *config.Config) *Server {
	cache := cache.NewCache(cfg.MaxKeySize, cfg.MaxValueSize)
	go cache.MonitorMemoryUsage()
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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := map[string]string{
			"status": "healthy",
			"time":   time.Now().String(),
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	})
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", s.config.Port),
		Handler:      handler,
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
