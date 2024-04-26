package command

import (
	"encoding/json"

	"github.com/JulianVidal/tagger/internal/serialize"
)

type Verb int

const (
	Add Verb = iota
	Delete
	Query
)

func (v Verb) String() string {
	return [...]string{"Add", "Delete", "Query"}[v]
}

func (v Verb) EnumIndex() int {
	return int(v)
}

type Noun int

const (
	Tag Noun = iota
	Object
)

func (n Noun) String() string {
	return [...]string{"Tag", "Object"}[n]
}

func (n Noun) EnumIndex() int {
	return int(n)
}

type Command struct {
	Verb    Verb
	Noun    Noun
	Subject []byte
}

func CreateCommand[S serialize.Tag | serialize.Obj](verb Verb, noun Noun, subject S) (Command, error) {
	sub, err := json.Marshal(subject)
	if err != nil {
		return Command{}, err
	}

	return Command{
			Verb:    verb,
			Noun:    noun,
			Subject: sub,
		},
		nil
}
