package main

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

type POW struct {
	//​区块：block
	Block *Block
	//目标值：target
	Target *big.Int
}

func NewPOW(block *Block) *POW {
	pow := POW{
		Block: block,
	}
	targetstr := "0001000000000000000000000000000000000000000000000000000000000000"
	targetInt := new(big.Int)
	targetInt.SetString(targetstr, 16)
	pow.Target = targetInt
	return &pow
}

//挖矿函数
func (p *POW) Run() ([]byte, uint64) {
	var nonce uint64
	var hash []byte
	fmt.Println("开始挖矿")
	for {
		fmt.Printf("%x\r", hash)
		hash = p.PrepareData(nonce)
		target := sha256.Sum256(hash)
		var targetInt = new(big.Int)
		targetInt.SetBytes(target[:])
		if targetInt.Cmp(p.Target) == -1 {
			fmt.Println("挖矿成功")
			break
		}
		nonce++
	}
	return hash, nonce
}

func (p *POW) PrepareData(nonce uint64) []byte {
	b := p.Block
	tmp := [][]byte{
		uintToByte(b.Version),
		b.MerkleRoot,
		uintToByte(b.TimeStamp),
		uintToByte(b.Bits),
		uintToByte(nonce),
		b.PrevHash,
	}
	data := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	return hash[:]
}

func (p *POW) IsValid() bool {
	data := p.PrepareData(p.Block.Nonce)
	hash := sha256.Sum256(data)
	tempInt := new(big.Int)
	tempInt.SetBytes(hash[:])
	return tempInt.Cmp(p.Target)==-1
}
