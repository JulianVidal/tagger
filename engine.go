package main

import (
	"errors"
	"fmt"
)

func delItemFromSlice[S ~[]I, I comparable](s *S, a I) error {
	for i, b := range *s {
		if a == b {
			(*s)[i] = (*s)[len(*s)-1]
			*s = (*s)[:len(*s)-1]
		}
	}
	return errors.New("Item to delete not found in slice")
}

type Object struct {
	name    string
	format  string
	parents []*Node
}

type Node struct {
	name     string
	children []*Node
	parents  []*Node
	objects  []*Object
}

func (node *Node) print() {
	fmt.Printf("node: %s\n", node.name)
	for _, doc := range node.objects {
		fmt.Printf("\tName:%s\n\tFormat:%s\n", doc.name, doc.format)
	}
	for _, child := range node.children {
		child.print()
	}
}

type Engine struct {
	table map[string]*Node
}

func newEngine() *Engine {
	root := Node{
		name: "Root",
	}
	table := make(map[string]*Node)
	table[root.name] = &root

	return &Engine{
		table: table,
	}
}

func (eng *Engine) addTag(name string, parents []string) {
	node := Node{
		name: name,
	}

	if len(parents) == 0 {
		parents = append(parents, "Root")
	}

	for _, parent_name := range parents {
		parent := eng.table[parent_name]
		parent.children = append(parent.children, &node)
		node.parents = append(node.parents, parent)
	}

	eng.table[name] = &node
}

func (eng *Engine) delTag(name string) error {
	node, exist := eng.table[name]
	if !exist {
		return errors.New("Tag not found in table")
	}

	for _, parent := range node.parents {
		// err := delNodeFromSlice(&parent.children, node)

		if err := delItemFromSlice(&parent.children, node); err != nil {
			panic(err)
		}
	}

	for _, child := range node.children {
		if err := delItemFromSlice(&child.parents, node); err != nil {
			panic(err)
		}
	}

	eng.table["Root"].objects = append(eng.table["Root"].objects, node.objects...)

	eng.table[name] = nil

	return nil
}

func (eng *Engine) addObj(obj *Object, tags []string) {
	if len(tags) == 0 {
		tags = append(tags, "Root")
	}
	for _, tag_name := range tags {
		tag := eng.table[tag_name]
		obj.parents = append(obj.parents, tag)
		tag.objects = append(tag.objects, obj)
	}
}

func (eng *Engine) delObj(obj *Object) {
	for _, parent := range obj.parents {
		if err := delItemFromSlice(&parent.objects, obj); err != nil {
			panic(err)
		}
	}
}

func (eng *Engine) print() {
	eng.table["Root"].print()
}
