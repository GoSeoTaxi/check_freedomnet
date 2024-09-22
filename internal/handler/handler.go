package handler

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/GoSeoTaxi/check_freedomnet/internal/service"
)

type FreedomNetHandler struct {
	Service *service.FreedomNetService
	Logger  *zap.Logger
}

func NewFreedomNetHandler(service *service.FreedomNetService, logger *zap.Logger) *FreedomNetHandler {
	return &FreedomNetHandler{
		Service: service,
		Logger:  logger,
	}
}

func (h *FreedomNetHandler) GetFreedomNetHandler(w http.ResponseWriter, r *http.Request) {
	_ = r
	result, err := h.Service.GetFreedomNet()
	if err != nil {
		h.Logger.Error("Failed to fetch data", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(result + "\n"))
}
