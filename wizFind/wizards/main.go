package wizards

type Wizard struct {
	Name   string `json:"name"`
	Origin string `json:"origin"`
	School string `json:"school"`
}

func Find(wizards []Wizard, test func(Wizard) bool) (ret []Wizard) {
	for i, w := range wizards {
		if test(w) {
			ret = append(ret, wizards[i])
		}
	}

	//make sure we always return something
	if ret == nil {
		ret = []Wizard{}
	}

	return
}

func GetAll() []Wizard {
	return []Wizard{
		Wizard{
			Name:   "Ged",
			Origin: "Gont",
			School: "Wizard School of Roke",
		},
		Wizard{
			Name:   "Ogion",
			Origin: "Gont",
			School: "Wizard School of Roke",
		},
		Wizard{
			Name:   "Estarriol",
			Origin: "Iffish",
			School: "Wizard School of Roke",
		},
		Wizard{
			Name:   "Harry Potter",
			Origin: "England",
			School: "Hogwarts",
		},
		Wizard{
			Name:   "Hermione",
			Origin: "England",
			School: "Hogwarts",
		},
		Wizard{
			Name:   "Merlin",
			Origin: "England",
			School: "Hard Nocks",
		},
		Wizard{
			Name:   "Gandalf",
			Origin: "Middle Earth",
			School: "Unknown",
		},
	}
}
