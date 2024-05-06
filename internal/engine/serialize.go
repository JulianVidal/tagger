package engine

import (
	"encoding/json"
)

type Engine struct {
	Tags    map[string]*Tag
	Objects map[string]*Object
}

func ToJson() []byte {
	data, err := json.Marshal(Engine{
		Tags:    tagMap,
		Objects: objectMap,
	})
	if err != nil {
		panic("Couldn't marshal data")
	}

	return data
}

func FromJson(data []byte) {
	var eng Engine
	err := json.Unmarshal(data, &eng)
	if err != nil {
		panic("Couldn't unmarshall data")
	}

	tagMap = eng.Tags
	objectMap = eng.Objects
}
