package model

// DTO structs

type (
	SampleDTO struct {
		ID      int64
		Message string
	}
)

// functions to convert to View models

func (dto SampleDTO) ToVm() SampleVm {
	return SampleVm{
		ID:      dto.ID,
		Message: dto.Message,
	}
}
