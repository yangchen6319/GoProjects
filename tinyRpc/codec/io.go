/*
 * @Author: YangChen
 * @Date: 2022-11-04 16:27:39
 * @Last Modified by: YangChen
 * @Last Modified time: 2022-11-04 19:48:45
 */

package codec

import (
	"encoding/binary"
	"io"
	"net"
)

func sendFrame(w io.Writer, data []byte) (err error) {
	var size [binary.MaxVarintLen64]byte
	if data == nil || len(data) == 0 {
		n := binary.PutUvarint(size[:], uint64(0))
		if err = write(w, size[:n]); err != nil {
			return
		}
		return
	}

	n := binary.PutUvarint(size[:], uint64(len(data)))
	if err = write(w, size[:n]); err != nil {
		return
	}
	if err = write(w, data[:]); err != nil {
		return
	}
	return
}

func write(w io.Writer, data []byte) error {
	for index := 0; index < len(data); {
		n, err := w.Write(data[index:])
		if _, ok := err.(net.Error); !ok {
			return err
		}
		index += n
	}
	return nil
}
