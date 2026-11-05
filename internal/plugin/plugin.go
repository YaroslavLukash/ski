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
	return hours * item.HourlyRate * p.Multiplier
}

func (p *SkiPlugin) ValidateItem(item *domain.RentalItem) error {
	if item.Type != domain.ItemTypeSki {
		return fmt.Errorf("invalid item type for ski plugin")
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
	return hours * item.HourlyRate * p.Multiplier
}

func (p *BikePlugin) ValidateItem(item *domain.RentalItem) error {
	if item.Type != domain.ItemTypeBike {
		return fmt.Errorf("invalid item type for bike plugin")
	}
	return nil
}
