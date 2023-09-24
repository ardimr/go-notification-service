package utils

import (
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateOTP(userEmail string) (string, string, error) {
	// perdiodDuration := time.Duration(5) * time.Minute
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "mata-duitan.org",
		AccountName: userEmail,
		Digits:      otp.DigitsSix,
		Period:      300,
	})

	if err != nil {
		return "", "", err
	}

	otp, err := totp.GenerateCode(key.Secret(), time.Now())

	return otp, key.Secret(), err
}

func VerifyOTP(userSecret string, otpCode string) bool {
	fmt.Println(otpCode, userSecret)
	passcode, err := totp.ValidateCustom(otpCode, userSecret, time.Now(), totp.ValidateOpts{
		Period: 300,
		Digits: otp.DigitsSix,
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(passcode)

	return totp.Validate(otpCode, userSecret)
}
