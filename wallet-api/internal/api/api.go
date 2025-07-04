package api

import (
	"context"
	"encoding/json"
	"net/http"

	"wallet-api/internal/domain"
	"wallet-api/internal/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type WalletAPI struct {
	service *service.WalletService
}

func NewWalletAPI(service *service.WalletService) *WalletAPI {
	return &WalletAPI{service: service}
}

func (a *WalletAPI) Routes() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Post("/api/v1/wallet", a.processOperation)
	r.Get("/api/v1/wallets/{walletId}", a.getWalletBalance)

	return r
}

func (a *WalletAPI) processOperation(w http.ResponseWriter, r *http.Request) {
	var req domain.OperationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.service.ProcessOperation(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *WalletAPI) getWalletBalance(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "walletId")
	id, err := domain.ParseUUID(walletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wallet, err := a.service.GetWalletByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}
