package HandlerOrders

import (
	"arch/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

func (o *OrderHandler) HandlerGetListOrder(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	filter := repository.Orderfilter{
		OwnerUuid:    query.Get("OwnerUuid"),
		Branchuuid:   query.Get("BranchUuid"),
		SellerUuid:   query.Get("SellerUuid"),
		InvoicesUuid: query.Get("InvoicesUuid"),
	}

	filter.Limit, _ = strconv.Atoi(query.Get("Limit"))
	filter.Offset, _ = strconv.Atoi(query.Get("Offset"))

	orders, err := o.OrderServices.GetListOrdersServices(r.Context(), filter)
	if err != nil {
		http.Error(w, "failed to get orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
