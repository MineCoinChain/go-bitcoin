package main

import (
	"bytes"
	"time"
	"encoding/gob"
	"log"
)

type Block struct {
	//版本号
	Version uint64
	//交易的根哈希值
	MerkleRoot []byte
	//时间戳
	TimeStamp uint64
	//难度值, 系统提供一个数据，用于计算出一个哈希值
	Bits uint64
	//随机数，挖矿要求的数值
	Nonce uint64
	// 前区块哈希
	PrevHash []byte
	// 哈希, 为了方便，我们将当前区块的哈希放入Block中
	Hash []byte
	//交易数据
	Data []byte
}

func NewBlock(data string, PrevHash []byte) *Block {
	b := Block{
		Version:    0,
		MerkleRoot: nil,
		TimeStamp:  uint64(time.Now().Unix()),
		Bits:       0,
		Nonce:      0,
		PrevHash:   PrevHash,
		Hash:       nil,
		Data:       []byte(data),
	}
	var pow = NewPOW(&b)
	hash, nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce
	return &b
}

func (b *Block)Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if err != nil {
		log.Fatal("encode err:", err)
	}
	return buffer.Bytes()
}

func Deserialize(src []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(src))
	err := decoder.Decode(&block)
	if err!=nil{
		log.Fatal("Deserialize err1:",err)
	}
	return &block
}
