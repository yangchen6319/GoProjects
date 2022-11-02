/*
 * @Author: YangChen
 * @Date: 2022-11-02 16:42:58
 * @Last Modified by: YangChen
 * @Last Modified time: 2022-11-02 17:05:13
 */

package header

import (
	"sync"

	"github.com/zehuamama/tinyrpc/compressor"
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
