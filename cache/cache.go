package cache

import (
	"Problem-solving/dto"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Config struct {
	KycNameAndPanNumber map[string]KycRecord `json:"kyc_name_and_pan_number"`
}

type KycRecord struct {
	Name       string    `json:"name"`
	PanNumber  string    `json:"pan_number"`
	MerchantId uuid.UUID `json:"merchant_id"`
}

var (
	configFilePath = "kycDetails.json" // Can be made configurable
	kycDetails     sync.Map            // In-memory cache of KYC data
	db             *gorm.DB            // Database connection instance
)

// -----------------------------------------------------------------------------
// InitializeDB sets up the database and optionally populates cache from file.
// -----------------------------------------------------------------------------
func InitializeDB(database *gorm.DB) error {
	db = database

	// Attempt to load cache from file
	if err := loadCacheFromFile(); err != nil {
		log.Printf("Cache file not found or invalid. Fetching from database: %v", err)
		return populateCacheFromDB()
	}
	log.Println("Cache initialized from file")
	return nil
}

// --------------------------------------------------
// InvalidateCache removes an entry from the cache.
// --------------------------------------------------
func InvalidateCache(key string) {

	kycDetails.Delete(key)
	log.Printf("Cache invalidated for key: %s", key)
}

// ---------------------------------------------------
// SetToCache adds or updates an entry in the cache.
// ---------------------------------------------------
func SetToCache(key string, kycData KycRecord) {

	if kycData.Name == "" || kycData.PanNumber == "" {
		log.Printf("Invalid KYC data for key: %s. Skipping cache update.", key)
		return
	}

	kycDetails.Store(key, kycData)
	log.Printf("Cache updated for key: %s", key)
}

// ---------------------------------------------------------
// SaveCacheToFile persists the in-memory cache to a file.
// ---------------------------------------------------------
func SaveCacheToFile() error {

	config := buildConfigFromCache()

	file, err := os.Create(configFilePath)
	if err != nil {
		log.Printf("Error creating cache file: %v", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	// Use indent for better readability
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		log.Printf("Error encoding cache to file: %v", err)
		return err
	}

	log.Printf("Cache successfully saved to file: %s", configFilePath)
	return nil
}

// ------------------------------------------------------
// FetchKycDetails retrieves KYC data from the database.
// ------------------------------------------------------
func FetchKycDetails() ([]dto.KYCDetails, error) {

	var kycDetailsList []dto.KYCDetails
	if err := db.Table("kycs").Select("name, pan_number, merchant_id").Find(&kycDetailsList).Error; err != nil {
		log.Printf("Error fetching KYC details: %v", err)
		return nil, err
	}
	return kycDetailsList, nil
}

// Private Helper Methods

// -------------------------------------------------------
// loadCacheFromFile loads the cache from the JSON file.
// -------------------------------------------------------
func loadCacheFromFile() error {

	file, err := os.Open(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to open cache file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	if err := decoder.Decode(&config); err != nil {
		return fmt.Errorf("failed to decode cache file: %w", err)
	}

	for key, value := range config.KycNameAndPanNumber {

		if value.MerchantId == uuid.Nil || value.Name == "" || value.PanNumber == "" {
			log.Printf("Skipping invalid cache entry: %+v", value)
			continue
		}

		kycDetails.Store(key, value)
	}
	return nil
}

// ----------------------------------------------------------------------------
// populateCacheFromDB fetches data from the database and populates the cache.
// ----------------------------------------------------------------------------
func populateCacheFromDB() error {

	kycDetailsList, err := FetchKycDetails()
	if err != nil {
		log.Printf("Error fetching KYC details from DB: %v", err)
		return err
	}

	for _, kyc := range kycDetailsList {

		if kyc.MerchantId == uuid.Nil || kyc.Name == "" || kyc.PanNumber == "" {
			log.Printf("Skipping invald KYC record: %+v", kyc)
			continue
		}
		cacheKey := kyc.MerchantId.String()
		kycRecord := KycRecord{
			Name:       kyc.Name,
			PanNumber:  kyc.PanNumber,
			MerchantId: kyc.MerchantId,
		}
		kycDetails.Store(cacheKey, kycRecord)
		log.Printf("Cache populated for merchant_id: %s", cacheKey)
	}

	// Save to file after population
	if err := SaveCacheToFile(); err != nil {
		log.Printf("Error populated for MerchantId: %s", err)
		return err
	}

	return nil
}

// -----------------------------------------------------------------
// buildConfigFromCache converts in-memory cache to Config struct.
// -----------------------------------------------------------------
func buildConfigFromCache() Config {

	config := Config{KycNameAndPanNumber: make(map[string]KycRecord)}
	kycDetails.Range(func(key, value interface{}) bool {
		strKey, _ := key.(string)
		kycRecord, _ := value.(KycRecord)
		config.KycNameAndPanNumber[strKey] = kycRecord
		return true
	})
	return config
}
