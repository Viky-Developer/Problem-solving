package models

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ----------------------------------------------------------------
// KYC represents the Know Your company information for a company.
// ----------------------------------------------------------------
type KYC struct {
	gorm.Model
	Name         string    `json:"name,omitempty"`
	PanNumber    string    `json:"panNumber,omitempty" gorm:"uniqueIndex"`
	AadharNumber string    `json:"aadharNumber,omitempty" gorm:"uniqueIndex"`
	MerchantId   uuid.UUID `json:"merchant_id,omitempty" gorm:"not null"`
}

// -----------------------------------------------------------------------------------------
// BeforeCreate is a GORM hook that runs before a new record is inserted into the database.
// -----------------------------------------------------------------------------------------
func (k *KYC) BeforeCreate(tx *gorm.DB) (err error) {

	if k.MerchantId == uuid.Nil {
		k.MerchantId, err = uuid.NewUUID()
		if err != nil {
			log.Println("UUID not generating for kyc")
		}
	}
	return nil
}
