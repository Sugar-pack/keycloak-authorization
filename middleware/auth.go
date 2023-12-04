package middleware

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"

	"test_iam/generated/swagger/models"
)

const (
	pubKeyBeginning = "-----BEGIN CERTIFICATE-----"
	pubKeyEnding    = "-----END CERTIFICATE-----"
)

type Authenticator struct {
	publicKey string
}

func NewAuthenticator(publicKey string) *Authenticator {
	return &Authenticator{
		publicKey: fmt.Sprintf("%s\n%s\n%s", pubKeyBeginning, publicKey, pubKeyEnding),
	}
}

func (a *Authenticator) Authenticate(token string) (*models.Principal, error) {
	err := a.verifyTokenSignature(token)
	if err != nil {
		return nil, err
	}

	return (*models.Principal)(&token), nil
}

func (a *Authenticator) verifyTokenSignature(token string) error {
	pub, err := jwt.ParseRSAPublicKeyFromPEM([]byte(a.publicKey))
	if err != nil {
		return fmt.Errorf("authenticator - verifyTokenSignature - jwt.ParseRSAPublicKeyFromPEM: %w", err)
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pub, nil
	})
	if err != nil {
		return fmt.Errorf("authenticator - verifyTokenSignature - jwt.Parse: %w", err)
	}
	if !parsedToken.Valid {
		return fmt.Errorf("authenticator - verifyTokenSignature: token invalid")
	}
	return nil
}
