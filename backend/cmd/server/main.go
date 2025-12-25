package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/what-writers-like/backend/internal/handler"
	"github.com/what-writers-like/backend/internal/infrastructure/config"
	"github.com/what-writers-like/backend/internal/infrastructure/database"
	"github.com/what-writers-like/backend/internal/repository/gorm"
	"github.com/what-writers-like/backend/internal/service"
)

func main() {
	fx.New(
		fx.Provide(
			config.NewConfig,
			database.NewDatabase,
			gorm.NewWriterRepository,
			gorm.NewWorkRepository,
			gorm.NewOpinionRepository,
			service.NewWriterService,
			service.NewWorkService,
			service.NewOpinionService,
			handler.NewWriterHandler,
			handler.NewWorkHandler,
			handler.NewOpinionHandler,
			handler.SetupRouter,
			NewHTTPServer,
		),
		fx.Invoke(RegisterLifecycle),
	).Run()
}

func NewHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}
}

func RegisterLifecycle(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(fmt.Sprintf("failed to start server: %v", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})
}
