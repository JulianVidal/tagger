package engine

import (
	"fmt"
)

type ObjectJSON struct {
	Name string
	Tags []string
}

type Object struct {
	name string
	tags []*Tag
}

func NewObject(name string) (*Object, error) {
	if _, exists := objectMap[name]; exists {
		return nil, fmt.Errorf("Object '%s' already exists", name)
	}
	objectMap[name] = &Object{
		name: name,
	}
	return objectMap[name], nil
}

func (o *Object) Print() {
	fmt.Println(o)
}

func (o *Object) String() string {
	var t string
	for _, tag := range o.tags {
		t += " " + tag.name
	}
	return fmt.Sprintf("Object:%s, Tags: %v", o.name, t)
}

func (o *Object) Name() string {
	return o.name
}

func (o *Object) Tags() []*Tag {
	return o.tags
}

func (o *Object) AddTags(tags ...*Tag) error {
	for _, tag := range tags {
		if _, exists := tagMap[tag.name]; !exists {
			return fmt.Errorf("Tag '%s' doesn't exist", tag.name)
		}
	}

	o.tags = append(o.tags, tags...)
	for _, tag := range tags {
		tag.objects = append(tag.objects, o)
	}

	return nil
}

func (o *Object) RemoveTag(tags ...*Tag) {
	o.tags, _ = delItemsFromSlice(o.tags, tags...)
	for _, tag := range tags {
		tag.objects, _ = delItemsFromSlice(tag.objects, o)
	}
}

func (o *Object) Delete() {

	for _, tag := range o.tags {
		tag.objects, _ = delItemsFromSlice(tag.objects, o)
	}

	delete(objectMap, o.name)
}

func (o *Object) SetName(name string) {
	delete(objectMap, o.name)
	o.name = name
	objectMap[name] = o
}

func (o *Object) json() ObjectJSON {
	var tags []string
	for _, tag := range o.tags {
		tags = append(tags, tag.name)
	}
	return ObjectJSON{
		Name: o.name,
		Tags: tags,
	}
}
