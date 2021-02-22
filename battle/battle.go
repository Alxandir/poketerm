package battle

import (
	"fmt"
	"strconv"

	"github.com/alxandir/poketerm/player"
	"github.com/alxandir/poketerm/pokemon"
	"github.com/alxandir/poketerm/term"
)

type battle struct {
	prize                 float64
	plyr                  *player.Player
	activePlayerPokemon   *pokemon.Pokemon
	activeOpponentPokemon *pokemon.Pokemon
	opponentPokemon       []*pokemon.Pokemon
	remainingPokemon      uint
	complete              bool
	playerWon             bool
}

func New(prize float64, plyr *player.Player, opponentPokemon []*pokemon.Pokemon) battle {
	activeOpponentPokemon := opponentPokemon[0]
	activePlayerPokemon := plyr.GetPokemon()[0]
	remainingPokemon := uint(len(opponentPokemon))
	complete := false
	playerWon := false
	b := battle{prize, plyr, activePlayerPokemon, activeOpponentPokemon, opponentPokemon, remainingPokemon, complete, playerWon}
	return b
}

func (b battle) Perform() {
	term.ShowNoResponseDialog("\n\nBattle started!")
	b.opponentThrow()
	b.playerThrow()
	b.DisplayCurrentBattleData()
	for !b.complete {
		b.PlayRound()
	}
	b.Conclude()
}

func (b battle) opponentThrow() {
	term.ShowNoResponseDialog("\tThe opponent sent out %v", b.activeOpponentPokemon.GetName())
}

func (b battle) playerThrow() {
	term.ShowNoResponseDialog("\t%v go get 'em!", b.activePlayerPokemon.GetName())
}

func (b battle) DisplayCurrentBattleData() {
	str := fmt.Sprintf("\n\t(Opponent) %v\n\t\tLevel: %v\n\t\tHP: %v", b.activeOpponentPokemon.GetName(), b.activeOpponentPokemon.GetLevel(), b.activeOpponentPokemon.GetHP())
	str += fmt.Sprintf("\n\t(%v) %v\n\t\tLevel: %v\n\t\tHP: %v", b.plyr.GetName(), b.activePlayerPokemon.GetName(), b.activePlayerPokemon.GetLevel(), b.activePlayerPokemon.GetHP())
	term.ShowNoResponseDialog(str)
}

func (b battle) GetBattleHeader() string {
	str := fmt.Sprintf("\n\t\t\t\t\t\t%v (lvl: %v) %v/%v", b.activeOpponentPokemon.GetName(), b.activeOpponentPokemon.GetLevel(), b.activeOpponentPokemon.GetHP(), b.activeOpponentPokemon.GetMaxHP())
	str += fmt.Sprintf("\n%v (lvl: %v) %v/%v\n", b.activePlayerPokemon.GetName(), b.activePlayerPokemon.GetLevel(), b.activePlayerPokemon.GetHP(), b.activePlayerPokemon.GetMaxHP())
	return str
}

func (b battle) GetAttackMenu() string {
	str := b.GetBattleHeader() + b.activePlayerPokemon.GetAttacksString()
	str += "\nSelect an attack:"
	return str
}

func (b battle) DisplayRoundData() string {
	str := b.GetAttackMenu()
	attacks := b.activePlayerPokemon.GetAttacks()
	return term.ShowInputDialogValidated(func(userInput string) string {
		err := term.ValidateNumericChoice(1, len(attacks))(userInput)
		if len(err) > 0 {
			return err
		}
		index, _ := strconv.Atoi(userInput)
		if attacks[index-1].GetPP() == 0 {
			return "Attack has no PP left"
		}
		return ""
	}, str)
}

func (b *battle) PlayRound() (battleOver bool, playerWon bool) {
	battleOver = false
	attacks := b.activePlayerPokemon.GetAttacks()
	response := b.DisplayRoundData()
	index, _ := strconv.Atoi(response)
	attack := attacks[index-1]
	oppenentAttack := chooseAttack(*b.activeOpponentPokemon, *b.activePlayerPokemon)

	b.activePlayerPokemon.UseAttack(index - 1)
	firstAttacker := b.activePlayerPokemon
	secondAttacker := b.activeOpponentPokemon
	firstAttack := attack.GetAttack()
	secondAttack := oppenentAttack
	if secondAttacker.GetAdjustedSpeed() > firstAttacker.GetAdjustedSpeed() {
		firstAttacker = secondAttacker
		firstAttack = secondAttack
		secondAttack = attack.GetAttack()
		secondAttacker = b.activePlayerPokemon
	}

	b.performAttackStep(firstAttacker, secondAttacker, firstAttack)
	if b.complete {
		return
	}
	b.performAttackStep(secondAttacker, firstAttacker, secondAttack)
	if b.complete {
		return
	}

	return
}

func (b *battle) performAttackStep(attacker *pokemon.Pokemon, defender *pokemon.Pokemon, attack pokemon.Attack) {
	hit, effectiveness, stageModifications := defender.ReceiveAttack(attacker, attack)
	displayAttackResult(attacker.GetName(), defender.GetName(), attack.GetName(), hit, effectiveness, stageModifications)
	if defender.GetHP() <= 0 {
		b.evaluateFaint(defender, attacker)
		if b.complete {
			return
		}
	}
	if attacker.GetHP() <= 0 {
		b.evaluateFaint(attacker, defender)
		if b.complete {
			return
		}
	}
	return
}

