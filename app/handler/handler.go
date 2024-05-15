package handler

import "github.com/JulianVidal/tagger/internal/engine"

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
