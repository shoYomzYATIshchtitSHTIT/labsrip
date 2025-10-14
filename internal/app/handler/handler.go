package handler

import (
	"Backend-RIP/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.Engine, repo *repository.Repository) {
	apiRouter := router.Group("/api")

	intervalHandler := NewIntervalHandler(repo)
	intervalRouter := apiRouter.Group("/intervals")
	{
		intervalRouter.GET("", intervalHandler.GetIntervals)
		intervalRouter.GET("/:id", intervalHandler.GetInterval)
		intervalRouter.POST("", intervalHandler.CreateInterval)
		intervalRouter.PUT("/:id", intervalHandler.UpdateInterval)
		intervalRouter.DELETE("/:id", intervalHandler.DeleteInterval)
		intervalRouter.POST("/add-to-composition", intervalHandler.AddIntervalToComposition)
		intervalRouter.POST("/:id/image", intervalHandler.UpdateIntervalPhoto)
	}

	compositionHandler := NewCompositionHandler(repo)
	compositionRouter := apiRouter.Group("/compositions")
	{
		compositionRouter.GET("/comp-cart", compositionHandler.GetCompositionCart)
		compositionRouter.GET("", compositionHandler.GetCompositions)
		compositionRouter.GET("/:id", compositionHandler.GetComposition)
		compositionRouter.PUT("/:id", compositionHandler.UpdateCompositionFields)
		compositionRouter.PUT("/:id/form", compositionHandler.FormComposition)
		compositionRouter.PUT("/:id/complete", compositionHandler.CompleteComposition)
		compositionRouter.PUT("/:id/reject", compositionHandler.RejectComposition)
		compositionRouter.DELETE("/:id", compositionHandler.DeleteComposition)
	}

	compositionIntervalHandler := NewCompositionIntervalHandler(repo)
	compositionIntervalRouter := apiRouter.Group("/composition-intervals")
	{
		compositionIntervalRouter.DELETE("", compositionIntervalHandler.RemoveFromComposition)
		compositionIntervalRouter.PUT("", compositionIntervalHandler.UpdateCompositionInterval)
	}

	userHandler := NewUserHandler(repo)
	userRouter := apiRouter.Group("/users")
	{
		userRouter.POST("/register", userHandler.Register)
		userRouter.GET("/profile", userHandler.GetProfile)
		userRouter.PUT("/profile", userHandler.UpdateProfile)
		userRouter.POST("/login", userHandler.Login)
		userRouter.POST("/logout", userHandler.Logout)
	}
}
