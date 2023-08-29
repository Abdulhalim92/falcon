package usermgmtuc

import (
	"bytes"
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type identityManager interface {
	CreateUser(ctx context.Context, user gocloak.User, password string, role string) (*gocloak.User, error)
	LoginUser(ctx context.Context, username, password string) (*gocloak.User, *gocloak.JWT, error)
	UpdateUserAttribute(ctx context.Context, userID string, attribute map[string][]string) error
	GenerateOTP(ctx context.Context) (string, bytes.Buffer, error)
	ValidateOTP(ctx context.Context, userID, passcode string) (bool, error)
}
