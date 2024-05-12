package engine

import (
	"encoding/json"
)

type Engine struct {
	Tags    map[string]*Tag
	Objects map[string]*Object
}

func ToJson() ([]byte, error) {
	data, err := json.Marshal(Engine{
		Tags:    tagMap,
		Objects: objectMap,
	})

	return data, err
}

func FromJson(data []byte) error {
	var eng Engine
	err := json.Unmarshal(data, &eng)

	tagMap = eng.Tags
	objectMap = eng.Objects
	return err
}
