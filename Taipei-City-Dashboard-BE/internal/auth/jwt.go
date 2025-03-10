package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"TaipeiCityDashboardBE/global"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// custom claims
type Claims struct {
	LoginType   string       `json:"login_type"`
	AccountID   int          `json:"account_id"`
	IsAdmin     bool         `json:"is_admin"`
	Permissions []Permission `json:"permissions"`
	jwt.StandardClaims
}

var jwtSecret = []byte(global.JwtSecret)

// validate JWT
func ValidateJWT(c *gin.Context) {
	const authPrefix = "Bearer "
	token, err := getAuthFromRequest(c, authPrefix)
	// logs.FError(err.Error())
	if err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// c.Abort()
		// If there is an error in extracting the token from the request,
		// set group:public role:viewer permission and proceed to the next middleware.
		c.Set("loginType", "no login")
		c.Set("accountID", 0)
		c.Set("isAdmin", false)
		permissions := []Permission{
			{GroupID: 1, RoleID: 3},
		}
		c.Set("permissions", permissions)
		c.Next()
		return
	}

	// parse and validate token for six things:
	// validationErrorMalformed => token is malformed
	// validationErrorUnverifiable => token could not be verified because of signing problems
	// validationErrorSignatureInvalid => signature validation failed
	// validationErrorExpired => exp validation failed
	// validationErrorNotValidYet => nbf validation failed
	// validationErrorIssuedAt => iat validation failed
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret, nil
	})
	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			// Handle different validation errors and set appropriate error messages.
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				message = "token is malformed"
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				message = "token could not be verified because of signing problems"
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				message = "signature validation failed"
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "token is expired"
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				message = "token is not yet valid before sometime"
			} else {
				message = "can not handle this token"
			}
		}
		// Respond with an unauthorized status and the error message.
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": message,
		})
		c.Abort()
		return
	}

	// If the token is valid, extract claims and set them in the context.
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		c.Set("loginType", claims.LoginType)
		c.Set("accountID", claims.AccountID)
		c.Set("isAdmin", claims.IsAdmin)
		c.Set("permissions", claims.Permissions)
		c.Set("expiresAt", claims.ExpiresAt)
		c.Next()
	} else {
		// If the token claims are not valid, abort the request.
		c.Abort()
		return
	}
}

// GenerateJWT generates a JWT token using the provided information.
// It includes user type, user ID, role list, group list, and expiration details in the JWT claims.
// The token is signed using HS256 and returned as a string.
func GenerateJWT(ExpiresAt time.Time, loginType string, userId int, isAdmin bool, permissions []Permission) (string, error) {
	// Create a unique user ID for JWT
	now := time.Now()
	uid := loginType + strconv.FormatInt(int64(userId), 10)
	jwtId := uid + strconv.FormatInt(now.Unix(), 10)

	// Set JWT claims and sign
	claims := Claims{
		LoginType:   loginType,
		AccountID:   userId,
		IsAdmin:     isAdmin,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			Audience:  uid,
			ExpiresAt: ExpiresAt.Unix(),
			Id:        jwtId,
			IssuedAt:  now.Unix(),
			Issuer:    "Taipei citydashboard",
			NotBefore: now.Add(global.NotBeforeDuration).Unix(),
			Subject:   uid,
		},
	}

	// Sign the claims using JWT signing method HS256 and obtain the token string
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("generate JWT token error: %v", err)
	}

	return token, nil
}
