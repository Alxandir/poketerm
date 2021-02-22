package pokemon

type ElemenalType interface {
	GetName() string
	Weaknesses() map[string]float64
}

type Normal struct{}

func (Normal) GetName() string {
	return "normal"
}
func (Normal) Weaknesses() map[string]float64 {
	return map[string]float64{
		"fighting": 2.0,
	}
}

type Fire struct{}

func (Fire) GetName() string {
	return "fire"
}
func (Fire) Weaknesses() map[string]float64 {
	return map[string]float64{
		"water": 2.0,
		"bug":   0.5,
		"ice":   0.5,
		"rock":  2.0,
	}
}

type Water struct{}

func (Water) GetName() string {
	return "water"
}
func (Water) Weaknesses() map[string]float64 {
	return map[string]float64{
		"grass": 2.0,
		"ice":   2.0,
		"rock":  0.5,
	}
}

type Grass struct{}

func (Grass) GetName() string {
	return "grass"
}
func (Grass) Weaknesses() map[string]float64 {
	return map[string]float64{
		"fire":   2.0,
		"flying": 2.0,
		"ice":    2.0,
		"poison": 2.0,
		"bug":    2.0,
		"water":  0.5,
	}
}

type Flying struct{}

func (Flying) GetName() string {
	return "flying"
}
func (Flying) Weaknesses() map[string]float64 {
	return map[string]float64{
		"electric": 2.0,
		"ice":      2.0,
		"rock":     2.0,
	}
}

type Poison struct{}

func (Poison) GetName() string {
	return "poison"
}
func (Poison) Weaknesses() map[string]float64 {
	return map[string]float64{
		"psychic": 2.0,
	}
}

type Bug struct{}

func (Bug) GetName() string {
	return "bug"
}
func (Bug) Weaknesses() map[string]float64 {
	return map[string]float64{
		"fire":   2.0,
		"flying": 2.0,
		"rock":   2.0,
	}
}

type Psychic struct{}

func (Psychic) GetName() string {
	return "psychic"
}
func (Psychic) Weaknesses() map[string]float64 {
	return map[string]float64{
		"fire": 2.0,
	}
}

type Ice struct{}

func (Ice) GetName() string {
	return "ice"
}
func (Ice) Weaknesses() map[string]float64 {
	return map[string]float64{
		"fire": 2.0,
	}
}
