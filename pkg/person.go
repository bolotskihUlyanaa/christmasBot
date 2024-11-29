package pkg

import "fmt"

type Person struct {
	Username string
	Free     bool
	Block    []string
	Friend   string //кому дарит подарок
	Sex      bool   //true - женский
}

// проверить не в блоке ли у p участник с именем name
func (p *Person) IsBlock(name string) bool {
	for _, val := range p.Block {
		if val == name {
			return true
		}
	}
	return false
}

// функция для распределения участников для игры в тайного санту
// p -> person
// можно ли чтобы p дарил подарок person
// проверка:
// не занят ли уже person
// не дарит ли p подарок самому себе
// не находится ли person в блоке у p
func (p *Person) Filter1(person *Person) bool {
	if person.Free && p.Username != person.Username && !p.IsBlock(person.Username) {
		person.Free = false
		p.Friend = person.Username
		return true
	}
	return false
}

// функция для распределения участников для игры в тайного санту
// p -> person
// проверки:
// не занят ли уже person
// пары назначаются только противоположного пола
// не находится ли person в блоке у p
func (p *Person) Filter2(person *Person) bool {
	if person.Free && p.Sex != person.Sex && !p.IsBlock(person.Username) {
		person.Free = false
		p.Friend = person.Username
		fmt.Println("TRUE")
		return true
	}
	fmt.Println("false")
	return false
}
