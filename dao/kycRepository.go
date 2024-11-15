package dao

import (
	"Problem-solving/cache"
	"Problem-solving/dto"
	"Problem-solving/models"
	"fmt"
	"log"

	"gorm.io/gorm"
)

var writeDB *gorm.DB

// -------------------------------------------------
// SetDB initializes the write database connection.
// -------------------------------------------------
func SetDB(write *gorm.DB) {

	if write == nil {
		log.Println("Received nil database connection")
	}
	writeDB = write
}

// -----------------------------------------------------------------------------
// InsertKyc inserts a new KYC record into the database and updates the cache.
// -----------------------------------------------------------------------------
func InsertKyc(kyc *models.KYC) (*dto.ResponseKyc, error) {

	// Insert the KYC record into the database
	if err := writeDB.Create(kyc).Error; err != nil {
		log.Printf("Error inserting KYC: %v", err)
		return nil, err
	}

	// Invalidate the cache for this merchant ID
	cacheKey := kyc.MerchantId.String() // Assuming MerchantId is of type string
	cache.InvalidateCache(cacheKey)

	// Create a new KYC info structure for the cache
	kycInfo := cache.KycRecord{
		Name:       kyc.Name,
		PanNumber:  kyc.PanNumber,
		MerchantId: kyc.MerchantId,
	}

	// Set the cache for this merchant ID
	cache.SetToCache(cacheKey, kycInfo)
	log.Printf("Cache set for merchant_id: %s", cacheKey)

	cache.SaveCacheToFile()

	// Return the response DTO
	return &dto.ResponseKyc{
		Name:         kyc.Name,
		PanNumber:    kyc.PanNumber,
		AadharNumber: kyc.AadharNumber,
		MerchantId:   kyc.MerchantId,
	}, nil
}

// ----------------------------------------------------------------------------------------------------
// UpdateKyc updates the KYC record for a given merchant ID with the provided PAN and Aadhar numbers.
// ----------------------------------------------------------------------------------------------------
func UpdateKyc(merchantId, panNumber, aadharNumber string) (string, error) {

	// Update the KYC record in the database
	var kyc models.KYC

	// Attempt to update the KYC record in the database with new PAN and Aadhar numbers.
	result := writeDB.Model(&models.KYC{}).Where("merchant_id = ?", merchantId).Updates(map[string]interface{}{
		"pan_number":    panNumber,
		"aadhar_number": aadharNumber,
	})

	// Check if the update affected any rows
	if result.RowsAffected == 0 {
		log.Printf("No record found for merchant_id: %s", merchantId)
		return "", fmt.Errorf("no record found for merchant_id: %s", merchantId)
	}

	// Handle any errors during the update operation
	if result.Error != nil {
		log.Printf("Error updating KYC: %v", result.Error)
		return "", result.Error
	}

	// Retrieve the updated KYC info
	if err := writeDB.Where("merchant_id = ?", merchantId).First(&kyc).Error; err != nil {
		log.Printf("Error retrieving updated KYC: %v", err)
		return "", err
	}

	// Prepare the cache info with the updated values
	kycInfo := cache.KycRecord{
		Name:       kyc.Name,
		PanNumber:  kyc.PanNumber,
		MerchantId: kyc.MerchantId,
	}

	// Invalidate the cache for this merchant ID
	cacheKey := merchantId
	cache.InvalidateCache(cacheKey)

	// Update the cache with the new KYC info
	cache.SetToCache(cacheKey, kycInfo)
	log.Printf("Cache updated for merchant_id: %s", cacheKey)

	cache.SaveCacheToFile()

	// Return the updated pan number
	return panNumber, nil
}
