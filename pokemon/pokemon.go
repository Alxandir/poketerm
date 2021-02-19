package pokemon

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/alxandir/poketerm/term"
)

type Pokemon struct {
	nickname         string
	level            uint
	hpReduction      uint
	hp               uint
	pokedexItem      PokedexItem
	attacks          []PokemonAttack
	evolutionPending bool
	attack           uint
	defence          uint
	spAttack         uint
	spDefence        uint
	accuracyStage    int
	evasivenessStage int
}

func New(pokedexId uint, nickname string, level uint) (p Pokemon, err error) {
	err = nil
	pokedexItem, ok := FindInPokedex(pokedexId)
	if !ok {
		return p, errors.New("Borked")
	}
	lvl := level
	if lvl == 0 {
		lvl = pokedexItem.baseLevel
	}
	hp := uint(float64(lvl) * 1.5)
	attack := uint(float64(lvl) * 1.0)
	defence := uint(float64(lvl) * 1.0)
	p = Pokemon{nickname: nickname, pokedexItem: pokedexItem, level: lvl, hp: hp, evolutionPending: false, attack: attack, defence: defence, accuracyStage: 0, evasivenessStage: 0}
	p.populateAttacks()
	return
}

func (p *Pokemon) populateAttacks() {
	allAttacks := p.pokedexItem.attacks
	appropriateAttacks := []PokedexItemAttack{}
	attacks := []PokemonAttack{}
	for _, attack := range allAttacks {
		if p.level >= attack.minLevel {
			appropriateAttacks = append(appropriateAttacks, attack)
		}
	}
	for len(appropriateAttacks) > 0 && len(attacks) < 4 {
		pokedexItemAttack := removePokedexItemAttack(&appropriateAttacks, RandomNumber(0, len(appropriateAttacks)-1))
		attack, _ := NewAttack(pokedexItemAttack.attack.id)
		attacks = append(attacks, attack)
	}
	p.attacks = attacks
}

func RandomNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func (p Pokemon) GetAccuracyModifier() float64 {
	if p.accuracyStage >= 0 {
		return (3.0 + float64(p.accuracyStage)) / 3.0
	}
	return 3.0 / (3.0 + float64(p.accuracyStage))
}

func (p Pokemon) GetEvasivenesssModifier() float64 {
	if p.evasivenessStage <= 0 {
		// 3/3 4/3 5/3
		return (3.0 + float64(p.accuracyStage)) / 3.0
	}
	return 3.0 / (3.0 + float64(p.accuracyStage))
}

func removePokedexItemAttack(array *[]PokedexItemAttack, index int) (pAttack PokedexItemAttack) {
	pAttack = (*array)[index]
	*array = append((*array)[0:index], (*array)[index+1:]...)
	return
}

func (p *Pokemon) replaceAttack(index int, newAttack PokemonAttack) {
	trailingAttacks := p.attacks[index+1:]
	p.attacks = append(p.attacks[0:index], newAttack)
	p.attacks = append(p.attacks, trailingAttacks...)
	return
}

func (p Pokemon) GetPokedexItem() PokedexItem {
	return p.pokedexItem
}

func (p Pokemon) GetName() string {
	if len(p.nickname) > 0 {
		return p.nickname
	}
	return p.pokedexItem.GetName()
}

func (p Pokemon) GetSpeciesName() string {
	return p.pokedexItem.GetName()
}

func (p *Pokemon) SetNickname(str string) {
	p.nickname = str
}

func (p *Pokemon) Nerf() {
	p.attacks[0].pp = 1
}

func (p Pokemon) GetLevel() uint {
	return p.level
}

func (p Pokemon) GetHP() uint {
	return p.GetMaxHP() - p.hpReduction
}

func (p Pokemon) GetMaxHP() uint {
	return p.hp + p.pokedexItem.baseHP
}

func (p Pokemon) GetWeaknesses() map[string]float64 {
	output := make(map[string]float64)
	var ok bool
	for _, pType := range p.pokedexItem.types {
		if _, ok = output[pType.GetName()]; !ok {
			output[pType.GetName()] = 1
		}
		for key, value := range pType.Weaknesses() {
			if _, ok = output[key]; !ok {
				output[key] = value
			}
		}
	}
	return output
}

func (p Pokemon) GetSpecialAttack() uint {
	return p.spAttack + p.pokedexItem.baseSpAttack
}

func (p Pokemon) GetSpecialDefence() uint {
	return p.spDefence + p.pokedexItem.baseSpDefence
}

func (p Pokemon) GetAttack() uint {
	return p.attack + p.pokedexItem.baseAttack
}

func (p Pokemon) GetDefence() uint {
	return p.defence + p.pokedexItem.baseDefence
}

