// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package csv

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/rowcodec"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/pkg/config"
	"github.com/stretchr/testify/require"
)

type csvTestColumnTuple struct {
	col     model.Column
	colInfo rowcodec.ColInfo
	want    interface{}
}

var csvTestColumnsGroup = [][]*csvTestColumnTuple{
	{
		{
			model.Column{Name: "tiny", Value: int64(1), Type: mysql.TypeTiny},
			rowcodec.ColInfo{
				ID:            1,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeTiny),
			},
			int64(1),
		},
		{
			model.Column{Name: "short", Value: int64(1), Type: mysql.TypeShort},
			rowcodec.ColInfo{
				ID:            2,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeShort),
			},
			int64(1),
		},
		{
			model.Column{Name: "int24", Value: int64(1), Type: mysql.TypeInt24},
			rowcodec.ColInfo{
				ID:            3,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeInt24),
			},
			int64(1),
		},
		{
			model.Column{Name: "long", Value: int64(1), Type: mysql.TypeLong},
			rowcodec.ColInfo{
				ID:            4,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeLong),
			},
			int64(1),
		},
		{
			model.Column{Name: "longlong", Value: int64(1), Type: mysql.TypeLonglong},
			rowcodec.ColInfo{
				ID:            5,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeLonglong),
			},
			int64(1),
		},
		{
			model.Column{
				Name:  "tinyunsigned",
				Value: uint64(1),
				Type:  mysql.TypeTiny,
				Flag:  model.UnsignedFlag,
			},
			rowcodec.ColInfo{
				ID:            6,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setFlag(types.NewFieldType(mysql.TypeTiny), uint(model.UnsignedFlag)),
			},
			uint64(1),
		},
		{
			model.Column{
				Name:  "shortunsigned",
				Value: uint64(1),
				Type:  mysql.TypeShort,
				Flag:  model.UnsignedFlag,
			},
			rowcodec.ColInfo{
				ID:            7,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setFlag(types.NewFieldType(mysql.TypeShort), uint(model.UnsignedFlag)),
			},
			uint64(1),
		},
		{
			model.Column{
				Name:  "int24unsigned",
				Value: uint64(1),
				Type:  mysql.TypeInt24,
				Flag:  model.UnsignedFlag,
			},
			rowcodec.ColInfo{
				ID:            8,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setFlag(types.NewFieldType(mysql.TypeInt24), uint(model.UnsignedFlag)),
			},
			uint64(1),
		},
		{
			model.Column{
				Name:  "longunsigned",
				Value: uint64(1),
				Type:  mysql.TypeLong,
				Flag:  model.UnsignedFlag,
			},
			rowcodec.ColInfo{
				ID:            9,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setFlag(types.NewFieldType(mysql.TypeLong), uint(model.UnsignedFlag)),
			},
			uint64(1),
		},
		{
			model.Column{
				Name:  "longlongunsigned",
				Value: uint64(1),
				Type:  mysql.TypeLonglong,
				Flag:  model.UnsignedFlag,
			},
			rowcodec.ColInfo{
				ID:            10,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft: setFlag(
					types.NewFieldType(mysql.TypeLonglong),
					uint(model.UnsignedFlag),
				),
			},
			uint64(1),
		},
	},
	{
		{
			model.Column{Name: "float", Value: float64(3.14), Type: mysql.TypeFloat},
			rowcodec.ColInfo{
				ID:            11,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeFloat),
			},
			float64(3.14),
		},
		{
			model.Column{Name: "double", Value: float64(3.14), Type: mysql.TypeDouble},
			rowcodec.ColInfo{
				ID:            12,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeDouble),
			},
			float64(3.14),
		},
	},
	{
		{
			model.Column{Name: "bit", Value: uint64(683), Type: mysql.TypeBit},
			rowcodec.ColInfo{
				ID:            13,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeBit),
			},
			uint64(683),
		},
	},
	{
		{
			model.Column{Name: "decimal", Value: "129012.1230000", Type: mysql.TypeNewDecimal},
			rowcodec.ColInfo{
				ID:            14,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeNewDecimal),
			},
			"129012.1230000",
		},
	},
	{
		{
			model.Column{Name: "tinytext", Value: []byte("hello world"), Type: mysql.TypeTinyBlob},
			rowcodec.ColInfo{
				ID:            15,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeBlob),
			},
			"hello world",
		},
		{
			model.Column{Name: "mediumtext", Value: []byte("hello world"), Type: mysql.TypeMediumBlob},
			rowcodec.ColInfo{
				ID:            16,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeMediumBlob),
			},
			"hello world",
		},
		{
			model.Column{Name: "text", Value: []byte("hello world"), Type: mysql.TypeBlob},
			rowcodec.ColInfo{
				ID:            17,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeBlob),
			},
			"hello world",
		},
		{
			model.Column{Name: "longtext", Value: []byte("hello world"), Type: mysql.TypeLongBlob},
			rowcodec.ColInfo{
				ID:            18,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeLongBlob),
			},
			"hello world",
		},
		{
			model.Column{Name: "varchar", Value: []byte("hello world"), Type: mysql.TypeVarchar},
			rowcodec.ColInfo{
				ID:            19,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeVarchar),
			},
			"hello world",
		},
		{
			model.Column{Name: "varstring", Value: []byte("hello world"), Type: mysql.TypeVarString},
			rowcodec.ColInfo{
				ID:            20,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeVarString),
			},
			"hello world",
		},
		{
			model.Column{Name: "string", Value: []byte("hello world"), Type: mysql.TypeString},
			rowcodec.ColInfo{
				ID:            21,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeString),
			},
			"hello world",
		},
		{
			model.Column{Name: "json", Value: `{"key": "value"}`, Type: mysql.TypeJSON},
			rowcodec.ColInfo{
				ID:            31,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeJSON),
			},
			`{"key": "value"}`,
		},
	},
	{
		{
			model.Column{
				Name:  "tinyblob",
				Value: []byte("hello world"),
				Type:  mysql.TypeTinyBlob,
				Flag:  model.BinaryFlag,
			},
			rowcodec.ColInfo{
				ID:            22,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setBinChsClnFlag(types.NewFieldType(mysql.TypeTinyBlob)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				Name:  "mediumblob",
				Value: []byte("hello world"),
				Type:  mysql.TypeMediumBlob,
				Flag:  model.BinaryFlag,
			},
			rowcodec.ColInfo{
				ID:            23,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setBinChsClnFlag(types.NewFieldType(mysql.TypeMediumBlob)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				Name:  "blob",
				Value: []byte("hello world"),
				Type:  mysql.TypeBlob,
				Flag:  model.BinaryFlag,
			},
			rowcodec.ColInfo{
				ID:            24,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setBinChsClnFlag(types.NewFieldType(mysql.TypeBlob)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				Name:  "longblob",
				Value: []byte("hello world"),
				Type:  mysql.TypeLongBlob,
				Flag:  model.BinaryFlag,
			},
			rowcodec.ColInfo{
				ID:            25,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setBinChsClnFlag(types.NewFieldType(mysql.TypeLongBlob)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				Name:  "varbinary",
				Value: []byte("hello world"),
				Type:  mysql.TypeVarchar,
				Flag:  model.BinaryFlag,
			},
			rowcodec.ColInfo{
				ID:            26,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setBinChsClnFlag(types.NewFieldType(mysql.TypeVarchar)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				Name:  "varbinary1",
				Value: []byte("hello world"),
				Type:  mysql.TypeVarString,
				Flag:  model.BinaryFlag,
			},
			rowcodec.ColInfo{
				ID:            27,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setBinChsClnFlag(types.NewFieldType(mysql.TypeVarString)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				Name:  "binary",
				Value: []byte("hello world"),
				Type:  mysql.TypeString,
				Flag:  model.BinaryFlag,
			},
			rowcodec.ColInfo{
				ID:            28,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setBinChsClnFlag(types.NewFieldType(mysql.TypeString)),
			},
			"aGVsbG8gd29ybGQ=",
		},
	},
	{
		{
			model.Column{Name: "enum", Value: uint64(1), Type: mysql.TypeEnum},
			rowcodec.ColInfo{
				ID:            29,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setElems(types.NewFieldType(mysql.TypeEnum), []string{"a,", "b"}),
			},
			"a,",
		},
	},
	{
		{
			model.Column{Name: "set", Value: uint64(9), Type: mysql.TypeSet},
			rowcodec.ColInfo{
				ID:            30,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            setElems(types.NewFieldType(mysql.TypeSet), []string{"a", "b", "c", "d"}),
			},
			"a,d",
		},
	},
	{
		{
			model.Column{Name: "date", Value: "2000-01-01", Type: mysql.TypeDate},
			rowcodec.ColInfo{
				ID:            32,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeDate),
			},
			"2000-01-01",
		},
		{
			model.Column{Name: "datetime", Value: "2015-12-20 23:58:58", Type: mysql.TypeDatetime},
			rowcodec.ColInfo{
				ID:            33,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeDatetime),
			},
			"2015-12-20 23:58:58",
		},
		{
			model.Column{Name: "timestamp", Value: "1973-12-30 15:30:00", Type: mysql.TypeTimestamp},
			rowcodec.ColInfo{
				ID:            34,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeTimestamp),
			},
			"1973-12-30 15:30:00",
		},
		{
			model.Column{Name: "time", Value: "23:59:59", Type: mysql.TypeDuration},
			rowcodec.ColInfo{
				ID:            35,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeDuration),
			},
			"23:59:59",
		},
	},
	{
		{
			model.Column{Name: "year", Value: int64(1970), Type: mysql.TypeYear},
			rowcodec.ColInfo{
				ID:            36,
				IsPKHandle:    false,
				VirtualGenCol: false,
				Ft:            types.NewFieldType(mysql.TypeYear),
			},
			int64(1970),
		},
	},
}

func setBinChsClnFlag(ft *types.FieldType) *types.FieldType {
	types.SetBinChsClnFlag(ft)
	return ft
}

func setFlag(ft *types.FieldType, flag uint) *types.FieldType {
	ft.SetFlag(flag)
	return ft
}

func setElems(ft *types.FieldType, elems []string) *types.FieldType {
	ft.SetElems(elems)
	return ft
}

func TestFormatWithQuotes(t *testing.T) {
	csvConfig := &config.CSVConfig{
		Quote: "\"",
	}

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "string does not contain quote mark",
			input:    "a,b,c",
			expected: `"a,b,c"`,
		},
		{
			name:     "string contains quote mark",
			input:    `"a,b,c`,
			expected: `"""a,b,c"`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: `""`,
		},
	}
	for _, tc := range testCases {
		csvMessage := newCSVMessage(csvConfig)
		strBuilder := new(strings.Builder)
		csvMessage.formatWithQuotes(tc.input, strBuilder)
		require.Equal(t, tc.expected, strBuilder.String(), tc.name)
	}
}

