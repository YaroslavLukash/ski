# System Verleih

An extensible rental framework prototype for skis and bikes, written in Go.

## Features
- Modular Plugin Architecture for different rental items (Ski, Bike).
- Dynamic hourly pricing calculations.
- Thread-safe in-memory data store using `sync.RWMutex`.
- HTTP REST API for statistics and customer ID verification.
- Integrated HTML/JS frontend dashboard.

## Running the Project

### Start the Server

```bash
go run ./cmd/server
```

This will start the HTTP server on `http://localhost:8080`.

### Endpoints

- **Frontend UI**: `GET /`
- **Dashboard Stats**: `GET /api/v1/stats`
- **List Items**: `GET /api/v1/items`
- **Checkout**: `POST /api/v1/checkout`
  - Body (JSON): `{"duration_hours": 24, "cart_items": [{"id": "ski-1", "voucher": "HALF"}]}`
- **Verify Customer ID**: `POST /api/v1/customers/verify-id`
  - Body (JSON): `{"customer_id": "...", "id_number": "...", "id_type": "...", "document_base64": "..."}`

### Pricing and Promotions

- **Base Pricing**: Skis are €6.00/hr, Bikes are €9.00/hr.
- **Pro Multipliers**: Pro items apply a multiplier (e.g., 1.5x) to the base rate.
- **3 + 1 Free Promotion**: Rent 3 items for a full day (>= 24 hours), and a 4th item is automatically free (zeroes out the cheapest item).

## Testing

To run the tests:

```bash
go test ./...
```
