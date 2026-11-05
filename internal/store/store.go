package store

import (
	"fmt"
	"sync"
	"system-verleih/internal/domain"
)

type MemoryStore struct {
	mu        sync.RWMutex
	items     map[string]*domain.RentalItem
	customers map[string]*domain.Customer
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items:     make(map[string]*domain.RentalItem),
		customers: make(map[string]*domain.Customer),
	}
}

func (s *MemoryStore) AddItem(item *domain.RentalItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[item.ID]; exists {
		return fmt.Errorf("item with id %s already exists", item.ID)
	}
	s.items[item.ID] = item
	return nil
}

func (s *MemoryStore) MarkRented(itemIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if all exist and are not already rented
	for _, id := range itemIDs {
		item, exists := s.items[id]
		if !exists {
			return fmt.Errorf("item %s not found", id)
		}
		if item.IsRented {
			return fmt.Errorf("item %s is already rented", id)
		}
	}

	// Mark as rented
	for _, id := range itemIDs {
		s.items[id].IsRented = true
	}

	return nil
}

func (s *MemoryStore) GetItem(id string) (*domain.RentalItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, exists := s.items[id]
	if !exists {
		return nil, fmt.Errorf("item %s not found", id)
	}
	return item, nil
}

func (s *MemoryStore) GetAllItems() []*domain.RentalItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var allItems []*domain.RentalItem
	for _, item := range s.items {
		allItems = append(allItems, item)
	}
	return allItems
}

func (s *MemoryStore) GetCustomerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.customers)
}

func (s *MemoryStore) AddCustomer(customer *domain.Customer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.customers[customer.CustomerID]; exists {
		return fmt.Errorf("customer with id %s already exists", customer.CustomerID)
	}
	s.customers[customer.CustomerID] = customer
	return nil
}
