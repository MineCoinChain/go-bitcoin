package main

import (
	"crypto/sha256"
	"bytes"
	"time"
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

func NewBlock(data []byte, PrevHash []byte) *Block {
	b := Block{
		Version:    0,
		MerkleRoot: nil,
		TimeStamp:  uint64(time.Now().Unix()),
		Bits:       0,
		Nonce:      0,
		PrevHash:   PrevHash,
		Hash:       nil,
		Data:       data,
	}
	b.SetHash()
	return &b
}

func (b *Block) SetHash() {
	tmp := [][]byte{
		uintToByte(b.Version),
		b.MerkleRoot,
		uintToByte(b.TimeStamp),
		uintToByte(b.Bits),
		uintToByte(b.Nonce),
		b.PrevHash,
		b.Hash,
		b.Data,
	}
	data := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	b.Hash = hash[:]
}
