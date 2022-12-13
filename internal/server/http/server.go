package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/julinserg/go_home_project/internal/app"
)

type Application interface {
	GetImagePreview(params app.InputParams, header http.Header) ([]byte, int, bool, error)
	ClearCache()
}

type Server struct {
	server   *http.Server
	logger   Logger
	endpoint string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func NewServer(logger Logger, app Application, endpoint string) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:              endpoint,
		Handler:           loggingMiddleware(mux, logger),
		ReadHeaderTimeout: 3 * time.Second,
	}
	ch := previewerHandler{logger, app}
	mux.HandleFunc("/", ch.hellowHandler)
	mux.HandleFunc("/fill/", ch.mainHandler)
	mux.HandleFunc("/clearcache/", ch.clearCacheHandler)
	return &Server{server, logger, endpoint}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("http server started on " + s.endpoint)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.logger.Info("http server stopped")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
