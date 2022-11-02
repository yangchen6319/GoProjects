/*
 * @Author: YangChen
 * @Date: 2022-11-02 16:42:58
 * @Last Modified by: YangChen
 * @Last Modified time: 2022-11-02 22:00:15
 */

package header

import (
	"encoding/binary"
	"sync"

	"github.com/zehuamama/tinyrpc/compressor"
)

const (
	// MaxHeaderSize = 2 + 10 + 10 + 10 + 4(10 refer to MaxVarintLen64)
	MaxHeaderSize = 36
	Uint16Size    = 2
	Uint32Size    = 4
)

// RequestHeader request header structure looks like:
// +--------------+----------------+----------+------------+----------+
// | CompressType |      Method    |    ID    | RequestLen | Checksum |
// +--------------+----------------+----------+------------+----------+
// |    uint16    | uvarint+string |  uvarint |   uvarint  |  uint32  |
// +--------------+----------------+----------+------------+----------+

type RequestHeader struct {
	sync.RWMutex
	CompressType compressor.CompressType
	Method       string
	ID           uint64
	RequestLen   uint32
	Checksum     uint32
}

func (r *RequestHeader) Marshal() []byte {
	//在对r进行操作前，锁住当前线程
	r.RLock()
	// 最后解封当前线程
	defer r.RUnlock()
	idx := 0
	header := make([]byte, MaxHeaderSize+len(r.Method))

	binary.LittleEndian.PutUint16(header[idx:], uint16(r.CompressType))
	idx += Uint16Size
	idx += writeString(header[idx:], r.Method)
	idx += binary.PutUvarint(header[idx:], r.ID)
	idx += binary.PutUvarint(header[idx:], uint64(r.RequestLen))

	binary.LittleEndian.AppendUint32(header[idx:], r.Checksum)
	idx += Uint32Size
	return header
}

// 将s的长度转为varint，将s的二进制形式放入header
func writeString(data []byte, s string) int {
	idx := 0
	idx += binary.PutUvarint(data[idx:], uint64(len(s)))
	copy(data[idx:], s)
	idx += len(s)
	return idx
}
