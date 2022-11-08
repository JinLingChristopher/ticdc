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

package open

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/sinkv2/codec"
	"github.com/pingcap/tiflow/cdc/sinkv2/codec/common"
	"github.com/pingcap/tiflow/cdc/sinkv2/eventsink"

	"github.com/pingcap/tiflow/pkg/config"
	cerror "github.com/pingcap/tiflow/pkg/errors"
	"go.uber.org/zap"
)

// BatchEncoder encodes the events into the byte of a batch into.
type BatchEncoder struct {
	messageBuf   []*common.Message
	callbackBuff []func()
	curBatchSize int

	// configs
	MaxMessageBytes int
	MaxBatchSize    int
}

// AppendRowChangedEvents implements the EventBatchEncoder interface
func (d *BatchEncoder) AppendRowChangedEvents(
	_ context.Context,
	_ string,
	events []*eventsink.RowChangeCallbackableEvent,
) error {
	for _, event := range events {
		keyMsg, valueMsg := rowChangeToMsg(event.Event)
		key, err := keyMsg.Encode()
		if err != nil {
			return errors.Trace(err)
		}
		value, err := valueMsg.encode()
		if err != nil {
			return errors.Trace(err)
		}

		var keyLenByte [8]byte
		binary.BigEndian.PutUint64(keyLenByte[:], uint64(len(key)))
		var valueLenByte [8]byte
		binary.BigEndian.PutUint64(valueLenByte[:], uint64(len(value)))

		// for single message that longer than max-message-size, do not send it.
		// 16 is the length of `keyLenByte` and `valueLenByte`, 8 is the length of `versionHead`
		length := len(key) + len(value) + common.MaxRecordOverhead + 16 + 8
		if length > d.MaxMessageBytes {
			log.Warn("Single message too large",
				zap.Int("max-message-size", d.MaxMessageBytes),
				zap.Int("length", length),
				zap.Any("table", event.Event.Table))
			return cerror.ErrOpenProtocolCodecRowTooLarge.GenWithStackByArgs()
		}

		if len(d.messageBuf) == 0 ||
			d.curBatchSize >= d.MaxBatchSize ||
			d.messageBuf[len(d.messageBuf)-1].Length()+len(key)+len(value)+16 > d.MaxMessageBytes {
			// Before we create a new message, we should handle the previous callbacks.
			d.tryBuildCallback()
			versionHead := make([]byte, 8)
			binary.BigEndian.PutUint64(versionHead, codec.BatchVersion1)
			msg := common.NewMsg(config.ProtocolOpen, versionHead, nil, 0, model.MessageTypeRow, nil, nil)
			d.messageBuf = append(d.messageBuf, msg)
			d.curBatchSize = 0
		}

		message := d.messageBuf[len(d.messageBuf)-1]
		message.Key = append(message.Key, keyLenByte[:]...)
		message.Key = append(message.Key, key...)
		message.Value = append(message.Value, valueLenByte[:]...)
		message.Value = append(message.Value, value...)
		message.Ts = event.Event.CommitTs
		message.Schema = &event.Event.Table.Schema
		message.Table = &event.Event.Table.Table
		message.IncRowsCount()

		if event.Callback != nil {
			d.callbackBuff = append(d.callbackBuff, event.Callback)
		}

		d.curBatchSize++
	}
	return nil
}

// AppendTxnEvent is no-op for ow
func (d *BatchEncoder) AppendTxnEvent(txn *eventsink.TxnCallbackableEvent) error {
	return nil
}

// EncodeDDLEvent implements the EventBatchEncoder interface
func (d *BatchEncoder) EncodeDDLEvent(e *model.DDLEvent) (*common.Message, error) {
	keyMsg, valueMsg := ddlEventToMsg(e)
	key, err := keyMsg.Encode()
	if err != nil {
		return nil, errors.Trace(err)
	}
	value, err := valueMsg.encode()
	if err != nil {
		return nil, errors.Trace(err)
	}

	var keyLenByte [8]byte
	binary.BigEndian.PutUint64(keyLenByte[:], uint64(len(key)))
	var valueLenByte [8]byte
	binary.BigEndian.PutUint64(valueLenByte[:], uint64(len(value)))

	keyBuf := new(bytes.Buffer)
	var versionByte [8]byte
	binary.BigEndian.PutUint64(versionByte[:], codec.BatchVersion1)
	keyBuf.Write(versionByte[:])
	keyBuf.Write(keyLenByte[:])
	keyBuf.Write(key)

	valueBuf := new(bytes.Buffer)
	valueBuf.Write(valueLenByte[:])
	valueBuf.Write(value)

	ret := common.NewDDLMsg(config.ProtocolOpen, keyBuf.Bytes(), valueBuf.Bytes(), e)
	return ret, nil
}

// EncodeCheckpointEvent implements the EventBatchEncoder interface
func (d *BatchEncoder) EncodeCheckpointEvent(ts uint64) (*common.Message, error) {
	keyMsg := newResolvedMessage(ts)
	key, err := keyMsg.Encode()
	if err != nil {
		return nil, errors.Trace(err)
	}

	var keyLenByte [8]byte
	binary.BigEndian.PutUint64(keyLenByte[:], uint64(len(key)))
	var valueLenByte [8]byte
	binary.BigEndian.PutUint64(valueLenByte[:], 0)

	keyBuf := new(bytes.Buffer)
	var versionByte [8]byte
	binary.BigEndian.PutUint64(versionByte[:], codec.BatchVersion1)
	keyBuf.Write(versionByte[:])
	keyBuf.Write(keyLenByte[:])
	keyBuf.Write(key)

	valueBuf := new(bytes.Buffer)
	valueBuf.Write(valueLenByte[:])

	ret := common.NewResolvedMsg(config.ProtocolOpen, keyBuf.Bytes(), valueBuf.Bytes(), ts)
	return ret, nil
}

// Build implements the EventBatchEncoder interface
func (d *BatchEncoder) Build() (messages []*common.Message) {
	d.tryBuildCallback()
	ret := d.messageBuf
	d.messageBuf = make([]*common.Message, 0)
	return ret
}

// tryBuildCallback will collect all the callbacks into one message's callback.
func (d *BatchEncoder) tryBuildCallback() {
	if len(d.messageBuf) != 0 && len(d.callbackBuff) != 0 {
		lastMsg := d.messageBuf[len(d.messageBuf)-1]
		callbacks := d.callbackBuff
		lastMsg.Callback = func() {
			for _, cb := range callbacks {
				cb()
			}
		}
		d.callbackBuff = make([]func(), 0)
	}
}

type batchEncoderBuilder struct {
	config *common.Config
}

// Build a BatchEncoder
func (b *batchEncoderBuilder) Build() codec.EventBatchEncoder {
	encoder := NewBatchEncoder()
	encoder.(*BatchEncoder).MaxMessageBytes = b.config.MaxMessageBytes
	encoder.(*BatchEncoder).MaxBatchSize = b.config.MaxBatchSize

	return encoder
}

// NewBatchEncoderBuilder creates an open-protocol batchEncoderBuilder.
func NewBatchEncoderBuilder(config *common.Config) codec.EncoderBuilder {
	return &batchEncoderBuilder{config: config}
}

// NewBatchEncoder creates a new BatchEncoder.
func NewBatchEncoder() codec.EventBatchEncoder {
	batch := &BatchEncoder{}
	return batch
}
