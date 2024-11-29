package pkg

type Persons struct {
	Persons []Person
}

// функция для поиска в массиве Persons чела с определенным именем
func (p *Persons) Find(name string) *Person {
	for i := range p.Persons {
		if p.Persons[i].Username == name {
			return &p.Persons[i]
		}
	}
	return nil
}
