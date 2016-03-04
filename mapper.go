package qb

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/serenize/snaker"
	"strings"
)

const tagPrefix = "qb"

// NewMapper instantiates a new mapper object and returns it as a mapper pointer
func NewMapper(driver string) *Mapper {
	return &Mapper{
		driver: driver,
	}
}

// Mapper is the generic struct for struct to table mapping
type Mapper struct {
	driver string
}

func (m *Mapper) extractValue(value string) string {

	hasParams := strings.Contains(value, "(") && strings.Contains(value, ")")

	if hasParams {
		startIndex := strings.Index(value, "(")
		endIndex := strings.Index(value, ")")
		return value[startIndex+1 : endIndex]
	}

	return ""
}

// ConvertMap converts a model struct to a map
func (m *Mapper) ConvertStructToMap(model interface{}) map[string]interface{} {
	return structs.Map(model)
}

// ModelName returns the table name of model
func (m *Mapper) ModelName(model interface{}) string {
	return snaker.CamelToSnake(structs.Name(model))
}

// ColName returns the column name of model
func (m *Mapper) ColName(col string) string {
	return snaker.CamelToSnake(col)
}

// ConvertType returns the type mapping of column.
// If tagType is, then colType would automatically be resolved.
// If tagType is not "", then automatic type resolving would be overridden by tagType
func (m *Mapper) ConvertType(colType string, tagType string) *Type {

	// convert tagType
	if tagType != "" {
		tagType = strings.ToUpper(tagType)
		return &Type{func() string { return tagType }}
	}

	// convert default type
	switch colType {
	case "string":
		return VarChar()
	case "int":
		return Int()
	case "int64":
		return BigInt()
	case "float32":
		return Float()
	case "float64":
		return Float()
	case "bool":
		return Boolean()
	case "time.Time":
		return Timestamp()
	case "*time.Time":
		return Timestamp()
	default:
		return VarChar()
	}
}

// Convert parses struct and converts it to a new table
func (m *Mapper) Convert(model interface{}) (*Table, error) {

	modelName := m.ModelName(model)

	table := &Table{
		name:        modelName,
		columns:     []Column{},
		constraints: []Constraint{},
		driver: m.driver,
	}

	//fmt.Printf("model name: %s\n\n", modelName)

	var col Column
	var rawTag string

	for _, f := range structs.Fields(model) {

		colName := m.ColName(f.Name())
		colType := fmt.Sprintf("%T", f.Value())

		rawTag = f.Tag(tagPrefix)

		constraints := []Constraint{}
		//fmt.Printf("field name: %s\n", colName)
		//fmt.Printf("field raw tag: %s\n", rawTag)
		//fmt.Printf("field type name: %T\n", f.Value())
		//fmt.Printf("field constraints: %v\n", constraints)

		// clean trailing spaces of tag
		rawTag = strings.Replace(f.Tag(tagPrefix), " ", "", -1)

		// parse tag
		tag, err := ParseTag(rawTag)
		if err != nil {
			return nil, err
		}

		// convert tag into constraints
		var constraint Constraint
		for _, v := range tag.Constraints {
			if v == "null" {
				constraint = Null()
			} else if v == "notnull" || v == "not_null" {
				constraint = NotNull()
			} else if v == "unique" {
				constraint = Constraint{
					Name: "UNIQUE",
				}
			} else if v == "auto_increment" || v == "autoincrement" {
				if m.driver == "mysql" {
					constraint = Constraint{
						Name: "AUTO_INCREMENT",
					}
				} else if m.driver == "sqlite" {
					constraint = Constraint{
						Name: "AUTOINCREMENT",
					}
				} else {
					continue
				}
			} else if strings.Contains(v, "default") {
				constraint = Default(m.extractValue(v))
			} else if strings.Contains(v, "primary_key") {
				table.AddPrimary(colName)
				continue
			} else if strings.Contains(v, "ref") && strings.Contains(v, "(") && strings.Contains(v, ")") {
				tc := strings.Split(m.extractValue(v), ".")
				table.AddRef(colName, tc[0], tc[1])
				continue
			} else {
				return nil, fmt.Errorf("Invalid constraint: %s", v)
			}
			constraints = append(constraints, constraint)
		}

		//fmt.Printf("field tag.Type: %s\n", tag.Type)
		//fmt.Printf("field tag.Constraints: %v\n", tag.Constraints)

		col = Column{
			Name:        colName,
			Constraints: constraints,
			Type:        m.ConvertType(colType, tag.Type),
		}

		table.AddColumn(col)

		//fmt.Println()
	}

	return table, nil
}
