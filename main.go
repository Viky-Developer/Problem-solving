package main

import (
	"Problem-solving/cache"
	"Problem-solving/config"
	"Problem-solving/dao"
	"Problem-solving/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

// add docs
func main() {

	//Load env file
	configure, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
	}

	//Create a db instance
	writeDB, err := config.DbCreate(configure)
	if err != nil {
		log.Println("Failed to connect database")
	}

	// create table
	config.CreateTable(writeDB)

	// Check the database connection
	dao.SetDB(writeDB)

	cache.InitializeDB(writeDB)

	// Gin router initialize
	router := gin.Default()

	router.POST("/create-kyc", handlers.NewKyc)
	router.PATCH("/kyc/:merchantId", handlers.UpdateKyc)

	log.Printf("Starting server on port %s", configure.Port)

	// get a port from env file
	router.Run(":" + configure.Port)

}
