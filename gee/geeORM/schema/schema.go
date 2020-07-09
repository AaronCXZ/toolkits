package schema

import (
	"go/ast"
	"reflect"

	"github.com/Muskchen/toolkits/gee/geeORM/dialect"
)

type Field struct {
	Name string // 字段名
	Type string // 类型
	Tag  string // 约束条件
}

type Schema struct {
	Model      interface{}       // 被映射的对象
	Name       string            // 表名
	Fields     []*Field          // 字段
	FieldNames []string          // 字段名
	fieldMap   map[string]*Field // 字段名和Field的映射关系
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func (schema *Schema) RecordValue(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(), // 结构体名称
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		// 根据下标获取字段
		p := modelType.Field(i)
		// @p.Name: 字段名
		// @p.Type：字段类型
		// @p.Tag：额外的约束条件
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
