package svc

import (
	"errors"
	"fmt"
	"github.com/jamesburns-rts/base-go-server/internal/call"
	"github.com/jamesburns-rts/base-go-server/internal/clients/example"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"github.com/jamesburns-rts/base-go-server/internal/db"
	"github.com/jamesburns-rts/base-go-server/internal/jwt"
	"github.com/jamesburns-rts/base-go-server/internal/vm"
	"github.com/jamesburns-rts/base-go-server/internal/vm/page"
	"github.com/jamesburns-rts/base-go-server/internal/vm/validator"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Client         *example.Client
	TokenGenerator func(params jwt.AccessTokenParams) (string, error)
}

func NewSampleService(props config.Application) *UserService {
	generator, err := jwt.NewTokenGenerator(props.JWT)
	if err != nil {
		panic(err) // todo handle better?
	}
	return &UserService{
		Client:         example.NewClient(props.ExampleClient),
		TokenGenerator: generator,
	}
}

func (s *UserService) Authenticate(ctx call.Context, login vm.Login) (*vm.Tokens, error) {
	if err := validator.Validate(login); err != nil {
		return nil, err // todo mapping
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u, err := db.FindUserByEmail(ctx, login.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if u == nil || u.PasswordHash != string(hashed) {
		return nil, errors.New("unauthorized") // todo
	}

	token, err := s.TokenGenerator(jwt.AccessTokenParams{
		UserID: u.ID,
	})
	if err != nil {
		return nil, err
	}
	return &vm.Tokens{AccessToken: token}, nil
}

func (s *UserService) GetUserExamplesPage(ctx call.Context, userID int, pageParams page.Request) ([]vm.Example, error) {
	// todo pagination

	// a contrived input validation example
	if userID < 0 {
		return nil, InvalidInputError{Message: "user id must be greater than zero"}
	}

	userExamples, err := db.GetExamplesByUser(ctx, userID, pageParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get examples by user: %w", err)
	}

	exampleVMs := make([]vm.Example, 0, len(userExamples))
	for _, e := range userExamples {
		p, err := s.Client.GetExample(ctx, e.ExampleName)
		if err != nil {
			return nil, fmt.Errorf("failed to get example %s: %w", e.ExampleName, err)
		}
		if p == nil {
			log.Warn("user has example that does not exist")
			continue
		}
		exampleVMs = append(exampleVMs, mapPokemonToExampleVM(*p))
	}
	return exampleVMs, nil
}

func mapPokemonToExampleVM(p example.Pokemon) vm.Example {
	moves := make([]string, len(p.Moves))
	for i, m := range p.Moves {
		moves[i] = m.Move.Name
	}
	return vm.Example{
		ID:    p.ID,
		Name:  p.Name,
		Moves: moves,
	}
}
