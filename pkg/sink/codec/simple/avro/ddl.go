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

type DDL struct {
	Version int32 `json:"version"`

	Type DDLType `json:"type"`

	Sql string `json:"sql"`

	CommitTs int64 `json:"commitTs"`

	BuildTs int64 `json:"buildTs"`

	TableSchema *UnionNullTableSchema `json:"tableSchema"`

	PreTableSchema *UnionNullTableSchema `json:"preTableSchema"`
}

const DDLAvroCRC64Fingerprint = "\x87\x99\xe5Ը\x9dl\x93"

func NewDDL() DDL {
	r := DDL{}
	r.TableSchema = nil
	r.PreTableSchema = nil
	return r
}

func DeserializeDDL(r io.Reader) (DDL, error) {
	t := NewDDL()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeDDLFromSchema(r io.Reader, schema string) (DDL, error) {
	t := NewDDL()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeDDL(r DDL, w io.Writer) error {
	var err error
	err = vm.WriteInt(r.Version, w)
	if err != nil {
		return err
	}
	err = writeDDLType(r.Type, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Sql, w)
	if err != nil {
		return err
	}
	err = vm.WriteLong(r.CommitTs, w)
	if err != nil {
		return err
	}
	err = vm.WriteLong(r.BuildTs, w)
	if err != nil {
		return err
	}
	err = writeUnionNullTableSchema(r.TableSchema, w)
	if err != nil {
		return err
	}
	err = writeUnionNullTableSchema(r.PreTableSchema, w)
	if err != nil {
		return err
	}
	return err
}

func (r DDL) Serialize(w io.Writer) error {
	return writeDDL(r, w)
}

func (r DDL) Schema() string {
	return "{\"docs\":\"the message format of the DDL event\",\"fields\":[{\"name\":\"version\",\"type\":\"int\"},{\"name\":\"type\",\"type\":{\"name\":\"DDLType\",\"symbols\":[\"CREATE\",\"ALTER\",\"ERASE\",\"RENAME\",\"TRUNCATE\",\"CINDEX\",\"DINDEX\",\"QUERY\"],\"type\":\"enum\"}},{\"name\":\"sql\",\"type\":\"string\"},{\"name\":\"commitTs\",\"type\":\"long\"},{\"name\":\"buildTs\",\"type\":\"long\"},{\"default\":null,\"name\":\"tableSchema\",\"type\":[\"null\",{\"docs\":\"table schema information\",\"fields\":[{\"name\":\"database\",\"type\":\"string\"},{\"name\":\"table\",\"type\":\"string\"},{\"name\":\"tableID\",\"type\":\"long\"},{\"name\":\"version\",\"type\":\"long\"},{\"name\":\"columns\",\"type\":{\"items\":{\"docs\":\"each column's schema information\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"dataType\",\"type\":{\"docs\":\"each column's mysql type information\",\"fields\":[{\"name\":\"mysqlType\",\"type\":\"string\"},{\"name\":\"charset\",\"type\":\"string\"},{\"name\":\"collate\",\"type\":\"string\"},{\"name\":\"length\",\"type\":\"long\"},{\"default\":null,\"name\":\"decimal\",\"type\":[\"null\",\"int\"]},{\"default\":null,\"name\":\"elements\",\"type\":[\"null\",{\"items\":\"string\",\"type\":\"array\"}]},{\"default\":null,\"name\":\"unsigned\",\"type\":[\"null\",\"boolean\"]},{\"default\":null,\"name\":\"zerofill\",\"type\":[\"null\",\"boolean\"]}],\"name\":\"DataType\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"}},{\"name\":\"nullable\",\"type\":\"boolean\"},{\"name\":\"default\",\"type\":[\"null\",\"string\"]}],\"name\":\"ColumnSchema\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"},\"type\":\"array\"}},{\"name\":\"indexes\",\"type\":{\"items\":{\"docs\":\"each index's schema information\",\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"unique\",\"type\":\"boolean\"},{\"name\":\"primary\",\"type\":\"boolean\"},{\"name\":\"nullable\",\"type\":\"boolean\"},{\"name\":\"columns\",\"type\":{\"items\":\"string\",\"type\":\"array\"}}],\"name\":\"IndexSchema\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"},\"type\":\"array\"}}],\"name\":\"TableSchema\",\"namespace\":\"com.pingcap.simple.avro\",\"type\":\"record\"}]},{\"default\":null,\"name\":\"preTableSchema\",\"type\":[\"null\",\"com.pingcap.simple.avro.TableSchema\"]}],\"name\":\"com.pingcap.simple.avro.DDL\",\"type\":\"record\"}"
}

