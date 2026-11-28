package HandlerOrders

import (
	"arch/internal/handler/RequestStruct"
	"arch/internal/handler/ResponseStruct"
	"arch/internal/middleware"
	"arch/internal/services"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type OrderHandler struct {
	OrderServices services.OrdersServices
}

func (o *OrderHandler) RegisterRouts(r *chi.Mux) {

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Post("/orders", o.CreateOrderHandler)
		r.Get("/orders", o.HandlerGetListOrder)
		r.Get("/orders/{id}", o.HandlerGetOrder)
		r.Delete("/orders/{id}", o.DeleteUserHandler)
	})

}

func (o *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {

	var req RequestStruct.OrdersRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		slog.Error("error parsing",
			"error", err.Error())
		http.Error(w, "invalid register", http.StatusBadRequest)
		return
	}

	orderUuidStr, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	orderUuid, err := uuid.Parse(orderUuidStr)
	if err != nil {
		http.Error(w, "failed uuid", http.StatusBadRequest)
		return
	}

	order, err := o.OrderServices.RegisterOrderServices(r.Context(),
		orderUuid,
		req.EquipmentUuid,
		req.BranchUuid,
		req.SellerUuid,
		req.InvoicesUuid,
		req.TotalAmount,
		req.IsPaid,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := ResponseStruct.OrdersResponse{
		OwnerUuid:     orderUuid,
		EquipmentUuid: order.EquipmentUuid,
		BranchUuid:    order.BranchUuid,
		SellerUuid:    order.SellerUuid,
		InvoicesUuid:  order.InvoicesUuid,
		TotalAmount:   order.TotalAmount,
		IsPaid:        order.IsPaid,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response register", http.StatusInternalServerError)
		return
	}
}
