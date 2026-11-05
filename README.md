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
- **Verify Customer ID**: `POST /api/v1/customers/verify-id`
  - Body (JSON): `{"customer_id": "...", "id_number": "...", "id_type": "...", "document_base64": "..."}`

## Testing

To run the tests:

```bash
go test ./...
```
