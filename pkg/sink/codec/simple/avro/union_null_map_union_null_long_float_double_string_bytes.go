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

type UnionNullMapUnionNullLongFloatDoubleStringBytesTypeEnum int

const (
	UnionNullMapUnionNullLongFloatDoubleStringBytesTypeEnumMapUnionNullLongFloatDoubleStringBytes UnionNullMapUnionNullLongFloatDoubleStringBytesTypeEnum = 1
)

type UnionNullMapUnionNullLongFloatDoubleStringBytes struct {
	Null                                   *types.NullVal
	MapUnionNullLongFloatDoubleStringBytes map[string]*UnionNullLongFloatDoubleStringBytes
	UnionType                              UnionNullMapUnionNullLongFloatDoubleStringBytesTypeEnum
}

func writeUnionNullMapUnionNullLongFloatDoubleStringBytes(r *UnionNullMapUnionNullLongFloatDoubleStringBytes, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(0, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullMapUnionNullLongFloatDoubleStringBytesTypeEnumMapUnionNullLongFloatDoubleStringBytes:
		return writeMapUnionNullLongFloatDoubleStringBytes(r.MapUnionNullLongFloatDoubleStringBytes, w)
	}
	return fmt.Errorf("invalid value for *UnionNullMapUnionNullLongFloatDoubleStringBytes")
}

func NewUnionNullMapUnionNullLongFloatDoubleStringBytes() *UnionNullMapUnionNullLongFloatDoubleStringBytes {
	return &UnionNullMapUnionNullLongFloatDoubleStringBytes{}
}

func (r *UnionNullMapUnionNullLongFloatDoubleStringBytes) Serialize(w io.Writer) error {
	return writeUnionNullMapUnionNullLongFloatDoubleStringBytes(r, w)
}

func DeserializeUnionNullMapUnionNullLongFloatDoubleStringBytes(r io.Reader) (*UnionNullMapUnionNullLongFloatDoubleStringBytes, error) {
	t := NewUnionNullMapUnionNullLongFloatDoubleStringBytes()
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

func DeserializeUnionNullMapUnionNullLongFloatDoubleStringBytesFromSchema(r io.Reader, schema string) (*UnionNullMapUnionNullLongFloatDoubleStringBytes, error) {
	t := NewUnionNullMapUnionNullLongFloatDoubleStringBytes()
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

func (r *UnionNullMapUnionNullLongFloatDoubleStringBytes) Schema() string {
	return "[\"null\",{\"default\":null,\"type\":\"map\",\"values\":[\"null\",\"long\",\"float\",\"double\",\"string\",\"bytes\"]}]"
}

func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) SetBoolean(v bool) {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) SetInt(v int32) {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) SetFloat(v float32) {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) SetDouble(v float64) {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) SetBytes(v []byte) {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) SetString(v string) {
	panic("Unsupported operation")
}

func (r *UnionNullMapUnionNullLongFloatDoubleStringBytes) SetLong(v int64) {

	r.UnionType = (UnionNullMapUnionNullLongFloatDoubleStringBytesTypeEnum)(v)
}

func (r *UnionNullMapUnionNullLongFloatDoubleStringBytes) Get(i int) types.Field {

	switch i {
	case 0:
		return r.Null
	case 1:
		r.MapUnionNullLongFloatDoubleStringBytes = make(map[string]*UnionNullLongFloatDoubleStringBytes)
		return &MapUnionNullLongFloatDoubleStringBytesWrapper{Target: (&r.MapUnionNullLongFloatDoubleStringBytes)}
	}
	panic("Unknown field index")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) NullField(i int) {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) HintSize(i int) {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) SetDefault(i int) {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) AppendMap(key string) types.Field {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) AppendArray() types.Field {
	panic("Unsupported operation")
}
func (_ *UnionNullMapUnionNullLongFloatDoubleStringBytes) Finalize() {}

func (r *UnionNullMapUnionNullLongFloatDoubleStringBytes) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionNullMapUnionNullLongFloatDoubleStringBytesTypeEnumMapUnionNullLongFloatDoubleStringBytes:
		return json.Marshal(map[string]interface{}{"map": r.MapUnionNullLongFloatDoubleStringBytes})
	}
	return nil, fmt.Errorf("invalid value for *UnionNullMapUnionNullLongFloatDoubleStringBytes")
}

func (r *UnionNullMapUnionNullLongFloatDoubleStringBytes) UnmarshalJSON(data []byte) error {

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	if value, ok := fields["map"]; ok {
		r.UnionType = 1
		return json.Unmarshal([]byte(value), &r.MapUnionNullLongFloatDoubleStringBytes)
	}
	return fmt.Errorf("invalid value for *UnionNullMapUnionNullLongFloatDoubleStringBytes")
}