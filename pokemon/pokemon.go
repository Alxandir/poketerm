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
	speed            uint
	accuracyStage    int
	evasivenessStage int
	attackStage      int
	defenceStage     int
	spAttackStage    int
	spDefenceStage   int
	speedStage       int
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

func (p *Pokemon) SetAccuracyStage(newValue int) {
	p.accuracyStage = newValue
}

func (p *Pokemon) SetEvasivenessStage(newValue int) {
	p.evasivenessStage = newValue
}

func (p *Pokemon) SetAttackStage(newValue int) {
	p.attackStage = newValue
}

func (p *Pokemon) SetDefenceStage(newValue int) {
	p.defenceStage = newValue
}

func (p *Pokemon) SetSpAttackStage(newValue int) {
	p.spAttackStage = newValue
}

func (p *Pokemon) SetSpDefenceStage(newValue int) {
	p.spDefenceStage = newValue
}

func (p *Pokemon) SetSpeedStage(newValue int) {
	p.speedStage = newValue
}

func (p Pokemon) GetAccuracyModifier() float64 {
	// https://bulbapedia.bulbagarden.net/wiki/Stat#Stage_multipliers
	if p.accuracyStage >= 0 {
		return (3.0 + float64(p.accuracyStage)) / 3.0
	}
	return 3.0 / (3.0 + (-1 * float64(p.accuracyStage)))
}

func (p Pokemon) GetEvasivenesssModifier() float64 {
	// https://bulbapedia.bulbagarden.net/wiki/Stat#Stage_multipliers
	if p.evasivenessStage <= 0 {
		// 3/3 4/3 5/3
		return (3.0 + (-1 * float64(p.evasivenessStage))) / 3.0
	}
	return 3.0 / (3.0 + float64(p.evasivenessStage))
}

func (p Pokemon) GetStageModifier(stage int) float64 {
	// https://bulbapedia.bulbagarden.net/wiki/Stat#Stage_multipliers
	if stage >= 0 {
		// 2/2 3/2 4/2
		return (2.0 + float64(stage)) / 2.0
	}
	return 2.0 / (2.0 + (-1 * float64(stage)))
}

func (p Pokemon) getAttackModifier() float64 {
	return p.GetStageModifier(p.attackStage)
}

func (p Pokemon) getDefenceModifier() float64 {
	return p.GetStageModifier(p.defenceStage)
}

func (p Pokemon) getSpAttackModifier() float64 {
	return p.GetStageModifier(p.spAttackStage)
}

func (p Pokemon) getSpDefenceModifier() float64 {
	return p.GetStageModifier(p.spDefenceStage)
}

