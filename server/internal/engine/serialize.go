package engine

import (
	"encoding/json"

	"github.com/JulianVidal/tagger/internal/serialize"
)

func (tag *TagNode) serialize() serialize.Tag {
	tagSer := serialize.Tag{
		Name: tag.name,
	}

	for _, parent := range tag.parents {
		tagSer.Tags = append(tagSer.Tags, parent.name)
	}

	return tagSer
}

func (obj *ObjNode) serialize() serialize.Obj {
	objSer := serialize.Obj{
		Name:   obj.name,
		Format: obj.format,
	}

	for _, parent := range obj.parents {
		objSer.Tags = append(objSer.Tags, parent.name)
	}

	return objSer
}

func (eng *Engine) serialize() serialize.Engine {
	tags := make(map[string]serialize.Tag)
	objs := make(map[string]serialize.Obj)

	for _, v := range eng.table {
		tags[v.name] = v.serialize()

		for _, object := range v.objects {
			if _, exist := objs[object.name]; !exist {
				objs[object.name] = object.serialize()
			}
		}
	}

	return serialize.Engine{
		Tags: tags,
		Objs: objs,
	}
}

func (eng *Engine) ToJson() []byte {
	data, err := json.Marshal(eng.serialize())
	if err != nil {
		panic("Couldn't marshal data")
	}

	return data
}

func (eng *Engine) addRecursiveTags(tags map[string]serialize.Tag, tag serialize.Tag) {
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

func deserialize(engSer serialize.Engine) *Engine {
	eng := NewEngine()

	for _, val := range engSer.Tags {
		eng.addRecursiveTags(engSer.Tags, val)
	}

	for _, val := range engSer.Objs {
		eng.AddObj(val.Name, val.Format, val.Tags)
	}

	return eng
}

func FromJson(data []byte) *Engine {
	var eng serialize.Engine
	err := json.Unmarshal(data, &eng)
	if err != nil {
		panic("Couldn't unmarshall data")
	}

	return deserialize(eng)
}
