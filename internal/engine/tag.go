package engine

import "fmt"

type Tag struct {
	name     string
	children []*Tag
	parents  []*Tag
	objects  []*Object
}

func NewTag(name string, parent_names []string) (*Tag, error) {
	var parents []*Tag
	for _, parent_name := range parent_names {
		parent, exists := tagMap[parent_name]
		if !exists {
			return nil, fmt.Errorf("Tag parent '%s' not found.\n", parent)
		}
		parents = append(parents, parent)
	}

	return &Tag{
		name:    name,
		parents: parents,
	}, nil
}

func (t Tag) Print() {
	fmt.Println(t)
}

func (t Tag) String() string {
	var str string
	str += fmt.Sprintf("Tag: %s\n", t.name)
	for _, object := range t.objects {
		str += fmt.Sprintf("\t%s\n", object)
	}
	return str
}

func (t *Tag) addParent(tag *Tag) {
	t.parents = append(t.parents, tag)
}

func (t *Tag) removeParent(tag *Tag) {
	t.parents, _ = delItemFromSlice(t.parents, tag)
}

func (t *Tag) addChild(tag *Tag) {
	t.children = append(t.children, tag)
}

func (t *Tag) removeChild(tag *Tag) {
	t.children, _ = delItemFromSlice(t.children, tag)
}

func (t *Tag) addObject(object *Object) {
	t.objects = append(t.objects, object)
}

func (t *Tag) removeObject(object *Object) {
	t.objects, _ = delItemFromSlice(t.objects, object)
}