func (r DDL) SchemaName() string {
	return "com.pingcap.simple.avro.DDL"
}

func (_ DDL) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ DDL) SetInt(v int32)       { panic("Unsupported operation") }
func (_ DDL) SetLong(v int64)      { panic("Unsupported operation") }
func (_ DDL) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ DDL) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ DDL) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ DDL) SetString(v string)   { panic("Unsupported operation") }
func (_ DDL) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *DDL) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.Int{Target: &r.Version}

		return w

	case 1:
		w := DDLTypeWrapper{Target: &r.Type}

		return w

	case 2:
		w := types.String{Target: &r.Sql}

		return w

	case 3:
		w := types.Long{Target: &r.CommitTs}

		return w

	case 4:
		w := types.Long{Target: &r.BuildTs}

		return w

	case 5:
		r.TableSchema = NewUnionNullTableSchema()

		return r.TableSchema
	case 6:
		r.PreTableSchema = NewUnionNullTableSchema()

		return r.PreTableSchema
	}
	panic("Unknown field index")
}

func (r *DDL) SetDefault(i int) {
	switch i {
	case 5:
		r.TableSchema = nil
		return
	case 6:
		r.PreTableSchema = nil
		return
	}
	panic("Unknown field index")
}

func (r *DDL) NullField(i int) {
	switch i {
	case 5:
		r.TableSchema = nil
		return
	case 6:
		r.PreTableSchema = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ DDL) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ DDL) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ DDL) HintSize(int)                     { panic("Unsupported operation") }
func (_ DDL) Finalize()                        {}

func (_ DDL) AvroCRC64Fingerprint() []byte {
	return []byte(DDLAvroCRC64Fingerprint)
}

func (r DDL) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["version"], err = json.Marshal(r.Version)
	if err != nil {
		return nil, err
	}
	output["type"], err = json.Marshal(r.Type)
	if err != nil {
		return nil, err
	}
	output["sql"], err = json.Marshal(r.Sql)
	if err != nil {
		return nil, err
	}
	output["commitTs"], err = json.Marshal(r.CommitTs)
	if err != nil {
		return nil, err
	}
	output["buildTs"], err = json.Marshal(r.BuildTs)
	if err != nil {
		return nil, err
	}
	output["tableSchema"], err = json.Marshal(r.TableSchema)
	if err != nil {
		return nil, err
	}
	output["preTableSchema"], err = json.Marshal(r.PreTableSchema)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *DDL) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
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
		if v, ok := fields["type"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Type); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for type")
	}
	val = func() json.RawMessage {
		if v, ok := fields["sql"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Sql); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for sql")
	}
	val = func() json.RawMessage {
		if v, ok := fields["commitTs"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.CommitTs); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for commitTs")
	}
	val = func() json.RawMessage {
		if v, ok := fields["buildTs"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.BuildTs); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for buildTs")
	}
	val = func() json.RawMessage {
		if v, ok := fields["tableSchema"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.TableSchema); err != nil {
			return err
		}
	} else {
		r.TableSchema = NewUnionNullTableSchema()

		r.TableSchema = nil
	}
	val = func() json.RawMessage {
		if v, ok := fields["preTableSchema"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.PreTableSchema); err != nil {
			return err
		}
	} else {
		r.PreTableSchema = NewUnionNullTableSchema()

		r.PreTableSchema = nil
	}
	return nil
}
