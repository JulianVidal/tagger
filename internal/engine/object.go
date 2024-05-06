package engine

import "fmt"

type Object struct {
	name    string
	format  string
	parents []string
}

func (object *Object) String() string {
	return fmt.Sprintf("Object:%s\nFormat:%s\n", object.name, object.format)
}

func AddObj(name string, format string, tags []string) error {
	obj := &Object{
		name:   name,
		format: format,
	}

	for _, tag_name := range tags {
		if _, exist := tagMap[tag_name]; !exist {
			return fmt.Errorf("Tag '%s' doesn't exist in engine", tag_name)
		}
	}

	for _, tag_name := range tags {
		tag, _ := tagMap[tag_name]
		obj.parents = append(obj.parents, tag_name)
		tag.objects = append(tag.objects, name)
	}

	objectMap[name] = obj

	return nil
}

func DelObj(name string) error {
	object, exist := objectMap[name]
	if !exist {
		return fmt.Errorf("Object %s not found", name)
	}

	for _, parent := range object.parents {
		parent := tagMap[parent]
		parent.objects, _ = delItemFromSlice(parent.objects, name)
	}

	objectMap[name] = nil

	return nil
}
