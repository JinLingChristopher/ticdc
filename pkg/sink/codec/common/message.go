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

package common

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/pkg/config"
	"github.com/tikv/client-go/v2/oracle"
)

// MaxRecordOverhead is used to calculate message size by sarama kafka client.
// reference: https://github.com/IBM/sarama/blob/
// 66521126c71c522c15a36663ae9cddc2b024c799/async_producer.go#L233
// For TiCDC, minimum supported kafka version is `0.11.0.2`,
// which will be treated as `version = 2` by sarama producer.
const MaxRecordOverhead = 5*binary.MaxVarintLen32 + binary.MaxVarintLen64 + 1

// Message represents an message to the sink
type Message struct {
	Key       []byte
	Value     []byte
	Ts        uint64 // reserved for possible output sorting
	TableID   int64
	HandleKey []string
	Schema    *string           // schema
	Table     *string           // table
	Type      model.MessageType // type
	Protocol  config.Protocol   // protocol
	rowsCount int               // rows in one Message
	Callback  func()            // Callback function will be called when the message is sent to the sink.

	// PartitionKey for pulsar, route messages to one or different partitions
	PartitionKey *string
}

func (m *Message) String() string {
	return fmt.Sprintf("Message{Ts: %d, Type: %d, Schema: %s, Table: %s, Key: %s, Value: %s, Protocol: %s}",
		m.Ts, m.Type, m.GetSchema(), m.GetTable(), m.Key, m.Value, m.Protocol)
}

// Length returns the expected size of the Kafka message
func (m *Message) Length() int {
	headers := m.Headers()
	var headerLen int
	for _, header := range headers {
		headerLen += len(header.Key) + len(header.Value) + 2*binary.MaxVarintLen32
	}
	return headerLen + len(m.Key) + len(m.Value) + MaxRecordOverhead
}

// Headers returns the headers of Kafka message
func (m *Message) Headers() []sarama.RecordHeader {
	var handleKey bytes.Buffer
	for idx, key := range m.HandleKey {
		handleKey.WriteString(key)
		if idx != len(m.HandleKey)-1 {
			handleKey.WriteString(",")
		}
	}
	handle := handleKey.Bytes()
	headers := []sarama.RecordHeader{
		{Key: []byte("tableID"), Value: []byte(strconv.FormatUint(uint64(m.TableID), 10))},
		{Key: []byte("commitTs"), Value: []byte(strconv.FormatUint(m.Ts, 10))},
		{Key: []byte("handleKey"), Value: handle},
	}
	return headers
}

// PhysicalTime returns physical time part of Ts in time.Time
func (m *Message) PhysicalTime() time.Time {
	return oracle.GetTimeFromTS(m.Ts)
}

// GetRowsCount returns the number of rows batched in one Message
func (m *Message) GetRowsCount() int {
	return m.rowsCount
}

// SetRowsCount set the number of rows
func (m *Message) SetRowsCount(cnt int) {
	m.rowsCount = cnt
}

// IncRowsCount increase the number of rows
func (m *Message) IncRowsCount() {
	m.rowsCount++
}

// GetSchema returns schema string
func (m *Message) GetSchema() string {
	if m.Schema == nil {
		return ""
	}
	return *m.Schema
}

// GetTable returns the Table string
func (m *Message) GetTable() string {
	if m.Table == nil {
		return ""
	}
	return *m.Table
}

// SetPartitionKey sets the PartitionKey for a message
// PartitionKey is used for pulsar producer, route messages to one or different partitions
func (m *Message) SetPartitionKey(key string) {
	m.PartitionKey = &key
}

// GetPartitionKey returns the GetPartitionKey
func (m *Message) GetPartitionKey() string {
	if m.PartitionKey == nil {
		return ""
	}
	return *m.PartitionKey
}

// NewDDLMsg creates a DDL message.
func NewDDLMsg(proto config.Protocol, key, value []byte, event *model.DDLEvent) *Message {
	return NewMsg(
		proto,
		key,
		value,
		event.CommitTs,
		model.MessageTypeDDL,
		&event.TableInfo.TableName.Schema,
		&event.TableInfo.TableName.Table,
	)
}

// NewResolvedMsg creates a resolved ts message.
func NewResolvedMsg(proto config.Protocol, key, value []byte, ts uint64) *Message {
	return NewMsg(proto, key, value, ts, model.MessageTypeResolved, nil, nil)
}

// NewMsg should be used when creating a Message struct.
// It copies the input byte slices to avoid any surprises in asynchronous MQ writes.
func NewMsg(
	proto config.Protocol,
	key []byte,
	value []byte,
	ts uint64,
	ty model.MessageType,
	schema, table *string,
) *Message {
	ret := &Message{
		Key:       nil,
		Value:     nil,
		Ts:        ts,
		Schema:    schema,
		Table:     table,
		Type:      ty,
		Protocol:  proto,
		rowsCount: 0,
	}

	if key != nil {
		ret.Key = make([]byte, len(key))
		copy(ret.Key, key)
	}

	if value != nil {
		ret.Value = make([]byte, len(value))
		copy(ret.Value, value)
	}

	return ret
}

// ClaimCheckMessage is the message sent to the claim-check external storage.
type ClaimCheckMessage struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

// UnmarshalClaimCheckMessage unmarshal bytes to ClaimCheckMessage.
func UnmarshalClaimCheckMessage(data []byte) (*ClaimCheckMessage, error) {
	var m ClaimCheckMessage
	err := json.Unmarshal(data, &m)
	return &m, err
}
