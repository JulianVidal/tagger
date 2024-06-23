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
		for _, t := range o.Tags() {
			if t.Name() == tag.Name() {
				return fmt.Errorf("Tag '%s' is already a there", tag.name)
			}
		}
	}

	o.tags = append(o.tags, tags...)
	for _, tag := range tags {
		tag.objects = append(tag.objects, o)
	}

	return nil
}

func (o *Object) RemoveTags(tags ...*Tag) error {
	for _, tag := range tags {
		if _, exists := tagMap[tag.name]; !exists {
			return fmt.Errorf("Tag '%s' doesn't exist", tag.name)
		}
		exists := false
		for _, t := range o.Tags() {
			if t.Name() == tag.Name() {
				exists = true
			}
		}
		if !exists {
			return fmt.Errorf("Tag '%s' is not there", tag.name)
		}
	}

	var err error
	o.tags, err = delItemsFromSlice(o.tags, tags...)
	if err != nil {
		panic("Couldn't delete tags from object")
	}
	for _, tag := range tags {
		tag.objects, err = delItemsFromSlice(tag.objects, o)
		if err != nil {
			panic("Couldn't delete object from tags")
		}
	}

	return nil
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
