package middlewares

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"log"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
)

type TokenRetrospector interface {
	RetrospectToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error)
}

func NewGinJwtMiddleware(tokenRetrospect TokenRetrospector) gin.HandlerFunc {
	return func(c *gin.Context) {
		successHandler(c, tokenRetrospect)
	}
}

func successHandler(c *gin.Context, tokenRetrospector TokenRetrospector) {
	userToken := c.GetHeader("user")
	if userToken == "" {
		log.Println("unable to get token")
		c.JSON(http.StatusUnauthorized, "cannot get token")
		return
	}

	base64Str := viper.GetString("KeyCloak.RealmRS256PublicKey")
	publicKey, err := parseKeycloakRSAPublicKey(base64Str)
	if err != nil {
		log.Println(err, "unable to get public key")
		c.JSON(http.StatusInternalServerError, "unable to get public key")
		return
	}

	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		log.Println(err, "unable to parse token")
		c.JSON(http.StatusInternalServerError, "unable to parse token")
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	var ctx = c.Request.Context()

	//c.Set(string(enums.ContextKeyClaims), claims)
	c.Set("key_claims", claims)

	rptResult, err := tokenRetrospector.RetrospectToken(ctx, token.Raw)
	if err != nil {
		panic(err)
	}
	if !*rptResult.Active {
		log.Println("token is not active")
		c.JSON(http.StatusUnauthorized, "token is not active")
		return
	}

	c.Next()
}

func parseKeycloakRSAPublicKey(base64Str string) (*rsa.PublicKey, error) {
	buf, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	parsedKey, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if ok {
		return publicKey, nil
	}
	return nil, fmt.Errorf("unexpected key type %T", publicKey)
}
