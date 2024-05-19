package handler

import "github.com/JulianVidal/tagger/internal/engine"

func TagParents(name string) []string {
	tag, exists := engine.FindTag(name)
	if !exists {
		panic("Tag doesn't exist")
	}

	var parents []string
	for _, parent := range tag.Tags() {
		parents = append(parents, parent.Name())
	}

	return parents
}

func ObjectTags(name string) []string {
	object, exists := engine.FindObject(name)
	if !exists {
		object, _ = engine.NewObject(name)
	}

	var tags []string
	for _, tag := range object.Tags() {
		tags = append(tags, tag.Name())
	}

	return tags
}

func SetTagParents(name string, parents []string) {
	tag, exists := engine.FindTag(name)
	if !exists {
		panic("Tag doesn't exist")
	}

	tag.RemoveTags(tag.Tags()...)
	for _, parentName := range parents {
		parent, exists := engine.FindTag(parentName)
		if !exists {
			panic("Tag doesn't exist")
		}
		tag.AddTags(parent)
	}
}

func SetObjectTags(name string, tags []string) {
	object, exists := engine.FindObject(name)
	if !exists {
		panic("Object doesn't exist")
	}

	object.RemoveTag(object.Tags()...)
	for _, tagName := range tags {
		tag, exists := engine.FindTag(tagName)
		if !exists {
			panic("Tag doesn't exist")
		}
		object.AddTags(tag)
	}
}

func QueryEngine(tagNames []string) []string {
	var tags []*engine.Tag
	for _, tagName := range tagNames {
		tag, exists := engine.FindTag(tagName)
		if !exists {
			panic("Couldn't find tag in engine")
		}
		tags = append(tags, tag)
	}
	objects, err := engine.Query(tags...)
	if err != nil {
		panic(err)
	}

	var files []string
	for _, object := range objects {
		files = append(files, object.Name())
	}
	return files
}

func EngineString() string {
	return engine.String()
}

func Tags() []string {
	return engine.Tags()
}

func ValidParenTags(tagName string) []string {
	t, exists := engine.FindTag(tagName)
	if !exists {
		panic("Couldn't find tag")
	}
	validParents := []string{}

OUTER:
	for _, parent := range engine.Tags() {
		if parent == tagName {
			continue
		}
		for _, child := range t.Children() {
			if parent == child.Name() {
				continue OUTER
			}
		}
		validParents = append(validParents, parent)
	}

	return validParents
}
