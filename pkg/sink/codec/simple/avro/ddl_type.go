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

	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
)

var _ = fmt.Printf

type DDLType int32

const (
	DDLTypeCREATE   DDLType = 0
	DDLTypeALTER    DDLType = 1
	DDLTypeERASE    DDLType = 2
	DDLTypeRENAME   DDLType = 3
	DDLTypeTRUNCATE DDLType = 4
	DDLTypeCINDEX   DDLType = 5
	DDLTypeDINDEX   DDLType = 6
	DDLTypeQUERY    DDLType = 7
)

func (e DDLType) String() string {
	switch e {
	case DDLTypeCREATE:
		return "CREATE"
	case DDLTypeALTER:
		return "ALTER"
	case DDLTypeERASE:
		return "ERASE"
	case DDLTypeRENAME:
		return "RENAME"
	case DDLTypeTRUNCATE:
		return "TRUNCATE"
	case DDLTypeCINDEX:
		return "CINDEX"
	case DDLTypeDINDEX:
		return "DINDEX"
	case DDLTypeQUERY:
		return "QUERY"
	}
	return "unknown"
}

func writeDDLType(r DDLType, w io.Writer) error {
	return vm.WriteInt(int32(r), w)
}

func NewDDLTypeValue(raw string) (r DDLType, err error) {
	switch raw {
	case "CREATE":
		return DDLTypeCREATE, nil
	case "ALTER":
		return DDLTypeALTER, nil
	case "ERASE":
		return DDLTypeERASE, nil
	case "RENAME":
		return DDLTypeRENAME, nil
	case "TRUNCATE":
		return DDLTypeTRUNCATE, nil
	case "CINDEX":
		return DDLTypeCINDEX, nil
	case "DINDEX":
		return DDLTypeDINDEX, nil
	case "QUERY":
		return DDLTypeQUERY, nil
	}

	return -1, fmt.Errorf("invalid value for DDLType: '%s'", raw)

}

func (b DDLType) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *DDLType) UnmarshalJSON(data []byte) error {
	var stringVal string
	err := json.Unmarshal(data, &stringVal)
	if err != nil {
		return err
	}
	val, err := NewDDLTypeValue(stringVal)
	*b = val
	return err
}

type DDLTypeWrapper struct {
	Target *DDLType
}

func (b DDLTypeWrapper) SetBoolean(v bool) {
	panic("Unable to assign boolean to int field")
}

func (b DDLTypeWrapper) SetInt(v int32) {
	*(b.Target) = DDLType(v)
}

func (b DDLTypeWrapper) SetLong(v int64) {
	panic("Unable to assign long to int field")
}

func (b DDLTypeWrapper) SetFloat(v float32) {
	panic("Unable to assign float to int field")
}

func (b DDLTypeWrapper) SetUnionElem(v int64) {
	panic("Unable to assign union elem to int field")
}

func (b DDLTypeWrapper) SetDouble(v float64) {
	panic("Unable to assign double to int field")
}

func (b DDLTypeWrapper) SetBytes(v []byte) {
	panic("Unable to assign bytes to int field")
}

func (b DDLTypeWrapper) SetString(v string) {
	panic("Unable to assign string to int field")
}

func (b DDLTypeWrapper) Get(i int) types.Field {
	panic("Unable to get field from int field")
}

func (b DDLTypeWrapper) SetDefault(i int) {
	panic("Unable to set default on int field")
}

func (b DDLTypeWrapper) AppendMap(key string) types.Field {
	panic("Unable to append map key to from int field")
}

func (b DDLTypeWrapper) AppendArray() types.Field {
	panic("Unable to append array element to from int field")
}

func (b DDLTypeWrapper) NullField(int) {
	panic("Unable to null field in int field")
}

func (b DDLTypeWrapper) HintSize(int) {
	panic("Unable to hint size in int field")
}

func (b DDLTypeWrapper) Finalize() {}