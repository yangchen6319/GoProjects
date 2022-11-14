/*
 * @Author: YangChen
 * @Date: 2022-11-05 14:38:06
 * @Last Modified by: YangChen
 * @Last Modified time: 2022-11-05 14:46:45
 */

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

// 写个map以将四种压缩器和压缩类型对应起来
var compressors = map[CompressType]Compressor{
	Raw:    RawCompressor{},
	Gzip:   GzipCompressor{},
	Snappy: SnappyCompressor{},
	Zlib:   ZlibCompressor{},
}
