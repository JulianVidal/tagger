package editor

import (
	"github.com/JulianVidal/tagger/internal/engine"
)

type EditorTag struct {
	tag *engine.Tag
}

func (e EditorTag) Parents() []string {
	parents := []string{}
	for _, tag := range e.tag.Tags() {
		parents = append(parents, tag.Name())
	}
	return parents
}

func (e EditorTag) SetParent(p string, set bool) {
	var err error
	parent, exists := engine.FindTag(p)
	if !exists {
		panic("Parent tag doesn't exists")
	}

	if set {
		err = e.tag.AddTags(parent)
	} else {
		err = e.tag.RemoveTags(parent)
	}

	if err != nil {
		panic(err)
	}
}

func (e EditorTag) PossibleParents() []string {
	parents := []string{}

OUTER:
	for _, parent := range engine.Tags() {
		if parent == e.tag.Name() {
			continue
		}
		for _, child := range e.tag.Children() {
			if parent == child.Name() {
				continue OUTER
			}
		}
		parents = append(parents, parent)
	}

	return parents
}

func NewEditorTag(name string) EditorTag {
	tag, exists := engine.FindTag(name)
	if !exists {
		panic("Can't edit tag that doesn't exist")
	}
	return EditorTag{tag}
}
