package service

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/pkg/errors"
)

const legacyServicesLabelKey = "legacy_servicesMetadata"

type LegacyServiceReference struct {
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
}

type labeler struct{}

func NewAppLabeler() *labeler {
	return &labeler{}
}

func (l *labeler) WriteServiceReference(appLabels externalschema.Labels, serviceReference LegacyServiceReference) (externalschema.LabelInput, error) {
	services, err := l.readLabel(appLabels)
	if err != nil {
		return externalschema.LabelInput{}, err
	}

	services[serviceReference.ID] = serviceReference

	return l.writeLabel(services)
}

func (l *labeler) ReadServiceReference(appLabels externalschema.Labels, serviceID string) (LegacyServiceReference, error) {
	services, err := l.readLabel(appLabels)
	if err != nil {
		return LegacyServiceReference{}, err
	}

	service, exists := services[serviceID]
	if !exists {
		return LegacyServiceReference{}, nil
	}

	return service, nil
}

func (l *labeler) DeleteServiceReference(appLabels externalschema.Labels, serviceID string) (externalschema.LabelInput, error) {
	services, err := l.readLabel(appLabels)
	if err != nil {
		return externalschema.LabelInput{}, err
	}

	delete(services, serviceID)

	return l.writeLabel(services)
}

func (l *labeler) ListServiceReferences(appLabels externalschema.Labels) ([]LegacyServiceReference, error) {
	services, err := l.readLabel(appLabels)
	if err != nil {
		return nil, err
	}
	var serviceReferences []LegacyServiceReference
	for serviceIDKey, _ := range services {
		value, err := l.ReadServiceReference(appLabels, serviceIDKey)
		if err != nil {
			return nil, err
		}
		serviceReferences = append(serviceReferences, value)
	}
	return serviceReferences, nil
}

func (l *labeler) readLabel(appLabels externalschema.Labels) (map[string]LegacyServiceReference, error) {
	value := appLabels[legacyServicesLabelKey]
	if value == nil {
		value = "{}"
	}

	strValue, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected: string; actual: %T", value)
	}

	var services map[string]LegacyServiceReference

	err := json.Unmarshal([]byte(strValue), &services)
	if err != nil {
		return nil, errors.Wrap(err, "while unmarshalling JSON value")
	}

	return services, nil
}

func (l *labeler) writeLabel(services map[string]LegacyServiceReference) (externalschema.LabelInput, error) {
	marshalledServices, err := json.Marshal(services)
	if err != nil {
		return externalschema.LabelInput{}, errors.Wrap(err, "while marshalling JSON value")
	}

	return externalschema.LabelInput{
		Key:   legacyServicesLabelKey,
		Value: strconv.Quote(string(marshalledServices)),
	}, nil
}
