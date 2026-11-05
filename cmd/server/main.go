package main

import (
	"fmt"
	"log"
	"net/http"
	"system-verleih/internal/domain"
	"system-verleih/internal/frontend"
	"system-verleih/internal/handler"
	"system-verleih/internal/plugin"
	"system-verleih/internal/store"
)

func main() {
	// 1. Initialize Store
	memStore := store.NewMemoryStore()

	// 2. Initialize Plugin Registry
	registry := plugin.NewRegistry()
	registry.Register(&plugin.SkiPlugin{Multiplier: 1.5})
	registry.Register(&plugin.BikePlugin{Multiplier: 1.2})

	// Seed Mock Data
	seedMockData(memStore)

	// 3. Initialize Handlers
	apiHandler := handler.NewHandler(memStore)

	// 4. Setup Routes
	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("/api/v1/customers/verify-id", apiHandler.VerifyID)
	mux.HandleFunc("/api/v1/stats", apiHandler.Stats)

	// Frontend Route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(frontend.HTMLTemplate))
	})

	// 5. Start Server
	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func seedMockData(s *store.MemoryStore) {
	// Seed Ski
	ski := &domain.Ski{
		RentalItem: domain.RentalItem{
			ID:           "ski-1",
			Type:         domain.ItemTypeSki,
			SerialNumber: "SN-SKI-1001",
			IsRented:     false,
			HourlyRate:   10.0,
		},
		LengthCM:    170,
		BindingType: "Alpine",
	}
	s.AddItem(&ski.RentalItem)

	// Seed Bike
	bike := &domain.Bike{
		RentalItem: domain.RentalItem{
			ID:           "bike-1",
			Type:         domain.ItemTypeBike,
			SerialNumber: "SN-BIKE-2001",
			IsRented:     false,
			HourlyRate:   15.0,
		},
		FrameSize:  "L",
		IsElectric: false,
	}
	s.AddItem(&bike.RentalItem)
}
