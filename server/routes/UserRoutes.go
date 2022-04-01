package routes

import (
	"github.com/DragonSov/smasher/server/domain/Users"
	"github.com/DragonSov/smasher/server/handlers"
	"github.com/DragonSov/smasher/server/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func UserRoutes(apiRouter *gin.RouterGroup, db *sqlx.DB) {
	h := handlers.UserHandlers{
		Service: services.NewUserService(Users.NewUserRepositoryDb(db)),
	}
	route := apiRouter.Group("/users")
	{
		route.POST("/", h.CreateUser)
		route.GET("/", h.SelectUserByLogin)
		route.GET("/:uuid", h.SelectUserByUUID)
		route.POST("/sign-in", h.SignIn)
	}
}
