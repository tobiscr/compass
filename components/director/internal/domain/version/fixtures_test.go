package version_test

import (
	"github.com/kyma-incubator/compass/components/director/internal/domain/version"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/internal/repo"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

func fixModelVersion(value string, deprecated bool, deprecatedSince string, forRemoval bool) *model.Version {
	return &model.Version{
		Value:           value,
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}
}

func fixGQLVersion(value string, deprecated bool, deprecatedSince string, forRemoval bool) *externalschema.Version {
	return &externalschema.Version{
		Value:           value,
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}
}

func fixModelVersionInput(value string, deprecated bool, deprecatedSince string, forRemoval bool) *model.VersionInput {
	return &model.VersionInput{
		Value:           value,
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}
}

func fixGQLVersionInput(value string, deprecated bool, deprecatedSince string, forRemoval bool) *externalschema.VersionInput {
	return &externalschema.VersionInput{
		Value:           value,
		Deprecated:      &deprecated,
		DeprecatedSince: &deprecatedSince,
		ForRemoval:      &forRemoval,
	}
}

func fixVersionEntity(value string, deprecated bool, deprecatedSince string, forRemoval bool) *version.Version {
	return &version.Version{
		VersionValue:           repo.NewNullableString(&value),
		VersionDepracated:      repo.NewNullableBool(&deprecated),
		VersionDepracatedSince: repo.NewNullableString(&deprecatedSince),
		VersionForRemoval:      repo.NewNullableBool(&forRemoval),
	}
}
