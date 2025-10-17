package cmd

import (
	"bytes"
	"context"
	"pokemon/internal/pokemon"
	"testing"

	svc "pokemon/internal/service/pokemon"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SearchPokemonCmdTestSuite struct {
	suite.Suite
	mockService *svc.MockService
}

func (s *SearchPokemonCmdTestSuite) SetupTest() {
	s.mockService = &svc.MockService{}
}

func (s *SearchPokemonCmdTestSuite) Test_newSearchPokemonCmd() {
	got := new(bytes.Buffer)
	ac := &pokemon.Context{PokemonService: s.mockService}
	ac.EC = cmd.NewExecutionContext("search", "pokemon", "Search for pokemons", nil, got, got)
	c := newSearchPokemonCmd(ac)

	s.Run("command should not be nil", func() {
		s.NotNil(c)
	})

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
	}{
		{
			name:           "name flag",
			args:           []string{"--name", "pikachu"},
			expectedOutput: "Searching for pokemon \"pikachu\"",
		},
		{
			name:           "short description",
			args:           []string{"--help"},
			expectedOutput: "Searches for pokemons",
		},
		{
			name:           "long description",
			args:           []string{"--help"},
			expectedOutput: "Searches for pokemons using pokeapi.co",
		},
		{
			name:           "example",
			args:           []string{"--help"},
			expectedOutput: "  # Search for a pokemon by name\n  $ pokemon search pokemon --name pikachu",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			c.SetArgs(tc.args)
			c.SetOut(got)

			s.mockService.On("SearchPokemon", mock.Anything).Return(&svc.Pokemon{}, nil)
			s.mockService.On("GetPokemons").Return(svc.Pokemons{}, nil)

			err := c.ExecuteContext(context.Background())
			s.NoError(err)

			s.Contains(got.String(), tc.expectedOutput)
		})
	}
}

func TestSearchPokemonCmdTestSuite(t *testing.T) {
	suite.Run(t, new(SearchPokemonCmdTestSuite))
}
