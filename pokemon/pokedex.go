package pokemon

type PokedexItem struct {
	id            uint
	name          string
	baseLevel     uint
	evolutionId   uint
	baseHP        uint
	baseAttack    uint
	baseDefence   uint
	baseSpAttack  uint
	baseSpDefence uint
	baseAccuracy  uint
	types         []ElemenalType
	attacks       []PokedexItemAttack
}

type PokedexItemAttack struct {
	minLevel uint
	attack   Attack
}

func (p PokedexItem) GetName() string {
	return p.name
}

func (p PokedexItem) GetID() uint {
	return p.id
}

func (p PokedexItem) GetEvolution() (PokedexItem, bool) {
	return FindInPokedex(p.evolutionId)
}

func (p PokedexItem) GetEvolutionLevel() (uint, bool) {
	evolution, exists := p.GetEvolution()
	if exists {
		return evolution.baseLevel, true
	}
	return 0, false
}

func (p PokedexItem) GetLevelAttacks(level uint) []Attack {
	levelAttacks := []Attack{}
	for _, attack := range p.attacks {
		if level == attack.minLevel {
			levelAttacks = append(levelAttacks, attack.attack)
		}
	}
	return levelAttacks
}

func (p PokedexItem) GetTypes() []ElemenalType {
	return p.types
}

func FindInPokedex(id uint) (PokedexItem, bool) {
	for _, item := range Pokedex {
		if item.id == id {
			return item, true
		}
	}
	return PokedexItem{}, false
}

// Contains tells whether a contains x.
func Contains(id uint) bool {
	for _, item := range Pokedex {
		if item.id == id {
			return true
		}
	}
	return false
}

var Pokedex = []PokedexItem{
	{
		id:            1,
		name:          "Bulbasaur",
		baseLevel:     1,
		evolutionId:   2,
		baseHP:        39,
		baseAttack:    52,
		baseDefence:   43,
		baseSpAttack:  60,
		baseSpDefence: 50,
		types: []ElemenalType{
			Grass{},
			Poison{},
		},
		attacks: []PokedexItemAttack{
			BasicTackle,
			BasicSandAttack,
			{
				attack:   Attacks[3],
				minLevel: 6,
			},
			{
				attack:   Attacks[5],
				minLevel: 6,
			},
		},
	},
	{
		id:          2,
		name:        "Ivysaur",
		baseLevel:   16,
		evolutionId: 3,
		types: []ElemenalType{
			Grass{},
			Poison{},
		},
	},
	{
		id:        3,
		name:      "Venusaur",
		baseLevel: 32,
		types: []ElemenalType{
			Grass{},
			Poison{},
		},
	},
	{ // https://bulbapedia.bulbagarden.net/wiki/Charmander
		id:            4,
		name:          "Charmander",
		baseLevel:     1,
		evolutionId:   5,
		baseHP:        39,
		baseAttack:    52,
		baseDefence:   43,
		baseSpAttack:  60,
		baseSpDefence: 50,
		types: []ElemenalType{
			Fire{},
		},
		attacks: []PokedexItemAttack{
			BasicTackle,
			{
				attack:   Attacks[2],
				minLevel: 5,
			},
			{
				attack:   Attacks[6],
				minLevel: 5,
			},
			{
				attack:   Attacks[5],
				minLevel: 6,
			},
			{
				attack:   Attacks[4],
				minLevel: 6,
			},
		},
	},
	{
		id:          5,
		name:        "Charmeleon",
		baseLevel:   6,
		evolutionId: 6,
		types: []ElemenalType{
			Fire{},
		},
	},
	{
		id:        6,
		name:      "Charizard",
		baseLevel: 32,
		types: []ElemenalType{
			Fire{},
			Flying{},
		},
	},
	{
		id:            7,
		name:          "Squirtle",
		baseLevel:     1,
		evolutionId:   8,
		baseHP:        39,
		baseAttack:    52,
		baseDefence:   43,
		baseSpAttack:  60,
		baseSpDefence: 50,
		types: []ElemenalType{
			Water{},
		},
		attacks: []PokedexItemAttack{
			BasicTackle,
			{
				attack:   Attacks[1],
				minLevel: 5,
			},
		},
	},
	{
		id:          8,
		name:        "Wartortle",
		baseLevel:   16,
		evolutionId: 9,
		types: []ElemenalType{
			Water{},
		},
	},
	{
		id:        9,
		name:      "Blastoise",
		baseLevel: 32,
		types: []ElemenalType{
			Water{},
		},
	},
	{ // https://bulbapedia.bulbagarden.net/wiki/Caterpie
		id:            10,
		name:          "Caterpie",
		baseLevel:     1,
		evolutionId:   11,
		baseHP:        45,
		baseAttack:    30,
		baseDefence:   35,
		baseSpAttack:  20,
		baseSpDefence: 20,
		types: []ElemenalType{
			Bug{},
		},
	},
	{
		id:          11,
		name:        "Metapod",
		baseLevel:   10,
		evolutionId: 12,
		types: []ElemenalType{
			Bug{},
		},
	},
	{
		id:        12,
		name:      "Butterfree",
		baseLevel: 24,
		types: []ElemenalType{
			Bug{},
			Flying{},
		},
	},
}