func TestFormatWithEscape(t *testing.T) {
	testCases := []struct {
		name      string
		csvConfig *config.CSVConfig
		input     string
		expected  string
	}{
		{
			name:      "string does not contain CR/LF/backslash/delimiter",
			csvConfig: &config.CSVConfig{Delimiter: ","},
			input:     "abcdef",
			expected:  "abcdef",
		},
		{
			name:      "string contains CRLF",
			csvConfig: &config.CSVConfig{Delimiter: ","},
			input:     "abc\r\ndef",
			expected:  "abc\\r\\ndef",
		},
		{
			name:      "string contains backslash",
			csvConfig: &config.CSVConfig{Delimiter: ","},
			input:     `abc\def`,
			expected:  `abc\\def`,
		},
		{
			name:      "string contains a single character delimiter",
			csvConfig: &config.CSVConfig{Delimiter: ","},
			input:     "abc,def",
			expected:  `abc\,def`,
		},
		{
			name:      "string contains multi-character delimiter",
			csvConfig: &config.CSVConfig{Delimiter: "***"},
			input:     "abc***def",
			expected:  `abc\*\*\*def`,
		},
		{
			name:      "string contains CR, LF, backslash and delimiter",
			csvConfig: &config.CSVConfig{Delimiter: "?"},
			input:     `abc\def?ghi\r\n`,
			expected:  `abc\\def\?ghi\\r\\n`,
		},
	}

	for _, tc := range testCases {
		csvMessage := newCSVMessage(tc.csvConfig)
		strBuilder := new(strings.Builder)
		csvMessage.formatWithEscapes(tc.input, strBuilder)
		require.Equal(t, tc.expected, strBuilder.String())
	}
}

