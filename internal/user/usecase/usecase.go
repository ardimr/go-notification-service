package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	queueclient "go_project_template/configs/queue_client"
	"go_project_template/internal/auth"
	"go_project_template/internal/user/model"
	"go_project_template/internal/user/repository"
	"go_project_template/internal/utils"
)

type IUserUseCase interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserById(ctx context.Context, id int64) (model.User, error)
	AddNewUser(ctx context.Context, newUser model.User) (int64, error)
	UpdateUser(ctx context.Context, user model.User) (int64, error)
	DeleteUser(ctx context.Context, id int64) error
	RegisterUser(ctx context.Context, newUser model.User) error
	VerifyOTP(ctx context.Context, otpCode string) (string, error)
	UpdateEmailVerificationStatus(ctx context.Context, email string) error
	RequestNewOTP(ctx context.Context, email string) error
}

type UserUseCase struct {
	userRepo  repository.IUserRepository
	userCache repository.IUserRedisRepository
	publisher *queueclient.Publisher
}

func NewUserUseCae(userRepo repository.IUserRepository, userCache repository.IUserRedisRepository, publisher *queueclient.Publisher) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		userCache: userCache,
		publisher: publisher,
	}
}

func (uc *UserUseCase) GetUsers(ctx context.Context) ([]model.User, error) {
	return nil, nil
}

func (uc *UserUseCase) GetUserById(ctx context.Context, id int64) (model.User, error) {

	userDetails, err := uc.userRepo.GetUserById(ctx, id)

	if err != nil {
		return userDetails, err
	}

	return userDetails, nil
}

func (uc *UserUseCase) AddNewUser(ctx context.Context, newUser model.User) (int64, error) {
	return 0, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user model.User) (int64, error) {
	return 0, nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id int64) error {
	return nil
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, newUser model.User) error {
	// Hash user's password
	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		return nil
	}
	newUser.Password = hashedPassword

	// Generate OTP
	otp, secret, err := utils.GenerateOTP(newUser.Email)

	if err != nil {
		return err
	}

	// User OTP Verification data
	userOTPVerification := model.UserOTPVerification{
		OTPCode: otp,
		Secret:  secret,
		Email:   newUser.Email,
	}

	// Store temporary in database for 5 minutes
	expiration := time.Duration(5) * time.Minute

	err = uc.userCache.SetUserOTP(ctx, otp, userOTPVerification, expiration)

	if err != nil {
		return err
	}
	// Insert user to database
	_, err = uc.userRepo.AddNewUser(ctx, newUser)

	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUseCase) RequestNewOTP(ctx context.Context, email string) error {
	// Check if email is already exist in database
	user, err := uc.userRepo.FindByEmail(ctx, email)

	if err != nil {
		return err
	}

	if user.IsVerified {
		return errors.New("[otp] user is already verified")
	}

	// Generate OTP
	otp, secret, err := utils.GenerateOTP(email)

	if err != nil {
		return err
	}

	// User OTP Verification data
	userOTPVerification := model.UserOTPVerification{
		OTPCode: otp,
		Secret:  secret,
		Email:   email,
	}

	// Store temporary in database for 5 minutes
	expiration := time.Duration(5) * time.Minute

	err = uc.userCache.SetUserOTP(ctx, otp, userOTPVerification, expiration)

	if err != nil {
		return err
	}

	verificationEmailPayload := model.OTPVerificationEmailContent{
		Email:   userOTPVerification.Email,
		OTPCode: userOTPVerification.OTPCode,
		Url:     fmt.Sprintf("http://localhost:8080/api/user-service/verify-otp?otp_code=%s", userOTPVerification.OTPCode),
	}

	publishPayload, err := json.Marshal(verificationEmailPayload)

	if err != nil {
		return err
	}

	err = uc.publisher.Publish(ctx, "mailQueue", publishPayload)

	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) VerifyOTP(ctx context.Context, otpCode string) (string, error) {
	// Get the user otp secret from cache
	userOTP, err := uc.userCache.GetUserOTP(ctx, otpCode)
	if err != nil {
		return userOTP.Email, err
	}

	if !utils.VerifyOTP(userOTP.Secret, otpCode) {
		return userOTP.Email, errors.New("invalid.otp.code")
	}

	// Remove otp from cache
	err = uc.userCache.RemoveUserOTP(ctx, otpCode)

	if err != nil {
		return userOTP.Email, err
	}

	return userOTP.Email, nil
}

func (uc *UserUseCase) UpdateEmailVerificationStatus(ctx context.Context, email string) error {

	err := uc.userRepo.UpdateEmailVerificationStatus(ctx, email)

	if err != nil {
		return err
	}
	return nil
}
