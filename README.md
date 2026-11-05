# Software Design Document (SDD): Extensible Rental Framework (`system-verleih`)

\---

## 1\. Executive Summary \& Goals

The objective of **`system-verleih`** is to build, modular prototype for rental systems (starting with Ski and Bike rentals) in Go. It builds directly on the core principles specified in `CLAUDE.md`: rapid iteration, prototype clarity, in-memory thread-safe execution, and a strict plugin architecture.

### Key Requirements \& Adaptations

* **SDD Implementation**: Strict specification-driven development aligning directly with `CLAUDE.md` and repository guidelines.



* **Modular Domain (`RentalItem`)**: Abstract entity handling base properties extended by concrete implementations (`Ski`, `Bike`).



* **Flexible Hourly Pricing**: Dynamic rate calculation based on duration rather than fixed daily rates.



* **ID Verification Support**: Ability to upload and parse customer ID details (ID number and ID type).



* **Plugin Architecture**: Clean interface abstraction allowing seamless swapping or adding of rental inventory types.



* **Public Repository Submission**: Fully version-controlled via Git/GitHub.



\---

## 2\. System Architecture \& Component Design

The backend uses a standard clean Go prototype layout:

```text
cmd/server/main.go      --> HTTP Server setup \& router initialization
internal/domain/        --> Core business entities (RentalItem, Ski, Bike, Customer)
internal/store/         --> In-memory thread-safe data store (Mutex + Maps)
internal/handler/       --> HTTP handlers mapping JSON requests to domain calls
internal/plugin/        --> Plugin registry \& concrete rental strategy implementations

```

### Domain Data Model (`RentalItem` Hierarchy)

The domain relies on polymorphism via Go interfaces and embedding:

```go
package domain

import "time"

type ItemType string

const (
    ItemTypeSki  ItemType = "SKI"
    ItemTypeBike ItemType = "BIKE"
)

type RentalItem struct {
    ID           string    `json:"id"`
    Type         ItemType  `json:"type"`
    SerialNumber string    `json:"serial\_number"`
    IsRented     bool      `json:"is\_rented"`
    HourlyRate   float64   `json:"hourly\_rate"`
}

type Ski struct {
    RentalItem
    LengthCM    int    `json:"length\_cm"`
    BindingType string `json:"binding\_type"`
}

type Bike struct {
    RentalItem
    FrameSize string `json:"frame\_size"`
    IsElectric bool  `json:"is\_electric"`
}

```

\---

## 3\. Plugin Architecture Framework

To achieve full decoupling between the core rental engine and specific rental items (Ski vs. Bike), a plugin registry interface is established.

### Interface Definition

```go
package plugin

import (
    "time"
    "system-verleih/internal/domain"
)

type RentalPlugin interface {
    GetType() domain.ItemType
    CalculatePrice(item \*domain.RentalItem, duration time.Duration) float64
    ValidateItem(item \*domain.RentalItem) error
}

type Registry struct {
    plugins map\[domain.ItemType]RentalPlugin
}

func NewRegistry() \*Registry {
    return \&Registry{plugins: make(map\[domain.ItemType]RentalPlugin)}
}

func (r \*Registry) Register(p RentalPlugin) {
    r.plugins\[p.GetType()] = p
}

func (r \*Registry) Get(t domain.ItemType) (RentalPlugin, bool) {
    p, ok := r.plugins\[t]
    return p, ok
}

```

\---

## 4\. Feature Specifications \& API Endpoints

### 4.1 Hourly Dynamic Pricing Calculation

Instead of standard flat daily rates, rentals are calculated dynamically per hour or fraction of an hour.

$$\\text{Total Price} = \\lceil \\text{Duration in Hours} \\rceil \\times \\text{Base Hourly Rate} \\times \\text{Plugin Multiplier}$$

### 4.2 Customer ID Verification Endpoint

Allows uploading identity verification document metadata alongside the rental record.

* **Endpoint**: `POST /api/v1/customers/verify-id`
* **Payload**:

```json
{
  "customer\_id": "cust\_98234",
  "id\_number": "A-83920192",
  "id\_type": "PASSPORT",
  "document\_base64": "data:image/png;base64,..."
}

```

\---

## 5\. In-Memory Thread-Safe Data Store

Per `CLAUDE.md` guidelines, no external databases (PostgreSQL, SQLite, ORMs) are used. Data consistency across concurrent HTTP handlers is guaranteed via `sync.RWMutex`.

```go
package store

import (
    "sync"
    "fmt"
    "system-verleih/internal/domain"
)

type MemoryStore struct {
    mu    sync.RWMutex
    items map\[string]\*domain.RentalItem
}

func NewMemoryStore() \*MemoryStore {
    return \&MemoryStore{
        items: make(map\[string]\*domain.RentalItem),
    }
}

func (s \*MemoryStore) AddItem(item \*domain.RentalItem) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    if \_, exists := s.items\[item.ID]; exists {
        return fmt.Errorf("item with id %s already exists", item.ID)
    }
    s.items\[item.ID] = item
    return nil
}

func (s \*MemoryStore) GetItem(id string) (\*domain.RentalItem, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    item, exists := s.items\[id]
    if !exists {
        return nil, fmt.Errorf("item %s not found", id)
    }
    return item, nil
}

```

\---

## 6\. Execution Commands \& Guidelines

Adhering strictly to repository command conventions:

|Step|Command|Description|
|-|-|-|
|**Run Server**|`go run ./cmd/server`|Starts the Go backend HTTP listener on port 8080|

|
| **Run Tests** | `go test ./...` | Runs unit tests across all domain and plugin packages

|
| **Code Formatting** | `go fmt ./...` | Standardizes all Go codebase files

|
| **Public Submission** | `git push origin main` | Pushes changes to the configured public GitHub repository

Dashboard Key UI Components

UI Section

Core Purpose

HTML/JS Implementation

Top Navigation Bar

Switch between view scopes (Dashboard, Items, Customers, Rentals)

<nav class="navbar"> with flexbox/grid routing tabs

Metric Summary Cards

Quick counts of Available, Rented, Total Items, Registered Users

Dynamic <div> metric cards updated via GET /api/v1/stats

Quick Action Buttons

Rapid entry points ("New Rental", "Manage Items", "Manage Persons")

Action buttons routing directly to modal dialogs or sub-pages

|

