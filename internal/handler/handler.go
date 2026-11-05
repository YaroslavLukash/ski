package handler

import (
	"encoding/json"
	"net/http"
	"system-verleih/internal/domain"
	"system-verleih/internal/plugin"
	"system-verleih/internal/store"
	"time"
)

type Handler struct {
	store    *store.MemoryStore
	registry *plugin.Registry
}

func NewHandler(store *store.MemoryStore, registry *plugin.Registry) *Handler {
	return &Handler{
		store:    store,
		registry: registry,
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

func (h *Handler) Items(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items := h.store.GetAllItems()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

type CartItem struct {
	ID      string `json:"id"`
	Voucher string `json:"voucher"`
}

type CheckoutRequest struct {
	DurationHours float64    `json:"duration_hours"`
	CartItems     []CartItem `json:"cart_items"`
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.CartItems) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	duration := time.Duration(req.DurationHours * float64(time.Hour))
	var itemIDs []string
	var prices []float64

	for _, cItem := range req.CartItems {
		itemIDs = append(itemIDs, cItem.ID)
	}

	// Validate items are available (MarkRented will do this, but we need to calculate first)
	// So we'll fetch them, calculate, then try to MarkRented
	var items []*domain.RentalItem
	for _, id := range itemIDs {
		item, err := h.store.GetItem(id)
		if err != nil || item.IsRented {
			http.Error(w, "Item not available: "+id, http.StatusBadRequest)
			return
		}
		items = append(items, item)
	}

	for i, item := range items {
		p, ok := h.registry.Get(item.Type)
		if !ok {
			http.Error(w, "Plugin not found for type", http.StatusInternalServerError)
			return
		}

		price := p.CalculatePrice(item, duration)

		// Very simple voucher logic: if voucher == "HALF", cut price in half
		if req.CartItems[i].Voucher == "HALF" {
			price = price / 2.0
		}

		prices = append(prices, price)
	}

	// Promotional Rental Rule (3 + 1 Free) for full day (>= 24 hours)
	if req.DurationHours >= 24.0 && len(items) >= 4 {
		// Find the cheapest item among the first 4 (or all) and zero it out
		// To be precise: If a user rents 3 items, they can add a 4th for free.
		// We'll zero out the cheapest item in the cart.
		minPriceIdx := 0
		minPrice := prices[0]
		for i, price := range prices {
			if price < minPrice {
				minPrice = price
				minPriceIdx = i
			}
		}
		prices[minPriceIdx] = 0
	}

	totalPrice := 0.0
	for _, price := range prices {
		totalPrice += price
	}

	if err := h.store.MarkRented(itemIDs); err != nil {
		http.Error(w, "Error marking items as rented: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Checkout successful",
		"total_price": totalPrice,
	})
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
