package main

import (
	"flight-booking/internal/database"
	"flight-booking/internal/router"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Setup Gin router
	r := router.SetupRouter(db)

	r.Run() // listen and serve on 0.0.0.0:8080
}
