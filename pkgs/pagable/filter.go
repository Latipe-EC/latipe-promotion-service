package pagable

import (
	"github.com/gofiber/fiber/v2/log"
	"net/url"
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

func decodeFilterURL(encodedUrl string) (string, error) {
	decodedUrl, err := url.QueryUnescape(encodedUrl)
	if err != nil {
		return "", err
	}

	log.Info("url:%v", decodedUrl)
	return decodedUrl, nil
}

func FilterBinding(uri string) ([]Filter, error) {
	urlDecode, err := decodeFilterURL(uri)
	if err != nil {
		return nil, err
	}

	var filters []Filter
	// Find all matches in the uri
	matches := filterRegex.FindAllStringSubmatch(urlDecode, -1)
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
