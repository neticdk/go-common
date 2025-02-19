package pokemon

import (
	"log/slog"

	"github.com/mtslzr/pokeapi-go"
	"github.com/mtslzr/pokeapi-go/structs"
)

type Pokemons []*Pokemon

type Pokemon struct {
	Name   string
	URL    string
	Types  []string
	Weight int
	Height int
}

type Service interface {
	GetPokemons() (Pokemons, error)
	SearchPokemon(name string) (*Pokemon, error)
}

type service struct {
	logger *slog.Logger
}

func (s *service) GetPokemons() (Pokemons, error) {
	s.Log("Getting all pokemons")
	var pokemons Pokemons
	offset := 0
	limit := 200

	for {
		pmons, err := pokeapi.Resource("pokemon", offset, limit)
		if err != nil {
			return nil, err
		}
		for _, p := range pmons.Results {
			pokemons = append(pokemons, &Pokemon{Name: p.Name, URL: p.URL})
		}
		if pmons.Next == "" {
			break
		}
		offset += limit
	}
	return pokemons, nil
}

func (s *service) SearchPokemon(name string) (*Pokemon, error) {
	p, err := pokeapi.Pokemon(name)
	if err != nil {
		return nil, err
	}
	return toPokemon(p), nil
}

func (s *service) Log(msg string) {
	if s.logger != nil {
		s.logger.Debug(msg)
	}
}

func NewService(opts ...Option) Service {
	s := &service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Option func(s *service)

func WithLogger(logger *slog.Logger) Option {
	return func(s *service) {
		s.logger = logger
	}
}

func toPokemon(p structs.Pokemon) *Pokemon {
	types := make([]string, 0)
	for _, t := range p.Types {
		types = append(types, t.Type.Name)
	}
	return &Pokemon{
		Name:   p.Name,
		Types:  types,
		Weight: p.Weight,
		Height: p.Height,
	}
}
