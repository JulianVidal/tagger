package engine

import (
	"errors"
	"fmt"
)

// TODO: Create a TUI that shows all the files, search with query or filename, like locate and everything
// TODO: Add a way to edit the children of tags
// TODO: Add a way to edit the parents of tags
// TODO: Add a way to edit the parents of objects

var tagMap map[string]*Tag
var objectMap map[string]*Object

func InitEngine() {
	tagMap = make(map[string]*Tag)
	objectMap = make(map[string]*Object)
}

func Print() {
	fmt.Println("Printing engine:")
	for _, v := range tagMap {
		v.print()
	}
}

func getAllObjectsFromTag(tagName string) map[string]string {
	results := make(map[string]string)
	tag := tagMap[tagName]
	for _, object := range tag.objects {
		results[object] = object
	}

	for _, childTag := range tag.children {
		mapUnion(results, getAllObjectsFromTag(childTag))
	}

	return results
}

func Query(tags []string) ([]*Object, error) {
	results := make(map[string]string)

	for _, tagName := range tags {
		if _, exist := tagMap[tagName]; !exist {
			return nil, fmt.Errorf("Tag '%s' not found", tagName)
		}
	}

	for _, tagName := range tags {
		mapUnion(results, getAllObjectsFromTag(tagName))
	}

	var resultList []*Object
	for _, objectName := range results {
		obj := objectMap[objectName]
		resultList = append(resultList, obj)
	}

	return resultList, nil
}
