package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"os/signal"
	"proxy_manager/config"
	v1 "proxy_manager/internal/controller/http/v1"
	"proxy_manager/internal/infrastructure/repository"
	"proxy_manager/internal/usecase"
	"proxy_manager/pkg/logger"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	rootCtx := context.Background()
	l := logger.New(cfg.LogLevel)
	errorChan := make(chan error)

	pgxPool, err := pgxpool.New(rootCtx, cfg.PostgresURL+
		fmt.Sprintf("?pool_max_conns=%d", cfg.PostgresMaxCons))
	if err != nil {
		log.Fatal(err)
	}

	proxyRepo := repository.NewPostgresProxyRepository(rootCtx, pgxPool, time.Minute*time.Duration(cfg.OccupiesExpireTime), l)
	u := usecase.New(proxyRepo)

	handler := gin.New()

	v1.NewRouter(handler, u, l, cfg.ServeSwagger)
	httpServer := serveHttpInBackground(errorChan, handler, fmt.Sprintf(":%s", cfg.HttpPort))

	// For graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		l.Info(fmt.Sprintf("Got signal: %s", sig))
	case err := <-errorChan:
		l.Error(fmt.Sprintf("Got error: %s", err.Error()))
	}

	l.Info("Shutdown httpServer...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		l.Fatal("httpServer shutdown:", err)
	}

	<-ctx.Done()
	l.Info("httpServer exited!")
}

func serveHttpInBackground(errorChan chan<- error, handler *gin.Engine, addr string) *http.Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		errorChan <- fmt.Errorf("httpServer: %w", srv.ListenAndServe())
	}()

	return srv
}
