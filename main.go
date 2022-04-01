package main

import (
	"github.com/DragonSov/smasher/db"
	"github.com/DragonSov/smasher/server/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	db.Init()
}

func main() {
	router := gin.Default()

	// Setting routes
	routes.StatusRoutes(&router.RouterGroup)
	routes.UserRoutes(&router.RouterGroup, db.DB)
	routes.WalletRoutes(&router.RouterGroup, db.DB)

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
