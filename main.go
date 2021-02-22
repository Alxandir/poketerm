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
	playerName := term.ShowInputDialog("\nHi there! What's your name?")
	playerGender := player.DisplayGenderChoice()
	term.ShowNoResponseDialog("\nGreat, so you're a %v!", playerGender.GetName())
	user := player.New(playerName, playerGender)
	starters := pokemon.GetStarterPokemon()
	choiceString := "\nI have some little critters for you to take a look at!"
	for index, starter := range starters {
		choiceString += fmt.Sprintf("\n\t(%v) %v", index+1, starter.GetName())
	}
	choiceString += "\nWhich of these pokemon takes your fancy?"
	response := term.ShowInputDialogValidated(validateStarterChoice(starters), choiceString)
	index, _ := strconv.Atoi(response)
	chosenStarter := starters[index-1]
	opponentPokemon := starters[index%len(starters)]
	helperPokemon := starters[(index+1)%len(starters)]

	response = term.ShowInputDialog("\nGreat choice %v! Would you like to give %v a nickname?", user.GetName(), chosenStarter.GetName())
	if term.UserAccepted(response) {
		nickname := term.ShowInputDialog("What would you like to nickname %v?", chosenStarter.GetName())
		chosenStarter.SetNickname(nickname)
	}

	user.AddPokemonToParty(&chosenStarter)
	term.ShowNoResponseDialog("\nThis battle will be tougher than you think, so I'll have %v here help you out!", helperPokemon.GetName())
	user.AddPokemonToParty(&helperPokemon)
	b := battle.New(64.5, &user, []*pokemon.Pokemon{&opponentPokemon})
	b.Perform()
}

func validateStarterChoice(starters []pokemon.Pokemon) func(string) string {
	return term.ValidateNumericChoice(1, len(starters))
}