func (p Pokemon) getSpeedModifier() float64 {
	return p.GetStageModifier(p.speedStage)
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

func (p Pokemon) GetSpeed() uint {
	return p.speed + p.pokedexItem.baseSpeed
}

func (p Pokemon) GetAdjustedAttack() float64 {
	return float64(p.attack+p.pokedexItem.baseAttack) * p.getAttackModifier()
}

func (p Pokemon) GetAdjustedDefence() float64 {
	return float64(p.defence+p.pokedexItem.baseDefence) * p.getDefenceModifier()
}

func (p Pokemon) GetAdjustedSpecialAttack() float64 {
	return float64(p.spAttack+p.pokedexItem.baseSpAttack) * p.getSpAttackModifier()
}

func (p Pokemon) GetAdjustedSpecialDefence() float64 {
	return float64(p.spDefence+p.pokedexItem.baseSpDefence) * p.getSpDefenceModifier()
}

func (p Pokemon) GetAdjustedSpeed() float64 {
	return float64(p.speed+p.pokedexItem.baseSpeed) * p.getSpeedModifier()
}

func (p *Pokemon) ReceiveAttack(attacker *Pokemon, attack Attack) (bool, float64, []string) {
	// https://bulbapedia.bulbagarden.net/wiki/Accuracy
	accuracy := float64(attack.GetAccuracyValue()) * attacker.GetAccuracyModifier() * p.GetEvasivenesssModifier()
	hit := RandomNumber(0, 255) < int(math.Round(accuracy))
	if !hit {
		return false, 1.0, []string{}
	}
	A := attacker.GetAdjustedAttack()
	D := p.GetAdjustedDefence()
	if attack.special {
		A = attacker.GetAdjustedSpecialAttack()
		D = p.GetAdjustedSpecialDefence()
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
	// https://bulbapedia.bulbagarden.net/wiki/Damage
	damage := 0.0
	if attack.power > 0 {
		damage = ((((((2.0 * float64(attacker.level)) / 5.0) + 2.0) * float64(attack.power) * (A / D)) / 50.0) + 2.0) * typeEffect * attackerTypeBonus
	}
	// term.ShowNoResponseDialog("Damage: %.2f", damage)
	damageInt := uint(math.Round(damage))
	stageModificationStrings := []string{}
	if p.GetHP() <= damageInt {
		p.hpReduction = p.GetMaxHP()
	} else {
		p.hpReduction += damageInt
		stageModificationStrings = PerformStageModifications(attack, attacker, p)
	}

	return true, typeEffect, stageModificationStrings
}

func PerformStageModifications(attack Attack, attacker *Pokemon, defender *Pokemon) (stageModificationStrings []string) {
	for _, stageModifier := range attack.modifiers {
		if stageModifier.modifyAmount != 0 && (stageModifier.chance >= 100 || stageModifier.chance >= uint(RandomNumber(1, 100))) {
			stageModifyTarget := defender.GetName()
			if stageModifier.toSelf {
				stageModifyTarget = attacker.GetName()
			}
			stageModifierChange := "fell"
			if stageModifier.modifyAmount < -1 {
				stageModifierChange = "harshly fell!"
			} else if stageModifier.modifyAmount > 1 {
				stageModifierChange = "rose greatly!"
			} else if stageModifier.modifyAmount > 0 {
				stageModifierChange = "rose"
			}
			stageModificationStrings = append(stageModificationStrings, fmt.Sprintf("%v's %v %v", stageModifyTarget, stageModifier.stageType, stageModifierChange))
			switch stageModifier.stageType {
			case STAGE_MODIFIER_ACCURACY:
				if stageModifier.toSelf {
					attacker.SetAccuracyStage(attacker.accuracyStage + stageModifier.modifyAmount)
				} else {
					defender.SetAccuracyStage(defender.accuracyStage + stageModifier.modifyAmount)
				}
			case STAGE_MODIFIER_EVASIVENESS:
				if stageModifier.toSelf {
					attacker.SetEvasivenessStage(attacker.evasivenessStage + stageModifier.modifyAmount)
				} else {
					defender.SetEvasivenessStage(defender.evasivenessStage + stageModifier.modifyAmount)
				}
			case STAGE_MODIFIER_ATTACK:
				if stageModifier.toSelf {
					attacker.SetAttackStage(attacker.attackStage + stageModifier.modifyAmount)
				} else {
					defender.SetAttackStage(defender.attackStage + stageModifier.modifyAmount)
				}
			case STAGE_MODIFIER_DEFENCE:
				if stageModifier.toSelf {
					attacker.SetDefenceStage(attacker.defenceStage + stageModifier.modifyAmount)
				} else {
					defender.SetDefenceStage(defender.defenceStage + stageModifier.modifyAmount)
				}
			case STAGE_MODIFIER_SP_ATTACK:
				if stageModifier.toSelf {
					attacker.SetSpAttackStage(attacker.spAttackStage + stageModifier.modifyAmount)
				} else {
					defender.SetSpAttackStage(defender.spAttackStage + stageModifier.modifyAmount)
				}
			case STAGE_MODIFIER_SP_DEFENCE:
				if stageModifier.toSelf {
					attacker.SetSpDefenceStage(attacker.spDefenceStage + stageModifier.modifyAmount)
				} else {
					defender.SetSpDefenceStage(defender.spDefenceStage + stageModifier.modifyAmount)
				}
			case STAGE_MODIFIER_SPEED:
				if stageModifier.toSelf {
					attacker.SetSpeedStage(attacker.speedStage + stageModifier.modifyAmount)
				} else {
					defender.SetSpeedStage(defender.speedStage + stageModifier.modifyAmount)
				}
			}

		}
	}
	return
}

func (p *Pokemon) LevelUp(inBattle bool) {
	if p.level >= 100 {
		return
	}
	p.level++
	hpIncrease := GenerateStatChange()
	attackIncrease := GenerateStatChange()
	defenceIncrease := GenerateStatChange()
	spAttackIncrease := GenerateStatChange()
	spDefenceIncrease := GenerateStatChange()
	speedIncrease := GenerateStatChange()
	p.hp += hpIncrease
	p.attack += attackIncrease
	p.defence += defenceIncrease
	p.spAttack += spAttackIncrease
	p.spDefence += spDefenceIncrease
	p.speed += speedIncrease
	evolutionLevel, hasEvolution := p.pokedexItem.GetEvolutionLevel()
	levelAttacks := p.pokedexItem.GetLevelAttacks(p.level)
	readyToEvolve := hasEvolution && p.level >= evolutionLevel
	term.ShowNoResponseDialog("\n\t%v is now level %v!", p.GetName(), p.GetLevel())
	statStr := fmt.Sprintf("\n\tHP: %v", p.GetMaxHP())
	statStr += statStringBracket(hpIncrease)
	statStr += fmt.Sprintf("\n\tAttack: %v", p.GetAttack())
	statStr += statStringBracket(attackIncrease)
	statStr += fmt.Sprintf("\n\tDefence: %v", p.GetDefence())
	statStr += statStringBracket(defenceIncrease)
	statStr += fmt.Sprintf("\n\tSP. Attack: %v", p.GetSpecialAttack())
	statStr += statStringBracket(spAttackIncrease)
	statStr += fmt.Sprintf("\n\tSP. Defence: %v", p.GetSpecialDefence())
	statStr += statStringBracket(spDefenceIncrease)
	statStr += fmt.Sprintf("\n\tSpeed: %v", p.GetSpeed())
	statStr += statStringBracket(speedIncrease)
	term.ShowNoResponseDialog(statStr)
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

func GenerateStatChange() uint {
	R := RandomNumber(1, 100)
	if R > 80 {
		return 2
	}
	if R > 40 {
		return 1
	}
	return 0
}

func statStringBracket(changeAmount uint) string {
	if changeAmount > 0 {
		return fmt.Sprintf(" (+%v)", changeAmount)
	}
	return ""
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
