package model

import (
	"errors"
	"fmt"

	schema "github.com/kyma-incubator/compass/components/connector/pkg/graphql/externalschema"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

const appNameLabel = "name"

type InfoProviderFunc func(application graphql.ApplicationExt, tenant string, configuration schema.Configuration) (interface{}, error)

func NewCSRInfoResponseProvider(connectivityAdapterBaseURL, connectivityAdapterMTLSBaseURL string) InfoProviderFunc {
	return func(application graphql.ApplicationExt, _ string, configuration schema.Configuration) (interface{}, error) {
		if configuration.Token == nil {
			return nil, errors.New("empty token returned from Connector")
		}

		csrURL := connectivityAdapterBaseURL + CertsEndpoint
		tokenParam := fmt.Sprintf(TokenFormat, configuration.Token.Token)

		api := Api{
			CertificatesURL: csrURL,
			InfoURL:         connectivityAdapterMTLSBaseURL + ManagementInfoEndpoint,
			RuntimeURLs:     makeRuntimeURLs(application.ID, connectivityAdapterMTLSBaseURL, application.EventingConfiguration.DefaultURL),
		}

		return CSRInfoResponse{
			CsrURL:          csrURL + tokenParam,
			API:             api,
			CertificateInfo: ToCertInfo(configuration.CertificateSigningRequestInfo),
		}, nil
	}
}

func NewManagementInfoResponseProvider(connectivityAdapterMTLSBaseURL string) InfoProviderFunc {
	return func(application graphql.ApplicationExt, tenant string, configuration schema.Configuration) (interface{}, error) {
		clientIdentity := ClientIdentity{
			Application:   retrieveAppName(application),
			ApplicationID: application.ID,
			Tenant:        tenant,
			Group:         "",
		}

		managementURLs := MgmtURLs{
			RuntimeURLs:   makeRuntimeURLs(application.ID, connectivityAdapterMTLSBaseURL, application.EventingConfiguration.DefaultURL),
			RenewCertURL:  fmt.Sprintf(RenewCertURLFormat, connectivityAdapterMTLSBaseURL),
			RevokeCertURL: fmt.Sprintf(RevocationCertURLFormat, connectivityAdapterMTLSBaseURL),
		}

		return MgmtInfoReponse{
			ClientIdentity:  clientIdentity,
			URLs:            managementURLs,
			CertificateInfo: ToCertInfo(configuration.CertificateSigningRequestInfo),
		}, nil
	}
}

func makeRuntimeURLs(appID, connectivityAdapterMTLSBaseURL string, eventServiceBaseURL string) *RuntimeURLs {
	return &RuntimeURLs{
		MetadataURL:   connectivityAdapterMTLSBaseURL + fmt.Sprintf(ApplicationRegistryEndpointFormat, appID),
		EventsURL:     eventServiceBaseURL,
		EventsInfoURL: eventServiceBaseURL + EventsInfoEndpoint,
	}
}

func retrieveAppName(application graphql.ApplicationExt) string {
	appName := application.Name

	labels := map[string]interface{}(application.Labels)
	if labels != nil || len(labels) != 0 {
		labelValue, exists := labels[appNameLabel]
		if exists {
			name, ok := labelValue.(string)
			if ok {
				appName = name
			}
		}
	}

	return appName
}
