package common

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"github.com/snluu/uuid"
	"io/ioutil"
	"log"
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

func UUID() string {
	return uuid.Rand().Hex()
}


func SaveFile(dir string, data []byte) bool {
	err := ioutil.WriteFile(dir,data,0666)
	log.Println("save file :",dir)
	if err != nil{
		log.Fatalln(err)
		return false
	}
	return true
}
