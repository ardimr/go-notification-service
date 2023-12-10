package controller

import (
	"database/sql"
	"go_project_template/internal/exception"
	"go_project_template/internal/user/model"
	"go_project_template/internal/user/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase usecase.IUserUseCase
}

func NewUserController(userUseCase usecase.IUserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

// Controller Implementation

func (controller *UserController) GetUsers(ctx *gin.Context) {
	// Get users data from db
	users, err := controller.userUseCase.GetUsers(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"Message": err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		users,
	)
}

func (controller *UserController) GetUserById(ctx *gin.Context) {
	var reqUri model.GetUserByIdReqUri

	// Request URI Binding
	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": err.Error(),
			},
		)
		return
	}

	user, err := controller.userUseCase.GetUserById(ctx, reqUri.ID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(
				http.StatusNotFound,
				gin.H{
					"Message": "Not Found",
				},
			)
		}
		return
	}

	// Success state
	ctx.JSON(
		http.StatusOK,
		user,
	)
}

func (controler *UserController) AddNewUser(ctx *gin.Context) {
	var newUser model.User

	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"Message": err.Error(),
			},
		)
		return
	}

	newId, err := controler.userUseCase.AddNewUser(ctx, newUser)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"Message": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"New Id": newId,
		},
	)
}

func (controller *UserController) UpdateUser(ctx *gin.Context) {
	var user model.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"Message": err.Error(),
			},
		)
		return
	}

	res, err := controller.userUseCase.UpdateUser(ctx, user)

	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"Message": err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"Rows affected": res,
		},
	)
}

func (controller *UserController) DeleteUser(ctx *gin.Context) {
	var reqUri model.DeleteUserReqUri

	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"Message": err.Error(),
			},
		)
		return
	}

	if err := controller.userUseCase.DeleteUser(ctx, reqUri.ID); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"Message": err.Error(),
			},
		)
	}

	ctx.Status(
		http.StatusOK,
	)
}

func (controller *UserController) UserRegistration(ctx *gin.Context) {
	var newUser model.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(exception.ErrorResponse(err))
		return
	}

	// Add user
	err := controller.userUseCase.RegisterUser(ctx, newUser)

	if err != nil {
		ctx.JSON(exception.ErrorResponse(err))
		return
	}

	ctx.Status(http.StatusAccepted)
}

func (controller *UserController) VerifyOTP(ctx *gin.Context) {
	// var reqBody model.UserOTPVerification
	var reqParam model.UserOTPVerificationParam
	if err := ctx.ShouldBindQuery(&reqParam); err != nil {
		ctx.AbortWithStatusJSON(exception.ErrorResponse(err))
		return
	}

	email, err := controller.userUseCase.VerifyOTP(ctx, reqParam.OTPCode)

	if err != nil {
		ctx.AbortWithStatusJSON(exception.ErrorResponse(err))
		return
	}

	// OTP is verified, Mark the user's email address as verified in the database

	err = controller.userUseCase.UpdateEmailVerificationStatus(ctx, email)
	if err != nil {
		ctx.AbortWithStatusJSON(exception.ErrorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *UserController) RequestOTP(ctx *gin.Context) {
	var reqBody model.UserOTPRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(exception.ErrorResponse(err))
		return
	}

	err := controller.userUseCase.RequestNewOTP(ctx, reqBody.Email)

	if err != nil {
		ctx.AbortWithStatusJSON(exception.ErrorResponse(err))
		return
	}

	ctx.Status(http.StatusAccepted)
}
