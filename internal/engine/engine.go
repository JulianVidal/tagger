package engine

import (
	"fmt"
	"os"
)

// TODO: Create a TUI that shows all the files, search with query or filename, like locate and everything
// TODO: Add a way to edit the children of tags
// TODO: Add a way to edit the parents of tags
// TODO: Add a way to edit the parents of objects

var tagMap map[string]*Tag
var objectMap map[string]*Object

func InitEngine() {
	data, err := os.ReadFile("engine.json")
	if err != nil {
		tagMap = make(map[string]*Tag)
		objectMap = make(map[string]*Object)
		return
	}

	err = FromJson(data)
	if err != nil {
		panic(err)
	}
}

func Print() {
	fmt.Println("Printing engine:")
	for _, v := range tagMap {
		v.Print()
	}
}

func FindTag(name string) (Tag, bool) {
	tag, exist := tagMap[name]
	return *tag, exist
}

func FindObject(name string) (Object, bool) {
	object, exist := objectMap[name]
	return *object, exist
}

func Tags() []string {
	return mapKeys(tagMap)
}

func Objects() []string {
	return mapKeys(objectMap)
}

func AddTag(tag *Tag) error {
	if _, exist := tagMap[tag.name]; exist {
		return fmt.Errorf("Tag '%s' already exists", tag.name)
	}

	for _, parent := range tag.parents {
		if _, exist := tagMap[parent.name]; !exist {
			return fmt.Errorf("Parent tag '%s' not found", parent)
		}
	}

	for _, parent := range tag.parents {
		parent.addChild(tag)
	}

	tagMap[tag.name] = tag

	return nil
}

func DelTag(tag *Tag) error {
	if _, exist := tagMap[tag.name]; !exist {
		return fmt.Errorf("Tag '%s' not found", tag.name)
	}

	for _, child := range tag.children {
		child.removeParent(tag)
	}

	for _, parent := range tag.parents {
		parent.removeChild(tag)
	}

	delete(tagMap, tag.name)

	return nil
}

func AddObject(object *Object) error {
	if _, exist := objectMap[object.name]; exist {
		return fmt.Errorf("Object '%s' already exists", object.name)
	}

	for _, tag := range object.tags {
		if _, exist := tagMap[tag.name]; !exist {
			return fmt.Errorf("Tag '%s' doesn't exist in engine", tag.name)
		}
	}

	for _, tag := range object.tags {
		tag.addObject(object)
	}

	objectMap[object.name] = object

	return nil
}

func DeleteObject(object *Object) error {
	if _, exist := objectMap[object.name]; !exist {
		return fmt.Errorf("Object %s not found", object.name)
	}

	for _, parent := range object.tags {
		parent.removeObject(object)
	}

	delete(objectMap, object.name)

	return nil
}

func getAllObjectsFromTag(tag *Tag) map[string]*Object {
	results := make(map[string]*Object)
	for _, object := range tag.objects {
		results[object.name] = object
	}

	for _, childTag := range tag.children {
		mapUnion(results, getAllObjectsFromTag(childTag))
	}

	return results
}

func Query(tags ...*Tag) ([]Object, error) {
	results := make(map[string]*Object)

	for _, tag := range tags {
		if _, exist := tagMap[tag.name]; !exist {
			return nil, fmt.Errorf("Tag '%s' not found", tag)
		}
	}

	for _, tag := range tags {
		mapUnion(results, getAllObjectsFromTag(tag))
	}

	var resultList []Object
	for _, object := range results {
		resultList = append(resultList, *object)
	}

	return resultList, nil
}
