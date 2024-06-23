package editor

import (
	"github.com/JulianVidal/tagger/internal/engine"
)

type EditorObject struct {
	object *engine.Object
}

func (e EditorObject) Parents() []string {
	parents := []string{}
	for _, tag := range e.object.Tags() {
		parents = append(parents, tag.Name())
	}
	return parents
}

func (e EditorObject) SetParent(p string, set bool) {
	var err error
	parent, exists := engine.FindTag(p)
	if !exists {
		panic("Parent tag doesn't exists")
	}

	if set {
		err = e.object.AddTags(parent)
	} else {
		err = e.object.RemoveTags(parent)
	}

	if err != nil {
		panic(err)
	}
}

func (e EditorObject) PossibleParents() []string {
	return engine.Tags()
}

func NewEditorObject(name string) EditorObject {
	object, exists := engine.FindObject(name)
	if !exists {
		object, _ = engine.NewObject(name)
		// panic("Can't edit object that doesn't exist")
	}
	return EditorObject{object}
}