func TestCSVMessageEncode(t *testing.T) {
	type fields struct {
		csvConfig  *config.CSVConfig
		opType     string
		tableName  string
		schemaName string
		commitTs   uint64
		columns    []any
	}
	testCases := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "csv encode with typical configurations",
			fields: fields{
				csvConfig: &config.CSVConfig{
					Delimiter:       ",",
					Quote:           "\"",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     insertOperation,
				tableName:  "table1",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{123, "hello,world"},
			},
			want: []byte("\"I\",\"table1\",\"test\",435661838416609281,123,\"hello,world\"\n"),
		},
		{
			name: "csv encode values containing single-character delimter string, without quote mark",
			fields: fields{
				csvConfig: &config.CSVConfig{
					Delimiter:       "!",
					Quote:           "",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     updateOperation,
				tableName:  "table2",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a!b!c", "def"},
			},
			want: []byte(`U!table2!test!435661838416609281!a\!b\!c!def` + "\n"),
		},
		{
			name: "csv encode values containing single-character delimter string, with quote mark",
			fields: fields{
				csvConfig: &config.CSVConfig{
					Delimiter:       ",",
					Quote:           "\"",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     updateOperation,
				tableName:  "table3",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a,b,c", "def", "2022-08-31 17:07:00"},
			},
			want: []byte(`"U","table3","test",435661838416609281,"a,b,c","def","2022-08-31 17:07:00"` + "\n"),
		},
		{
			name: "csv encode values containing multi-character delimiter string, without quote mark",
			fields: fields{
				csvConfig: &config.CSVConfig{
					Delimiter:       "[*]",
					Quote:           "",
					Terminator:      "\r\n",
					NullString:      "\\N",
					IncludeCommitTs: false,
				},
				opType:     deleteOperation,
				tableName:  "table4",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a[*]b[*]c", "def"},
			},
			want: []byte(`D[*]table4[*]test[*]a\[\*\]b\[\*\]c[*]def` + "\r\n"),
		},
		{
			name: "csv encode with values containing multi-character delimiter string, with quote mark",
			fields: fields{
				csvConfig: &config.CSVConfig{
					Delimiter:       "[*]",
					Quote:           "'",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: false,
				},
				opType:     insertOperation,
				tableName:  "table5",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a[*]b[*]c", "def", nil, 12345.678},
			},
			want: []byte(`'I'[*]'table5'[*]'test'[*]'a[*]b[*]c'[*]'def'[*]\N[*]12345.678` + "\n"),
		},
		{
			name: "csv encode with values containing backslash and LF, without quote mark",
			fields: fields{
				csvConfig: &config.CSVConfig{
					Delimiter:       ",",
					Quote:           "",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     updateOperation,
				tableName:  "table6",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a\\b\\c", "def\n"},
			},
			want: []byte(`U,table6,test,435661838416609281,a\\b\\c,def\n` + "\n"),
		},
		{
			name: "csv encode with values containing backslash and CR, with quote mark",
			fields: fields{
				csvConfig: &config.CSVConfig{
					Delimiter:       ",",
					Quote:           "'",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: false,
				},
				opType:     insertOperation,
				tableName:  "table7",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"\\", "\\\r", "\\\\"},
			},
			want: []byte("'I','table7','test','\\','\\\r','\\\\'" + "\n"),
		},
		{
			name: "csv encode with values containing unicode characters",
			fields: fields{
				csvConfig: &config.CSVConfig{
					Delimiter:       "\t",
					Quote:           "\"",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     deleteOperation,
				tableName:  "table8",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a\tb", 123.456, "你好，世界"},
			},
			want: []byte("\"D\"\t\"table8\"\t\"test\"\t435661838416609281\t\"a\tb\"\t123.456\t\"你好，世界\"\n"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &csvMessage{
				csvConfig:  tc.fields.csvConfig,
				opType:     tc.fields.opType,
				tableName:  tc.fields.tableName,
				schemaName: tc.fields.schemaName,
				commitTs:   tc.fields.commitTs,
				columns:    tc.fields.columns,
				newRecord:  true,
			}

			require.Equal(t, tc.want, c.encode())
		})
	}
}

func TestConvertToCSVType(t *testing.T) {
	for _, group := range csvTestColumnsGroup {
		for _, c := range group {
			val, _ := convertToCSVType(&c.col, c.colInfo.Ft)
			require.Equal(t, c.want, val, c.col.Name)
		}
	}
}

func TestBuildRowData(t *testing.T) {
	for idx, group := range csvTestColumnsGroup {
		row := &model.RowChangedEvent{}
		var cols []*model.Column = make([]*model.Column, 0)
		var colInfos []rowcodec.ColInfo = make([]rowcodec.ColInfo, 0)
		for _, c := range group {
			cols = append(cols, &c.col)
			colInfos = append(colInfos, c.colInfo)
		}
		row.ColInfos = colInfos
		row.Table = &model.TableName{
			Table:  fmt.Sprintf("table%d", idx),
			Schema: "test",
		}

		if idx%3 == 0 { // delete operation
			row.PreColumns = cols
		} else if idx%3 == 1 { // insert operation
			row.Columns = cols
		} else { // update operation
			row.PreColumns = cols
			row.Columns = cols
		}
		data, err := buildRowData(&config.CSVConfig{
			Delimiter:       "\t",
			Quote:           "\"",
			Terminator:      "\n",
			NullString:      "\\N",
			IncludeCommitTs: true,
		}, row)
		require.NotNil(t, data)
		require.Nil(t, err)
	}
}
