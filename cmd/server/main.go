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
	registry.Register(&plugin.HelmetPlugin{})
	registry.Register(&plugin.BootsPlugin{})
	registry.Register(&plugin.GogglesPlugin{})

	// Seed Mock Data
	seedMockData(memStore)

	// 3. Initialize Handlers
	apiHandler := handler.NewHandler(memStore, registry)

	// 4. Setup Routes
	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("/api/v1/customers", apiHandler.Customers)
	mux.HandleFunc("/api/v1/customers/verify-id", apiHandler.VerifyID)
	mux.HandleFunc("/api/v1/stats", apiHandler.Stats)
	mux.HandleFunc("/api/v1/items", apiHandler.Items)
	mux.HandleFunc("/api/v1/checkout", apiHandler.Checkout)

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
	// Standard Ski
	ski1 := &domain.Ski{
		RentalItem: domain.RentalItem{
			ID:           "ski-1",
			Type:         domain.ItemTypeSki,
			Description:  "Standard Ski",
			IsPro:        false,
			SerialNumber: "SN-SKI-1001",
			IsRented:     false,
			HourlyRate:   6.0,
		},
		LengthCM:    170,
		BindingType: "Alpine",
	}
	s.AddItem(&ski1.RentalItem)

	// Pro Carving Ski
	ski2 := &domain.Ski{
		RentalItem: domain.RentalItem{
			ID:           "ski-2",
			Type:         domain.ItemTypeSki,
			Description:  "Pro Carving Ski",
			IsPro:        true,
			SerialNumber: "SN-SKI-1002",
			IsRented:     false,
			HourlyRate:   6.0,
		},
		LengthCM:    165,
		BindingType: "Performance",
	}
	s.AddItem(&ski2.RentalItem)

	// City Bike
	bike1 := &domain.Bike{
		RentalItem: domain.RentalItem{
			ID:           "bike-1",
			Type:         domain.ItemTypeBike,
			Description:  "City Bike",
			IsPro:        false,
			SerialNumber: "SN-BIKE-2001",
			IsRented:     false,
			HourlyRate:   9.0,
		},
		FrameSize:  "M",
		IsElectric: false,
	}
	s.AddItem(&bike1.RentalItem)

	// Pro Enduro Mountain Bike
	bike2 := &domain.Bike{
		RentalItem: domain.RentalItem{
			ID:           "bike-2",
			Type:         domain.ItemTypeBike,
			Description:  "Pro Enduro Mountain Bike",
			IsPro:        true,
			SerialNumber: "SN-BIKE-2002",
			IsRented:     false,
			HourlyRate:   9.0,
		},
		FrameSize:  "L",
		IsElectric: false,
	}
	s.AddItem(&bike2.RentalItem)

	helmet := &domain.Helmet{
		RentalItem: domain.RentalItem{
			ID:           "helmet-1",
			Type:         domain.ItemTypeHelmet,
			Description:  "Safety Helmet",
			IsPro:        false,
			SerialNumber: "SN-HELMET-3001",
			IsRented:     false,
			HourlyRate:   2.0,
		},
		Size: "M",
	}
	s.AddItem(&helmet.RentalItem)

	boots := &domain.Boots{
		RentalItem: domain.RentalItem{
			ID:           "boots-1",
			Type:         domain.ItemTypeBoots,
			Description:  "Ski Boots",
			IsPro:        false,
			SerialNumber: "SN-BOOTS-4001",
			IsRented:     false,
			HourlyRate:   4.0,
		},
		Size: 42.5,
	}
	s.AddItem(&boots.RentalItem)

	goggles := &domain.Goggles{
		RentalItem: domain.RentalItem{
			ID:           "goggles-1",
			Type:         domain.ItemTypeGoggles,
			Description:  "Snow Goggles",
			IsPro:        false,
			SerialNumber: "SN-GOGGLES-5001",
			IsRented:     false,
			HourlyRate:   1.5,
		},
	}
	s.AddItem(&goggles.RentalItem)
}
