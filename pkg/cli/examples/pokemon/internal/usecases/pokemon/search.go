package pokemon

import (
	"fmt"
	"strings"

	pkmn "pokemon/internal/service/pokemon"
)

func SearchPokemon(svc pkmn.Service, name string) (*pkmn.Pokemon, error) {
	return svc.SearchPokemon(name)
}

func ListPokemons(svc pkmn.Service) (pkmn.Pokemons, error) {
	return svc.GetPokemons()
}

func PrintPokemon(p *pkmn.Pokemon) {
	fmt.Printf("Name    : %s\n", p.Name)
	fmt.Printf("Types   : %s\n", strings.Join(p.Types, ", "))
	fmt.Printf("Height  : %d\n", p.Height)
	fmt.Printf("Weight  : %d\n", p.Weight)
}
