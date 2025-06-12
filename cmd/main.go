package main

import (
	"log"
	"net/http"

	"github.com/jithesh-wq/oolio-food-cart/db/postgres"
	"github.com/jithesh-wq/oolio-food-cart/logger"
	"github.com/jithesh-wq/oolio-food-cart/routes"
	"github.com/jithesh-wq/oolio-food-cart/store"
)

func main() {
	//initialize the logger
	if err := logger.Init(); err != nil {
		logger.Log.Infoln("Error initializing logger:", err)
		return
	}
	logger.Log.Info("Logger initialized successfully")
	defer logger.Close()

	//initialize the memory store
	memoryStore := store.NewMemoryStore()
	logger.Log.Infow("Memory store initialized successfully")

	//initialize the database connection
	db, err := postgres.CreatePostgres()
	if err != nil {
		log.Fatal(err.Error())
	}

	//create the memory store
	routes := routes.CreateRoutes(memoryStore, db)

	//add server configs and start server
	server := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}
	logger.Log.Infoln("Server is starting on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		logger.Log.Infoln("Error starting server:", err)
	} else {
		logger.Log.Infoln("Server started successfully")
	}
	logger.Log.Infoln("Exiting main function")
}
