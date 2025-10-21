package http

import (
	"em-test/internal/domain"
	"em-test/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type SumSubscriptionsRequest struct {
    UserID       uuid.UUID `json:"user_id"`                
    ServiceName  *string   `json:"service_name,omitempty"` 
    StartDate    string    `json:"start_date"`             
    EndDate      string    `json:"end_date"`               
}

type Handler struct {
	subscriptionService *service.SubscriptionService
	log                 *slog.Logger
}

func NewHandler(subscriptionService *service.SubscriptionService, log *slog.Logger) *Handler {
	return &Handler{subscriptionService: subscriptionService, log: log}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /subscriptions/{user_id}", h.listSubscriptions)
	mux.HandleFunc("GET /subscriptions/{user_id}/{service_name}", h.getSubscription)
	mux.HandleFunc("POST /subscriptions/total-price", h.sumSubscriptionsPrice)
	mux.HandleFunc("POST /subscriptions", h.createSubscription)
	mux.HandleFunc("PUT /subscriptions/{user_id}/{service_name}", h.updateSubscription)
	mux.HandleFunc("DELETE /subscriptions/{user_id}/{service_name}", h.deleteSubscription)
}

func (h *Handler) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.log.Error("failed to list subscriptions, invalid user_id", "user_id", userIDStr, "error", err)
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	subs, err := h.subscriptionService.ListSubscriptions(r.Context(), userID)
	if err != nil {
		h.log.Error("failed to list subscriptions", "user_id", userIDStr, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

func (h *Handler) getSubscription(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	serviceName := r.PathValue("service_name")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.log.Error("failed to get subscription, invalid user_id", "user_id", userIDStr, "error", err)
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	sub, err := h.subscriptionService.GetSubscription(r.Context(), userID, serviceName)
	if err != nil {
		h.log.Error("failed to get subscription", "user_id", userIDStr, "service_name", serviceName, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	var sub domain.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		h.log.Error("failed to create subscription, invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.subscriptionService.CreateSubscription(r.Context(), &sub); err != nil {
		h.log.Error("failed to create subscription", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	serviceName := r.PathValue("service_name")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.log.Error("failed to update subscription, invalid user_id", "user_id", userIDStr, "error", err)
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	var sub domain.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		h.log.Error("failed to update subscription, invalid request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure the path and body match
	if sub.UserID != userID || sub.ServiceName != serviceName {
		h.log.Error("failed to update subscription, user_id or service_name mismatch", "user_id", userID, "service_name", serviceName)
		http.Error(w, "user_id or service_name mismatch", http.StatusBadRequest)
		return
	}

	if err := h.subscriptionService.UpdateSubscription(r.Context(), &sub); err != nil {
		h.log.Error("failed to update subscription", "user_id", userIDStr, "service_name", serviceName, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	serviceName := r.PathValue("service_name")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.log.Error("failed to delete subscription, invalid user_id", "user_id", userIDStr, "error", err)
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.subscriptionService.DeleteSubscription(r.Context(), userID, serviceName); err != nil {
		h.log.Error("failed to delete subscription", "user_id", userIDStr, "service_name", serviceName, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) sumSubscriptionsPrice(w http.ResponseWriter, r *http.Request) {
    var req SumSubscriptionsRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.log.Error("invalid request body", "err", err)
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if req.UserID == uuid.Nil {
        h.log.Error("user_id is required")
        http.Error(w, "user_id is required", http.StatusBadRequest)
        return
    }

    if req.StartDate == "" || req.EndDate == "" {
        h.log.Error("start_date and end_date are required")
        http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
        return
    }

    totalPrice, err := h.subscriptionService.SumSubscriptionsPrice(
        r.Context(),
        req.UserID,
        req.ServiceName,
        req.StartDate,
        req.EndDate,
    )
    if err != nil {
        h.log.Error("failed to calculate total price", "err", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(map[string]int{"total_price": totalPrice}); err != nil {
        h.log.Error("failed to encode response", "err", err)
    }
}