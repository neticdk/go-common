package pokemon

import (
	pkmn "pokemon/internal/service/pokemon"

	"github.com/neticdk/go-common/pkg/cli/cmd"
)

type Context struct {
	EC *cmd.ExecutionContext

	PokemonService pkmn.Service
}

// SetupPokemonService will setup the PokemonService in the context
// It should be called from rootCmd.PersistentPreRunE
func (ac *Context) SetupPokemonService() {
	// Makes it possible to mock the service in tests
	if ac.PokemonService != nil {
		return
	}
	ac.PokemonService = pkmn.NewService(pkmn.WithLogger(ac.EC.Logger))
}

func NewContext() *Context {
	return &Context{}
}
