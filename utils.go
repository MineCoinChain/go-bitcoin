package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
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

func IsFileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	//根据数据中的error值进行判断，不存在返回false，存在返回true
	return !os.IsNotExist(err)
}
