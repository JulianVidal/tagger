package engine

import (
	"encoding/json"
	"fmt"
)

type ObjectJSON struct {
	name string
	tags []string
}

type Object struct {
	name string
	tags []*Tag
}

func NewObject(name string, tag_names []string) (*Object, error) {
	var parents []*Tag
	for _, tag_name := range tag_names {
		tag, exists := tagMap[tag_name]
		if !exists {
			return nil, fmt.Errorf("Tag '%s' not found.\n", tag)
		}
		parents = append(parents, tag)
	}

	return &Object{
		name: name,
		tags: parents,
	}, nil
}

func (o *Object) Print() {
	fmt.Println(o)
}

func (o *Object) String() string {
	return fmt.Sprintf("Object:%s", o.name)
}

func (o *Object) addTag(parent *Tag) {
	o.tags = append(o.tags, parent)
}

func (o *Object) removeTag(parent *Tag) {
	o.tags, _ = delItemFromSlice(o.tags, parent)
}

func (o *Object) MarshalJSON() ([]byte, error) {
	var tags []string
	for _, tag := range o.tags {
		tags = append(tags, tag.name)
	}
	return json.Marshal(ObjectJSON{
		name: o.name,
		tags: tags,
	})
}
