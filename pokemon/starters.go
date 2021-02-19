package pokemon

func GetStarterPokemon() []Pokemon {
	bulbasaur, _ := New(1, "", 5)
	squirtle, _ := New(4, "", 5)
	charmander, _ := New(7, "", 5)
	return []Pokemon{bulbasaur, squirtle, charmander}
}
