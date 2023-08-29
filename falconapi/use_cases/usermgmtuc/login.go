package usermgmtuc

import (
	"context"
	"falconapi/domain/entities"
	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	Username string `validate:"required,min=3,max=15"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	Message      string
	OtpGenerated bool
	UserID       string
}

type LoginUseCase struct {
	identityManager identityManager
}

func NewLoginUseCase(im identityManager) *LoginUseCase {
	return &LoginUseCase{
		identityManager: im,
	}
}

func (uc *LoginUseCase) Login(ctx context.Context, request LoginRequest) (*LoginResponse, *entities.ErrorModel) {
	var (
		validate   = validator.New()
		response   LoginResponse
		errorModel entities.ErrorModel
		attributes map[string][]string
	)

	err := validate.Struct(request)
	if err != nil {
		errorModel.Message = err.Error()
		errorModel.Err = err
		errorModel.Code = 1 // todo
		return nil, &errorModel
	}

	user, _, err := uc.identityManager.LoginUser(ctx, request.Username, request.Password)
	if err != nil {
		errorModel.Message = err.Error()
		errorModel.Err = err
		errorModel.Code = 1 // todo
		return nil, &errorModel
	}

	response.UserID = *user.ID

	if user.Attributes == nil {
		response.Message = "please generate OTP"
		response.OtpGenerated = false
	} else if user.Attributes != nil {
		attributes = *user.Attributes
		_, ok := attributes["otp_secret"]
		if !ok {
			response.Message = "please generate OTP"
			response.OtpGenerated = false
		} else if ok {
			response.Message = "please validate OTP"
			response.OtpGenerated = true
		}
	}

	return &response, nil
}
