package service

import (
	"Problem-solving/dao"
	"Problem-solving/dto"
	"Problem-solving/models"
)

// ----------------------------------------------------------------
// CreateKyc creates a new KYC record in the database and returns
// a ResponseKyc DTO containing the saved details.
// ----------------------------------------------------------------
func CreateKyc(kycToCreate *models.KYC) (*dto.ResponseKyc, error) {

	newKycRecord, err := dao.InsertKyc(kycToCreate)
	if err != nil {
		return nil, err
	}

	// Return a ResponseKyc DTO with the newly created KYC details.
	return &dto.ResponseKyc{
		Name:         newKycRecord.Name,
		PanNumber:    newKycRecord.PanNumber,
		AadharNumber: newKycRecord.AadharNumber,
		MerchantId:   newKycRecord.MerchantId,
	}, nil
}

// -----------------------------------------------------------------------
// UpdateKycDetails updates the KYC information for a specific merchant.
// It accepts the merchantId and panNumber as input, updates the record,
// and returns the updated pan number or an error.
// -----------------------------------------------------------------------

func UpdateKycDetails(merchantId, panNumber, aadharNumber string) (string, error) {

	updatePanAndAadhar, err := dao.UpdateKyc(merchantId, panNumber, aadharNumber)

	if err != nil {
		return "", err
	}

	// Return the updated PAN number if the update was successful.
	return updatePanAndAadhar, nil
}