func (b *battle) evaluateFaint(faintedPokemon *pokemon.Pokemon, victor *pokemon.Pokemon) {
	playerVictorious := victor == b.activePlayerPokemon
	if playerVictorious {
		term.ShowNoResponseDialog("Opponent's %v fainted", faintedPokemon.GetName())
		if victor.GetHP() > 0 {
			victor.LevelUp(true)
		}
		nextOpponent, hasValidPokemon := b.GetNextOpponentPokemon()
		if !hasValidPokemon {
			term.ShowNoResponseDialog("The opponent has no more Pokemon")
			b.complete = true
			b.playerWon = true
			return
		}
		b.activeOpponentPokemon = nextOpponent
		b.opponentThrow()
		b.DisplayCurrentBattleData()
	} else {
		term.ShowNoResponseDialog("%v fainted", faintedPokemon.GetName())
		nextPlayerPokemon, hasValidPokemon := b.GetNextPokemon()
		if !hasValidPokemon {
			term.ShowNoResponseDialog("You have no more Pokemon")
			b.complete = true
			b.playerWon = false
			return
		}
		b.activePlayerPokemon = nextPlayerPokemon
		b.playerThrow()
		b.DisplayCurrentBattleData()
	}
	return
}

func (b *battle) GetNextOpponentPokemon() (*pokemon.Pokemon, bool) {
	validPokemon := []*pokemon.Pokemon{}
	for _, opponentPokemon := range b.opponentPokemon {
		if opponentPokemon.GetHP() > 0 {
			validPokemon = append(validPokemon, opponentPokemon)
		}
	}
	b.remainingPokemon = uint(len(validPokemon))
	if b.remainingPokemon > 0 {
		return validPokemon[0], true
	}
	return nil, false
}

func (b *battle) GetNextPokemon() (*pokemon.Pokemon, bool) {
	validPokemon := []*pokemon.Pokemon{}
	for _, opponentPokemon := range b.plyr.GetPokemon() {
		if opponentPokemon.GetHP() > 0 {
			validPokemon = append(validPokemon, opponentPokemon)
		}
	}
	b.remainingPokemon = uint(len(validPokemon))
	if b.remainingPokemon > 0 {
		return validPokemon[0], true
	}
	return nil, false
}

func displayAttackResult(attackerName string, targetName string, attackName string, hit bool, effectiveness float64, stageModifications []string) {
	term.ShowNoResponseDialog("\n\t%v used %v...", attackerName, attackName)
	if hit {
		switch {
		case effectiveness > 1.0:
			term.ShowNoResponseDialog("\tIt was super effective!")
		case effectiveness == 0.0:
			term.ShowNoResponseDialog("\tIt had no effect")
		case effectiveness < 1.0:
			term.ShowNoResponseDialog("\tIt was not very effective")
		default:
			term.ShowNoResponseDialog("\t%v was hit with %v", targetName, attackName)
		}
	}
	if !hit {
		term.ShowNoResponseDialog("\tIt missed")
	}
	for _, stageMstageModification := range stageModifications {
		term.ShowNoResponseDialog(stageMstageModification)
	}
}

func (b battle) Conclude() {
	if b.playerWon {
		term.ShowNoResponseDialog("\n\nCongratulations on winning the batte!")
		for _, pokemon := range b.plyr.GetPokemon() {
			pokemon.EvaluateEvolution()
		}
		b.plyr.AddMoney(b.prize)
		term.ShowNoResponseDialog("\n\n£%.2f was added to your Monzo!", b.prize)
	} else {
		term.ShowNoResponseDialog("\n\nLooks like today wasn't your day %v", b.plyr.GetName())
	}
}

type attackChance struct {
	index  int
	attack pokemon.Attack
	chance float64
}

func chooseAttack(attacker pokemon.Pokemon, defender pokemon.Pokemon) pokemon.Attack {
	potentialAttacks := []attackChance{}
	totalChances := 0.0
	for i, pAttack := range attacker.GetAttacks() {
		if pAttack.GetPP() > 0 {
			attack := pAttack.GetAttack()
			chance := defender.GetAttackEffectiveness(attack) * attacker.GetAttackTypeBonus(attack) * float64(attack.GetPower())
			totalChances += chance
			potentialAttacks = append(potentialAttacks, attackChance{
				index:  i,
				attack: attack,
				chance: chance,
			})
		}
	}
	lastUpperBound := 0.0
	r := float64(pokemon.RandomNumber(0, 100))
	for _, potentialAttack := range potentialAttacks {
		upperBound := lastUpperBound + ((potentialAttack.chance / totalChances) * 100.0)
		if r >= lastUpperBound && r < upperBound {
			attacker.UseAttack(potentialAttack.index)
			return potentialAttack.attack
		}
		lastUpperBound = upperBound
	}
	return pokemon.Struggle
}
