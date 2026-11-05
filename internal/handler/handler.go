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
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var customer domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	if customer.CustomerID == "" || customer.IDNumber == "" || customer.IDType == "" || customer.DocumentBase64 == "" {
		http.Error(w, "Fehlende Pflichtfelder", http.StatusBadRequest)
		return
	}

	if err := h.store.AddCustomer(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kunden-ID erfolgreich verifiziert und gespeichert"})
}

func (h *Handler) Customers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customers := h.store.GetAllCustomers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
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
	ID            string  `json:"id"`
	Voucher       string  `json:"voucher"`
	DurationHours float64 `json:"duration_hours"`
}

type CheckoutRequest struct {
	CustomerID string     `json:"customer_id"`
	CartItems  []CartItem `json:"cart_items"`
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		return
	}

	if len(req.CartItems) == 0 {
		http.Error(w, "Warenkorb ist leer", http.StatusBadRequest)
		return
	}

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
			http.Error(w, "Artikel nicht verfügbar: "+id, http.StatusBadRequest)
			return
		}
		items = append(items, item)
	}

	fullDayCount := 0

	for i, item := range items {
		p, ok := h.registry.Get(item.Type)
		if !ok {
			http.Error(w, "Plugin nicht gefunden", http.StatusInternalServerError)
			return
		}

		cItem := req.CartItems[i]
		if cItem.DurationHours >= 24.0 {
			fullDayCount++
		}

		duration := time.Duration(cItem.DurationHours * float64(time.Hour))
		price := p.CalculatePrice(item, duration)

		// Scarcity Logic: if 2 or fewer available, multiply base price by 1.5
		if h.store.GetAvailableCountByType(item.Type) <= 2 {
			price *= 1.5
		}

		// Very simple voucher logic: if voucher == "HALF", cut price in half
		if cItem.Voucher == "HALF" {
			price = price / 2.0
		}

		prices = append(prices, price)
	}

	// Promotional Rental Rule (3 + 1 Free) for full day
	// If at least 3 items have duration >= 24h, and total items >= 4, the cheapest item is free
	if fullDayCount >= 3 && len(items) >= 4 {
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

	// Verified User 10% Discount
	if req.CustomerID != "" {
		if _, err := h.store.GetCustomer(req.CustomerID); err == nil {
			totalPrice = totalPrice * 0.9
		}
	}

	if err := h.store.MarkRented(itemIDs); err != nil {
		http.Error(w, "Fehler beim Aktualisieren der Artikel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Checkout erfolgreich",
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
