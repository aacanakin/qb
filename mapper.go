package qbit

import (
	"fmt"
	"github.com/fatih/structs"
)

func NewMapper() *Mapper {
	return &Mapper{}
}

type Mapper struct {
}

func (m *Mapper) Convert(model interface{}) *Table {

	//	fmt.Printf("%s\n", structs.Fields(model))

	for _, f := range structs.Fields(model) {
		fmt.Printf("field name: %s, type: %s, tag: %s;\n", f.Name(), f.Kind(), f.Tag("qbit"))
	}
	//	modelType := reflect.TypeOf(model)
	//	for i := 0; i < modelType.NumField(); i++ {
	//		fmt.Printf("field name: %s\n", modelType.Field(i).Name)
	//		fmt.Printf("field tag: %s\n", string(modelType.Field(i).Tag))
	//	}
	return NewTable("-", []Column{}, []Constraint{})
}
