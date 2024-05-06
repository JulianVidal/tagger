package command

import (
	"encoding/json"
	"errors"

	"github.com/JulianVidal/tagger/internal/serialized"
)

type Command int

const (
	AddTag Command = iota
	DelTag
	AddObj
	DelObj
	Query
)

func (c Command) String() string {
	return [...]string{"AddTag", "DelTag", "AddObj", "DelObj", "Query"}[c]
}

func (c Command) EnumIndex() int {
	return int(c)
}

type AddTagData struct {
	Tag serialized.Tag
}

type DelTagData struct {
	Tag string
}

type QueryData struct {
	Tags []string
}

type AddObjData struct {
	Obj serialized.Obj
}

type DelObjData struct {
	Obj string
}

type Packet struct {
	Type Command
	Data interface{}
}

func (p *Packet) UnmarshalJSON(b []byte) (err error) {
	var packet map[string]interface{}

	err = json.Unmarshal(b, &packet)
	if err != nil {
		return err
	}

	p.Type = Command((packet["Type"].(float64)))
	data := packet["Data"].(map[string]interface{})

	switch p.Type {
	case AddTag:
		tag := data["Tag"].(map[string]interface{})

		tagSer := serialized.Tag{
			Name: tag["Name"].(string),
		}
		if tag["Tags"] != nil {
			var tags []string
			for _, t := range tag["Tags"].([]interface{}) {
				tags = append(tags, t.(string))
			}
			tagSer.Tags = tags
		}

		p.Data = AddTagData{
			Tag: tagSer,
		}
	case DelTag:
		p.Data = DelTagData{
			Tag: data["Tag"].(string),
		}

	case AddObj:
		obj := data["Obj"].(map[string]interface{})

		objSer := serialized.Obj{
			Name:   obj["Name"].(string),
			Format: obj["Format"].(string),
		}
		if obj["Tags"] != nil {
			var tags []string
			for _, t := range obj["Tags"].([]interface{}) {
				tags = append(tags, t.(string))
			}
			objSer.Tags = tags
		}

		p.Data = AddObjData{
			Obj: objSer,
		}

	case DelObj:
		p.Data = DelObjData{
			Obj: data["Obj"].(string),
		}

	case Query:
		var tags []string
		for _, t := range data["Tags"].([]interface{}) {
			tags = append(tags, t.(string))
		}
		p.Data = QueryData{
			Tags: tags,
		}

	default:
		return errors.New("Unknown command")
	}

	return nil
}