func (p *Pokemon) ReceiveAttack(attacker Pokemon, attack Attack) (bool, float64) {
	accuracy := float64(attack.GetAccuracy()) * attacker.GetAccuracyModifier() * p.GetEvasivenesssModifier()
	term.ShowNoWaitDialog("Accuracy: %.2f", accuracy)
	hit := RandomNumber(0, 255) < int(math.Round(accuracy))
	if !hit {
		return false, 1.0
	}
	A := float64(attacker.GetAttack())
	D := float64(p.GetDefence())
	if attack.special {
		A = float64(attacker.GetSpecialAttack())
		D = float64(p.GetSpecialDefence())
	}
	typeWeaknesses := p.GetWeaknesses()
	typeEffect := 1.0
	attackerTypeBonus := 1.0
	if _, ok := typeWeaknesses[attack.attackType.GetName()]; ok {
		typeEffect = typeWeaknesses[attack.attackType.GetName()]
	}
	for _, attackerType := range attacker.pokedexItem.types {
		if attackerType.GetName() == attack.attackType.GetName() {
			attackerTypeBonus = 1.5
			break
		}
	}
	damage := ((((((2.0 * float64(attacker.level)) / 5.0) + 2.0) * float64(attack.power) * (A / D)) / 50.0) + 2.0) * typeEffect * attackerTypeBonus
	term.ShowNoResponseDialog("Damage: %v", damage)
	damageInt := uint(math.Round(damage))
	if p.GetHP() <= damageInt {
		p.hpReduction = p.GetMaxHP()
	} else {
		p.hpReduction += damageInt
	}
	return true, typeEffect
}

func (p *Pokemon) LevelUp(inBattle bool) {
	if p.level >= 100 {
		return
	}
	p.level++
	p.hp += 1
	p.attack += 1
	p.defence += 1
	evolutionLevel, hasEvolution := p.pokedexItem.GetEvolutionLevel()
	levelAttacks := p.pokedexItem.GetLevelAttacks(p.level)
	readyToEvolve := hasEvolution && p.level >= evolutionLevel
	term.ShowNoResponseDialog("\n\t%v is now level %v!", p.GetName(), p.GetLevel())
	if len(levelAttacks) > 0 {
		for _, attack := range levelAttacks {
			p.LearnAttack(attack)
		}
	}
	if inBattle {
		p.evolutionPending = readyToEvolve
	} else if readyToEvolve {
		p.EvaluateEvolution()
	}
}

func (p *Pokemon) EvaluateEvolution() {
	if p.evolutionPending {
		p.evolutionPending = false
		evolutionPokedexItem, _ := p.GetPokedexItem().GetEvolution()
		oldName := p.GetName()
		response := term.ShowInputDialog("\t%v is about to evolve into %v. Is that alright?", oldName, evolutionPokedexItem.GetName())
		if term.UserAccepted(response) {
			p.Evolve()
			term.ShowNoResponseDialog("\t%v successfully evolved into %v", oldName, p.GetPokedexItem().GetName())
		} else {
			term.ShowNoResponseDialog("\t%v stopped evolving", oldName)
		}
	}
}

func (p *Pokemon) LearnAttack(newAttack Attack) {
	attack, _ := NewAttack(newAttack.id)
	if len(p.attacks) < 4 {
		p.attacks = append(p.attacks, attack)
		term.ShowNoResponseDialog("\t%v learned %v", p.GetName(), newAttack.name)
		return
	}
	response := term.ShowInputDialog("\n%v would like to learn %v, but would need to forget an old move. Should %v forget an old move?", p.GetName(), newAttack.name, p.GetName())
	learnedNewAttack := false
	if term.UserAccepted(response) {
		str := fmt.Sprintf("\nWhich move should %v forget?\n", p.GetName())
		str += p.GetAttacksSummaryString()
		str += "\n" + newAttack.GetAttackSummaryString("\t", fmt.Sprintf("(%v - Do Not Learn) ", len(p.attacks)+1))
		str += "\nSelect an option:"
		response := term.ShowInputDialogValidated(term.ValidateNumericChoice(1, len(p.attacks)+1), str)
		index, _ := strconv.Atoi(response)
		if index < len(p.attacks)+1 {
			learnedNewAttack = true
			forgottenAttack := p.attacks[index-1]
			p.replaceAttack(index-1, attack)
			term.ShowNoResponseDialog("\t%v forgot %v...", p.GetName(), forgottenAttack.attack.name)
			term.ShowNoResponseDialog("\tAnd learned %v!", newAttack.name)
		}
	}
	if !learnedNewAttack {
		term.ShowNoResponseDialog("\t%v did not learn %v", p.GetName(), newAttack.name)
	}
}

func (p Pokemon) GetAttacksSummaryString() string {
	str := ""
	for i, pAttack := range p.attacks {
		if len(str) > 0 {
			str += "\n"
		}
		str += pAttack.GetAttackSummaryString("\t", fmt.Sprintf("(%v) ", i+1))
		// str += fmt.Sprintf("\t(%v) %v\n\t\tType: %v\n\t\tPP: %v", i+1, pAttack.attack.name, pAttack.attack.attackType.GetName(), pAttack.basePP)
	}
	return str
}

func (p Pokemon) GetAttacksString() string {
	str := ""
	for i, pAttack := range p.attacks {
		if i%2 == 0 {
			str += "\n\t"
		} else {
			str += "\t\t\t"
		}
		str += pAttack.GetAttackString(fmt.Sprintf("(%v) ", i+1))
	}
	return str
}

func (p Pokemon) GetAttacks() []PokemonAttack {
	return p.attacks
}

func (p *Pokemon) UseAttack(attackIndex int) {
	p.attacks[attackIndex].pp--
	return
}

func (p *Pokemon) Evolve() (ok bool, err error) {
	newPokedexItem, ok := p.pokedexItem.GetEvolution()
	if !ok {
		return ok, errors.New("Invalid evolution")
	}
	p.pokedexItem = newPokedexItem
	return true, nil
}
