package main

import (
	"bytes"
	"encoding/binary"
	"github.com/micro/go-log"
)

func uintToByte(num uint64) []byte {
	var buffer bytes.Buffer
	//使用二进制编码
	err := binary.Write(&buffer, binary.LittleEndian, num)
	if err != nil {
		log.Fatal("binary write error:", err)
	}
	//返回byte切片
	return buffer.Bytes()
}

