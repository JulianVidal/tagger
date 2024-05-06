package engine

import "fmt"

type Tag struct {
	name     string
	children []string
	parents  []string
	objects  []string
}

func (tag *Tag) print() {
	fmt.Println(tag)
}

func (tag *Tag) String() string {
	str := fmt.Sprintf("----------------------------------\n")
	str += fmt.Sprintf("Tag: %s\n", tag.name)
	for _, objectName := range tag.objects {
		obj := objectMap[objectName]
		str += fmt.Sprintf("%s\n", obj)
	}
	str += fmt.Sprintf("----------------------------------\n")
	return str
}

func AddTag(name string, parents []string) error {
	if _, exist := tagMap[name]; exist {
		return fmt.Errorf("Tag '%s' already exists", name)
	}

	node := &Tag{
		name: name,
	}

	for _, parent_name := range parents {
		if _, exist := tagMap[parent_name]; !exist {
			return fmt.Errorf("Parent tag '%s' not found", parent_name)
		}
	}

	for _, parent_name := range parents {
		parent, _ := tagMap[parent_name]
		parent.children = append(parent.children, name)
		node.parents = append(node.parents, parent.name)
	}

	tagMap[name] = node

	return nil
}

func DelTag(name string) error {
	tag, exist := tagMap[name]
	if !exist {
		return fmt.Errorf("Tag not found in engine: %s", name)
	}

	for _, childName := range tag.children {
		child := tagMap[childName]
		child.parents, _ = delItemFromSlice(child.parents, tag.name)
	}

	for _, parentName := range tag.parents {
		parent := tagMap[parentName]
		parent.children, _ = delItemFromSlice(parent.children, tag.name)
	}

	tagMap[name] = nil

	return nil
}
