package main

import (
	"os"

	"github.com/JulianVidal/tagger/internal/engine"
)

func main() {
	engine.InitEngine()
	data, err := engine.ToJson()
	if err != nil {
		panic(err)
	}

	bt_tag, err := engine.NewTag("BT", []string{})
	if err != nil {
		panic(err)
	}
	err = engine.AddTag(bt_tag)
	if err != nil {
		panic(err)
	}

	sat_tag, err := engine.NewTag("Sat", []string{"BT"})
	if err != nil {
		panic(err)
	}
	err = engine.AddTag(sat_tag)
	if err != nil {
		panic(err)
	}

	sat_obj, err := engine.NewObject("/home/julian/temp/satshelf.txt", []string{"Sat"})
	if err != nil {
		panic(err)
	}
	err = engine.AddObject(sat_obj)
	if err != nil {
		panic(err)
	}

	start_obj, err := engine.NewObject("/home/julian/temp/starter.txt", []string{"BT"})
	if err != nil {
		panic(err)
	}
	err = engine.AddObject(start_obj)
	if err != nil {
		panic(err)
	}

	engine.Print()

	objs, err := engine.Query(bt_tag)
	if err != nil {
		panic(err)
	}
	for _, obj := range objs {
		obj.Print()
	}

	err = os.WriteFile("engine.json", data, 0644)
	if err != nil {
		panic(err)
	}
}
