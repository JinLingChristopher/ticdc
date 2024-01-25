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

var _ = fmt.Printf

type TableSchema struct {
	Database string `json:"database"`

	Table string `json:"table"`

	TableID int64 `json:"tableID"`

	Version int64 `json:"version"`

	Columns []ColumnSchema `json:"columns"`

	Indexes []IndexSchema `json:"indexes"`
}

const TableSchemaAvroCRC64Fingerprint = "/\xaa\xf8\x96\xed.w!"

func NewTableSchema() TableSchema {
	r := TableSchema{}
	r.Columns = make([]ColumnSchema, 0)

	r.Indexes = make([]IndexSchema, 0)

	return r
}

func DeserializeTableSchema(r io.Reader) (TableSchema, error) {
	t := NewTableSchema()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeTableSchemaFromSchema(r io.Reader, schema string) (TableSchema, error) {
	t := NewTableSchema()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeTableSchema(r TableSchema, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Database, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Table, w)
	if err != nil {
		return err
	}
	err = vm.WriteLong(r.TableID, w)
	if err != nil {
		return err
	}
	err = vm.WriteLong(r.Version, w)
	if err != nil {
		return err
	}
	err = writeArrayColumnSchema(r.Columns, w)
	if err != nil {
		return err
	}
	err = writeArrayIndexSchema(r.Indexes, w)
	if err != nil {
		return err
	}
	return err
}

func (r TableSchema) Serialize(w io.Writer) error {
	return writeTableSchema(r, w)
}

func (r TableSchema) Schema() string {
	return "{\"docs\":\"table schema information\",\"fields\":[{\"name\":\"database\",\"type\":\"string\"},{\"name\":\"table\",\"type\":\"string\"},{\"name\":\"tableID\",\"type\":\"long\"},{\"name\":\"version\",\"type\":\"long\"},{\"name\":\"columns\",\"type\":{\"items\":{\"docs\":\"each column's schema information\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"dataType\",\"type\":{\"docs\":\"each column's mysql type information\",\"fields\":[{\"name\":\"mysqlType\",\"type\":\"string\"},{\"name\":\"charset\",\"type\":\"string\"},{\"name\":\"collate\",\"type\":\"string\"},{\"name\":\"length\",\"type\":\"long\"},{\"default\":null,\"name\":\"decimal\",\"type\":[\"null\",\"int\"]},{\"default\":null,\"name\":\"elements\",\"type\":[\"null\",{\"items\":\"string\",\"type\":\"array\"}]},{\"default\":null,\"name\":\"unsigned\",\"type\":[\"null\",\"boolean\"]},{\"default\":null,\"name\":\"zerofill\",\"type\":[\"null\",\"boolean\"]}],\"name\":\"DataType\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"}},{\"name\":\"nullable\",\"type\":\"boolean\"},{\"name\":\"default\",\"type\":[\"null\",\"string\"]}],\"name\":\"ColumnSchema\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"},\"type\":\"array\"}},{\"name\":\"indexes\",\"type\":{\"items\":{\"docs\":\"each index's schema information\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"unique\",\"type\":\"boolean\"},{\"name\":\"primary\",\"type\":\"boolean\"},{\"name\":\"nullable\",\"type\":\"boolean\"},{\"name\":\"columns\",\"type\":{\"items\":\"string\",\"type\":\"array\"}}],\"name\":\"IndexSchema\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"},\"type\":\"array\"}}],\"name\":\"com.pingcap.simple.avro.TableSchema\",\"type\":\"record\"}"
}

func (r TableSchema) SchemaName() string {
	return "com.pingcap.simple.avro.TableSchema"
}

func (_ TableSchema) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ TableSchema) SetInt(v int32)       { panic("Unsupported operation") }
func (_ TableSchema) SetLong(v int64)      { panic("Unsupported operation") }
func (_ TableSchema) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ TableSchema) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ TableSchema) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ TableSchema) SetString(v string)   { panic("Unsupported operation") }
func (_ TableSchema) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *TableSchema) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Database}

		return w

	case 1:
		w := types.String{Target: &r.Table}

		return w

	case 2:
		w := types.Long{Target: &r.TableID}

		return w

	case 3:
		w := types.Long{Target: &r.Version}

		return w

	case 4:
		r.Columns = make([]ColumnSchema, 0)

		w := ArrayColumnSchemaWrapper{Target: &r.Columns}

		return w

	case 5:
		r.Indexes = make([]IndexSchema, 0)

		w := ArrayIndexSchemaWrapper{Target: &r.Indexes}

		return w

	}
	panic("Unknown field index")
}

func (r *TableSchema) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *TableSchema) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ TableSchema) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ TableSchema) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ TableSchema) HintSize(int)                     { panic("Unsupported operation") }
func (_ TableSchema) Finalize()                        {}

func (_ TableSchema) AvroCRC64Fingerprint() []byte {
	return []byte(TableSchemaAvroCRC64Fingerprint)
}

func (r TableSchema) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["database"], err = json.Marshal(r.Database)
	if err != nil {
		return nil, err
	}
	output["table"], err = json.Marshal(r.Table)
	if err != nil {
		return nil, err
	}
	output["tableID"], err = json.Marshal(r.TableID)
	if err != nil {
		return nil, err
	}
	output["version"], err = json.Marshal(r.Version)
	if err != nil {
		return nil, err
	}
	output["columns"], err = json.Marshal(r.Columns)
	if err != nil {
		return nil, err
	}
	output["indexes"], err = json.Marshal(r.Indexes)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *TableSchema) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["database"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Database); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for database")
	}
	val = func() json.RawMessage {
		if v, ok := fields["table"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Table); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for table")
	}
	val = func() json.RawMessage {
		if v, ok := fields["tableID"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.TableID); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for tableID")
	}
	val = func() json.RawMessage {
		if v, ok := fields["version"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Version); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for version")
	}
	val = func() json.RawMessage {
		if v, ok := fields["columns"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Columns); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for columns")
	}
	val = func() json.RawMessage {
		if v, ok := fields["indexes"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Indexes); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for indexes")
	}
	return nil
}