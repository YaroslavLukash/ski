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
	svgStandardSki := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#e0f7fa"/><path d="M 10 90 L 90 10" stroke="#00bcd4" stroke-width="6" stroke-linecap="round"/><text x="10" y="20" font-size="12" fill="#0097a7">Standard Ski</text></svg>`
	svgProSki := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#ffebee"/><path d="M 10 90 Q 50 50 90 10" stroke="#f44336" stroke-width="8" fill="none" stroke-linecap="round"/><text x="10" y="20" font-size="12" fill="#c62828">Pro Ski</text></svg>`
	svgFreestyleSki := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#fff3e0"/><path d="M 20 80 Q 50 20 80 80" stroke="#ff9800" stroke-width="6" fill="none" stroke-linecap="round"/><text x="10" y="20" font-size="12" fill="#e65100">Freestyle</text></svg>`

	svgCityBike := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#f1f8e9"/><circle cx="30" cy="70" r="15" stroke="#8bc34a" stroke-width="4" fill="none"/><circle cx="70" cy="70" r="15" stroke="#8bc34a" stroke-width="4" fill="none"/><path d="M 30 70 L 50 45 L 70 70 M 50 45 L 25 45" stroke="#8bc34a" stroke-width="4" fill="none"/><text x="10" y="20" font-size="12" fill="#33691e">City Bike</text></svg>`
	svgEnduroBike := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#eceff1"/><circle cx="25" cy="75" r="15" stroke="#607d8b" stroke-width="6" fill="none"/><circle cx="75" cy="75" r="15" stroke="#607d8b" stroke-width="6" fill="none"/><path d="M 25 75 L 50 40 L 75 75 M 50 40 L 20 40" stroke="#607d8b" stroke-width="6" fill="none"/><text x="10" y="20" font-size="12" fill="#263238">Enduro</text></svg>`
	svgGravelBike := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#fff8e1"/><circle cx="30" cy="70" r="18" stroke="#ffc107" stroke-width="3" fill="none" stroke-dasharray="2,2"/><circle cx="70" cy="70" r="18" stroke="#ffc107" stroke-width="3" fill="none" stroke-dasharray="2,2"/><path d="M 30 70 L 50 45 L 70 70" stroke="#ffc107" stroke-width="4" fill="none"/><text x="10" y="20" font-size="12" fill="#ff8f00">Gravel</text></svg>`

	svgHelmet := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#e8eaf6"/><path d="M 20 70 Q 50 10 80 70 Z" fill="#3f51b5"/><text x="25" y="90" font-size="12" fill="#1a237e">Helmet</text></svg>`
	svgBoots := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#fce4ec"/><path d="M 30 20 L 30 80 L 70 80 L 70 60 L 50 60 L 50 20 Z" fill="#e91e63"/><text x="35" y="95" font-size="12" fill="#880e4f">Boots</text></svg>`
	svgGoggles := `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg"><rect width="100" height="100" fill="#e0f2f1"/><ellipse cx="35" cy="50" rx="20" ry="15" fill="#009688"/><ellipse cx="65" cy="50" rx="20" ry="15" fill="#009688"/><path d="M 15 50 Q 50 20 85 50" stroke="#004d40" stroke-width="4" fill="none"/><text x="25" y="85" font-size="12" fill="#004d40">Goggles</text></svg>`

	// Seed multiple quantities of skis
	for i := 1; i <= 3; i++ {
		item := &domain.Ski{
			RentalItem: domain.RentalItem{
				ID:           fmt.Sprintf("%d", i),
				Type:         domain.ItemTypeSki,
				Description:  "Standard Ski",
				IsPro:        false,
				SerialNumber: fmt.Sprintf("SN-SKI-100%d", i),
				IsRented:     false,
				HourlyRate:   6.0,
				ImageSVG:     svgStandardSki,
			},
			LengthCM:    170,
			BindingType: "Alpine",
		}
		s.AddItem(&item.RentalItem)
	}

	for i := 4; i <= 5; i++ {
		item := &domain.Ski{
			RentalItem: domain.RentalItem{
				ID:           fmt.Sprintf("%d", i),
				Type:         domain.ItemTypeSki,
				Description:  "Pro Carving Ski",
				IsPro:        true,
				SerialNumber: fmt.Sprintf("SN-SKI-200%d", i),
				IsRented:     false,
				HourlyRate:   6.0,
				ImageSVG:     svgProSki,
			},
			LengthCM:    165,
			BindingType: "Performance",
		}
		s.AddItem(&item.RentalItem)
	}

	item6 := &domain.Ski{
		RentalItem: domain.RentalItem{
			ID:           "6",
			Type:         domain.ItemTypeSki,
			Description:  "Freestyle Ski",
			IsPro:        false,
			SerialNumber: "SN-SKI-3001",
			IsRented:     false,
			HourlyRate:   6.0,
			ImageSVG:     svgFreestyleSki,
		},
		LengthCM:    160,
		BindingType: "Freestyle",
	}
	s.AddItem(&item6.RentalItem)

	// Seed multiple quantities of bikes
	for i := 7; i <= 8; i++ {
		item := &domain.Bike{
			RentalItem: domain.RentalItem{
				ID:           fmt.Sprintf("%d", i),
				Type:         domain.ItemTypeBike,
				Description:  "City Bike",
				IsPro:        false,
				SerialNumber: fmt.Sprintf("SN-BIKE-100%d", i),
				IsRented:     false,
				HourlyRate:   9.0,
				ImageSVG:     svgCityBike,
			},
			FrameSize:  "M",
			IsElectric: false,
		}
		s.AddItem(&item.RentalItem)
	}

	item9 := &domain.Bike{
		RentalItem: domain.RentalItem{
			ID:           "9",
			Type:         domain.ItemTypeBike,
			Description:  "Pro Enduro Mountain Bike",
			IsPro:        true,
			SerialNumber: "SN-BIKE-2001",
			IsRented:     false,
			HourlyRate:   9.0,
			ImageSVG:     svgEnduroBike,
		},
		FrameSize:  "L",
		IsElectric: false,
	}
	s.AddItem(&item9.RentalItem)

	item10 := &domain.Bike{
		RentalItem: domain.RentalItem{
			ID:           "10",
			Type:         domain.ItemTypeBike,
			Description:  "Gravel Bike",
			IsPro:        false,
			SerialNumber: "SN-BIKE-3001",
			IsRented:     false,
			HourlyRate:   9.0,
			ImageSVG:     svgGravelBike,
		},
		FrameSize:  "L",
		IsElectric: false,
	}
	s.AddItem(&item10.RentalItem)

	// Seed Helmets
	for i := 11; i <= 12; i++ {
		item := &domain.Helmet{
			RentalItem: domain.RentalItem{
				ID:           fmt.Sprintf("%d", i),
				Type:         domain.ItemTypeHelmet,
				Description:  "Safety Helmet",
				IsPro:        false,
				SerialNumber: fmt.Sprintf("SN-HELMET-%d", i),
				IsRented:     false,
				HourlyRate:   2.0,
				ImageSVG:     svgHelmet,
			},
			Size: "M",
		}
		s.AddItem(&item.RentalItem)
	}

	// Seed Boots
	item13 := &domain.Boots{
		RentalItem: domain.RentalItem{
			ID:           "13",
			Type:         domain.ItemTypeBoots,
			Description:  "Ski Boots",
			IsPro:        false,
			SerialNumber: "SN-BOOTS-4001",
			IsRented:     false,
			HourlyRate:   4.0,
			ImageSVG:     svgBoots,
		},
		Size: 42.5,
	}
	s.AddItem(&item13.RentalItem)

	// Seed Goggles
	item14 := &domain.Goggles{
		RentalItem: domain.RentalItem{
			ID:           "14",
			Type:         domain.ItemTypeGoggles,
			Description:  "Snow Goggles",
			IsPro:        false,
			SerialNumber: "SN-GOGGLES-5001",
			IsRented:     false,
			HourlyRate:   1.5,
			ImageSVG:     svgGoggles,
		},
	}
	s.AddItem(&item14.RentalItem)
}
