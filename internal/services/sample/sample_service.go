package sample

import (
	"github.com/jamesburns-rts/base-go-server/internal/clients/example"
	. "github.com/jamesburns-rts/base-go-server/internal/model"
)

type (
	// Service what methods can be used to access the stuff
	Service interface {
		GetSample() (sample *SampleDTO, err error)
	}

	sampleService struct {
		client example.Client
	}
)

func NewService(client example.Client) Service {
	return &sampleService{
		client: client,
	}
}
func (s sampleService) GetSample() (sample *SampleDTO, err error) {
	// do stuff with s.client
	return &SampleDTO{
		ID:      1,
		Message: "hello world",
	}, nil
}
