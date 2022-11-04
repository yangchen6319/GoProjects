// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compressor

// CompressType type of compressions supported by rpc
type CompressType uint16

const (
	Raw CompressType = iota
	Gzip
	Snappy
	Zlib
)

// Compressor is interface, each compressor has Zip and Unzip functions
type Compressor interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}
