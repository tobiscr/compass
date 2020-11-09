package labelfilter

import (
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)

type LabelFilter struct {
	Key   string
	Query *string
}

func FromGraphQL(in *externalschema.LabelFilter) *LabelFilter {
	return &LabelFilter{
		Key:   in.Key,
		Query: in.Query,
	}
}

func MultipleFromGraphQL(in []*externalschema.LabelFilter) []*LabelFilter {
	var filters []*LabelFilter

	for _, f := range in {
		filters = append(filters, FromGraphQL(f))
	}

	return filters
}

func NewForKey(key string) *LabelFilter {
	return &LabelFilter{key, nil}
}

func NewForKeyWithQuery(key, query string) *LabelFilter {
	return &LabelFilter{key, &query}
}
