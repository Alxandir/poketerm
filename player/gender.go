package player

import (
	"fmt"
	"strconv"

	"github.com/alxandir/poketerm/term"
)

type Gender struct {
	name             string
	pronoun          string
	ownershipPronoun string
}

func (g Gender) GetName() string {
	return g.name
}
func (g Gender) GetPronoun() string {
	return g.pronoun
}
func (g Gender) GetOwnershipPronoun() string {
	return g.ownershipPronoun
}

var GenderMale = Gender{
	name:             "boy",
	pronoun:          "he",
	ownershipPronoun: "his",
}
var GenderFemale = Gender{
	name:             "girl",
	pronoun:          "she",
	ownershipPronoun: "hers",
}
var GenderNonBinary = Gender{
	name:             "thing",
	pronoun:          "they",
	ownershipPronoun: "theirs",
}

var Genders = []Gender{
	GenderMale,
	GenderFemale,
	GenderNonBinary,
}

func DisplayGenderChoice() Gender {
	choiceString := "\nI'd like to get to know you a little better'"
	for index, gender := range Genders {
		choiceString += fmt.Sprintf("\n\t(%v) %v", index+1, gender.GetName())
	}
	choiceString += "\nWhat gender does your snowflake self identify as?"
	response := term.ShowInputDialogValidated(validateGenderChoice(), choiceString)
	index, _ := strconv.Atoi(response)
	return Genders[index-1]
}

func validateGenderChoice() func(string) string {
	return term.ValidateNumericChoice(1, len(Genders))
}
