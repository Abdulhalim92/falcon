package otp

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
	"image/png"
)

func GenerateOTP() (string, bytes.Buffer, error) {
	var qrCode bytes.Buffer

	// генерация OTP
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "humopay",
		AccountName: "humopay",
		SecretSize:  15,
	})
	if err != nil {
		return "", bytes.Buffer{}, errors.Wrap(err, "unable to generate OTP")
	}

	// создание QR - кода
	img, err := key.Image(200, 200)
	if err != nil {
		return "", bytes.Buffer{}, errors.Wrap(err, "unable to create image for QR - code")
	}

	err = png.Encode(&qrCode, img)
	if err != nil {
		return "", bytes.Buffer{}, errors.Wrap(err, "unable to save image from created QR - code")
	}

	return key.Secret(), qrCode, nil
}

func ValidateOTP(passcode, secret string) bool {
	return totp.Validate(passcode, secret)
}
