package main

import (
	"crypto/sha256"
	"bytes"
)

type Block struct {
	PrevHash []byte
	Hash     []byte
	Data     []byte
}

func NewBlock(data []byte, PrevHash []byte) *Block {
	b := Block{
		PrevHash: PrevHash,
		Hash:     nil,
		Data:     data,
	}
	b.SetHash()
	return &b
}

func (b *Block) SetHash() {
	tmp := [][]byte{
		b.PrevHash,
		b.Hash,
		b.Data,
	}
	data := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	b.Hash = hash[:]
}
