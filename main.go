package main

// import "github.com/JulianVidal/tagger/cmd/tagger"
import "github.com/JulianVidal/tagger/internal/engine"

func main() {
	// tagger.Execute()
	eng := engine.NewEngine()

	eng.AddTag("BT", []string{})
	eng.AddTag("Sat", []string{"BT"})

	eng.AddObj("lab_doc.doc", "Word Document", []string{"Sat"})
	eng.AddObj("start.ptx", "Powerpoint Presentation", []string{"BT"})

	println("Query: BT")
	for _, obj := range eng.Query([]string{"BT"}) {
		obj.Print()
	}

	println("Query: Sat")
	for _, obj := range eng.Query([]string{"Sat"}) {
		obj.Print()
	}

	println("All:")

	eng.Print()
}
