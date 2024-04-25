package engine

import (
	"encoding/json"
	"os"
)

type EngineSer struct {
	Tags map[string]TagSer
	Objs map[string]ObjSer
}

type TagSer struct {
	Name string
	Tags []string
}

type ObjSer struct {
	Name   string
	Format string
	Tags   []string
}

// TODO: This serializes the tags out of order, meaning if a tag is added when
// TODO: it's parent is not there, the add function will panic.
// TODO: Maybe store the serialized objects as they are added, so order is mantained
// TODO: Idea: When De-serializing turn slice into map, when adding tag add parents first
func (eng *Engine) serialize() EngineSer {
	tags := make(map[string]TagSer)
	objs := make(map[string]ObjSer)

	for _, v := range eng.table {
		tagSer := TagSer{
			Name: v.name,
		}
		for _, parent := range v.parents {
			tagSer.Tags = append(tagSer.Tags, parent.name)
		}
		for _, object := range v.objects {
			if _, exist := objs[object.name]; !exist {
				objSer := ObjSer{
					Name:   object.name,
					Format: object.format,
				}

				for _, parent := range object.parents {
					objSer.Tags = append(objSer.Tags, parent.name)
				}
				objs[object.name] = objSer
			}
		}
		tags[tagSer.Name] = tagSer
	}

	return EngineSer{
		Tags: tags,
		Objs: objs,
	}
}

func (eng *Engine) Json(pathname string) {
	if pathname == "" {
		pathname = "engine.json"
	}

	data, err := json.Marshal(eng.serialize())
	if err != nil {
		panic("Marshalling failed")
	}

	err = os.WriteFile(pathname, data, 0644)

	if err != nil {
		panic("Couldn't write engine to file")
	}
}

func (eng *Engine) addRecursiveTags(tags map[string]TagSer, tag TagSer) {
	_, exist := eng.table[tag.Name]

	if exist {
		return
	}

	for _, parentName := range tag.Tags {
		parent, exist := tags[parentName]

		if !exist {
			panic("Child claims missing parent")
		}

		eng.addRecursiveTags(tags, parent)
	}

	eng.AddTag(tag.Name, tag.Tags)
}

func deserialize(engSer EngineSer) *Engine {
	eng := NewEngine()

	for _, val := range engSer.Tags {
		eng.addRecursiveTags(engSer.Tags, val)
	}

	for _, val := range engSer.Objs {
		eng.AddObj(val.Name, val.Format, val.Tags)
	}

	return eng
}

func FromJson(pathname string) *Engine {
	file, err := os.ReadFile(pathname)
	if err != nil {
		panic("Couldn't read file")
	}

	var eng EngineSer
	err = json.Unmarshal(file, &eng)
	if err != nil {
		panic("Couldn't unmarshall file")
	}

	return deserialize(eng)
}
