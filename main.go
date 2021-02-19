package main

import (
	"fmt"
	"strconv"

	"github.com/alxandir/poketerm/battle"
	"github.com/alxandir/poketerm/player"
	"github.com/alxandir/poketerm/pokemon"
	"github.com/alxandir/poketerm/term"
)

func main() {
	term.PrettyPrint("Pokemon.txt")
	response := term.ShowInputDialog("\nHi there! What's your name?")
	user := player.New(response, player.GenderMale{})
	starters := pokemon.GetStarterPokemon()
	choiceString := "\nI have some little critters for you to take a look at!"
	for index, starter := range starters {
		choiceString += fmt.Sprintf("\n\t(%v) %v", index+1, starter.GetName())
	}
	choiceString += "\nWhich of these pokemon takes your fancy?"
	response = term.ShowInputDialogValidated(validateStarterChoice(starters), choiceString)
	index, _ := strconv.Atoi(response)
	chosenStarter := starters[index-1]

	response = term.ShowInputDialog("\nGreat choice %v! Would you like to give %v a nickname?", user.GetName(), chosenStarter.GetName())
	if term.UserAccepted(response) {
		nickname := term.ShowInputDialog("What would you like to nickname %v?", chosenStarter.GetName())
		chosenStarter.SetNickname(nickname)
	}

	user.AddPokemonToParty(&chosenStarter)
	term.ShowNoResponseDialog("\n\t%v was added to your %v's party!", chosenStarter.GetName(), user.GetName())
	chosenStarter.Nerf()
	enemyPokemon, _ := pokemon.New(10, "", 4)
	enemyPokemon2, _ := pokemon.New(1, "", 6)
	b := battle.New(64.5, &user, []*pokemon.Pokemon{&enemyPokemon, &enemyPokemon2})
	b.Perform()
	term.ShowNoResponseDialog("\n\n%#v", user)
}

func validateStarterChoice(starters []pokemon.Pokemon) func(string) string {
	return term.ValidateNumericChoice(1, len(starters))
}
