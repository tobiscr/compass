package authentication

import (
	"context"

	"github.com/kyma-incubator/compass/components/connector/internal/apperrors"

	log "github.com/sirupsen/logrus"
)

//go:generate mockery -name=Authenticator
type Authenticator interface {
	AuthenticateToken(context context.Context) (string, error)
	Authenticate(context context.Context) (string, error)
	AuthenticateCertificate(context context.Context) (string, string, error)
}

func NewAuthenticator() Authenticator {
	return &authenticator{}
}

type authenticator struct {
}

func (a *authenticator) Authenticate(context context.Context) (string, error) {
	clientId, tokenAuthErr := a.AuthenticateToken(context)
	if tokenAuthErr == nil {
		log.Debugf("Client with id %s successfully authenticated with token", clientId)
		return clientId, nil
	}

	clientId, _, certAuthErr := a.AuthenticateCertificate(context)
	if certAuthErr != nil {
		return "", apperrors.NotAuthenticated("Failed to authenticate request. Token authentication error: %s. Certificate authentication error: %s",
			tokenAuthErr.Error(), certAuthErr.Error())
	}

	log.Debugf("Client with id %s successfully authenticated with certificate", clientId)
	return clientId, nil
}

func (a *authenticator) AuthenticateToken(context context.Context) (string, error) {
	clientId, err := GetStringFromContext(context, ClientIdFromTokenKey)
	if err != nil {
		return "", err.Append("Failed to authenticate request, token not provided")
	}

	if clientId == "" {
		return "", apperrors.NotAuthenticated("Failed to authenticate with one time token.")
	}

	return clientId, nil
}

func (a *authenticator) AuthenticateCertificate(context context.Context) (string, string, error) {
	clientId, err := GetStringFromContext(context, ClientIdFromCertificateKey)
	if err != nil {
		return "", "", err.Append("Failed to authenticate with Certificate. Invalid subject.")
	}

	if clientId == "" {
		return "", "", apperrors.NotAuthenticated("Failed to authenticate with Certificate. Invalid subject.")
	}

	certificateHash, err := GetStringFromContext(context, ClientCertificateHashKey)
	if err != nil {
		return "", "", err.Append("Failed to authenticate with Certificate. Invalid certificate hash.")
	}

	return clientId, certificateHash, nil
}
