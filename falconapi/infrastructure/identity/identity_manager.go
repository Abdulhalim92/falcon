package identity

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/viper"
)

type identityManager struct {
	baseUrl             string
	realm               string
	restApiClientId     string
	restApiClientSecret string
}

func NewIdentityManager() *identityManager {
	return &identityManager{
		baseUrl:             viper.GetString("KeyCloak.BaseUrl"),
		realm:               viper.GetString("KeyCloak.Realm"),
		restApiClientId:     viper.GetString("KeyCloak.RestApi.ClientId"),
		restApiClientSecret: viper.GetString("KeyCloak.RestApi.ClientSecret"),
	}
}

func (im *identityManager) loginRestApiClient(ctx context.Context) (*gocloak.JWT, error) {
	client := gocloak.NewClient(im.baseUrl)

	token, err := client.LoginClient(ctx, im.restApiClientId, im.restApiClientSecret, im.realm)
	if err != nil {
		return nil, errors.Wrap(err, "unable to login the rest client")
	}
	return token, nil
}

func (im *identityManager) loginRestApiUser(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	client := gocloak.NewClient(im.baseUrl)

	token, err := client.Login(ctx, im.restApiClientId, im.restApiClientSecret, im.realm, username, password)
	if err != nil {
		return nil, errors.Wrap(err, "unable to login rest user")
	}

	return token, nil
}

func (im *identityManager) CreateUser(ctx context.Context, user gocloak.User, password string, role string) (*gocloak.User, error) {

	token, err := im.loginRestApiClient(ctx)
	if err != nil {
		return nil, err
	}

	client := gocloak.NewClient(im.baseUrl)

	userId, err := client.CreateUser(ctx, token.AccessToken, im.realm, user)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create the user")
	}

	err = client.SetPassword(ctx, token.AccessToken, userId, im.realm, password, false)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set the password for the user")
	}

	var roleNameLowerCase = strings.ToLower(role)
	roleKeycloak, err := client.GetRealmRole(ctx, token.AccessToken, im.realm, roleNameLowerCase)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to get role by name: '%v'", roleNameLowerCase))
	}
	err = client.AddRealmRoleToUser(ctx, token.AccessToken, im.realm, userId, []gocloak.Role{
		*roleKeycloak,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to add a realm role to user")
	}

	userKeycloak, err := client.GetUserByID(ctx, token.AccessToken, im.realm, userId)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get recently created user")
	}

	return userKeycloak, nil
}

func (im *identityManager) LoginUser(ctx context.Context, username, password string) (*gocloak.User, *gocloak.JWT, error) {
	client := gocloak.NewClient(im.baseUrl)

	userToken, err := im.loginRestApiUser(ctx, username, password)
	if err != nil {
		return nil, nil, err
	}

	clientToken, err := im.loginRestApiClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	userInfo, err := client.GetUserInfo(ctx, userToken.AccessToken, im.realm)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get user information")
	}

	user, err := client.GetUserByID(ctx, clientToken.AccessToken, im.realm, *userInfo.Sub)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get user by ID")
	}

	return user, userToken, nil
}

func (im *identityManager) UpdateUserAttribute(ctx context.Context, userID string, attribute map[string][]string) error {
	client := gocloak.NewClient(im.baseUrl)

	clientToken, err := im.loginRestApiClient(ctx)
	if err != nil {
		return err
	}

	user, err := client.GetUserByID(ctx, clientToken.AccessToken, im.realm, userID)
	if err != nil {
		return errors.Wrap(err, "unable to ge user")
	}

	user.Attributes = &attribute

	err = client.UpdateUser(ctx, clientToken.AccessToken, im.realm, *user)
	if err != nil {
		return errors.Wrap(err, "unable to update user")
	}

	return nil
}

func (im *identityManager) GenerateOTP(ctx context.Context) (string, bytes.Buffer, error) {
	var (
		qrCode bytes.Buffer
	)

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

func (im *identityManager) ValidateOTP(ctx context.Context, userID, passcode string) (bool, error) {
	client := gocloak.NewClient(im.baseUrl)

	clientToken, err := im.loginRestApiClient(ctx)
	if err != nil {
		return false, err
	}

	user, err := client.GetUserByID(ctx, clientToken.AccessToken, im.realm, userID)
	if err != nil {
		return false, errors.Wrap(err, "unable to get user by ID")
	}

	var otpSecret string

	if user.Attributes == nil {
		return false, errors.New("for this user didn't generated OTP")
	} else if user.Attributes != nil {
		otpSecret = (*user.Attributes)["otp_secret"][0]
		if otpSecret == "" {
			return false, errors.New("for this user didn't generated OTP")
		}
	}

	otpSecret = (*user.Attributes)["otp_secret"][0]

	validate := totp.Validate(passcode, otpSecret)

	return validate, nil
}

func (im *identityManager) RetrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {

	client := gocloak.NewClient(im.baseUrl)

	rptResult, err := client.RetrospectToken(ctx, accessToken, im.restApiClientId, im.restApiClientSecret, im.realm)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrospect token")
	}
	return rptResult, nil
}
