package page

import "math"

type (
	Request struct {
		PageSize      int
		PageNumber    int
		SortProperty  string
		SortDirection string
	}

	ResponseMeta struct {
		TotalPages       int   `json:"totalPages"`
		TotalElements    int64 `json:"totalElements"`
		NumberOfElements int   `json:"numberOfElements"`
		First            bool  `json:"first"`
		Last             bool  `json:"last"`
		Size             int   `json:"size"`
		Number           int   `json:"number"`
	}
)

func CreateResponseMeta(NumberOfElements, pageSize, pageNumber int, totalElements int64) ResponseMeta {

	totalPages := int(math.Ceil(float64(totalElements) / float64(pageSize)))

	return ResponseMeta{
		TotalPages:       totalPages,
		TotalElements:    totalElements,
		NumberOfElements: NumberOfElements,
		First:            pageNumber == 0,
		Last:             pageNumber == (totalPages-1) || totalPages == 0,
		Size:             pageSize,
		Number:           pageNumber,
	}
}
