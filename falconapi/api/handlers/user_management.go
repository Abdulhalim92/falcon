package handlers

import (
	"context"
	"falconapi/domain/entities"
	"falconapi/use_cases/usermgmtuc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type RegisterUseCase interface {
	Register(context.Context, usermgmtuc.RegisterRequest) (*usermgmtuc.RegisterResponse, *entities.ErrorModel)
}

type LoginUseCase interface {
	Login(ctx context.Context, request usermgmtuc.LoginRequest) (*usermgmtuc.LoginResponse, *entities.ErrorModel)
}

type OtpUseCase interface {
	GenerateOTP(ctx context.Context, request usermgmtuc.GenerateOtpRequest) (*usermgmtuc.GenerateOtpResponse, *entities.ErrorModel)
	ValidateOTP(ctx context.Context, request usermgmtuc.ValidateOtpRequest) (*usermgmtuc.ValidateOtpResponse, *entities.ErrorModel)
}

// @Summary Метод регистрации пользователя
// @Description Регистрация пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param RegisterInput body usermgmtuc.RegisterRequest true "Login data"
// @OperationId login
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user [post]
func Register(useCase RegisterUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx     = c.Request.Context()
			request = usermgmtuc.RegisterRequest{}
		)

		err := c.BindJSON(&request)
		if err != nil {
			log.Println(err, "unable to parse incoming request")
			c.JSON(http.StatusBadRequest, "unable to parse incoming request")
			return
		}

		response, errModel := useCase.Register(ctx, request)
		if errModel != nil {
			log.Println(errModel.Err, "unable to register user")
			if errModel.Code == 1 {
				c.JSON(http.StatusBadRequest, "please send valid data")
				return
			} else if errModel.Code == 2 {
				c.JSON(http.StatusInternalServerError, "something went wrong")
				return
			}
		}

		c.JSON(http.StatusCreated, response.User.ID)
	}
}

// @Summary Метод входа пользователя
// @Description Вход пользователя под логином и паролем
// @Tags Auth
// @Accept json
// @Produce json
// @Param LoginInput body usermgmtuc.LoginRequest true "Login data"
// @OperationId login
// @Success 200 {object} usermgmtuc.LoginResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /login [post]
func Login(useCase LoginUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx     = c.Request.Context()
			request = usermgmtuc.LoginRequest{}
		)

		err := c.BindJSON(&request)
		if err != nil {
			log.Println(err, "unable to parse incoming request")
			c.JSON(http.StatusBadRequest, "unable to parse incoming request")
			return
		}

		response, errModel := useCase.Login(ctx, request)
		if errModel != nil {
			log.Println(errModel.Err, "unable to login user")
			if errModel.Code == 1 {
				c.JSON(http.StatusBadRequest, "incorrect login or password")
				return
			} else if errModel.Code == 2 {
				c.JSON(http.StatusInternalServerError, "something went wrong")
				return
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Метод генарации OTP
// @Description Генерация OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param GenerateOtp body usermgmtuc.GenerateOtpRequest true "Generate OTP data"
// @Success 200 {string} binary "PNG image data"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /generate-otp [post]
func GenerateOtp(useCase OtpUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx     = c.Request.Context()
			request = usermgmtuc.GenerateOtpRequest{}
		)

		err := c.BindJSON(&request)
		if err != nil {
			log.Println(err, "unable to parse incoming request")
			c.JSON(http.StatusBadRequest, "unable to parse incoming request")
			return
		}

		response, errModel := useCase.GenerateOTP(ctx, request)
		if errModel != nil {
			log.Println(errModel.Err, "unable to generate OTP")
			if errModel.Code == 1 {
				c.JSON(http.StatusBadRequest, "please send valid data")
				return
			} else if errModel.Code == 2 {
				c.JSON(http.StatusInternalServerError, "something went wrong")
				return
			}
		}

		c.Set("Content-Type", "image/png")
		c.Set("Content-Length", strconv.Itoa(len(response.QrCode.Bytes())))

		c.Status(http.StatusOK)
		_, err = c.Writer.Write(response.QrCode.Bytes())
		if err != nil {
			log.Println("unable to write image.")
		}

	}
}

// @Summary Метод валидации OTP
// @Description Валидация OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param ValidateOtp body usermgmtuc.ValidateOtpRequest true "Validate OTP data"
// @Success 200 {object} usermgmtuc.ValidateOtpResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /validate-otp [post]
func ValidateOtp(useCase OtpUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ctx     = c.Request.Context()
			request = usermgmtuc.ValidateOtpRequest{}
		)

		err := c.BindJSON(&request)
		if err != nil {
			log.Println(err, "unable to parse incoming request")
			c.JSON(http.StatusBadRequest, "unable to parse incoming request")
			return
		}

		response, errModel := useCase.ValidateOTP(ctx, request)
		if errModel != nil {
			log.Println(err, "unable to validate OTP")
			if errModel.Code == 1 {
				c.JSON(http.StatusBadRequest, "please send valid data")
				return
			} else if errModel.Code == 2 {
				c.JSON(http.StatusInternalServerError, "something went wrong")
				return
			}
		}

		c.JSON(http.StatusOK, response)
	}
}
