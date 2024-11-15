package handlers

import (
	"Problem-solving/dto"
	"Problem-solving/models"
	"Problem-solving/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --------------------------------------------------------------------
// This function is responsible for receiving a KYC creation request,
// validating it, and creating the KYC record.
// --------------------------------------------------------------------
func NewKyc(context *gin.Context) {

	var kycRequest models.KYC

	// Bind the incoming JSON body to the kyc variable.
	if err := context.ShouldBindBodyWithJSON(&kycRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service layer to create a new KYC record using the parsed data.
	savedKyc, err := service.CreateKyc(&kycRequest)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a successful response with the created KYC record details
	context.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "KYC created successfully",
		"data":       savedKyc,
	})
}

// -------------------------------------------------------------------------
// UpdateKyc is a handler function for updating KYC details of a merchant.
// -------------------------------------------------------------------------
func UpdateKyc(context *gin.Context) {

	// Get merchantId from URL params
	merchantId := context.Param("merchantId")

	// Check if merchantId is provided
	if merchantId == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "MerchantId is required"})
		return
	}

	// Declare the request body variable using the RequestBodyKyc struct
	var requestBody dto.RequestBodyKyc

	// Bind the request body to the struct
	if err := context.ShouldBindBodyWithJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if must two of PanNumber or AadharNumber is provided
	if requestBody.PanNumber == "" && requestBody.AadharNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Both panNumber and aadharNumber are required"})
		return
	}

	// Call the service function to update KYC details
	updatedPanNumber, err := service.UpdateKycDetails(merchantId, requestBody.PanNumber, requestBody.AadharNumber)

	// Check for errors in updating KYC
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	context.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "KYC updated successfully",
		"panNumber":  updatedPanNumber,
	})
}
