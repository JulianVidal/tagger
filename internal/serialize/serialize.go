package serialize

type Tag struct {
	Name string
	Tags []string
}

type Obj struct {
	Name   string
	Format string
	Tags   []string
}

type Engine struct {
	Tags map[string]Tag
	Objs map[string]Obj
}
