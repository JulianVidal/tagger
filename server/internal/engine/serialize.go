package engine

import (
	"encoding/json"

	"github.com/JulianVidal/tagger/internal/serialized"
)

func (tag *TagNode) serialize() serialized.Tag {
	tagSer := serialized.Tag{
		Name: tag.name,
	}

	for _, parent := range tag.parents {
		tagSer.Tags = append(tagSer.Tags, parent.name)
	}

	return tagSer
}

func (obj *ObjNode) serialize() serialized.Obj {
	objSer := serialized.Obj{
		Name:   obj.name,
		Format: obj.format,
	}

	for _, parent := range obj.parents {
		objSer.Tags = append(objSer.Tags, parent.name)
	}

	return objSer
}

func serialize() serialized.Engine {
	tags := make(map[string]serialized.Tag)
	objs := make(map[string]serialized.Obj)

	for _, v := range table {
		tags[v.name] = v.serialize()

		for _, object := range v.objects {
			if _, exist := objs[object.name]; !exist {
				objs[object.name] = object.serialize()
			}
		}
	}

	return serialized.Engine{
		Tags: tags,
		Objs: objs,
	}
}

func ToJson() []byte {
	data, err := json.Marshal(serialize())
	if err != nil {
		panic("Couldn't marshal data")
	}

	return data
}

func addRecursiveTags(tags map[string]serialized.Tag, tag serialized.Tag) {
	_, exist := table[tag.Name]

	if exist {
		return
	}

	for _, parentName := range tag.Tags {
		parent, exist := tags[parentName]

		if !exist {
			panic("Parent Tag doesn't exist")
		}

		addRecursiveTags(tags, parent)
	}

	AddTag(tag.Name, tag.Tags)
}

func deserialize(engSer serialized.Engine) {
	InitEngine()

	for _, val := range engSer.Tags {
		addRecursiveTags(engSer.Tags, val)
	}

	for _, val := range engSer.Objs {
		AddObj(val.Name, val.Format, val.Tags)
	}
}

func FromJson(data []byte) {
	var eng serialized.Engine
	err := json.Unmarshal(data, &eng)
	if err != nil {
		panic("Couldn't unmarshall data")
	}

	deserialize(eng)
}
