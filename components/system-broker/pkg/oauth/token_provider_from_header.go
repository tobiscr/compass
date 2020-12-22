package oauth

import (
	"context"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	httputils "github.com/kyma-incubator/compass/components/system-broker/pkg/http"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/log"
	"github.com/pkg/errors"
)

const AuthzHeader = "Authorization"

type TokenProviderFromHeader struct {
}

func NewTokenProviderFromHeader() *TokenProviderFromHeader {
	return &TokenProviderFromHeader{}
}

func (c *TokenProviderFromHeader) Name() string {
	return "TokenProviderFromHeader"
}

func (c *TokenProviderFromHeader) Matches(ctx context.Context) bool {
	if _, err := getBearerToken(ctx); err != nil {
		log.C(ctx).WithError(err).Errorf("while obtaining bearer token")
		return false
	}

	return true
}

func (c *TokenProviderFromHeader) GetAuthorizationToken(ctx context.Context) (httputils.Token, error) {
	token, err := getBearerToken(ctx)
	if err != nil {
		return httputils.Token{}, errors.Wrapf(err, "while obtaining bearer token from header %s", AuthzHeader)
	}

	tokenResponse := httputils.Token{
		AccessToken: token,
		Expiration:  0,
	}

	log.C(ctx).Info("Successfully unmarshal response oauth token for accessing Director")
	return tokenResponse, nil
}

func getBearerToken(ctx context.Context) (string, error) {
	headers, err := httputils.LoadFromContext(ctx)
	if err != nil {
		return "", errors.New("cannot read headers from context")
	}
	reqToken, ok := headers[AuthzHeader]
	if !ok {
		return "", errors.Errorf("cannot read header %s from context", AuthzHeader)
	}

	if reqToken == "" {
		return "", apperrors.NewUnauthorizedError("missing bearer token")
	}

	if !strings.HasPrefix(strings.ToLower(reqToken), "bearer ") {
		return "", apperrors.NewUnauthorizedError("invalid bearer token prefix")
	}

	return strings.TrimPrefix(reqToken, "Bearer "), nil
}
