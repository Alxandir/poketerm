package player

import (
	"github.com/alxandir/poketerm/pokemon"
	"github.com/alxandir/poketerm/term"
)

type Player struct {
	name         string
	gender       Gender
	wallet       float64
	pokemonParty []*pokemon.Pokemon
}

func New(name string, gender Gender) Player {
	pokemonParty := []*pokemon.Pokemon{}
	wallet := 50.0
	player := Player{name, gender, wallet, pokemonParty}
	return player
}

func (p *Player) AddMoney(amount float64) {
	p.wallet += amount
}

func (p *Player) removeMoney(amount float64) (canAfford bool) {
	canAfford = p.wallet >= amount
	if canAfford {
		p.wallet -= amount
	}
	return
}

func (p *Player) AddPokemonToParty(pokemn *pokemon.Pokemon) {
	p.pokemonParty = append(p.pokemonParty, pokemn)
	term.ShowNoResponseDialog("\n\t%v was added to %v's party!", pokemn.GetName(), p.GetName())
}

func (p Player) GetPokemon() []*pokemon.Pokemon {
	return p.pokemonParty
}

func (p Player) GetName() string {
	return p.name
}
func (p Player) GetGenderName() string {
	return p.gender.GetName()
}
func (p Player) GetGenderPronoun() string {
	return p.gender.GetPronoun()
}
func (p Player) GetGenderOwnershipPronoun() string {
	return p.gender.GetOwnershipPronoun()
}
