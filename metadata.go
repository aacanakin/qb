package qb

import (
	"fmt"
	"strings"
)

// NewMetaData creates a new MetaData object and returns
func NewMetaData(engine *Engine) *MetaData {

	return &MetaData{
		tables: []*Table{},
		engine: engine,
		mapper: NewMapper(engine.Driver()),
	}
}

// MetaData is the container for database structs and tables
type MetaData struct {
	tables []*Table
	engine *Engine
	mapper *Mapper
}

// Add retrieves the struct and converts it using mapper and appends to tables slice
func (m *MetaData) Add(model interface{}) {

	table, err := m.mapper.Convert(model)
	if err != nil {
		panic(err)
	}

	m.AddTable(table)
}

// AddTable appends table to tables slice
func (m *MetaData) AddTable(table *Table) {
	m.tables = append(m.tables, table)
}

func (m *MetaData) Table(name string) *Table {

	if m.engine.Driver() != "postgres" {
		name = fmt.Sprintf("`%s`", name)
	}

	for _, t := range m.tables {
		if t.name == name {
			return t
		}
	}

	return nil
}

// Tables returns the current tables slice
func (m *MetaData) Tables() []*Table {
	return m.tables
}

// CreateAll creates all the tables added to metadata
func (m *MetaData) CreateAll(engine *Engine) error {

	sqls := []string{}
	for _, t := range m.tables {
		sqls = append(sqls, t.SQL())
	}

	_, err := engine.DB().Exec(strings.Join(sqls, "\n"))
	return err
}

// DropAll drops all the tables which is added to metadata
func (m *MetaData) DropAll(engine *Engine) error {

	sqls := []string{}
	for i := len(m.tables) - 1; i >= 0; i-- {
		sqls = append(sqls, fmt.Sprintf("DROP TABLE %s;", m.tables[i].Name()))
	}

	_, err := engine.DB().Exec(strings.Join(sqls, "\n"))

	return err
}
