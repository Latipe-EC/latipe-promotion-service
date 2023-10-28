package pagable

import (
	"regexp"
)

const (
	filterPattern = "filters\\[(.*?)\\]\\[(.*?)\\]=(.*?)(?:&|\\z)"
)

// Compile the regex pattern
var filterRegex = regexp.MustCompile(filterPattern)

// ExpressionFilter is struct for add filtering in slice or array
// ref: https://docs.strapi.io/dev-docs/api/rest/filters-locale-publication#filtering
type Filter struct {
	Field     string      `json:"field"`
	Value     interface{} `json:"value"`
	Operation Operation   `json:"operation"`
}

func FilterBinding(uri string) ([]Filter, error) {
	var filters []Filter
	// Find all matches in the uri
	matches := filterRegex.FindAllStringSubmatch(uri, -1)
	for _, match := range matches {
		comp, err := OperationMapping(match[2])
		if err != nil {
			return nil, err
		}
		filter := Filter{
			Field:     match[1],
			Value:     match[3],
			Operation: comp,
		}
		filters = append(filters, filter)
	}
	return filters, nil
}
