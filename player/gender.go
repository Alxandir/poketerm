package player

type Gender interface {
	GetName() string
	GetPronoun() string
	GetOwnershipPronoun() string
}

type GenderMale struct {
	name string
}

func (GenderMale) GetName() string {
	return "boy"
}
func (GenderMale) GetPronoun() string {
	return "he"
}
func (GenderMale) GetOwnershipPronoun() string {
	return "his"
}

type GenderFemale struct {
	name string
}

func (GenderFemale) GetName() string {
	return "girl"
}
func (GenderFemale) GetPronoun() string {
	return "she"
}
func (GenderFemale) GetOwnershipPronoun() string {
	return "hers"
}
