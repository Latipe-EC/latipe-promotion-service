package pagable

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	defaultSize = 10
	maxSize     = 100
	defaultPage = 1
)

type PageableQuery struct {
	Page string `json:"page"`
	Size string `json:"size"`
}

type Query struct {
	Page              int         `json:"page"`
	Size              int         `json:"size"`
	ExpressionFilters []Filter    `json:"filters"`
	ormConditions     interface{} `json:"-"`
}

type ListResponse struct {
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	HasMore bool        `json:"has_more"`
}

// SetSize Set page size
func (q *Query) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}

	n, err := strconv.ParseUint(sizeQuery, 10, 32)
	if err != nil {
		return err
	}

	q.Size = int(n)
	if q.Size > maxSize {
		q.Size = maxSize
	}

	return nil
}

// SetPage Set page number
func (q *Query) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Page = defaultPage
		return nil
	}
	n, err := strconv.ParseUint(pageQuery, 10, 32)
	if err != nil {
		return err
	}
	q.Page = int(n)

	return nil
}

// GetOffset Get offset
func (q *Query) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit Get limit
func (q *Query) GetLimit() int {
	return q.Size
}

// GetPage Get OrderBy
func (q *Query) GetPage() int {
	return q.Page
}

// Get OrderBy
func (q *Query) GetSize() int {
	return q.Size
}

// Get total pages int
func (q *Query) GetTotalPages(totalCount int) int {
	d := float64(totalCount) / float64(q.Size)
	return int(math.Ceil(d))
}

// Get has more
func (q *Query) GetHasMore(total int) bool {
	return q.Page < total/q.Size
}

func (q *Query) ORMConditions() interface{} {
	if q.ormConditions != nil {
		return q.ormConditions
	}
	var conditions []string
	for _, filter := range q.ExpressionFilters {
		var condition string
		switch filter.Operation {
		case Equal:
			condition = filter.Field + " = " + fmt.Sprintf("'%s'", filter.Value)
		case NotEqual:
			condition = filter.Field + " <> " + fmt.Sprintf("'%s'", filter.Value)
		case LT:
			condition = filter.Field + " < " + fmt.Sprintf("'%s'", filter.Value)
		case LTE:
			condition = filter.Field + " <= " + fmt.Sprintf("'%s'", filter.Value)
		case GT:
			condition = filter.Field + " > " + fmt.Sprintf("'%s'", filter.Value)
		case GTE:
			condition = filter.Field + " >= " + fmt.Sprintf("'%s'", filter.Value)
		case In:
			condition = filter.Field + " IN " + fmt.Sprintf("'%s'", filter.Value)
		case NotIn:
			condition = filter.Field + " NOT IN " + fmt.Sprintf("'%s'", filter.Value)
		case Contains:
			condition = filter.Field + " LIKE " + "%" + fmt.Sprintf("'%s'", filter.Value) + "%"
		case NotContains:
			condition = filter.Field + " NOT LIKE " + "%" + fmt.Sprintf("'%s'", filter.Value) + "%"
		case IsNull:
			condition = filter.Field + " IS NULL"
		case IsNotNull:
			condition = filter.Field + " IS NOT NULL"
		case StartsWith:
			condition = filter.Field + " LIKE " + fmt.Sprintf(`'%s%s'`, filter.Value, "%")
		case EndsWith:
			condition = filter.Field + " LIKE " + fmt.Sprintf("'%s%s'", "%", filter.Value)
		}
		conditions = append(conditions, condition)

	}

	q.ormConditions = strings.Join(conditions, " AND ")
	return q.ormConditions
}

func (q *Query) ParseQueryParams() (map[string]string, error) {
	conditions := map[string]string{}
	for _, filter := range q.ExpressionFilters {
		if filter.Operation != Equal {
			return conditions, errors.New("only $eq filtering is supported")
		}
		conditions[filter.Field] = fmt.Sprint(filter.Value)
	}
	return conditions, nil
}
