package middlewares

import (
	"falconapi/shared/jwt"
	"github.com/gin-gonic/gin"
	golangJwt "github.com/golang-jwt/jwt/v5"
	"net/http"
)

func NewRequiresRealmRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsX := c.Value("key_claims")
		claims := claimsX.(golangJwt.MapClaims)
		jwtHelper := jwt.NewJwtHelper(claims)
		if !jwtHelper.IsUserInRealmRole(role) {
			c.JSON(http.StatusUnauthorized, "role authorization failed")
			return
		}
		c.Next()
	}
}

func NewRequiresRealmRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsX := c.Value("key_claims")
		claims := claimsX.(golangJwt.MapClaims)
		jwtHelper := jwt.NewJwtHelper(claims)

		for _, role := range roles {
			if !jwtHelper.IsUserInRealmRole(role) {
				c.JSON(http.StatusUnauthorized, "roles authorization failed")
				return
			}
		}

		c.Next()
	}
}

func NewRequiresRealmRolesWithLogins(roles, logins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsX := c.Value("key_claims")
		claims := claimsX.(golangJwt.MapClaims)
		jwtHelper := jwt.NewJwtHelper(claims)

		for _, role := range roles {
			if !jwtHelper.IsUserInRealmRole(role) {
				c.JSON(http.StatusUnauthorized, "roles authorization failed")
				return

			}
		}

		for _, login := range logins {
			if !jwtHelper.TokenHasScope(login) {
				c.JSON(http.StatusUnauthorized, "logins authorization failed")
				return
			}
		}

		c.Next()
	}
}
