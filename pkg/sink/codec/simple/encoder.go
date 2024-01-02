// Copyright 2023 PingCAP, Inc.
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

package simple

import (
	"context"
	"os"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/pkg/config"
	cerror "github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/sink/codec"
	"github.com/pingcap/tiflow/pkg/sink/codec/common"
	"github.com/pingcap/tiflow/pkg/sink/kafka/claimcheck"
	"go.uber.org/zap"
)

type encoder struct {
	messages []*common.Message

	config     *common.Config
	claimCheck *claimcheck.ClaimCheck
	marshaller Marshaller
}

// AppendRowChangedEvent implement the RowEventEncoder interface
func (e *encoder) AppendRowChangedEvent(
	ctx context.Context, _ string, event *model.RowChangedEvent, callback func(),
) error {
	var (
		m   interface{}
		err error
	)
	switch e.config.EncodingFormat {
	case common.EncodingFormatJSON:
		m, err = newDMLMessage(event, e.config, false)
	case common.EncodingFormatAvro:
		m, err = newDMLMessageMap(event, e.config, false)
	}
	if err != nil {
		return err
	}

	value, err := e.marshaller.Marshal(m)
	if err != nil {
		return cerror.WrapError(cerror.ErrEncodeFailed, err)
	}

	value, err = common.Compress(e.config.ChangefeedID,
		e.config.LargeMessageHandle.LargeMessageHandleCompression, value)
	if err != nil {
		return err
	}

	result := &common.Message{
		Value:    value,
		Ts:       event.CommitTs,
		Schema:   &event.Table.Schema,
		Table:    &event.Table.Table,
		Type:     model.MessageTypeRow,
		Protocol: config.ProtocolSimple,
		Callback: callback,
	}

	result.IncRowsCount()
	if result.Length() <= e.config.MaxMessageBytes {
		e.messages = append(e.messages, result)
		return nil
	}

	if e.config.LargeMessageHandle.Disabled() {
		log.Error("Single message is too large for simple",
			zap.Int("maxMessageBytes", e.config.MaxMessageBytes),
			zap.Int("length", result.Length()),
			zap.Any("table", event.Table))
		return cerror.ErrMessageTooLarge.GenWithStackByArgs()
	}

	switch e.config.EncodingFormat {
	case common.EncodingFormatJSON:
		m, err = newDMLMessage(event, e.config, true)
		if err != nil {
			return err
		}
		if e.config.LargeMessageHandle.EnableClaimCheck() {
			fileName := claimcheck.NewFileName()
			m.(*message).ClaimCheckLocation = e.claimCheck.FileNameWithPrefix(fileName)
			if err = e.claimCheck.WriteMessage(ctx, result.Key, result.Value, fileName); err != nil {
				return errors.Trace(err)
			}
		}
	case common.EncodingFormatAvro:
		m, err = newDMLMessageMap(event, e.config, true)
		if err != nil {
			return err
		}
		if e.config.LargeMessageHandle.EnableClaimCheck() {
			fileName := claimcheck.NewFileName()
			m.(map[string]interface{})["com.pingcap.simple.avro.DML"].(map[string]interface{})["claimCheckLocation"] = e.claimCheck.FileNameWithPrefix(fileName)
			if err = e.claimCheck.WriteMessage(ctx, result.Key, result.Value, fileName); err != nil {
				return errors.Trace(err)
			}
		}
	}

	value, err = e.marshaller.Marshal(m)
	if err != nil {
		return cerror.WrapError(cerror.ErrEncodeFailed, err)
	}
	value, err = common.Compress(e.config.ChangefeedID,
		e.config.LargeMessageHandle.LargeMessageHandleCompression, value)
	if err != nil {
		return err
	}
	result.Value = value

	if result.Length() <= e.config.MaxMessageBytes {
		log.Warn("Single message is too large for simple, only encode handle key columns",
			zap.Int("maxMessageBytes", e.config.MaxMessageBytes),
			zap.Int("originLength", result.Length()),
			zap.Int("length", result.Length()),
			zap.Any("table", event.Table))
		e.messages = append(e.messages, result)
		return nil
	}

	log.Error("Single message is still too large for simple after only encode handle key columns",
		zap.Int("maxMessageBytes", e.config.MaxMessageBytes),
		zap.Int("length", result.Length()),
		zap.Any("table", event.Table))
	return cerror.ErrMessageTooLarge.GenWithStackByArgs()
}

