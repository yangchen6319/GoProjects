/*
 * @Author: YangChen
 * @Date: 2022-11-04 16:02:37
 * @Last Modified by: YangChen
 * @Last Modified time: 2022-11-04 16:03:03
 */

// 写个消息对象池
package header

import "sync"

var (
	RequestPool  sync.Pool
	ResponsePool sync.Pool
)

func init() {
	RequestPool = sync.Pool{New: func() any {
		return &RequestHeader{}
	}}

	ResponsePool = sync.Pool{New: func() any {
		return &ResponseHeader{}
	}}
}
