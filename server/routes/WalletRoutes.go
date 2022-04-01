package routes

import (
	"github.com/DragonSov/smasher/server/domain/Users"
	"github.com/DragonSov/smasher/server/domain/Wallets"
	"github.com/DragonSov/smasher/server/handlers"
	"github.com/DragonSov/smasher/server/middlewares"
	"github.com/DragonSov/smasher/server/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func WalletRoutes(apiRouter *gin.RouterGroup, db *sqlx.DB) {
	h := handlers.WalletHandlers{
		Service: services.NewWalletService(Wallets.NewWalletRepositoryDb(db)),
	}
	m := middlewares.UserMiddlewares{
		Service: services.NewUserService(Users.NewUserRepositoryDb(db)),
	}
	route := apiRouter.Group("/wallets")
	route.Use(m.CheckAuthorization())
	{
		route.POST("/", h.CreateWallet)
		route.GET("/", h.SelectUserWallets)
		route.GET("/:uuid", h.SelectWalletByUUID)
		route.PUT("/:uuid/replenish", h.ReplenishWalletByUUID)
		route.PUT("/:uuid/transfer", h.TransferWalletByUUID)
		route.DELETE("/:uuid", h.DeleteWalletByUUID)
	}
}
