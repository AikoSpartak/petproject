package HandlerOrders

import (
	"arch/internal/handler/ResponseStruct"
	"arch/internal/middleware"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (o *OrderHandler) HandlerGetOrder(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "id")
	orderUuid, err := uuid.Parse(idStr)
	if err != nil {
		slog.Error("error receiving the ID",
			"error", err.Error())
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	userId, userRole, ok := middleware.GetUserRole(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	order, err := o.OrderServices.GetOrderServices(r.Context(), orderUuid)
	if err != nil {
		slog.Error("error receiving order", "error", err.Error())
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	userIdUuid, err := uuid.Parse(userId)
	if err != nil {
		slog.Error("invalid user id in token", "error", err.Error())
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	if userRole != "admin" && order.OwnerUuid != userIdUuid {
		http.Error(w, "access is denied", http.StatusForbidden)
		return
	}

	response := ResponseStruct.GetOrderResponse{
		OrderUuid:     order.Uuid,
		OwnerUuid:     order.OwnerUuid,
		EquipmentUuid: order.EquipmentUuid,
		BranchUuid:    order.BranchUuid,
		SellerUuid:    order.SellerUuid,
		InvoicesUuid:  order.InvoicesUuid,
		TotalAmount:   order.TotalAmount,
		IsPaid:        order.IsPaid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
