package utils

import (
	"gopkg.in/yaml.v2"
	"strings"
)

func ToYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}
