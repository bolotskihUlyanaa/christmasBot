package pkg

type persons struct {
	Persons []person
}

// функция для поиска в массиве Persons чела с определенным именем
func (p *persons) find(name string) *person {
	for i := range p.Persons {
		if p.Persons[i].Username == name {
			return &p.Persons[i]
		}
	}
	return nil
}
