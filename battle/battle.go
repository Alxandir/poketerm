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
}

func New(prize float64, plyr *player.Player, opponentPokemon []*pokemon.Pokemon) battle {
	activeOpponentPokemon := opponentPokemon[0]
	activePlayerPokemon := plyr.GetPokemon()[0]
	remainingPokemon := uint(len(opponentPokemon))
	b := battle{prize, plyr, activePlayerPokemon, activeOpponentPokemon, opponentPokemon, remainingPokemon}
	return b
}

func (b battle) Perform() {
	term.ShowNoResponseDialog("\n\nBattle started!")
	b.opponentThrow()
	b.playerThrow()
	b.DisplayCurrentBattleData()
	battleOver := false
	for !battleOver {
		battleOver = b.PlayRound()
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

func (b *battle) PlayRound() (battleOver bool) {
	battleOver = false
	attacks := b.activePlayerPokemon.GetAttacks()
	response := b.DisplayRoundData()
	index, _ := strconv.Atoi(response)
	attack := attacks[index-1]
	term.ShowNoResponseDialog("\t%v used %v...", b.activePlayerPokemon.GetName(), attack.GetName())
	b.activePlayerPokemon.UseAttack(index - 1)
	hit, effectiveness, stageModifications := b.activeOpponentPokemon.ReceiveAttack(b.activePlayerPokemon, attack.GetAttack())
	if hit {
		switch {
		case effectiveness > 1.0:
			term.ShowNoResponseDialog("\tIt was super effective!")
		case effectiveness < 1.0:
			term.ShowNoResponseDialog("\tIt was not very effective")
		default:
			term.ShowNoResponseDialog("\t%v was hit with %v", b.activeOpponentPokemon.GetName(), attack.GetName())
		}
	}
	if !hit {
		term.ShowNoResponseDialog("\tIt missed")
	}
	for _, stageMstageModification := range stageModifications {
		term.ShowNoResponseDialog(stageMstageModification)
	}
	if b.activeOpponentPokemon.GetHP() <= 0 {
		term.ShowNoResponseDialog("%v fainted", b.activeOpponentPokemon.GetName())
		b.activePlayerPokemon.LevelUp(true)
		nextOpponent, hasValidPokemon := b.GetNextOpponentPokemon()
		if !hasValidPokemon {
			term.ShowNoResponseDialog("There are no more pokemon")
			battleOver = true
			return
		}
		b.activeOpponentPokemon = nextOpponent
		b.opponentThrow()
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

func (b battle) Conclude() {
	term.ShowNoResponseDialog("\n\nCongratulations on winning the batte!")
	for _, pokemon := range b.plyr.GetPokemon() {
		pokemon.EvaluateEvolution()
	}
	b.plyr.AddMoney(b.prize)
	term.ShowNoResponseDialog("\n\n£%.2f was added to your Monzo!", b.prize)
}
