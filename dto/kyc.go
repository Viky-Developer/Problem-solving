package dto

import "github.com/google/uuid"

// ----------------------------------------------------------------
// ResponseKyc defines the structure for sending the KYC response.
// This structure is used when sending back the details of a KYC
// entity after an operation (e.g., update, fetch).
// ----------------------------------------------------------------
type ResponseKyc struct {
	Name         string    `json:"name"`
	PanNumber    string    `json:"panNumber"`    //PAN (Permanent Account Number)
	AadharNumber string    `json:"aadharNumber"` //Aadhar Number (Indian National ID)
	MerchantId   uuid.UUID `json:"merchant_id"`
}

// -------------------------------------------------------------
// RequestBodyKyc defines the structure for receiving KYC data
// in the request body during operations like update.
// -------------------------------------------------------------
type RequestBodyKyc struct {
	PanNumber    string `json:"panNumber,omitempty"`
	AadharNumber string `json:"aadharNumber,omitempty"`
}

// ----------------------------------------------------------------
// KYCDetails defines the structure for the KYC data to be stored
// or manipulated in the database.
// ----------------------------------------------------------------
type KYCDetails struct {
	Name       string    `json:"name"`
	PanNumber  string    `json:"panNumber"`
	MerchantId uuid.UUID `json:"merchant_id"`
}
