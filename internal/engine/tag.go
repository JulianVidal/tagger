package engine

import (
	"fmt"
)

type TagJSON struct {
	Name    string
	Parents []string
}

type Tag struct {
	name     string
	children []*Tag
	parents  []*Tag
	objects  []*Object
}

func NewTag(name string) (*Tag, error) {
	if _, exists := tagMap[name]; exists {
		return nil, fmt.Errorf("Tag '%s' already exists", name)
	}
	tagMap[name] = &Tag{
		name: name,
	}
	return tagMap[name], nil
}

func (t *Tag) Print() {
	fmt.Println(t)
}

func (t *Tag) String() string {
	var str string
	str += fmt.Sprintf("Tag: %s\n", t.name)
	for _, object := range t.objects {
		str += fmt.Sprintf("\t%s\n", object)
	}
	return str
}

func (t *Tag) Name() string {
	return t.name
}

func (t *Tag) Tags() []*Tag {
	return t.parents
}

func (t *Tag) AddTags(tags ...*Tag) error {
	for _, tag := range tags {
		if _, exists := tagMap[tag.name]; !exists {
			return fmt.Errorf("Tag '%s' doesn't exist", tag.name)
		}
	}

	t.parents = append(t.parents, tags...)
	for _, tag := range tags {
		tag.children = append(tag.children, t)
	}

	return nil
}

func (t *Tag) RemoveTags(tags ...*Tag) {
	t.parents, _ = delItemsFromSlice(t.parents, tags...)
	for _, tag := range tags {
		tag.children, _ = delItemsFromSlice(tag.children, t)
	}
}

func (t *Tag) Delete() {

	for _, parent := range t.parents {
		parent.children, _ = delItemsFromSlice(parent.children, t)
	}

	for _, child := range t.children {
		child.parents, _ = delItemsFromSlice(child.parents, t)
	}

	delete(tagMap, t.name)
}

func (t *Tag) SetName(name string) {
	delete(tagMap, t.name)
	t.name = name
	tagMap[name] = t
}

func (t *Tag) json() TagJSON {
	var parents []string
	for _, parent := range t.parents {
		parents = append(parents, parent.name)
	}
	return TagJSON{
		Name:    t.name,
		Parents: parents,
	}
}
