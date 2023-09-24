package user

import (
	"go_project_template/internal/user/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	controller *controller.UserController
}

func NewRouter(controller *controller.UserController) *Router {
	return &Router{
		controller: controller,
	}
}

func (router *Router) AddRoute(superRoute *gin.RouterGroup) {
	router.userRoutes(superRoute)
	// router.sfmRoutes(superRoute)
}

func (router *Router) userRoutes(superRoute *gin.RouterGroup) {
	userRouter := superRoute.Group("/user-service")
	userRouter.GET("/users", router.controller.GetUsers)
	userRouter.GET("/users/:id", router.controller.GetUserById)
	userRouter.POST("/users", router.controller.AddNewUser)
	userRouter.PATCH("/users", router.controller.UpdateUser)
	userRouter.DELETE("/users/:id", router.controller.DeleteUser)

	userRouter.POST("/signup", router.controller.UserRegistration)
	userRouter.POST("/verify-otp", router.controller.VerifyOTP)
	userRouter.POST("/request-otp", router.controller.RequestOTP)
}
