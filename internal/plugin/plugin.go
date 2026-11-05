package plugin

import (
	"fmt"
	"math"
	"system-verleih/internal/domain"
	"time"
)

type RentalPlugin interface {
	GetType() domain.ItemType
	CalculatePrice(item *domain.RentalItem, duration time.Duration) float64
	ValidateItem(item *domain.RentalItem) error
}

type Registry struct {
	plugins map[domain.ItemType]RentalPlugin
}

func NewRegistry() *Registry {
	return &Registry{plugins: make(map[domain.ItemType]RentalPlugin)}
}

func (r *Registry) Register(p RentalPlugin) {
	r.plugins[p.GetType()] = p
}

func (r *Registry) Get(t domain.ItemType) (RentalPlugin, bool) {
	p, ok := r.plugins[t]
	return p, ok
}

type SkiPlugin struct {
	Multiplier float64
}

func (p *SkiPlugin) GetType() domain.ItemType {
	return domain.ItemTypeSki
}

func (p *SkiPlugin) CalculatePrice(item *domain.RentalItem, duration time.Duration) float64 {
	hours := math.Ceil(duration.Hours())
	multiplier := 1.0
	if item.IsPro {
		multiplier = p.Multiplier
	}
	return hours * item.HourlyRate * multiplier
}

func (p *SkiPlugin) ValidateItem(item *domain.RentalItem) error {
	if item.Type != domain.ItemTypeSki {
		return fmt.Errorf("invalid item type for ski plugin")
	}
	return nil
}

type HelmetPlugin struct{}

func (p *HelmetPlugin) GetType() domain.ItemType {
	return domain.ItemTypeHelmet
}

func (p *HelmetPlugin) CalculatePrice(item *domain.RentalItem, duration time.Duration) float64 {
	hours := math.Ceil(duration.Hours())
	return hours * item.HourlyRate
}

func (p *HelmetPlugin) ValidateItem(item *domain.RentalItem) error {
	if item.Type != domain.ItemTypeHelmet {
		return fmt.Errorf("invalid item type for helmet plugin")
	}
	return nil
}

type BootsPlugin struct{}

func (p *BootsPlugin) GetType() domain.ItemType {
	return domain.ItemTypeBoots
}

func (p *BootsPlugin) CalculatePrice(item *domain.RentalItem, duration time.Duration) float64 {
	hours := math.Ceil(duration.Hours())
	return hours * item.HourlyRate
}

func (p *BootsPlugin) ValidateItem(item *domain.RentalItem) error {
	if item.Type != domain.ItemTypeBoots {
		return fmt.Errorf("invalid item type for boots plugin")
	}
	return nil
}

type GogglesPlugin struct{}

func (p *GogglesPlugin) GetType() domain.ItemType {
	return domain.ItemTypeGoggles
}

func (p *GogglesPlugin) CalculatePrice(item *domain.RentalItem, duration time.Duration) float64 {
	hours := math.Ceil(duration.Hours())
	return hours * item.HourlyRate
}

func (p *GogglesPlugin) ValidateItem(item *domain.RentalItem) error {
	if item.Type != domain.ItemTypeGoggles {
		return fmt.Errorf("invalid item type for goggles plugin")
	}
	return nil
}

type BikePlugin struct {
	Multiplier float64
}

func (p *BikePlugin) GetType() domain.ItemType {
	return domain.ItemTypeBike
}

func (p *BikePlugin) CalculatePrice(item *domain.RentalItem, duration time.Duration) float64 {
	hours := math.Ceil(duration.Hours())
	multiplier := 1.0
	if item.IsPro {
		multiplier = p.Multiplier
	}
	return hours * item.HourlyRate * multiplier
}

func (p *BikePlugin) ValidateItem(item *domain.RentalItem) error {
	if item.Type != domain.ItemTypeBike {
		return fmt.Errorf("invalid item type for bike plugin")
	}
	return nil
}
