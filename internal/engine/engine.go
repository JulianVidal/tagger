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

func Init() {
	data, err := os.ReadFile("engine.json")
	if err != nil {
		New()
		return
	}

	err = FromJson(data)
	if err != nil {
		panic(err)
	}
}

func New() {
	tagMap = make(map[string]*Tag)
	objectMap = make(map[string]*Object)
}

func Print() {
	fmt.Println("Printing engine:")
	for _, v := range tagMap {
		v.Print()
	}
}

func String() string {
	str := ""
	str += fmt.Sprintln("Printing engine:")
	for _, v := range tagMap {
		str += v.String()
	}
	for _, v := range objectMap {
		str += v.String() + "\n"
	}
	return str
}

func FindTag(name string) (*Tag, bool) {
	tag, exist := tagMap[name]
	return tag, exist
}

func FindObject(name string) (*Object, bool) {
	object, exist := objectMap[name]
	if !exist {
		return nil, exist
	}
	return object, exist
}

func Tags() []string {
	return mapKeys(tagMap)
}

func Objects() []string {
	return mapKeys(objectMap)
}

func getAllObjectsFromTag(tag *Tag) map[string]*Object {
	results := make(map[string]*Object)
	for _, object := range tag.objects {
		results[object.Name()] = object
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

func TagExists(tag string) bool {
	_, ok := tagMap[tag]
	return ok
}

func ObjectExists(object string) bool {
	_, ok := objectMap[object]
	return ok
}
