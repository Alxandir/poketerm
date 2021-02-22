package pokemon

import (
	"errors"
	"fmt"
)

const (
	STAGE_MODIFIER_ATTACK      = "attack"
	STAGE_MODIFIER_DEFENCE     = "defence"
	STAGE_MODIFIER_SP_ATTACK   = "special attack"
	STAGE_MODIFIER_SP_DEFENCE  = "special defence"
	STAGE_MODIFIER_SPEED       = "speed"
	STAGE_MODIFIER_ACCURACY    = "accuracy"
	STAGE_MODIFIER_EVASIVENESS = "evasiveness"
)

type Attack struct {
	id         uint
	name       string
	attackType ElemenalType
	basePP     uint
	power      int
	special    bool
	accuracy   uint
	modifiers  []AttackStageModifier
}

type PokemonAttack struct {
	attack Attack
	pp     uint
	basePP uint
}

type AttackStageModifier struct {
	stageType    string
	modifyAmount int
	toSelf       bool
	chance       uint
}

func (pAttack PokemonAttack) GetName() string {
	return pAttack.attack.name
}

func (pAttack PokemonAttack) GetAttackSummaryString(spacer string, prefix string) (output string) {
	output = ""
	output += spacer
	output += prefix
	output += pAttack.attack.name + "\n"
	output += spacer
	output += "\tType: " + pAttack.attack.attackType.GetName() + fmt.Sprintf("\tPP: %v\n", pAttack.basePP)
	output += spacer
	output += fmt.Sprintf("\tPower: %v\tAccuracy: %v", pAttack.attack.power, pAttack.attack.accuracy)
	output += "%%"
	return
}

func (attack Attack) GetAttackSummaryString(spacer string, prefix string) (output string) {
	output = ""
	output += spacer
	output += prefix
	output += attack.name + "\n"
	output += spacer
	output += "\tType: " + attack.attackType.GetName() + fmt.Sprintf("\tPP: %v\n", attack.basePP)
	output += spacer
	output += fmt.Sprintf("\tPower: %v\tAccuracy: %v", attack.power, attack.accuracy)
	output += "%%"
	return
}

func (attack Attack) GetAccuracyValue() uint {
	return 255 * attack.accuracy
}

func (pAttack PokemonAttack) GetAttackString(prefix string) (output string) {
	output = ""
	output += prefix
	output += pAttack.attack.name
	output += fmt.Sprintf(" (%v) %v/%v", pAttack.attack.attackType.GetName(), pAttack.pp, pAttack.basePP)
	return
}

func (pAttack PokemonAttack) GetPP() uint {
	return pAttack.pp
}

func (pAttack PokemonAttack) GetAttack() Attack {
	return pAttack.attack
}

func NewAttack(attackId uint) (pAttack PokemonAttack, err error) {
	err = nil
	attack, ok := FindAttack(attackId)
	if !ok {
		err = errors.New("Borked")
	}
	pAttack = PokemonAttack{attack: attack, pp: attack.basePP, basePP: attack.basePP}
	return
}

func FindAttack(id uint) (Attack, bool) {
	for _, item := range Attacks {
		if item.id == id {
			return item, true
		}
	}
	return Attack{}, false
}

var BasicTackle = PokedexItemAttack{
	attack:   Attacks[0],
	minLevel: 1,
}
var BasicSandAttack = PokedexItemAttack{
	attack:   Attacks[4],
	minLevel: 1,
}

var Attacks = []Attack{
	{
		id:         1,
		name:       "Tackle",
		attackType: Normal{},
		basePP:     30,
		power:      30,
		accuracy:   100,
	},
	{
		id:         2,
		name:       "Water Gun",
		attackType: Water{},
		basePP:     25,
		power:      40,
		special:    true,
		accuracy:   100,
	},
	{
		id:         3,
		name:       "Ember",
		attackType: Fire{},
		basePP:     25,
		power:      40,
		special:    true,
		accuracy:   100,
	},
	{
		id:         4,
		name:       "Vine Whip",
		attackType: Grass{},
		basePP:     25,
		power:      40,
		accuracy:   100,
	},
	{
		id:         5,
		name:       "Sand Attack",
		attackType: Normal{},
		basePP:     30,
		power:      0,
		accuracy:   100,
		modifiers: []AttackStageModifier{
			{
				stageType:    STAGE_MODIFIER_ACCURACY,
				modifyAmount: -1,
				chance:       100,
			},
		},
	},
	{
		id:         6,
		name:       "String Shot",
		attackType: Bug{},
		basePP:     20,
		power:      40,
		accuracy:   100,
	},
	{
		id:         7,
		name:       "Scratch",
		attackType: Normal{},
		basePP:     15,
		power:      30,
		accuracy:   100,
	},
}
