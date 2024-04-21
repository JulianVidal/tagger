package main

import (
	"errors"
	"fmt"
)

type Object struct {
	name   string
	format string
}

type Node struct {
	name     string
	children []*Node
	parents  []*Node
	objects  []*Object
}

type Engine struct {
	graph *Node
	table map[string]*Node
}

func newEngine() *Engine {
	root := Node{
		name: "Root",
	}
	table := make(map[string]*Node)
	table[root.name] = &root

	return &Engine{
		graph: &root,
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

func (eng *Engine) addObj(obj *Object, tags []string) {
	if len(tags) == 0 {
		tags = append(tags, "Root")
	}
	for _, tag_name := range tags {
		tag := eng.table[tag_name]
		tag.objects = append(tag.objects, obj)
	}
}

func delNodeFromSlice(nodes *[]*Node, nodeA *Node) error {
	for i, nodeB := range *nodes {
		if nodeA == nodeB {
			(*nodes)[i] = (*nodes)[len(*nodes)-1]
			*nodes = (*nodes)[:len(*nodes)-1]
		}
	}
	return errors.New("Node not found in slice")
}

func (eng *Engine) delTag(name string) error {
	node, exist := eng.table[name]
	if !exist {
		return errors.New("Tag not found in table")
	}

	for _, parent := range node.parents {
		err := delNodeFromSlice(&parent.children, node)
		if err != nil {
			panic(err)
		}
	}

	for _, child := range node.children {
		err := delNodeFromSlice(&child.parents, node)
		if err != nil {
			panic(err)
		}
	}

	eng.table[name] = nil

	return nil
}

func (eng *Engine) print() {
	eng.graph.print()
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

func main() {
	engine := newEngine()

	engine.addTag("BT", []string{})
	engine.addTag("Sat", []string{"BT"})

	lab_poc := Object{
		name:   "lab_doc.doc",
		format: "Word Document",
	}

	start_doc := Object{
		name:   "start.ptx",
		format: "Powerpoint Presentation",
	}

	engine.addObj(&lab_poc, []string{"Sat"})
	engine.addObj(&start_doc, []string{"BT"})
	engine.print()
}
