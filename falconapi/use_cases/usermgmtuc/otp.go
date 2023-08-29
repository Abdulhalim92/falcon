package usermgmtuc

import (
	"bytes"
	"context"
	"falconapi/domain/entities"
	"github.com/pkg/errors"
)

type ValidateOtpRequest struct {
	UserID   string `json:"user_id"`
	OtpToken string `json:"otp_token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ValidateOtpResponse struct {
	Message string `json:"message,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	//QrCode       bytes.Buffer `json:"qr_code,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type GenerateOtpRequest struct {
	UserID string `json:"user_id"`
}

type GenerateOtpResponse struct {
	UserID string       `json:"user_id,omitempty"`
	QrCode bytes.Buffer `json:"qr_code,omitempty"`
}

type OtpUseCase struct {
	identityManager identityManager
}

func NewOtpUseCase(im identityManager) *OtpUseCase {
	return &OtpUseCase{
		identityManager: im,
	}
}

func (uc *OtpUseCase) GenerateOTP(ctx context.Context, request GenerateOtpRequest) (*GenerateOtpResponse, *entities.ErrorModel) {
	var (
		errorModel entities.ErrorModel
	)

	otpSecret, buf, err := uc.identityManager.GenerateOTP(ctx)
	if err != nil {
		errorModel.Message = err.Error()
		errorModel.Err = err
		errorModel.Code = 2 // todo
		return nil, &errorModel
	}

	attribute := make(map[string][]string)
	attribute["otp_secret"] = []string{otpSecret}

	err = uc.identityManager.UpdateUserAttribute(ctx, request.UserID, attribute)
	if err != nil {
		errorModel.Message = err.Error()
		errorModel.Err = err
		errorModel.Code = 2 // todo
		return nil, &errorModel
	}

	var response = &GenerateOtpResponse{
		UserID: request.UserID,
		QrCode: buf,
	}

	return response, nil
}

func (uc *OtpUseCase) ValidateOTP(ctx context.Context, request ValidateOtpRequest) (*ValidateOtpResponse, *entities.ErrorModel) {
	var (
		errorModel entities.ErrorModel
	)

	validateOTP, err := uc.identityManager.ValidateOTP(ctx, request.UserID, request.OtpToken)
	if err != nil {
		errorModel.Message = err.Error()
		errorModel.Err = err
		errorModel.Code = 2 // todo
		return nil, &errorModel
	}
	if !validateOTP {
		errorModel.Message = "invalid otp"
		errorModel.Err = errors.New("invalid otp")
		errorModel.Code = 1
		return nil, &errorModel
	}

	user, token, err := uc.identityManager.LoginUser(ctx, request.Username, request.Password)
	if err != nil {
		errorModel.Message = err.Error()
		errorModel.Err = err
		errorModel.Code = 1 // todo
		return nil, &errorModel
	}

	var response = &ValidateOtpResponse{
		Message:      "valid OTP",
		UserID:       *user.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	return response, nil
}
