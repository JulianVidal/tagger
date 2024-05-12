package engine

import (
	"encoding/json"
)

type EngineJSON struct {
	Tags    map[string]TagJSON
	Objects map[string]ObjectJSON
}

func ToJson() ([]byte, error) {
	tags := make(map[string]TagJSON)
	objects := make(map[string]ObjectJSON)

	for k, v := range tagMap {
		tags[k] = v.json()
	}

	for k, v := range objectMap {
		objects[k] = v.json()
	}

	data, err := json.Marshal(EngineJSON{
		Tags:    tags,
		Objects: objects,
	})

	return data, err
}

func addTagsRecursive(tags map[string]TagJSON, tagJSON TagJSON) {
	_, exist := tagMap[tagJSON.Name]
	if exist {
		return
	}

	for _, parentName := range tagJSON.Parents {
		parent, exist := tags[parentName]
		if !exist {
			panic("Parent tag doesn't exist")
		}

		addTagsRecursive(tags, parent)
	}

	tag, err := NewTag(tagJSON.Name)
	if err != nil {
		panic(err)
	}

	for _, parentName := range tagJSON.Parents {
		parent, exists := FindTag(parentName)
		if !exists {
			panic("Tag doesn't exist")
		}
		tag.AddTags(parent)
	}
}

func FromJson(data []byte) error {
	var eng EngineJSON
	err := json.Unmarshal(data, &eng)

	New()

	for _, t := range eng.Tags {
		addTagsRecursive(eng.Tags, t)
	}

	for _, o := range eng.Objects {
		object, err := NewObject(o.Name)
		if err != nil {
			panic(err)
		}
		for _, tagName := range o.Tags {
			tag, exists := FindTag(tagName)
			if !exists {
				panic("Tag doesn't exist")
			}
			object.AddTags(tag)
		}
	}

	return err
}
