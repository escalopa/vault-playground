package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/escalopa/vault-playground/internal/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	server        *mux.Router
	ordersService service.OrderService

	loggger *zap.Logger
}

func New(ordersService service.OrderService) *Handler {
	return &Handler{
		ordersService: ordersService,
		loggger:       zap.L().Named("handler"),
	}
}

func (h *Handler) Start(port string) error {
	h.server = mux.NewRouter()

	h.server.HandleFunc("/order/{orderId}", h.GetOrderItems).Methods(http.MethodGet)

	return http.ListenAndServe(port, h.server)
}

func (h *Handler) GetOrderItems(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(mux.Vars(r)["orderId"])
	if err != nil || orderId < 1 {
		h.loggger.Error("invalid order ID")

		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	items, err := h.ordersService.GetOrderItems(r.Context(), orderId)
	if err != nil {
		h.loggger.Error("failed to get order", zap.Error(err))

		http.Error(w, "failed to get order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		h.loggger.Error("failed to encode order", zap.Error(err))
	}
}
