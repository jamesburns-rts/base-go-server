package example

type (
	Pokemon struct {
		ID        int         `json:"id"`
		Abilities []Abilities `json:"abilities"`
		Height    int         `json:"height"`
		Moves     []Moves     `json:"moves"`
		Name      string      `json:"name"`
		Order     int         `json:"order"`
		Types     []Types     `json:"types"`
	}

	Ability struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	Abilities struct {
		Ability  Ability `json:"ability"`
		IsHidden bool    `json:"is_hidden"`
		Slot     int     `json:"slot"`
	}

	Move struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	Moves struct {
		Move Move `json:"move"`
	}

	Type struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	Types struct {
		Slot int  `json:"slot"`
		Type Type `json:"type"`
	}
)
