/*
 * @Author: YangChen
 * @Date: 2022-11-02 16:42:58
 * @Last Modified by: YangChen
 * @Last Modified time: 2022-11-04 14:50:42
 */

package header

import (
	"encoding/binary"
	"errors"
	"github.com/yangchen6319/GoProjects/tinyRpc/compressor"
	"sync"
)

const (
	// MaxHeaderSize = 2 + 10 + 10 + 10 + 4(10 refer to MaxVarintLen64)
	MaxHeaderSize = 36
	Uint16Size    = 2
	Uint32Size    = 4
)

var ErrorUnmarshal = errors.New("a error occurred in unmarshal")

type CompressType uint16

// RequestHeader request header structure looks like:
// +--------------+----------------+----------+------------+----------+
// | CompressType |      Method    |    ID    | RequestLen | Checksum |
// +--------------+----------------+----------+------------+----------+
// |    uint16    | uvarint+string |  uvarint |   uvarint  |  uint32  |
// +--------------+----------------+----------+------------+----------+

type RequestHeader struct {
	sync.RWMutex
	CompressType CompressType
	Method       string
	ID           uint64
	RequestLen   uint32
	Checksum     uint32
}

// 请求头编码
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

	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size
	return header[:idx]
}

// 请求头解码
func (r *RequestHeader) Unmarshal(data []byte) (err error) {
	r.Lock()
	defer r.Unlock()
	if len(data) == 0 {
		return ErrorUnmarshal
	}
	// 当前函数出现异常的处理
	defer func() {
		if r := recover(); r != nil {
			err = ErrorUnmarshal
		}
	}()
	// 解码data
	index, size := 0, 0
	r.CompressType = CompressType(binary.LittleEndian.Uint16(data[index:]))
	index += Uint16Size

	r.Method, size = readString(data[index:])
	index += size

	r.ID, size = binary.Uvarint(data[index:])
	index += size

	context, size := binary.Uvarint(data[index:])
	r.RequestLen = uint32(context)
	index += size

	r.Checksum = binary.LittleEndian.Uint32(data[index:])

	return

}

// 得到compressType
func (r *RequestHeader) GetCompressType() compressor.CompressType {
	r.RLock()
	defer r.RUnlock()
	return compressor.CompressType(r.CompressType)
}

// 重置请求头
func (r *RequestHeader) ResetHeader() {
	r.Lock()
	defer r.Unlock()
	r.CompressType = 0
	r.Method = ""
	r.ID = 0
	r.RequestLen = 0
	r.Checksum = 0
}

// ResponseHeader request header structure looks like:
// +--------------+---------+----------------+-------------+----------+
// | CompressType |    ID   |      Error     | ResponseLen | Checksum |
// +--------------+---------+----------------+-------------+----------+
// |    uint16    | uvarint | uvarint+string |    uvarint  |  uint32  |
// +--------------+---------+----------------+-------------+----------+
type ResponseHeader struct {
	sync.RWMutex
	CompressType CompressType
	ID           uint64
	Error        string
	ResponseLen  uint32
	Checksum     uint32
}

// 响应头编码
func (r *ResponseHeader) Marshal() []byte {
	//在对r进行操作前，锁住当前线程
	r.RLock()
	// 最后解封当前线程
	defer r.RUnlock()
	idx := 0
	header := make([]byte, MaxHeaderSize+len(r.Error))

	binary.LittleEndian.PutUint16(header[idx:], uint16(r.CompressType))
	idx += Uint16Size
	idx += binary.PutUvarint(header[idx:], r.ID)
	idx += writeString(header[idx:], r.Error)
	idx += binary.PutUvarint(header[idx:], uint64(r.ResponseLen))

	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size
	return header[:idx]
}

// 请求头解码
func (r *ResponseHeader) Unmarshal(data []byte) (err error) {
	r.Lock()
	defer r.Unlock()
	if len(data) == 0 {
		return ErrorUnmarshal
	}
	// 当前函数出现异常的处理
	defer func() {
		if r := recover(); r != nil {
			err = ErrorUnmarshal
		}
	}()
	// 解码data
	index, size := 0, 0
	r.CompressType = CompressType(binary.LittleEndian.Uint16(data[index:]))
	index += Uint16Size

	r.ID, size = binary.Uvarint(data[index:])
	index += size

	r.Error, size = readString(data[index:])
	index += size

	context, size := binary.Uvarint(data[index:])
	r.ResponseLen = uint32(context)
	index += size

	r.Checksum = binary.LittleEndian.Uint32(data[index:])

	return

}

// 得到compressType
func (r *ResponseHeader) GetCompressType() compressor.CompressType {
	r.RLock()
	defer r.RUnlock()
	return compressor.CompressType(r.CompressType)
}

// 重置请求头
func (r *ResponseHeader) ResetHeader() {
	r.Lock()
	defer r.Unlock()
	r.CompressType = 0
	r.Error = ""
	r.ID = 0
	r.ResponseLen = 0
	r.Checksum = 0
}

func readString(data []byte) (string, int) {
	index := 0
	length, size := binary.Uvarint(data[index:])
	index += size
	s := string(data[index : index+int(length)])
	index += len(s)
	return s, index
}

// 将s的长度转为varint，将s的二进制形式放入header
func writeString(data []byte, s string) int {
	idx := 0
	idx += binary.PutUvarint(data[idx:], uint64(len(s)))
	copy(data[idx:], s)
	idx += len(s)
	return idx
}
