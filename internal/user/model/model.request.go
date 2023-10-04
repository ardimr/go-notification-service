package model

type GetUserByIdReqUri struct {
	ID int64 `uri:"id"`
}

type AddNewUserReqBody struct {
	Name string `json:"name"`
}

type DeleteUserReqUri struct {
	ID int64 `uri:"id"`
}

type UserOTPVerification struct {
	Email   string `json:"email" binding:"required,email"`
	OTPCode string `json:"otp_code" binding:"required"`
	Secret  string `json:"secret" binding:"required"`
}

type UserOTPVerificationParam struct {
	OTPCode string `form:"otp_code"`
}

type UserOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type OTPVerificationEmailContent struct {
	Email   string `json:"email" binding:"required,email"`
	OTPCode string `json:"otp_code" binding:"required"`
	Url     string `json:"url"`
}