// Build implement the RowEventEncoder interface
func (e *encoder) Build() []*common.Message {
	if len(e.messages) == 0 {
		return nil
	}
	result := e.messages
	e.messages = nil
	return result
}

// EncodeCheckpointEvent implement the DDLEventBatchEncoder interface
func (e *encoder) EncodeCheckpointEvent(ts uint64) (*common.Message, error) {
	var m interface{}
	switch e.config.EncodingFormat {
	case common.EncodingFormatJSON:
		m = newResolvedMessage(ts)
	case common.EncodingFormatAvro:
		m = newResolvedMessageMap(ts)
	}
	value, err := e.marshaller.Marshal(m)
	if err != nil {
		return nil, cerror.WrapError(cerror.ErrEncodeFailed, err)
	}

	value, err = common.Compress(e.config.ChangefeedID,
		e.config.LargeMessageHandle.LargeMessageHandleCompression, value)
	if err != nil {
		return nil, err
	}
	return common.NewResolvedMsg(config.ProtocolSimple, nil, value, ts), nil
}

// EncodeDDLEvent implement the DDLEventBatchEncoder interface
func (e *encoder) EncodeDDLEvent(event *model.DDLEvent) (*common.Message, error) {
	var (
		m   interface{}
		err error
	)
	switch e.config.EncodingFormat {
	case common.EncodingFormatJSON:
		m, err = newDDLMessage(event)
		if err != nil {
			return nil, err
		}
	case common.EncodingFormatAvro:
		if event.IsBootstrap {
			m = newBootstrapMessageMap(event.TableInfo)
		} else {
			m = newDDLMessageMap(event)
		}
		if m == nil {
			return nil, nil
		}
	}
	value, err := e.marshaller.Marshal(m)
	if err != nil {
		return nil, cerror.WrapError(cerror.ErrEncodeFailed, err)
	}

	value, err = common.Compress(e.config.ChangefeedID,
		e.config.LargeMessageHandle.LargeMessageHandleCompression, value)
	if err != nil {
		return nil, err
	}
	result := common.NewDDLMsg(config.ProtocolSimple, nil, value, event)

	if result.Length() > e.config.MaxMessageBytes {
		log.Error("DDL message is too large for simple",
			zap.Int("maxMessageBytes", e.config.MaxMessageBytes),
			zap.Int("length", result.Length()),
			zap.Any("table", event.TableInfo.TableName))
		return nil, cerror.ErrMessageTooLarge.GenWithStackByArgs()
	}
	return result, nil
}

type builder struct {
	config     *common.Config
	claimCheck *claimcheck.ClaimCheck
	marshaller Marshaller
}

// NewBuilder returns a new builder
func NewBuilder(ctx context.Context, config *common.Config) (*builder, error) {
	var (
		claimCheck *claimcheck.ClaimCheck
		err        error
	)
	if config.LargeMessageHandle.EnableClaimCheck() {
		claimCheck, err = claimcheck.New(ctx,
			config.LargeMessageHandle.ClaimCheckStorageURI, config.ChangefeedID)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	var marshaller Marshaller
	switch config.EncodingFormat {
	case common.EncodingFormatJSON:
		marshaller = newJSONMarshaller()
	case common.EncodingFormatAvro:
		schema, err := os.ReadFile("message.json")
		if err != nil {
			return nil, errors.Trace(err)
		}
		marshaller, err = newAvroMarshaller(string(schema))
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return &builder{
		config:     config,
		claimCheck: claimCheck,
		marshaller: marshaller,
	}, nil
}

// Build implement the RowEventEncoderBuilder interface
func (b *builder) Build() codec.RowEventEncoder {
	return &encoder{
		messages:   make([]*common.Message, 0, 1),
		config:     b.config,
		claimCheck: b.claimCheck,
		marshaller: b.marshaller,
	}
}

// CleanMetrics implement the RowEventEncoderBuilder interface
func (b *builder) CleanMetrics() {
	if b.claimCheck != nil {
		b.claimCheck.CleanMetrics()
	}
}
