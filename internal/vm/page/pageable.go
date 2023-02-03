package page

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type (
	// Request input params for pagination
	Request struct {
		PageSize   int    `query:"size"`
		PageNumber int    `query:"page"`
		Sort       []Sort `query:"sort"`
	}

	Sort struct {
		Property  string
		Direction string
	}

	// Response a paginated container
	Response[T any] struct {
		Content          []T  `json:"content"`
		TotalPages       int  `json:"totalPages"`
		TotalElements    int  `json:"totalElements"`
		NumberOfElements int  `json:"numberOfElements"`
		First            bool `json:"first"`
		Last             bool `json:"last"`
		Size             int  `json:"size"`
		Number           int  `json:"number"`
	}
)

const Ascending = "ASC"
const Descending = "DESC"

// CreateResponse fill in meta data fields from min data
func CreateResponse[T any](content []T, pageSize, pageNumber, totalElements int) Response[T] {

	var totalPages int
	if totalElements > 0 {
		totalPages = int(math.Ceil(float64(totalElements) / float64(pageSize)))
	}

	if content == nil {
		content = make([]T, 0)
	}

	return Response[T]{
		Content:          content,
		TotalPages:       totalPages,
		TotalElements:    totalElements,
		NumberOfElements: len(content),
		First:            pageNumber == 0,
		Last:             pageNumber == (totalPages-1) || totalPages == 0,
		Size:             pageSize,
		Number:           pageNumber,
	}
}

func RequestFromParams(size, number, sort string) (Request, error) {
	params := Request{
		PageSize:   20,
		PageNumber: 0,
		Sort:       nil,
	}
	if size != "" {
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			return params, fmt.Errorf("invalid integer for size: %w", err)
		}
		if sizeInt <= 0 {
			return params, errors.New("size must be greater than zero")
		}
		params.PageSize = sizeInt
	}
	if number != "" {
		numberInt, err := strconv.Atoi(number)
		if err != nil {
			return params, fmt.Errorf("invalid integer for number: %w", err)
		}
		if numberInt < 0 {
			return params, errors.New("number must be greater than or equal to zero")
		}
		params.PageNumber = numberInt
	}

	params.Sort = SortSliceFromString(sort)

	return params, nil
}

// SortSliceFromString converts a comma-separated list of properties (with optional directions) to slice of Sort
func SortSliceFromString(str string) []Sort {
	sort := make([]Sort, 0)
	for _, sortStr := range strings.Split(str, ",") {
		sortProperty := strings.TrimSpace(sortStr)
		sortDir := ""

		lastSpaceIndex := strings.LastIndex(sortProperty, " ")
		if lastSpaceIndex > 0 {
			lastField := sortProperty[lastSpaceIndex+1:]
			if strings.ToUpper(lastField) == Ascending || strings.ToUpper(lastField) == Descending {
				sortProperty = sortProperty[:lastSpaceIndex]
				sortDir = lastField
			}
		}

		sort = append(sort, Sort{
			Property:  sortProperty,
			Direction: sortDir,
		})
	}
	return sort
}

// SortSliceToString converts sort slice to string that can be used as an "order by" statement
func SortSliceToString(sort []Sort) string {
	sb := strings.Builder{}
	for i, s := range sort {
		sb.WriteString(s.Property)
		if s.Direction != "" {
			sb.WriteRune(' ')
			sb.WriteString(s.Direction)
		}
		if i != len(sort)-1 {
			sb.WriteRune(',')
		}
	}
	return sb.String()
}
