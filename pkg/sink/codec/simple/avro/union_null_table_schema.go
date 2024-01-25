// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCE:
 *     schema.avsc
 */
package avro

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
)

type UnionNullTableSchemaTypeEnum int

const (
	UnionNullTableSchemaTypeEnumTableSchema UnionNullTableSchemaTypeEnum = 1
)

type UnionNullTableSchema struct {
	Null        *types.NullVal
	TableSchema TableSchema
	UnionType   UnionNullTableSchemaTypeEnum
}

func writeUnionNullTableSchema(r *UnionNullTableSchema, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(0, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullTableSchemaTypeEnumTableSchema:
		return writeTableSchema(r.TableSchema, w)
	}
	return fmt.Errorf("invalid value for *UnionNullTableSchema")
}

func NewUnionNullTableSchema() *UnionNullTableSchema {
	return &UnionNullTableSchema{}
}

func (r *UnionNullTableSchema) Serialize(w io.Writer) error {
	return writeUnionNullTableSchema(r, w)
}

func DeserializeUnionNullTableSchema(r io.Reader) (*UnionNullTableSchema, error) {
	t := NewUnionNullTableSchema()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func DeserializeUnionNullTableSchemaFromSchema(r io.Reader, schema string) (*UnionNullTableSchema, error) {
	t := NewUnionNullTableSchema()
	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func (r *UnionNullTableSchema) Schema() string {
	return "[\"null\",{\"docs\":\"table schema information\",\"fields\":[{\"name\":\"database\",\"type\":\"string\"},{\"name\":\"table\",\"type\":\"string\"},{\"name\":\"tableID\",\"type\":\"long\"},{\"name\":\"version\",\"type\":\"long\"},{\"name\":\"columns\",\"type\":{\"items\":{\"docs\":\"each column's schema information\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"dataType\",\"type\":{\"docs\":\"each column's mysql type information\",\"fields\":[{\"name\":\"mysqlType\",\"type\":\"string\"},{\"name\":\"charset\",\"type\":\"string\"},{\"name\":\"collate\",\"type\":\"string\"},{\"name\":\"length\",\"type\":\"long\"},{\"default\":null,\"name\":\"decimal\",\"type\":[\"null\",\"int\"]},{\"default\":null,\"name\":\"elements\",\"type\":[\"null\",{\"items\":\"string\",\"type\":\"array\"}]},{\"default\":null,\"name\":\"unsigned\",\"type\":[\"null\",\"boolean\"]},{\"default\":null,\"name\":\"zerofill\",\"type\":[\"null\",\"boolean\"]}],\"name\":\"DataType\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"}},{\"name\":\"nullable\",\"type\":\"boolean\"},{\"name\":\"default\",\"type\":[\"null\",\"string\"]}],\"name\":\"ColumnSchema\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"},\"type\":\"array\"}},{\"name\":\"indexes\",\"type\":{\"items\":{\"docs\":\"each index's schema information\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"unique\",\"type\":\"boolean\"},{\"name\":\"primary\",\"type\":\"boolean\"},{\"name\":\"nullable\",\"type\":\"boolean\"},{\"name\":\"columns\",\"type\":{\"items\":\"string\",\"type\":\"array\"}}],\"name\":\"IndexSchema\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"},\"type\":\"array\"}}],\"name\":\"TableSchema\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"}]"
}

func (_ *UnionNullTableSchema) SetBoolean(v bool)   { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) SetInt(v int32)      { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) SetFloat(v float32)  { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) SetBytes(v []byte)   { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) SetString(v string)  { panic("Unsupported operation") }

func (r *UnionNullTableSchema) SetLong(v int64) {

	r.UnionType = (UnionNullTableSchemaTypeEnum)(v)
}

func (r *UnionNullTableSchema) Get(i int) types.Field {

	switch i {
	case 0:
		return r.Null
	case 1:
		r.TableSchema = NewTableSchema()
		return &types.Record{Target: (&r.TableSchema)}
	}
	panic("Unknown field index")
}
func (_ *UnionNullTableSchema) NullField(i int)                  { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) HintSize(i int)                   { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) SetDefault(i int)                 { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *UnionNullTableSchema) Finalize()                        {}

func (r *UnionNullTableSchema) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionNullTableSchemaTypeEnumTableSchema:
		return json.Marshal(map[string]interface{}{"com.pingcap.simple.avro.TableSchema": r.TableSchema})
	}
	return nil, fmt.Errorf("invalid value for *UnionNullTableSchema")
}

func (r *UnionNullTableSchema) UnmarshalJSON(data []byte) error {

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	if value, ok := fields["com.pingcap.simple.avro.TableSchema"]; ok {
		r.UnionType = 1
		return json.Unmarshal([]byte(value), &r.TableSchema)
	}
	return fmt.Errorf("invalid value for *UnionNullTableSchema")
}