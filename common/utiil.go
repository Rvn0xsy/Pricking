package common

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
)

// UnGzipData Gzip 解压
func UnGzipData(data []byte) []byte {
	b := new(bytes.Buffer)
	_ = binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return data
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			return data
		}
		return undatas
	}
}
