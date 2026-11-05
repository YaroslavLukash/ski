package domain

type ItemType string

const (
	ItemTypeSki  ItemType = "SKI"
	ItemTypeBike ItemType = "BIKE"
)

type RentalItem struct {
	ID           string   `json:"id"`
	Type         ItemType `json:"type"`
	Description  string   `json:"description"`
	IsPro        bool     `json:"is_pro"`
	SerialNumber string   `json:"serial_number"`
	IsRented     bool     `json:"is_rented"`
	HourlyRate   float64  `json:"hourly_rate"`
}

type Ski struct {
	RentalItem
	LengthCM    int    `json:"length_cm"`
	BindingType string `json:"binding_type"`
}

type Bike struct {
	RentalItem
	FrameSize  string `json:"frame_size"`
	IsElectric bool   `json:"is_electric"`
}

type Customer struct {
	CustomerID     string `json:"customer_id"`
	IDNumber       string `json:"id_number"`
	IDType         string `json:"id_type"`
	DocumentBase64 string `json:"document_base64"`
}
