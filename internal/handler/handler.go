package handler

import (
	"encoding/json"
	"net/http"
	"system-verleih/internal/domain"
	"system-verleih/internal/store"
)

type Handler struct {
	store *store.MemoryStore
}

func NewHandler(store *store.MemoryStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) VerifyID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var customer domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if customer.CustomerID == "" || customer.IDNumber == "" || customer.IDType == "" || customer.DocumentBase64 == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.store.AddCustomer(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer ID verified and saved"})
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items := h.store.GetAllItems()

	totalItems := len(items)
	availableItems := 0
	rentedItems := 0

	for _, item := range items {
		if item.IsRented {
			rentedItems++
		} else {
			availableItems++
		}
	}

	// Assuming a simpler way to get registered users since we don't have a GetAllCustomers
	// Or we can just add a GetAllCustomers

	registeredUsers := h.store.GetCustomerCount()

	stats := map[string]int{
		"total_items":      totalItems,
		"available_items":  availableItems,
		"rented_items":     rentedItems,
		"registered_users": registeredUsers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
