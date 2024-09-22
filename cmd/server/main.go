// cmd/main.go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/GoSeoTaxi/check_freedomnet/internal/config"
	"github.com/GoSeoTaxi/check_freedomnet/internal/handler"
	"github.com/GoSeoTaxi/check_freedomnet/internal/repository"
	"github.com/GoSeoTaxi/check_freedomnet/internal/service"
)

func main() {
	fx.New(
		// Провайдим зависимости
		fx.Provide(
			config.LoadConfig, // Загрузка конфигурации
			zap.NewProduction, // Логгер
			func(cfg *config.Config) *repository.FreedomNetRepo {
				return repository.NewFreedomNetRepo(cfg.Servers)
			},
			func(repo *repository.FreedomNetRepo, logger *zap.Logger, cfg *config.Config) *service.FreedomNetService {
				return service.NewFreedomNetService(repo, cfg.MaxRetries, logger)
			},
			func(svc *service.FreedomNetService, logger *zap.Logger) *handler.FreedomNetHandler {
				return handler.NewFreedomNetHandler(svc, logger)
			},
		),
		fx.Invoke(registerRoutes),
	).Run()
}

func registerRoutes(h *handler.FreedomNetHandler, cfg *config.Config) {
	r := chi.NewRouter()
	r.Get("/get_me_new_freedomnet", h.GetFreedomNetHandler)
	_ = http.ListenAndServe(":"+cfg.Port, r)
}
