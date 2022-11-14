/*
 * @Author: YangChen
 * @Date: 2022-11-05 14:49:11
 * @Last Modified by: YangChen
 * @Last Modified time: 2022-11-05 14:51:46
 */

package serializer

type serializer interface {
	Marshal(message interface{}) ([]byte, error)
	Unmarshal(data []byte, message interface{}) error
}
