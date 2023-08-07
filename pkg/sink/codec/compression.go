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

package codec

import (
	"bytes"
	"compress/gzip"

	snappy "github.com/eapache/go-xerial-snappy"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/pingcap/errors"
)

type CompressionCodec uint8

const (
	// CompressionNone no compression
	CompressionNone CompressionCodec = iota
	// CompressionGZIP compression using GZIP
	CompressionGZIP
	// CompressionSnappy compression using snappy
	CompressionSnappy
	// CompressionLZ4 compression using LZ4
	CompressionLZ4
	// CompressionZSTD compression using ZSTD
	CompressionZSTD
)

var (
	compressionNames = []string{"none", "gzip", "snappy", "lz4", "zstd"}
)

func (c CompressionCodec) String() string {
	return compressionNames[c]
}

func Compress(cc CompressionCodec, data []byte) ([]byte, error) {
	switch cc {
	case CompressionNone:
		return data, nil
	case CompressionGZIP:
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		if _, err := writer.Write(data); err != nil {
			return nil, errors.Trace(err)
		}
		if err := writer.Close(); err != nil {
			return nil, errors.Trace(err)
		}
		return buf.Bytes(), nil
	case CompressionSnappy:
		return snappy.Encode(data), nil
	case CompressionLZ4:
		var buf bytes.Buffer
		writer := lz4.NewWriter(&buf)
		if _, err := writer.Write(data); err != nil {
			return nil, errors.Trace(err)
		}
		if err := writer.Close(); err != nil {
			return nil, errors.Trace(err)
		}
		return buf.Bytes(), nil
	case CompressionZSTD:
		encoder, err := zstd.NewWriter(nil, zstd.WithZeroFrames(true),
			zstd.WithEncoderLevel(zstd.SpeedDefault),
			zstd.WithEncoderConcurrency(1))
		if err != nil {
			return nil, errors.Trace(err)
		}

		out := encoder.EncodeAll(data, nil)
		return out, nil
	default:
	}
	return nil, errors.New("unsupported compression codec")
}
